# Kubernetes VFIO Device Plugin
This plugin is used to present VFIO devices as allocatable resources in Kubernetes. While the plugin was written specifically to allow GPU passthrough for [Kata Containers](https://katacontainers.io), in practice the plugin is completely generic and it can handle any device controlled by the `vfio-pci` driver on the host machine. The configuration of devices via the `vfio-pci` driver must be done manually before starting the plugin.


## Credits
This plugin is a heavily modified version of the [smarter-device-plugin](https://gitlab.com/arm-research/smarter/smarter-device-manager) by the ARM Research group. Especially the gRPC part of the code was largely adopted vanilla from this project.


## Configuration
In the `manifest` directory there is a YAML file for starting the plugin as a container in Kubernetes. This manifest also contains a configration file for the plugin, where you can specify what VFIO devices you want to register as a resource.

```yaml
- resourceName: nvidia/a100
  vendorId: 10de
  deviceId: ['20f1']
- resourceName: nvidia/v100
  vendorId: 10de
  deviceId: ['1db5']
```

The `resourceName` is the name of the resource in Kubernetes and it must be of the format `xxxx/yyyy`. The next two variables are the PCI IDs of the devices you want to register as part of this resource. You can only specify a single vendor ID, but multiple device IDs.
