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

# create auto-completion file
RUN mkdir /workdir; cp ./test/testConfig.yaml /workdir/; \
    /out/baomon completion bash > /out/baomon.completion

FROM scratch AS bin

COPY --from=build /out/* /