# openbao-manager-go
Golang code for the Openbao Manager

Manager runs as a Kubernetes pod or host service and performs the following actions:
* Discover Openbao servers
* Initialize the Openbao cluster, raft backend
* Store and retrieve token and shards
* Add/remove Openbao servers from the raft
* Unseal Openbao servers

[Openbao](https://openbao.org/)] OpenBao is an open source, community-driven fork of Vault managed by the Linux Foundation.  Openbao Manager is a [Starlingx](https://www.starlingx.io/) project.

## Project Status
_Concept and early development_

The openbao-manager-go repo is a pre-alpha status project.  
