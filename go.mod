module github.com/michel-thebeau-WR/openbao-manager-go/baomon

go 1.18

replace github.com/michel-thebeau-WR/openbao-manager-go/baomon/config => ./config

require github.com/michel-thebeau-WR/openbao-manager-go/baomon/config v0.0.0-00010101000000-000000000000

require (
	github.com/go-yaml/yaml v2.1.0+incompatible // indirect
	github.com/kr/pretty v0.3.1 // indirect
)
