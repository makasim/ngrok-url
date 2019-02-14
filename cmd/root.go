// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
	"time"
)

var ngrokHostCmd string

type Tunnel struct {
	PublicUrl string `json:"public_url"`
	Proto string `json:"proto"`
}

type Tunnels struct {
	Tunnels []Tunnel `json:"tunnels"`
}

var rootCmd = &cobra.Command{
	Use:   "ngrok-url",
	Short: "Sets Ngrok HTTPS tunnel url as Telegram updates hook",
	Long: ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ngrokHost, err := getNgrokHost(cmd.Flag("api-host-cmd").Value.String())
		if err != nil {
			return fmt.Errorf("Failed to execute docker-compose port cmd. %s", err);
		}

		publicUri, err := getNgrokPublicUri(ngrokHost)
		if err != nil {
			return fmt.Errorf("Failed to guess Ngrok public URI. %s", err);
		}

		fmt.Println(publicUri)

		return nil
	},
}

func getNgrokHost(cmd string) (string, error) {
	var out []byte
	var err error

	cmdParts := strings.Split(cmd, " ")
	first := cmdParts[0]
	args := cmdParts[1:]

	endAt := time.Now().UnixNano() + int64(6 * time.Second)
	for time.Now().UnixNano() < endAt {
		out, err = exec.Command(first, args...).Output()
		if err == nil {
			return strings.TrimSpace(string(out)), nil
		}

		time.Sleep(100 * time.Millisecond)
	}

	return "", err
}

func getNgrokPublicUri(ngrokHost string) (string, error) {
	var out []byte
	var err error

	endAt := time.Now().UnixNano() + int64(2 * time.Second)

	for time.Now().UnixNano() < endAt {
		out, err = exec.Command("curl", ngrokHost + "/api/tunnels").Output()
		if err != nil {
			time.Sleep(100 * time.Millisecond)

			continue
		}

		var tunnels Tunnels
		if err = json.Unmarshal(out, &tunnels); err != nil {
			panic(err)
		}

		var publicUri string
		for _, tunnel := range tunnels.Tunnels {
			if tunnel.Proto == "https" {
				publicUri = tunnel.PublicUrl
			}
		}

		if publicUri == "" {
			err = errors.New("Public URI could not be parsed, or has not yet set")

			time.Sleep(100 * time.Millisecond)

			continue
		}

		return publicUri, nil
	}

	return "", err
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&ngrokHostCmd, "api-host-cmd", "", "", "A command will be used to find out ngrok host. Possibly 'docker-compose port ngrok 4040'")
	rootCmd.MarkFlagRequired("ngrok-host-cmd")
}
