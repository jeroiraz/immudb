/*
Copyright 2019-2020 vChain, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
	http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package immuadmin

import (
	"context"
	"fmt"
	"github.com/codenotary/immudb/cmd/docs/man"
	"github.com/codenotary/immudb/cmd/immuadmin/command/service"
	"strings"

	c "github.com/codenotary/immudb/cmd/helper"
	"github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/gw"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Init ...
func Init(rootCmd *cobra.Command, cmdName string, cfgFn *string) *cobra.Command {
	if err := configureFlags(rootCmd, cfgFn); err != nil {
		c.QuitToStdErr(err)
	}
	cl := new(commandline)
	cl.options = Options()
	cl.passwordReader = c.DefaultPasswordReader
	cl.context = context.Background()
	cl.hds = client.NewHomedirService()

	cl.user(rootCmd)
	cl.login(rootCmd)
	cl.logout(rootCmd)
	cl.status(rootCmd)
	cl.stats(rootCmd)
	cl.serverConfig(rootCmd)

	clb, err := NewCommandlineBck()
	if err != nil {
		c.QuitToStdErr(err)
	}
	clb.dumpToFile(rootCmd)
	clb.backup(rootCmd)
	clb.restore(rootCmd)
	cl.printTree(rootCmd)

	cld := service.NewCommandLine()
	cld.Service(rootCmd)

	rootCmd.AddCommand(man.Generate(rootCmd, cmdName, "./cmd/docs/man/"+cmdName))

	return rootCmd
}

func Options() *client.Options {
	port := viper.GetInt("immudb-port")
	address := viper.GetString("immudb-address")
	tokenFileName := viper.GetString("tokenfile")
	if !strings.HasSuffix(tokenFileName, client.AdminTokenFileSuffix) {
		tokenFileName += client.AdminTokenFileSuffix
	}
	mtls := viper.GetBool("mtls")
	certificate := viper.GetString("certificate")
	servername := viper.GetString("servername")
	pkey := viper.GetString("pkey")
	clientcas := viper.GetString("clientcas")
	options := client.DefaultOptions().
		WithPort(port).
		WithAddress(address).
		WithAuth(true).
		WithTokenFileName(tokenFileName).
		WithMTLs(mtls)
	if mtls {
		// todo https://golang.org/src/crypto/x509/root_linux.go
		options.MTLsOptions = client.DefaultMTLsOptions().
			WithServername(servername).
			WithCertificate(certificate).
			WithPkey(pkey).
			WithClientCAs(clientcas)
	}
	return options
}

func configureFlags(cmd *cobra.Command, cfgFn *string) error {
	cmd.PersistentFlags().IntP("immudb-port", "p", gw.DefaultOptions().ImmudbPort, "immudb port number")
	cmd.PersistentFlags().StringP("immudb-address", "a", gw.DefaultOptions().ImmudbAddress, "immudb host address")
	cmd.PersistentFlags().String(
		"tokenfile",
		client.DefaultOptions().TokenFileName,
		fmt.Sprintf(
			"authentication token file (default path is $HOME or binary location; the supplied "+
				"value will be automatically suffixed with %s; default filename is %s%s)",
			client.AdminTokenFileSuffix,
			client.DefaultOptions().TokenFileName,
			client.AdminTokenFileSuffix))
	cmd.PersistentFlags().StringVar(cfgFn, "config", "", "config file (default path is configs or $HOME; default filename is immuadmin.toml)")
	cmd.PersistentFlags().BoolP("mtls", "m", client.DefaultOptions().MTLs, "enable mutual tls")
	cmd.PersistentFlags().String("servername", client.DefaultMTLsOptions().Servername, "used to verify the hostname on the returned certificates")
	cmd.PersistentFlags().String("certificate", client.DefaultMTLsOptions().Certificate, "server certificate file path")
	cmd.PersistentFlags().String("pkey", client.DefaultMTLsOptions().Pkey, "server private key path")
	cmd.PersistentFlags().String("clientcas", client.DefaultMTLsOptions().ClientCAs, "clients certificates list. Aka certificate authority")
	if err := viper.BindPFlag("immudb-port", cmd.PersistentFlags().Lookup("immudb-port")); err != nil {
		return err
	}
	if err := viper.BindPFlag("immudb-address", cmd.PersistentFlags().Lookup("immudb-address")); err != nil {
		return err
	}
	if err := viper.BindPFlag("tokenfile", cmd.PersistentFlags().Lookup("tokenfile")); err != nil {
		return err
	}
	if err := viper.BindPFlag("mtls", cmd.PersistentFlags().Lookup("mtls")); err != nil {
		return err
	}
	if err := viper.BindPFlag("servername", cmd.PersistentFlags().Lookup("servername")); err != nil {
		return err
	}
	if err := viper.BindPFlag("certificate", cmd.PersistentFlags().Lookup("certificate")); err != nil {
		return err
	}
	if err := viper.BindPFlag("pkey", cmd.PersistentFlags().Lookup("pkey")); err != nil {
		return err
	}
	if err := viper.BindPFlag("clientcas", cmd.PersistentFlags().Lookup("clientcas")); err != nil {
		return err
	}
	viper.SetDefault("immudb-port", gw.DefaultOptions().ImmudbPort)
	viper.SetDefault("immudb-address", gw.DefaultOptions().ImmudbAddress)
	viper.SetDefault("tokenfile", client.DefaultOptions().TokenFileName)
	viper.SetDefault("mtls", client.DefaultOptions().MTLs)
	viper.SetDefault("servername", client.DefaultMTLsOptions().Servername)
	viper.SetDefault("certificate", client.DefaultMTLsOptions().Certificate)
	viper.SetDefault("pkey", client.DefaultMTLsOptions().Pkey)
	viper.SetDefault("clientcas", client.DefaultMTLsOptions().ClientCAs)

	return nil
}
