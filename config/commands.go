package config

const (
	//DOCKER will point to the docker executable on a UNIX/linux OS
	DOCKER = "/usr/local/bin/docker"
	//TATUM is the command used to generate secure wallets for docker executable
	TATUM = "tatum-kms"
	//NODE will point to the node executable on a UNIX/Linux OS
	NODE = "/usr/local/opt/node@16/bin/node"
	//WORKDIR points to the distributon of the tatym kms
	WORKDIR = "/usr/src/app"
)
