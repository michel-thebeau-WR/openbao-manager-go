#!/bin/bash

# unittest should already have created rootca.crt
text="$(
    curl --cacert /rootca.crt \
        https://Openbao:8200/v1/sys/health \
        2>/tmp/unittest_health.err \
    | jq . )"


