FROM alpine
LABEL maintainer="hansen@imada.sdu.dk"

RUN apk update
RUN apk upgrade
COPY vfio-device-plugin /usr/bin/vfio-device-plugin

ENTRYPOINT /usr/bin/vfio-device-plugin
