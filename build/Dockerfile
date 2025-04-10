#
# Copyright (c) 2025 Wind River Systems, Inc.
#
# SPDX-License-Identifier: Apache-2.0
#

# Simple environment image to build the go openbao monitor
# The build result should be in the output folder
# Example: docker build --target bin --output bin/ -f ./build/Dockerfile .

FROM golang:1.24.1-bookworm AS build

WORKDIR /src

# copy source code
COPY . ./

RUN go mod download

RUN go mod tidy

RUN go build -o /out/baomon .

# save the module sources
RUN mkdir /src/mod; GOPATH=/src/mod go get; \
    tar zcf /out/mod.src.tar.gz /src/mod

# create auto-completion file
RUN mkdir /workdir; cp ./test/testConfig.yaml /workdir/

# copy the ca files briefly so baomon doesn't error out during
# generating completion script. Will be deleted along with the image after.
COPY ./bin/OpenBaoCA/ /workdir/OpenBaoCA/
COPY ./bin/OpenBaoServerCert/ /workdir/OpenBaoServerCert/

RUN /out/baomon completion bash > /out/baomon.completion

FROM scratch AS bin

COPY --from=build /out/* /