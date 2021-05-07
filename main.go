// When compiled for an Alpine container use
// CGO_ENABLED=0 go build

package main

import (
	"log"
	"syscall"
	"strings"
	"flag"
	"github.com/fsnotify/fsnotify"
	api "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type vfioInstance struct {
	devicePlugin *vfioDevicePlugin
	resourceName string
	iommuGroups []string
	socketName string
}

func main() {
	var instances []vfioInstance
	var configFile string

	flag.StringVar(&configFile, "config", "/root/config/config.yml", "path to the configuration file")
	flag.Parse()

	log.Print("Starting VFIO device plugin for Kubernetes")
	config := readConfigFile(configFile)
	devices := scanDevices()
	groups := groupDevices(devices, config)

	for _,group := range groups {
		var instance vfioInstance
		instance.devicePlugin = nil
		instance.iommuGroups = group.iommuGroups
		instance.resourceName = group.resourceName
		instance.socketName = api.DevicePluginPath + strings.ReplaceAll(group.resourceName, "/", "-") + ".sock"
		instances = append(instances, instance)
	}

	log.Print("Starting new FS watcher")
	watcher, err := newFSWatcher(api.DevicePluginPath)
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	log.Print("Starting new OS watcher")
	sigs := newOSWatcher(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	restart := true

L:
	for {
		if restart {
			var err error

			for _,instance := range instances {
				if instance.devicePlugin != nil {
					instance.devicePlugin.Stop()
				}
			}

			for _,instance := range instances {
				instance.devicePlugin = NewDevicePlugin(instance.iommuGroups, instance.resourceName, instance.socketName)
				err = instance.devicePlugin.Serve()
				if err != nil {
					log.Print("Failed to contact Kubelet, retrying")
					break
				}
			}

			if err != nil {
				continue
			}

			restart = false
		}

		select {
			case event := <-watcher.Events:
				if (event.Name == api.KubeletSocket) && (event.Op & fsnotify.Create) == fsnotify.Create {
					log.Printf("inotify: %s created, restarting", api.KubeletSocket)
					restart = true
				}
			case err := <-watcher.Errors:
				log.Printf("inotify: %s", err)
			case s := <-sigs:
				switch s {
					case syscall.SIGHUP:
						log.Print("Received SIGHUP, restarting.")
						restart = true
					default:
						log.Printf("Received signal '%v', shutting down", s)
						for _,instance := range instances {
							if instance.devicePlugin != nil {
								instance.devicePlugin.Stop()
							}
						}
						break L
				}
		}
	}
}
