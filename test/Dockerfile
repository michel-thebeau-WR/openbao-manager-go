#
# Copyright (c) 2025 Wind River Systems, Inc.
#
# SPDX-License-Identifier: Apache-2.0
#

# Test image for go openbao monitor
# Debian image with openbao and go openbao monitor installed
# Example command to build: docker build -f test/Dockerfile .

FROM debian:stable-slim

ENV PKG_LIST="mawk bash coreutils curl grep sed jq uuid-runtime golang \
    wget procps less file vim bash-completion gzip tar sudo"

USER root

# install packages
RUN set -ex; \
    apt-get update && apt-get install -y $PKG_LIST \
    && apt-get clean && rm -r /var/lib/apt/lists/*

# Download openbao
RUN mkdir -p /tmp && \
    wget -P /tmp/ https://github.com/openbao/openbao/releases/download/v2.1.0/bao_2.1.0_linux_amd64.deb && \
    dpkg -i /tmp/bao_2.1.0_linux_amd64.deb

# Copy over go openbao monitor
# Use the output folder used in the build
COPY ./bin/baomon /usr/bin/

# create a non-root user/group for openbao-manager
RUN groupadd --gid 1000 manager \
    && adduser --uid 1000 --gid 1000 manager \
        --home /workdir --shell /bin/bash

RUN echo "root:root" | chpasswd
RUN echo "manager:manager" | chpasswd
RUN usermod -aG sudo manager

USER manager

# Copy and run auto-completion
COPY ./bin/baomon.completion /usr/share/bash-completion/completions/baomon

RUN echo "source /etc/bash_completion" >> /workdir/.bashrc

# Copy testing files
COPY ./bin/OpenBaoCA/ /workdir/OpenBaoCA/
COPY ./bin/OpenBaoServerCert/ /workdir/OpenBaoServerCert/
COPY ./test/* /workdir/

CMD ["bash"]
