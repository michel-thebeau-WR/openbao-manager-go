//
// Copyright (c) 2025 Wind River Systems, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//

module github.com/michel-thebeau-WR/openbao-manager-go/baomon

go 1.23.0

toolchain go1.24.2

replace github.com/michel-thebeau-WR/openbao-manager-go/baomon/config => ./config

replace github.com/michel-thebeau-WR/openbao-manager-go/baomon/commands => ./commands

require (
	github.com/michel-thebeau-WR/openbao-manager-go/baomon/commands v0.0.0-00010101000000-000000000000
	github.com/michel-thebeau-WR/openbao-manager-go/baomon/config v0.0.0-00010101000000-000000000000 // indirect
)

require (
	github.com/go-yaml/yaml v2.1.0+incompatible // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/spf13/cobra v1.9.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
)
