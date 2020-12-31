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
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
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

// FinancialStats is a struct for tracking some metrics relevant to financial health of mining operations
type FinancialStats struct {
	EthereumValue         float64
	EthereumValueDelta    float64
	PoolBalanceValueTotal float64
	PoolBalanceValueDelta float64
	DeltaWindow           int64 // in nano seconds
}

// InfluxDBLine will convert the struct to a byte slice for delivery as a network payload
func (fs *FinancialStats) InfluxDBLine(table string) (payload []byte) {
	payload = []byte(fmt.Sprintf("%s,EthereumValue=%f EthereumValueDelta=%f PoolBalanceValueTotal=%f PoolBalanceValueDelta=%f DeltaWindow=%d %d\n",
		table, fs.EthereumValue, fs.EthereumValueDelta, fs.PoolBalanceValueTotal, fs.PoolBalanceValueDelta, fs.DeltaWindow, time.Now().UTC().UnixNano()))
	return
}

// QuestDBResponse is the expected shape of a successful response to a query
type QuestDBResponse struct {
	Query   string          `json:"query"`
	Columns []interface{}   `json:"columns"`
	Dataset [][]interface{} `json:"dataset"`
	Count   int64           `json:"count"`
	Timings interface{}     `json:"timings"`
}

// QuestDBErrorResponse is the expected shape of an error response to a query
type QuestDBErrorResponse struct {
	Query    string `json:"query"`
	Error    string `json:"error"`
	Position int64  `json:"position"`
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
	Run: metricsCmdRun,
}

func metricsCmdRun(cmd *cobra.Command, args []string) {
	log.Debugln("metricsCmdRun called")
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
	metricsCmd.Flags().BoolP("dryrun", "d", false, "Print metrics instead of shipping them to a timeseries DB")
	metricsCmd.Flags().BoolP("first-run", "f", false, "First run, use this flag to automate some setup like create questDB tables")
}

func collectPoolStats() (poolStats PoolStats, err error) {
	return
}

func collectFinancialStats() (financialStats FinancialStats, err error) {
	return
}

func createTables() {
	tables := []string{
		"pool",
		"financial",
	}
	queryFmt := "CREATE TABLE %s"
	for _, t := range tables {
		queryQuestDB("http://localhost:9000", fmt.Sprintf(queryFmt, t), new(Metrics))
	}
}

func getLastTimeSeries(table string, metrics *Metrics) {
	query := fmt.Sprintf("%s LATEST BY timestamp", table)
	queryQuestDB("http://localhost:9000", query, metrics)
}

func queryQuestDB(apiRoot string, query string, metrics *Metrics) (err error) {
	u, err := url.Parse(apiRoot)
	if err != nil {
		fmt.Println(err)
		log.Errorf("queryQuestDB: url.Parse(%s); returned err=%s\n", apiRoot, err.Error())
		// TODO: handle error
		return
	}

	u.Path += "exec"
	params := url.Values{}
	params.Add("query", query)
	u.RawQuery = params.Encode()
	url := fmt.Sprintf("%v", u)

	resp, err := apiClient.Get(url)
	if err != nil {
		fmt.Println(err)
		log.Errorf("queryQuestDB: apiClient.Get(%s); returned err=%s\n", url, err.Error())
		// TODO: handle error
		return
	}

	defer resp.Body.Close()

	questDBResponse := new(QuestDBResponse)
	err = json.NewDecoder(resp.Body).Decode(&questDBResponse)
	if err != nil {
		fmt.Println(err)
		log.Errorf("queryQuestDB: json.NewDecoder(resp.Body).Decode(&minerBalance); returned err=%s\n", err.Error())
		// TODO: handle error
		return
	}
	return
}
