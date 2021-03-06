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
	"fmt"

	c "github.com/codenotary/immudb/cmd/helper"
	"github.com/codenotary/immudb/pkg/auth"
	"github.com/codenotary/immudb/pkg/client"
	"github.com/spf13/cobra"
)

func (cl *commandline) login(cmd *cobra.Command) {
	ccmd := &cobra.Command{
		Use:               "login username (you will be prompted for password)",
		Short:             fmt.Sprintf("Login using the specified username and password (admin username is %s)", auth.SysAdminUsername),
		Aliases:           []string{"l"},
		PersistentPreRunE: cl.connect,
		PersistentPostRun: cl.disconnect,
		RunE: func(cmd *cobra.Command, args []string) error {
			tokenFileName := cl.immuClient.GetOptions().TokenFileName
			ctx := cl.context
			user := args[0]
			if user != auth.SysAdminUsername {
				c.QuitToStdErr(fmt.Errorf("Permission denied: user %s has no admin rights", user))
			}
			pass, err := cl.passwordReader.Read("Password:")
			if err != nil {
				c.QuitToStdErr(err)
			}
			response, err := cl.immuClient.Login(ctx, []byte(user), pass)
			if err != nil {
				c.QuitWithUserError(err)
			}
			if err = cl.hds.WriteFileToUserHomeDir(response.Token, tokenFileName); err != nil {
				c.QuitToStdErr(err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "logged in\n")
			if cl.immuClient, err = client.NewImmuClient(cl.immuClient.GetOptions()); err != nil {
				c.QuitWithUserError(err)
			}
			if string(response.Warning) == auth.WarnDefaultAdminPassword {
				c.PrintfColorW(cmd.OutOrStdout(), c.Yellow, "SECURITY WARNING: %s\n", response.Warning)
				cl.changePassword(user, pass)
			}
			return nil
		},
		Args: cobra.ExactArgs(1),
	}
	cmd.AddCommand(ccmd)
}

func (cl *commandline) logout(cmd *cobra.Command) {
	ccmd := &cobra.Command{
		Use:               "logout",
		Aliases:           []string{"x"},
		PersistentPreRunE: cl.connect,
		PersistentPostRun: cl.disconnect,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cl.immuClient.Logout(cl.context); err != nil {
				c.QuitWithUserError(err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "logged out\n")
			return nil
		},
		Args: cobra.NoArgs,
	}
	cmd.AddCommand(ccmd)
}
