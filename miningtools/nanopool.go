/*
Package miningtools contains the various supported CLI commands for mining-tools
Copyright © 2020 Keith Olenchak <kenjin.domini@gmail.com>

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
package miningtools

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NanopoolError is a struct for marshaling json of any error coming from nanopool's API
type NanopoolError struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

// nanopoolCmd represents the nanopool command
var nanopoolCmd = &cobra.Command{
	Use:   "nanopool",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("nanopool called")
	},
}

func init() {
	rootCmd.AddCommand(nanopoolCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nanopoolCmd.PersistentFlags().String("foo", "", "A help for foo")
	nanopoolCmd.PersistentFlags().String("address", "", "miner account")
	viper.BindPFlag("nanopool.address", nanopoolCmd.PersistentFlags().Lookup("address"))
	nanopoolCmd.PersistentFlags().String("apiRoot", "https://api.nanopool.org/v1/eth/", "base URL for nanopool ethereum API")
	viper.BindPFlag("nanopool.apiRoot", nanopoolCmd.PersistentFlags().Lookup("apiRoot"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nanopoolCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
