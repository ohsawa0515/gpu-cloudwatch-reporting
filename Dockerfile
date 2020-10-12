FROM nvidia/cuda:10.0-base
USER root

ARG VERSION
ENV PATH $PATH:/work

RUN mkdir /work
WORKDIR /work

RUN apt-get update \
    && apt-get install -y curl \
    && curl -L -O https://github.com/ohsawa0515/gpu-cloudwatch-reporting/releases/download/${VERSION}/gpu-cloudwatch-reporting_linux_amd64.tar.gz \
    && tar zxf ./gpu-cloudwatch-reporting_linux_amd64.tar.gz \
    && chmod +x /work/gpu-cloudwatch-reporting

CMD ["/work/gpu-cloudwatch-reporting"]