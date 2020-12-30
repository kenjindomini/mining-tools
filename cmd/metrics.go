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
	"time"

	"github.com/spf13/cobra"
)

// Metrics is an interface to cover the various stats and metrics structs that may be created over time
type Metrics interface {
	InfluxDBLine() []byte
}

// PoolStats is a struct for tracking some metrics gathered and calculated from the mining pool
type PoolStats struct {
	BalanceCurrent float64
	BalanceDelta   float64
	SharesCurrent  int64
	SharesDelta    int64
	DeltaWindow    int64 // in nano seconds
}

// InfluxDBLine will convert the struct to a byte slice for delivery as a network payload
func (ps *PoolStats) InfluxDBLine(table string) (payload []byte) {
	payload = []byte(fmt.Sprintf("%s,BalanceCurrent=%f BalanceDelta=%f SharesCurrent=%d SharesDelta=%d DeltaWindow=%d %d\n",
		table, ps.BalanceCurrent, ps.BalanceDelta, ps.SharesCurrent, ps.SharesDelta, ps.DeltaWindow, time.Now().UTC().UnixNano()))
	return
}

// metricsCmd represents the metrics command
var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("metrics called")
	},
}

func init() {
	rootCmd.AddCommand(metricsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// metricsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// metricsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
