package command

import (
	"flag"
	"github.com/pefish/go-coin-eth"
	"github.com/pefish/go-commander"
	go_config "github.com/pefish/go-config"
	go_crypto "github.com/pefish/go-crypto"
	go_logger "github.com/pefish/go-logger"
)

type GeneAddressCommand struct {
}

func NewGeneAddressCommand() *GeneAddressCommand {
	return &GeneAddressCommand{}
}

func (dc *GeneAddressCommand) DecorateFlagSet(flagSet *flag.FlagSet) error {
	flagSet.String("mnemonic", "test", "mnemonic")
	flagSet.String("pass", "test", "pass")
	flagSet.String("path", "m/0/0", "path")
	return nil
}

func (dc *GeneAddressCommand) OnExited(data *commander.StartData) error {
	return nil
}

func (dc *GeneAddressCommand) Start(data *commander.StartData) error {
	mnemonic := go_config.ConfigManagerInstance.MustGetString("mnemonic")
	pass := go_config.ConfigManagerInstance.MustGetString("pass")
	path := go_config.ConfigManagerInstance.MustGetString("path")

	wallet := go_coin_eth.NewWallet()
	seed := wallet.SeedHexByMnemonic(mnemonic, "")
	result, err := wallet.DeriveFromPath(seed, path)
	if err != nil {
		return err
	}
	go_logger.Logger.InfoF("address: %s", result.Address)
	go_logger.Logger.InfoF("priv: %s", result.PrivateKey)
	go_logger.Logger.InfoF("encrypted priv: %s", go_crypto.Crypto.MustAesCbcEncrypt(pass, result.PrivateKey))
	return nil
}
