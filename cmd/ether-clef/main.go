package main

import (
	"github.com/pefish/ether-clef/cmd/ether-clef/command"
	"github.com/pefish/ether-clef/version"
	"github.com/pefish/go-commander"
	go_logger "github.com/pefish/go-logger"
)

func main() {
	commanderInstance := commander.NewCommander(version.AppName, version.Version, version.AppName + " is a substitute for official clef。Author：pefish")
	commanderInstance.RegisterDefaultSubcommand("", command.NewDefaultCommand())
	commanderInstance.RegisterSubcommand("gene-address", "", command.NewGeneAddressCommand())
	err := commanderInstance.Run()
	if err != nil {
		go_logger.Logger.Error(err)
	}
}
