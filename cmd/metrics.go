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
	"io/ioutil"
	"mining-tools/nanopool"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Metrics is an interface to cover the various stats and metrics structs that may be created over time
type Metrics interface {
	InfluxDBLine() []byte
}

// PoolStats is a struct for tracking some metrics gathered and calculated from the mining pool
type PoolStats struct {
	Location string
	Balance  float64
	Shares   int64
}

// InfluxDBLine will convert the struct to a byte slice for delivery as a network payload
func (ps *PoolStats) InfluxDBLine(table string) (payload []byte) {
	balance := floatToStringNoTrail(ps.Balance)
	payload = []byte(fmt.Sprintf("%s,Location=%s Balance=%s,Shares=%d %d\n",
		table, ps.Location, balance, ps.Shares, time.Now().UTC().UnixNano()))
	return
}

// FinancialStats is a struct for tracking some metrics relevant to financial health of mining operations
type FinancialStats struct {
	Location    string
	EthereumUSD float64
	BalanceETH  float64
	BalanceUSD  float64
	BalanceBTC  float64
}

// InfluxDBLine will convert the struct to a byte slice for delivery as a network payload
func (fs *FinancialStats) InfluxDBLine(table string) (payload []byte) {
	eth := floatToStringNoTrail(fs.EthereumUSD)
	balETH := floatToStringNoTrail(fs.BalanceETH)
	balUSD := floatToStringNoTrail(fs.BalanceUSD)
	balBTC := floatToStringNoTrail(fs.BalanceBTC)
	payload = []byte(fmt.Sprintf("%s,Location=%s EthereumUSD=%s,BalanceETH=%s,BalanceUSD=%s,BalanceBTC=%s %d\n",
		table, fs.Location, eth, balETH, balUSD, balBTC, time.Now().UTC().UnixNano()))
	return
}

// QuestDBSuccessResponse is the expected shape of a successful response to a query
type QuestDBSuccessResponse struct {
	Query   string           `json:"query"`
	Columns []QuestDBColumns `json:"columns"`
	Dataset [][]interface{}  `json:"dataset"`
	Count   int64            `json:"count"`
	Timings interface{}      `json:"timings"`
}

func (qsr *QuestDBSuccessResponse) isErrorResponse() bool {
	return false
}

// QuestDBColumns is the expected shape of the Columns field from a successful QuestDB success response
type QuestDBColumns struct {
	Name string
	Type string
}

// QuestDBErrorResponse is the expected shape of an error response to a query
type QuestDBErrorResponse struct {
	Query    string `json:"query"`
	Error    string `json:"error"`
	Position int64  `json:"position"`
}

func (qer *QuestDBErrorResponse) isErrorResponse() bool {
	return true
}

// QuestDBResponse is an interface to simplify the questDBQuery return values
type QuestDBResponse interface {
	isErrorResponse() bool
}

// EtherscanAccountBalance is the expected shape of the response from the etherscan API for an account balance plus "Balance" which is calaculated here based on result
type EtherscanAccountBalance struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
	Balance float64
}

// metricsCmd represents the metrics command
var (
	apiClient  = &http.Client{Timeout: 10 * time.Second}
	dryRunFlag bool

	metricsCmd = &cobra.Command{
		Use:   "metrics",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Run: metricsCmdRun,
	}
)

func metricsCmdRun(cmd *cobra.Command, args []string) {
	log.Debugln("metricsCmdRun called")
	poolStats, err := collectPoolStats()
	if err != nil {
		fmt.Println(err)
		log.Errorf("metricsCmdRun: collectPoolStats(); returned err=%s\n", err.Error())
		// TODO: handle error
		return
	}
	payload := poolStats.InfluxDBLine("pool")
	nanoStats, err := collectNanopoolFinancialStats()
	if err != nil {
		fmt.Println(err)
		log.Errorf("metricsCmdRun: collectFinancialStats(); returned err=%s\n", err.Error())
		// TODO: handle error
		return
	}
	payload = append(payload, nanoStats.InfluxDBLine("financial")...)
	walletStats, err := collectWalletFinancialStats()
	if err != nil {
		fmt.Println(err)
		log.Errorf("metricsCmdRun: collectWalletFinancialStats(); returned err=%s\n", err.Error())
		// TODO: handle error
		return
	}
	payload = append(payload, walletStats.InfluxDBLine("financial")...)
	if !dryRunFlag {
		insertQuestDB("127.0.0.1:9009", payload)
	} else {
		fmt.Printf("DRYRUN: Metrics in InfluxDB Line format - %s", payload)
	}
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
	metricsCmd.Flags().BoolVarP(&dryRunFlag, "dryrun", "d", false, "Print metrics instead of shipping them to a timeseries DB")
}

func collectPoolStats() (poolStats PoolStats, err error) {
	poolStats.Location = "nanopool"
	nanoAPIRoot := viper.GetString("miningtools.nanopool.apiroot")
	nanoAddress := viper.GetString("miningtools.nanopool.address")
	mb, err := nanopool.GetMinerBalance(nanoAPIRoot, nanoAddress)
	if err != nil {
		fmt.Println(err)
		log.Errorf("collectPoolStats: getMinerBalance(%s, %s); returned err=%s\n", nanoAPIRoot, nanoAddress, err.Error())
		// TODO: handle error
		return
	}
	poolStats.Balance = mb.Data

	d := time.Duration(10 * time.Minute)
	now := time.Now().UTC().Truncate(d).Unix()
	sr, err := nanopool.GetMinerShareRate(nanoAPIRoot, nanoAddress)
	if err != nil {
		fmt.Println(err)
		log.Errorf("collectPoolStats: getMinerShareRate(%s, %s); returned err=%s\n", nanoAPIRoot, nanoAddress, err.Error())
		// TODO: handle error
		return
	}
	found := false
	for i := range sr.Data {
		if sr.Data[i].Date == now {
			poolStats.Shares = sr.Data[i].Shares
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("Date %d not found in share rate history", now)
	}
	return
}

func collectNanopoolFinancialStats() (financialStats FinancialStats, err error) {
	financialStats.Location = "nanopool"
	nanoAPIRoot := viper.GetString("miningtools.nanopool.apiroot")
	nanoAddress := viper.GetString("miningtools.nanopool.address")
	mb, err := nanopool.GetMinerBalance(nanoAPIRoot, nanoAddress)
	if err != nil {
		fmt.Println(err)
		log.Errorf("collectFinancialStats: nanopool.GetMinerBalance(%s, %s); returned err=%s\n", nanoAPIRoot, nanoAddress, err.Error())
		// TODO: handle error
		return
	}
	bal := mb.Data
	financialStats.BalanceETH = bal
	p, err := nanopool.GetOtherPrices(nanoAPIRoot)
	if err != nil {
		fmt.Println(err)
		log.Errorf("collectFinancialStats: nanopool.GetOtherPrices(%s, %s); returned err=%s\n", nanoAPIRoot, nanoAddress, err.Error())
		// TODO: handle error
		return
	}
	financialStats.EthereumUSD = p.Data.PriceUSD
	financialStats.BalanceUSD = p.Data.PriceUSD * bal
	financialStats.BalanceBTC = p.Data.PriceBTC * bal
	return
}

func collectWalletFinancialStats() (financialStats FinancialStats, err error) {
	financialStats.Location = "wallet"
	nanoAPIRoot := viper.GetString("miningtools.nanopool.apiroot")
	etherscanAPIRoot := viper.GetString("miningtools.etherscan.apiroot")
	walletAddress := viper.GetString("miningtools.etherscan.address")
	ab, err := getWalletBalance(etherscanAPIRoot, walletAddress)
	if err != nil {
		fmt.Println(err)
		log.Errorf("collectFinancialStats: getWalletBalance(%s, %s); returned err=%s\n", etherscanAPIRoot, walletAddress, err.Error())
		// TODO: handle error
		return
	}
	bal := ab.Balance
	financialStats.BalanceETH = bal
	p, err := nanopool.GetOtherPrices(nanoAPIRoot)
	if err != nil {
		fmt.Println(err)
		log.Errorf("collectFinancialStats: nanopool.GetOtherPrices(%s, %s); returned err=%s\n", nanoAPIRoot, walletAddress, err.Error())
		// TODO: handle error
		return
	}
	financialStats.EthereumUSD = p.Data.PriceUSD
	financialStats.BalanceUSD = p.Data.PriceUSD * bal
	financialStats.BalanceBTC = p.Data.PriceBTC * bal
	return
}

func getLastTimeSeries(table string, metrics *Metrics) {
	query := fmt.Sprintf("%s LATEST BY timestamp", table)
	queryQuestDB("http://localhost:9000", query)
}

func queryQuestDB(apiRoot string, query string) (response QuestDBResponse, err error) {
	log.Debugln("queryQuestDB called")
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

	respBody, err := ioutil.ReadAll(resp.Body)

	questDBErrorResponse := new(QuestDBErrorResponse)
	err = json.Unmarshal(respBody, questDBErrorResponse)
	if err != nil {
		fmt.Println(err)
		log.Errorf("queryQuestDB: json.Unmarshal(respBody, questDBErrorResponse); returned err=%s\n", err.Error())
		// TODO: handle error
		return
	}
	response = questDBErrorResponse
	err = fmt.Errorf("Query '%s' failed with error: %s", query, questDBErrorResponse.Error)
	if questDBErrorResponse.Error == *new(string) {
		questDBSuccessResponse := new(QuestDBSuccessResponse)
		err = json.Unmarshal(respBody, questDBSuccessResponse)
		if err != nil {
			fmt.Println(err)
			log.Errorf("queryQuestDB: json.Unmarshal(respBody, questDBSuccessResponse); returned err=%s\n", err.Error())
			// TODO: handle error
			return
		}
		response = questDBSuccessResponse
	}
	return
}

func insertQuestDB(host string, payload []byte) {
	conn, err := net.DialTimeout("tcp", host, time.Second*10)
	if err != nil {
		fmt.Println(err)
		log.Errorf("insertQuestDB: net.DialTCP(tcp, nil, tcpAddr); returned err=%s\n", err.Error())
		// TODO: handle error
		return
	}
	defer conn.Close()

	log.Debugf("insertQuestDB: Sending query to QuestDB - %s", payload)
	_, err = conn.Write(payload)
	if err != nil {
		fmt.Println(err)
		log.Errorf("insertQuestDB: conn.Write(payload); returned err=%s\n", err.Error())
		// TODO: handle error
		return
	}
	// QuestDB does not appear to respond.
}

func getWalletBalance(apiRoot string, address string) (accountBalance EtherscanAccountBalance, err error) {
	log.Debugln("getWalletBalance called")
	u, err := url.Parse(apiRoot)
	if err != nil {
		fmt.Println(err)
		log.Errorf("getWalletBalance: url.Parse(%s); returned err=%s\n", apiRoot, err.Error())
		// TODO: handle error
		return
	}

	params := url.Values{}
	params.Add("module", "account")
	params.Add("action", "balance")
	params.Add("address", address)
	params.Add("tag", "latest")
	params.Add("apikey", viper.GetString("miningtools.etherscan.apikey"))
	u.RawQuery = params.Encode()
	url := fmt.Sprintf("%v", u)

	resp, err := apiClient.Get(url)
	if err != nil {
		fmt.Println(err)
		log.Errorf("getWalletBalance: apiClient.Get(%s); returned err=%s\n", url, err.Error())
		// TODO: handle error
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(respBody, &accountBalance)
	if err != nil {
		fmt.Println(err)
		log.Errorf("getWalletBalance: json.Unmarshal(respBody, accountBalance); returned err=%s\n", err.Error())
		// TODO: handle error
		return
	}

	if accountBalance.Status == "1" {
		bal, _ := strconv.ParseFloat(accountBalance.Result, 64)
		// The balance is provided as a whole number and must be divided to get the correct decimal
		accountBalance.Balance = bal / 1000000000000000000
	}
	return
}

func floatToStringNoTrail(number float64) (noTrail string) {
	noTrail = fmt.Sprintf("%.12f", number)
	noTrail = strings.TrimRight(noTrail, "0")
	return
}
