# gpu-cloudwatch-reporting

This repository provides a tool that sends metrics on GPU utilization on Amazon ECS to CloudWatch.
This tool is able to supports Linux only.

* Ubuntu 16.04/18.04/20.04 LTS
* Amazon Linux 2 (ECS GPU-optimized AMI)

## Installation

### Docker pull

You can get Docker image from [GitHub Container Registry](https://github.com/users/ohsawa0515/packages/container/package/gpu-cloudwatch-reporting)

```console
$ docker pull ghcr.io/ohsawa0515/gpu-cloudwatch-reporting:<TAG_NAME>
```

### Download binary

Download it from [releases page](https://github.com/ohsawa0515/gpu-cloudwatch-reporting/releases) and extract it to `/usr/local/bin`.

```console
$ curl -L -O https://github.com/ohsawa0515/gpu-cloudwatch-reporting/releases/download/<version>/gpu-cloudwatch-reporting_linux_amd64.tar.gz
$ tar zxf gpu-cloudwatch-reporting_linux_amd64.tar.gz
$ mv ./gpu-cloudwatch-reporting /usr/local/bin/
$ chmod +x /usr/local/bin/gpu-cloudwatch-reporting
```

### go get

```console
$ go get github.com/ohsawa0515/gpu-cloudwatch-reporting
$ mv $GOPATH/gpu-cloudwatch-reporting /usr/local/bin/
$ chmod +x /usr/local/bin/gpu-cloudwatch-reporting
```

## Run as systemd

```console
$ cat <<-EOH > /lib/systemd/system/gpu-cloudwatch-reporting.service
[Unit]
Description=GPU Utilization Metric Reporting
[Service]
Type=simple
PIDFile=/run/gpu-cloudwatch-reporting.pid
ExecStart=/usr/local/bin/gpu-cloudwatch-reporting
User=root
Group=root
WorkingDirectory=/
Restart=always
[Install]
WantedBy=multi-user.target
EOH
$ systemctl daemon-reload
$ systemctl enable gpu-cloudwatch-reporting.service
$ systemctl start gpu-cloudwatch-reporting.service
```

## Run as docker container

NVIDIA driver is required. Please install from [here](https://github.com/NVIDIA/nvidia-docker#quickstart).

```console
$ docker pull ghcr.io/ohsawa0515/gpu-cloudwatch-reporting:<TAG_NAME>
$ docker run -d --gpus=all --rm \
      -e REGION=us-east-1 \
      -e NAMESPACE=GPUMonitor \
      -e SEND_INTERVAL_SECOND=60s \
      -e COLLECT_INTERVAL_SECOND=5s \
      ghcr.io/ohsawa0515/gpu-cloudwatch-reporting:<TAG_NAME>
```

