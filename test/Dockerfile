# Test image for go openbao monitor
# Debian image with openbao and go openbao monitor installed
# Example command to build: docker build -f test/Dockerfile .

FROM debian:stable-slim

ENV PKG_LIST="mawk bash coreutils curl grep sed jq uuid-runtime golang \
    wget procps less file vim"

USER root

# install packages
RUN set -ex; \
    apt-get update && apt-get install -y $PKG_LIST \
    && apt-get clean && rm -r /var/lib/apt/lists/*

# Download openbao
RUN mkdir -p /tmp && \
    wget -P /tmp/ https://github.com/openbao/openbao/releases/download/v2.1.0/bao_2.1.0_linux_amd64.deb && \
    dpkg -i /tmp/bao_2.1.0_linux_amd64.deb

# create a non-root user/group for openbao-manager
RUN groupadd --gid 1000 manager \
    && adduser --uid 1000 --gid 1000 manager \
        --home /workdir --shell /bin/bash

USER manager

# Copy over go openbao monitor
# Use the output folder used in the build
COPY ./bin/baomon .

CMD ["bash"]
