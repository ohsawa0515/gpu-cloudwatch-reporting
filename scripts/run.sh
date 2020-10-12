#!/bin/bash

version=$1
region_param="${REGION:-us-east-1}"
namespace_param="${NAMESPACE:-GPUMonitor}"
send_interval_second_param="${SEND_INTERVAL_SECOND:-60s}"
collect_interval_second_param="${COLLECT_INTERVAL_SECOND:-5s}"

docker run -d --gpus=all --rm \
    -e REGION=${region_param} \
    -e NAMESPACE=${namespace_param} \
    -e SEND_INTERVAL_SECOND=${send_interval_second_param} \
    -e COLLECT_INTERVAL_SECOND=${collect_interval_second_param} \
    ohsawa0515/gpu-cloudwatch-reporting:$version
