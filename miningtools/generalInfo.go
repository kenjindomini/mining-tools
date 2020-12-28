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
	ID        string      `json:"id"`
	UID       json.Number `json:"uid"`
	Hashrate  string      `json:"hashrate"`
	Lastshare json.Number `json:"lastshare"`
	Rating    json.Number `json:"rating"`
	H1        string      `json:"h1"`
	H3        string      `json:"h3"`
	H6        string      `json:"h6"`
	H12       string      `json:"h12"`
	H24       string      `json:"h24"`
}

// MinerShareRate is for decoding json from a successful response of the nanopool miner general info api endpoint
type MinerShareRate struct {
	Status bool                 `json:"status"`
	Data   []MinerShareRateData `json:"data"`
}

// MinerShareRateData is for decoding json from a successful response of the nanopool miner general info api endpoint
type MinerShareRateData struct {
	Date   json.Number `json:"date"`
	Shares json.Number `json:"shares"`
}

// generalInfoCmd represents the generalInfo sub-command of nanopool
var (
	apiClient          = &http.Client{Timeout: 10 * time.Second}
	rewardPerShareFlag bool
	sharesPerHour      bool
	generalInfoCmd     = &cobra.Command{
		Use:   "generalInfo",
		Short: "Gets general info of nanopool ethereum miner account",
		Long: `Gets general info of nanopool ethereum miner account, with
				options for calculating additional values`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("generalInfo called")
			address := viper.GetString("nanopool.address")
			apiRoot := viper.GetString("nanopool.apiRoot")
			info, err := getMinerGeneralInfo(apiRoot, address)
			if err != nil {
				// TODO: Log error
				fmt.Println(err)
				// TODO: handle error
				return
			}
			if rewardPerShareFlag {
				balance, _ := strconv.ParseFloat(info.Data.Balance, 64)
				totalPayouts, _ := strconv.ParseFloat("0", 64)
				totalShares := calcTotalShares(info.Data.Workers)
				info.Data.RewardPerShare = calcRPS(balance, totalPayouts, totalShares)
			}
			if sharesPerHour {
				shareRate, err := getMinerShareRate(apiRoot, address)
				if err != nil {
					// TODO: Log error
					fmt.Println(err)
					// TODO: handle error
					return
				}
				info.Data.SharesPerHour = calcSharesPerHour(shareRate.Data, 24)
			}
			prettyPrint(info)
		},
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
	generalInfoCmd.Flags().BoolVarP(&rewardPerShareFlag, "rewardPerShare", "r", false, "Include calculated attribute rewardPerShare (lifetime average)")
	generalInfoCmd.Flags().BoolVarP(&sharesPerHour, "sharesPerHour", "s", false, "Include calculated attribute sharesPerHour (24h rolling average)")
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
	resp, err := apiClient.Get(fmt.Sprintf("%s%s%s", apiRoot, "user/", address))
	if err != nil {
		// TODO: Log error
		fmt.Println(err)
		// TODO: handle error
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		// TODO: Log error
		fmt.Println(err)
		// TODO: handle error
		return
	}
	return
}

func getMinerPayments(apiRoot string, address string) (resp string, err error) {
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

func calcRPS(balance float64, totalPayOuts float64, totalShares int64) (rps string) {
	lifetimeEarnings := balance + totalPayOuts
	rps = fmt.Sprintf("%f", lifetimeEarnings/float64(totalShares))
	return
}

func calcTotalShares(workers []MinerGeneralInfoWorker) (totalShares int64) {
	totalShares = 0
	for _, w := range workers {
		rating, _ := w.Rating.Int64()
		totalShares += rating
	}
	return
}

func calcSharesPerHour(shareRate []MinerShareRateData, hours int64) (sharesPerHour int64) {
	shares := int64(0)
	srLen := int64(len(shareRate))
	// Share Rate is returned in 10 minute segments and is ordered oldest to newest
	// len - 1 gets us the index of the last element, then there are 6 elements per hour
	sliceLow := srLen - 1 - (6 * hours)
	if hours <= 0 {
		sliceLow = 0
		hours = int64(math.Round(float64(srLen / 6)))
	}
	for _, sr := range shareRate[sliceLow:] {
		srShares, _ := sr.Shares.Int64()
		shares += srShares
	}
	sharesPerHour = int64(math.Round(float64(shares / hours)))
	return
}
