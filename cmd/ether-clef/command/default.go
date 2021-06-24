package command

import (
	"context"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pefish/ether-clef/pkg/global"
	"github.com/pefish/ether-clef/pkg/http"
	"github.com/pefish/go-commander"
	go_config "github.com/pefish/go-config"
	go_logger "github.com/pefish/go-logger"
	go_mysql "github.com/pefish/go-mysql"
	"github.com/pkg/errors"
)

type DefaultCommand struct {
}

func NewDefaultCommand() *DefaultCommand {
	return &DefaultCommand{}
}

func (dc *DefaultCommand) DecorateFlagSet(flagSet *flag.FlagSet) error {
	flagSet.String("http.vhosts", "*", "Comma separated list of virtual hostnames from which to accept requests (server enforced)")
	flagSet.String("http.addr", "0.0.0.0", "HTTP-RPC server listening interface")
	flagSet.Int("http.port", 8550, "HTTP-RPC server listening port")
	flagSet.Int("chainid", 1, "Chain id to use for signing (1=mainnet, 3=Ropsten, 4=Rinkeby, 5=Goerli)")
	flagSet.String("db.host", "0.0.0.0", "Host of mysql server to connect")
	flagSet.Int("db.port", 3306, "Port of mysql server to connect")
	flagSet.String("db.database", "wallet", "Database to connect")
	flagSet.String("db.username", "test", "Username to connect mysql server")
	flagSet.String("db.password", "", "Password to connect mysql server")
	flagSet.String("password", "", "Password to decrypt private key")
	return nil
}

func (dc *DefaultCommand) OnExited(data *commander.StartData) error {
	return nil
}

func (dc *DefaultCommand) Start(data *commander.StartData) error {
	global.Password = go_config.ConfigManagerInstance.MustGetString("password")

	err := go_mysql.MysqlInstance.ConnectWithConfiguration(go_mysql.Configuration{
		Host:     go_config.ConfigManagerInstance.MustGetString("db.host"),
		Port:     go_config.ConfigManagerInstance.MustGetInt("db.port"),
		Username: go_config.ConfigManagerInstance.MustGetString("db.username"),
		Password: go_config.ConfigManagerInstance.MustGetString("db.password"),
		Database: go_config.ConfigManagerInstance.MustGetString("db.database"),
	})
	if err != nil {
		return err
	}
	go_mysql.MysqlInstance.SetLogger(go_logger.Logger)

	vhosts := go_config.ConfigManagerInstance.MustGetString("http.vhosts")
	chainId := go_config.ConfigManagerInstance.MustGetInt64("chainid")
	rpcAPI := []rpc.API{
		{
			Namespace: "account",
			Public:    true,
			Service:   http.NewSignerAPI(chainId),
			Version:   "1.0",
		},
	}
	srv := rpc.NewServer()
	err = node.RegisterApisFromWhitelist(rpcAPI, []string{"account"}, srv, false)
	if err != nil {
		return errors.Wrap(err, "Could not register API")
	}
	handler := node.NewHTTPHandlerStack(srv, []string{"*"}, []string{vhosts})

	// set port
	httpAddr := go_config.ConfigManagerInstance.MustGetString("http.addr")
	port := go_config.ConfigManagerInstance.MustGetInt("http.port")

	// start http server
	httpEndpoint := fmt.Sprintf("%s:%d", httpAddr, port)
	httpServer, addr, err := node.StartHTTPEndpoint(httpEndpoint, rpc.DefaultHTTPTimeouts, handler)
	if err != nil {
		return errors.Wrap(err, "Could not start RPC api")
	}
	extapiURL := fmt.Sprintf("http://%v/", addr)
	go_logger.Logger.InfoF("HTTP endpoint opened. url: %s", extapiURL)

	defer func() {
		// Don't bother imposing a timeout here.
		httpServer.Shutdown(context.Background())
		go_logger.Logger.InfoF("HTTP endpoint closed. url: %s", extapiURL)
	}()

	<-data.ExitCancelCtx.Done()
	return nil
}
