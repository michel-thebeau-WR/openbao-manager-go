# Simple environment image to build the go openbao monitor
# The build result should be in the output folder
# Example: docker build --target bin --output bin/ -f ./build/Dockerfile .

FROM golang:1.24.1-bookworm AS build

WORKDIR /src

# Copy submodules
COPY ./config ./config

# Copy main module
COPY go.mod go.sum main.go ./

RUN go mod download

RUN go mod tidy

RUN go build -o /out/baomon .

FROM scratch AS bin

COPY --from=build /out/baomon /