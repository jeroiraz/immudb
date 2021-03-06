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

package immuc

import (
	"log"
	"strings"
	"testing"

	"github.com/codenotary/immudb/cmd/helper"
	c "github.com/codenotary/immudb/cmd/helper"
	"github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/server"
	"github.com/codenotary/immudb/pkg/server/servertest"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func newClient(pr helper.PasswordReader, dialer servertest.BuffDialer) Client {
	dialOptions := []grpc.DialOption{
		grpc.WithContextDialer(dialer), grpc.WithInsecure(),
	}

	c.DefaultPasswordReader = pr
	imc, err := Init(Options().WithDialOptions(&dialOptions))
	if err != nil {
		log.Fatal(err)
	}
	err = imc.Connect([]string{""})
	if err != nil {
		log.Fatal(err)
	}
	return imc
}
func login(username string, password string, dialer servertest.BuffDialer) Client {
	viper.Set("tokenfile", client.DefaultOptions().TokenFileName)
	imc := newClient(&testPasswordReader{
		pass: []string{password},
	}, dialer)
	msg, err := imc.Login([]string{username})
	if err != nil {
		log.Fatal(err)
	}
	if !strings.Contains(msg, "Successfully logged in.") {
		log.Fatal("Login error")
	}

	return imc
}

type testPasswordReader struct {
	pass       []string
	callNumber int
}

func (pr *testPasswordReader) Read(msg string) ([]byte, error) {
	pass := []byte(pr.pass[pr.callNumber])
	pr.callNumber++
	return pass, nil
}

func TestLogin(t *testing.T) {
	viper.Set("tokenfile", "token")
	options := server.DefaultOptions().WithAuth(true).WithInMemoryStore(true)
	bs := servertest.NewBufconnServer(options)
	bs.Start()

	dialOptions := []grpc.DialOption{
		grpc.WithContextDialer(bs.Dialer), grpc.WithInsecure(),
	}
	c.DefaultPasswordReader = &testPasswordReader{
		pass: []string{"immudb"},
	}
	imc, err := Init(Options().WithDialOptions(&dialOptions))
	if err != nil {
		t.Fatal(err)
	}
	err = imc.Connect([]string{""})
	if err != nil {
		t.Fatal(err)
	}

	msg, err := imc.Login([]string{"immudb"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(msg, "Successfully logged in.") {
		t.Fatal("Login error")
	}
}
func TestLogout(t *testing.T) {
	viper.Set("tokenfile", client.DefaultOptions().TokenFileName)
	options := server.DefaultOptions().WithAuth(true).WithInMemoryStore(true)
	bs := servertest.NewBufconnServer(options)
	bs.Start()

	dialOptions := []grpc.DialOption{
		grpc.WithContextDialer(bs.Dialer), grpc.WithInsecure(),
	}
	pr := new(testPasswordReader)
	pr.pass = []string{"immudb"}
	c.DefaultPasswordReader = pr

	imc, err := Init(Options().WithDialOptions(&dialOptions))
	if err != nil {
		t.Fatal(err)
	}
	err = imc.Connect([]string{""})
	if err != nil {
		t.Fatal(err)
	}
	_, err = imc.Logout([]string{""})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserList(t *testing.T) {
	options := server.DefaultOptions().WithAuth(true).WithInMemoryStore(true)
	bs := servertest.NewBufconnServer(options)
	bs.Start()
	imc := login("immudb", "immudb", bs.Dialer)
	_, err := imc.UserList([]string{""})
	if err != nil {
		t.Fatal("Userlist fail", err)
	}
}
func TestUserCreate(t *testing.T) {
	options := server.DefaultOptions().WithAuth(true).WithInMemoryStore(true)
	bs := servertest.NewBufconnServer(options)
	bs.Start()
	imc := login("immudb", "immudb", bs.Dialer)

	var userCreateTests = []struct {
		name     string
		args     []string
		password string
		expected string
		test     func(*testing.T, string, []string, string)
	}{
		{
			"Create user",
			[]string{"create", "myuser", "readwrite", "defaultdb"},
			"MyUser@9",
			"Created user",
			func(t *testing.T, password string, args []string, exp string) {
				imc = newClient(&testPasswordReader{
					pass: []string{password, password},
				}, bs.Dialer)
				msg, err := imc.UserOperations(args)
				if err != nil {
					t.Fatal("TestUserCreate fail", err)
				}
				if !strings.Contains(msg, exp) {
					t.Fatalf("TestUserCreate failed to create user: %s", msg)
				}
			},
		},
		{
			"Create user read",
			[]string{"create", "myuserRead", "read", "defaultdb"},
			"MyUser@9",
			"Created user",
			func(t *testing.T, password string, args []string, exp string) {
				imc = newClient(&testPasswordReader{
					pass: []string{password, password},
				}, bs.Dialer)
				msg, err := imc.UserOperations(args)
				if err != nil {
					t.Fatal("TestUserCreate fail", err)
				}
				if !strings.Contains(msg, exp) {
					t.Fatalf("TestUserCreate failed to create user: %s", msg)
				}
			},
		},
		{
			"Create user admin",
			[]string{"create", "myuseradmin", "admin", "defaultdb"},
			"MyUser@9",
			"Created user",
			func(t *testing.T, password string, args []string, exp string) {
				imc = newClient(&testPasswordReader{
					pass: []string{password, password},
				}, bs.Dialer)
				msg, err := imc.UserOperations(args)
				if err != nil {
					t.Fatal("TestUserCreate fail", err)
				}
				if !strings.Contains(msg, exp) {
					t.Fatalf("TestUserCreate failed to create user: %s", msg)
				}
			},
		},
		{
			"Create user wrong permission",
			[]string{"create", "myuserguard", "guard", "defaultdb"},
			"MyUser@9",
			"Permission value not recognized.",
			func(t *testing.T, password string, args []string, exp string) {
				imc = newClient(&testPasswordReader{
					pass: []string{password, password},
				}, bs.Dialer)
				msg, err := imc.UserOperations(args)
				if err != nil {
					t.Fatal("TestUserCreate fail", err)
				}
				if !strings.Contains(msg, exp) {
					t.Fatalf("TestUserCreate failed to create user: %s", msg)
				}
			},
		},
		{
			"Create duplicate user",
			[]string{"create", "myuser", "readwrite", "defaultdb"},
			"MyUser@9",
			"user already exists",
			func(t *testing.T, password string, args []string, exp string) {
				imc = newClient(&testPasswordReader{
					pass: []string{password, password},
				}, bs.Dialer)
				msg, err := imc.UserOperations(args)
				if err == nil {
					t.Fatal("TestUserCreate fail", err)
				}
				if !strings.Contains(err.Error(), exp) {
					t.Fatalf("TestUserCreate failed to create user: %s", err)
				}
				if msg != "" {
					t.Fatalf("TestUserCreate %s", msg)
				}
			},
		},
	}
	for _, tt := range userCreateTests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t, tt.password, tt.args, tt.expected)
		})
	}
}

func TestUserChangePassword(t *testing.T) {
	options := server.DefaultOptions().WithAuth(true).WithInMemoryStore(true)
	bs := servertest.NewBufconnServer(options)
	bs.Start()
	imc := login("immudb", "immudb", bs.Dialer)

	var userCreateTests = []struct {
		name     string
		args     []string
		password string
		expected string
		test     func(*testing.T, string, []string, string)
	}{
		{
			"Change user password",
			[]string{"immudb"},
			"MyUser@9",
			"Password of immudb was changed successfuly",
			func(t *testing.T, password string, args []string, exp string) {
				imc = newClient(&testPasswordReader{
					pass: []string{"immudb", password, password},
				}, bs.Dialer)
				msg, err := imc.ChangeUserPassword(args)
				if err != nil {
					t.Fatal("TestUserChangePassword fail", err)
				}
				if !strings.Contains(msg, exp) {
					t.Fatalf("TestUserChangePassword failed to change password: %s", msg)
				}
			},
		},
		{
			"Change user password wrong old password",
			[]string{"immudb"},
			"MyUser@9",
			"old password is incorrect",
			func(t *testing.T, password string, args []string, exp string) {
				imc := login("immudb", "MyUser@9", bs.Dialer)
				imc = newClient(&testPasswordReader{
					pass: []string{"immudb", password, password},
				}, bs.Dialer)
				msg, err := imc.ChangeUserPassword(args)
				if err == nil {
					t.Fatal("TestUserChangePassword fail", err)
				}
				if !strings.Contains(err.Error(), exp) {
					t.Fatalf("TestUserChangePassword failed to change password: %s", msg)
				}
			},
		},
	}
	for _, tt := range userCreateTests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t, tt.password, tt.args, tt.expected)
		})
	}
}

func TestUserSetActive(t *testing.T) {
	options := server.DefaultOptions().WithAuth(true).WithInMemoryStore(true)
	bs := servertest.NewBufconnServer(options)
	bs.Start()
	imc := login("immudb", "immudb", bs.Dialer)

	imc = newClient(&testPasswordReader{
		pass: []string{"MyUser@9", "MyUser@9"},
	}, bs.Dialer)
	_, err := imc.UserOperations([]string{"create", "myuser", "readwrite", "defaultdb"})
	if err != nil {
		t.Fatal("TestUserCreate fail", err)
	}
	var userCreateTests = []struct {
		name     string
		args     []string
		password string
		expected string
		test     func(*testing.T, string, []string, string)
	}{
		{
			"SetActiveUser",
			[]string{"myuser"},
			"",
			"User status changed successfully",
			func(t *testing.T, password string, args []string, exp string) {
				imc = newClient(&testPasswordReader{
					pass: []string{"immudb", password, password},
				}, bs.Dialer)
				msg, err := imc.SetActiveUser(args, true)
				if err != nil {
					t.Fatal("TestUserChangePassword fail", err)
				}
				if !strings.Contains(msg, exp) {
					t.Fatalf("SetActiveUser failed to change status: %s", msg)
				}
			},
		},
		{
			"Deactivate user",
			[]string{"myuser"},
			"",
			"User status changed successfully",
			func(t *testing.T, password string, args []string, exp string) {
				imc = newClient(&testPasswordReader{
					pass: []string{"immudb", password, password},
				}, bs.Dialer)
				msg, err := imc.SetActiveUser(args, false)
				if err != nil {
					t.Fatal("TestUserChangePassword fail", err)
				}
				if !strings.Contains(msg, exp) {
					t.Fatalf("SetActiveUser failed to change status: %s", msg)
				}
			},
		},
	}
	for _, tt := range userCreateTests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t, tt.password, tt.args, tt.expected)
		})
	}
}

func TestSetUserPermission(t *testing.T) {
	options := server.DefaultOptions().WithAuth(true).WithInMemoryStore(true)
	bs := servertest.NewBufconnServer(options)
	bs.Start()
	imc := login("immudb", "immudb", bs.Dialer)

	imc = newClient(&testPasswordReader{
		pass: []string{"MyUser@9", "MyUser@9"},
	}, bs.Dialer)
	_, err := imc.UserOperations([]string{"create", "myuser", "readwrite", "defaultdb"})
	if err != nil {
		t.Fatal("TestUserCreate fail", err)
	}
	var userCreateTests = []struct {
		name     string
		args     []string
		password string
		expected string
		test     func(*testing.T, string, []string, string)
	}{
		{
			"SetUserPermission user",
			[]string{"grant", "myuser", "admin", "defaultdb"},
			"MyUser@9",
			"Permission changed successfully",
			func(t *testing.T, password string, args []string, exp string) {
				imc = newClient(&testPasswordReader{
					pass: []string{password, password},
				}, bs.Dialer)
				msg, err := imc.SetUserPermission(args)
				if err != nil {
					t.Fatal("TestUserCreate fail", err)
				}
				if !strings.Contains(msg, exp) {
					t.Fatalf("TestUserCreate failed to create user: %s", msg)
				}
			},
		},
		{
			"SetUserPermission user",
			[]string{"revoke", "myuser", "admin", "defaultdb"},
			"MyUser@9",
			"Permission changed successfully",
			func(t *testing.T, password string, args []string, exp string) {
				imc = newClient(&testPasswordReader{
					pass: []string{password, password},
				}, bs.Dialer)
				msg, err := imc.SetUserPermission(args)
				if err != nil {
					t.Fatal("TestUserCreate fail", err)
				}
				if !strings.Contains(msg, exp) {
					t.Fatalf("TestUserCreate failed to create user: %s", msg)
				}
			},
		},
		{
			"SetUserPermission user",
			[]string{"grant", "myuser", "readwrite", "defaultdb"},
			"MyUser@9",
			"Permission changed successfully",
			func(t *testing.T, password string, args []string, exp string) {
				imc = newClient(&testPasswordReader{
					pass: []string{password, password},
				}, bs.Dialer)
				msg, err := imc.SetUserPermission(args)
				if err != nil {
					t.Fatal("TestUserCreate fail", err)
				}
				if !strings.Contains(msg, exp) {
					t.Fatalf("TestUserCreate failed to create user: %s", msg)
				}
			},
		},
		{
			"SetUserPermission user",
			[]string{"grant", "myuser", "read", "defaultdb"},
			"MyUser@9",
			"Permission changed successfully",
			func(t *testing.T, password string, args []string, exp string) {
				imc = newClient(&testPasswordReader{
					pass: []string{password, password},
				}, bs.Dialer)
				msg, err := imc.SetUserPermission(args)
				if err != nil {
					t.Fatal("TestUserCreate fail", err)
				}
				if !strings.Contains(msg, exp) {
					t.Fatalf("TestUserCreate failed to create user: %s", msg)
				}
			},
		},
	}
	for _, tt := range userCreateTests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t, tt.password, tt.args, tt.expected)
		})
	}
}
