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
	"math"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// MinerGeneralInfo is for decoding json from a successful response of the nanopool miner general info api endpoint
type MinerGeneralInfo struct {
	Status bool                 `json:"status"`
	Data   MinerGeneralInfoData `json:"data"`
}

// MinerGeneralInfoData is for decoding json from a successful response of the nanopool miner general info api endpoint
type MinerGeneralInfoData struct {
	Account            string                   `json:"account"`
	UnconfirmedBalance string                   `json:"unconfirmed_balance"`
	Balance            string                   `json:"balance"`
	Hashrate           string                   `json:"hashrate"`
	AvgHashrate        MinerAvgHashrate         `json:"avgHashrate"`
	Workers            []MinerGeneralInfoWorker `json:"workers"`
	RewardPerShare     string                   `json:"rewardPerShare,omitempty"`
	SharesPerHour      int64                    `json:"sharesPerHour,omitempty"`
	RewardPerHour      string                   `json:"rewardPerHour,omitempty"`
}

// MinerAvgHashrate is for decoding json from a successful response of the nanopool miner general info api endpoint
type MinerAvgHashrate struct {
	H1  string `json:"h1"`
	H3  string `json:"h3"`
	H6  string `json:"h6"`
	H12 string `json:"h12"`
	H24 string `json:"h24"`
}

// MinerGeneralInfoWorker is for decoding json from a successful response of the nanopool miner general info api endpoint
type MinerGeneralInfoWorker struct {
	ID        string `json:"id"`
	UID       int64  `json:"uid"`
	Hashrate  string `json:"hashrate"`
	Lastshare int64  `json:"lastshare"`
	Rating    int64  `json:"rating"`
	H1        string `json:"h1"`
	H3        string `json:"h3"`
	H6        string `json:"h6"`
	H12       string `json:"h12"`
	H24       string `json:"h24"`
}

// MinerShareRate is for decoding json from a successful response of the nanopool miner general info api endpoint
type MinerShareRate struct {
	Status bool                 `json:"status"`
	Data   []MinerShareRateData `json:"data"`
}

// MinerShareRateData is for decoding json from a successful response of the nanopool miner general info api endpoint
type MinerShareRateData struct {
	Date   int64 `json:"date"`
	Shares int64 `json:"shares"`
}

// MinerPayments is for decoding json from a successful response of the nanopool miner general info api endpoint
type MinerPayments struct {
	Status bool                `json:"status"`
	Data   []MinerPaymentsData `json:"data"`
}

// MinerPaymentsData is for decoding json from a successful response of the nanopool miner general info api endpoint
type MinerPaymentsData struct {
	Date      int64   `json:"date"`
	TXHash    string  `json:"txHash"`
	Amount    float64 `json:"amount"`
	Confirmed bool    `json:"confirmed"`
}

// HTTPClient is an interface to abstract http.client to support testing using mocks
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

// generalInfoCmd represents the generalInfo sub-command of nanopool
var (
	apiClient          = HTTPClient(&http.Client{Timeout: 10 * time.Second})
	rewardPerShareFlag bool
	sharesPerHourFlag  bool
	generalInfoCmd     = &cobra.Command{
		Use:   "generalInfo",
		Short: "Gets general info of nanopool ethereum miner account",
		Long: `Gets general info of nanopool ethereum miner account, with
				options for calculating additional values`,
		Run: generalInfoCmdRun,
	}
)

func init() {
	nanopoolCmd.AddCommand(generalInfoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generalInfoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generalInfoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	generalInfoCmd.Flags().BoolVarP(&rewardPerShareFlag, "rewardPerShare", "r", false,
		"Include calculated attribute rewardPerShare (lifetime average)")
	generalInfoCmd.Flags().BoolVarP(&sharesPerHourFlag, "sharesPerHour", "s", false,
		"Include calculated attribute sharesPerHour (24h rolling average)")
}

func generalInfoCmdRun(cmd *cobra.Command, args []string) {
	log.Debugln("generalInfoCmdRun called")
	address := viper.GetString("miningtools.nanopool.address")
	apiRoot := viper.GetString("miningtools.nanopool.apiRoot")
	log.Debugf("generalInfoCmdRun: address=%s; apiRoot=%s\n", address, apiRoot)
	info, err := getMinerGeneralInfo(apiRoot, address)
	if err != nil {
		fmt.Println(err)
		log.Errorf("generalInfoCmdRun: getMinerGeneralInfo(apiRoot=%s, address=%s); returned err=%s\n",
			apiRoot, address, err.Error())
		// TODO: handle error
		return
	}
	if rewardPerShareFlag {
		payments, err := getMinerPayments(apiRoot, address)
		if err != nil {
			fmt.Println(err)
			log.Errorf("generalInfoCmdRun: getMinerPayments(apiRoot=%s, address=%s); returned err=%s\n",
				apiRoot, address, err.Error())
			// TODO: handle error
			return
		}
		totalPayouts := calcTotalPayout(payments.Data)
		totalShares := calcTotalShares(info.Data.Workers)
		info.Data.RewardPerShare = calcRPS(info.Data.Balance, totalPayouts, totalShares)
		log.Debugf("generalInfoCmdRun: totalPayouts=%s; totalShares=%d; info.Data.RewardPerShare=%s\n",
			totalPayouts, totalShares, info.Data.RewardPerShare)
	}
	if sharesPerHourFlag {
		hours := int64(24)
		shareRate, err := getMinerShareRate(apiRoot, address)
		if err != nil {
			fmt.Println(err)
			log.Errorf("generalInfoCmdRun: getMinerShareRate(apiRoot=%s, address=%s); returned err=%s\n",
				apiRoot, address, err.Error())
			// TODO: handle error
			return
		}
		info.Data.SharesPerHour = calcSharesPerHour(shareRate.Data, &hours)
		log.Debugf("generalInfoCmdRun: info.Data.SharesPerHour=%s\n", info.Data.SharesPerHour)
	}
	if rewardPerShareFlag && sharesPerHourFlag {
		info.Data.RewardPerHour = calcRewardPerHour(info.Data.RewardPerShare, info.Data.SharesPerHour)
		log.Debugf("generalInfoCmdRun: info.Data.RewardPerHour=%s\n", info.Data.RewardPerHour)
	}
	prettyPrint(info)
	// TODO: updating the saved config should be optional
	// TODO: break this out to check if writing is desired and
	//	create the file if viper fails to, then retry
	if err = viper.WriteConfig(); err != nil {
		fmt.Println(err)
		log.Errorf("viper.WriteConfig(); returned err=%s\n", err.Error())
		// TODO: handle error
		return
	}
}

func prettyPrint(v interface{}) {
	pretty, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		// TODO: Log error
		fmt.Println(err)
		// TODO: handle error
		return
	}
	out := string(pretty)
	fmt.Println(out)
}

func getMinerGeneralInfo(apiRoot string, address string) (info MinerGeneralInfo, err error) {
	log.Debugln("getMinerGeneralInfo called")
	resp, err := apiClient.Get(fmt.Sprintf("%s%s%s", apiRoot, "user/", address))
	if err != nil {
		fmt.Println(err)
		log.Errorf("getMinerGeneralInfo: apiClient.Get(%suser/%s); returned err=%s\n", apiRoot, address, err.Error())
		// TODO: handle error
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		fmt.Println(err)
		log.Errorf("getMinerGeneralInfo: json.NewDecoder(resp.Body).Decode(&info); returned err=%s\n", err.Error())
		// TODO: handle error
		return
	}
	return
}

func getMinerPayments(apiRoot string, address string) (payments MinerPayments, err error) {
	resp, err := apiClient.Get(fmt.Sprintf("%s%s%s", apiRoot, "payments/", address))
	if err != nil {
		// TODO: Log error
		fmt.Println(err)
		// TODO: handle error
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&payments)
	if err != nil {
		// TODO: Log error
		fmt.Println(err)
		// TODO: handle error
		return
	}
	return
}

func getMinerShareRate(apiRoot string, address string) (shareRate MinerShareRate, err error) {
	resp, err := apiClient.Get(fmt.Sprintf("%s%s%s", apiRoot, "shareratehistory/", address))
	if err != nil {
		// TODO: Log error
		fmt.Println(err)
		// TODO: handle error
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&shareRate)
	if err != nil {
		// TODO: Log error
		fmt.Println(err)
		// TODO: handle error
		return
	}
	return
}

func calcRPS(balance string, totalPayouts string, totalShares int64) (rps string) {
	bal, _ := strconv.ParseFloat(balance, 64)
	payouts, _ := strconv.ParseFloat(totalPayouts, 64)
	lifetimeEarnings := bal + payouts
	rps = fmt.Sprintf("%f", lifetimeEarnings/float64(totalShares))
	return
}

func calcTotalShares(workers []MinerGeneralInfoWorker) (totalShares int64) {
	totalShares = 0
	for _, w := range workers {
		totalShares += w.Rating
	}
	return
}

func calcSharesPerHour(shareRate []MinerShareRateData, hours *int64) (sharesPerHour int64) {
	shares := int64(0)
	now := time.Now().UTC()
	d := time.Duration(10 * time.Minute)
	hoursAgo := now.Add(time.Duration(*hours) * -1 * time.Hour)
	thePast := hoursAgo.Truncate(d).Unix()
	oldestEntry := now.Unix()
	for _, sr := range shareRate {
		if oldestEntry > thePast {
			if sr.Date < oldestEntry {
				oldestEntry = sr.Date
			}
		}
		if sr.Date > thePast {
			shares += sr.Shares
		}
	}
	if oldestEntry > thePast {
		*hours = int64(math.Round(float64(now.Unix()-oldestEntry) / 60 / 60))
	}
	sharesPerHour = int64(math.Round(float64(shares / *hours)))
	return
}

func calcRewardPerHour(rewardPerShare string, sharesPerHour int64) (rewardPerHour string) {
	rps, _ := strconv.ParseFloat(rewardPerShare, 64)
	rph := rps * float64(sharesPerHour)
	rewardPerHour = fmt.Sprintf("%f", rph)
	return
}

func calcTotalPayout(payments []MinerPaymentsData) (totalPayout string) {
	if len(payments) == 0 {
		return "0"
	}
	payout := float64(0)
	for _, p := range payments {
		payout += p.Amount
	}
	totalPayout = fmt.Sprintf("%f", payout)
	return
}
