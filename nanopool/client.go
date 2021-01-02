package nanopool

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// HTTPClient is an interface to abstract http.client to support testing using mocks
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

var (
	apiClient = HTTPClient(&http.Client{Timeout: 10 * time.Second})
)

func get(fullPath string, output interface{}) (err error) {
	log.Debugf("get(fullPath=%s, output interface{}) called\n", fullPath)
	resp, err := apiClient.Get(fullPath)
	if err != nil {
		fmt.Println(err)
		log.Errorf("get: apiClient.Get(%s); returned err=%s\n", fullPath, err.Error())
		// TODO: handle error
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(output)
	if err != nil {
		fmt.Println(err)
		log.Errorf("get: json.NewDecoder(resp.Body).Decode(output); returned err=%s\n", err.Error())
		// TODO: handle error
		return
	}
	return
}

// GetMinerGeneralInfo calls the Miner:General Info endpoint user/:address and forms the response in to a usable Struct
func GetMinerGeneralInfo(apiRoot string, address string) (info MinerGeneralInfo, err error) {
	log.Debugln("GetMinerGeneralInfo called")
	resp, err := apiClient.Get(fmt.Sprintf("%s%s%s", apiRoot, "user/", address))
	if err != nil {
		fmt.Println(err)
		log.Errorf("GetMinerGeneralInfo: apiClient.Get(%suser/%s); returned err=%s\n", apiRoot, address, err.Error())
		// TODO: handle error
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		fmt.Println(err)
		log.Errorf("GetMinerGeneralInfo: json.NewDecoder(resp.Body).Decode(&info); returned err=%s\n", err.Error())
		// TODO: handle error
		return
	}
	return
}

// GetMinerPayments calls the Miner:Payments endpoint payments/:address and forms the response in to a usable Struct
func GetMinerPayments(apiRoot string, address string) (payments MinerPayments, err error) {
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

// GetMinerShareRate calls the Miner:Share Rate History endpoint shareratehistory/:address and forms the response in to a usable Struct
func GetMinerShareRate(apiRoot string, address string) (shareRate MinerShareRate, err error) {
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

// GetMinerBalance calls the Miner:Balance endpoint balance/:address and forms the response in to a usable Struct
func GetMinerBalance(apiRoot string, address string) (minerBalance MinerBalance, err error) {
	resp, err := apiClient.Get(fmt.Sprintf("%s%s%s", apiRoot, "balance/", address))
	if err != nil {
		fmt.Println(err)
		log.Errorf("GetMinerBalance: apiClient.Get(%sbalance/%s); returned err=%s\n", apiRoot, address, err.Error())
		// TODO: handle error
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&minerBalance)
	if err != nil {
		fmt.Println(err)
		log.Errorf("GetMinerBalance: json.NewDecoder(resp.Body).Decode(&minerBalance); returned err=%s\n", err.Error())
		// TODO: handle error
		return
	}
	return
}

// GetOtherPrices calls the Miner:Balance endpoint prices/ and forms the response in to a usable Struct
func GetOtherPrices(apiRoot string) (prices OtherPrices, err error) {
	fullPath := fmt.Sprintf("%s%s", apiRoot, "prices/")
	err = get(fullPath, &prices)
	// resp, err := apiClient.Get(fmt.Sprintf("%s%s", apiRoot, "prices/"))
	// if err != nil {
	// 	fmt.Println(err)
	// 	log.Errorf("GetOtherPrices: apiClient.Get(%sbalance/); returned err=%s\n", apiRoot, err.Error())
	// 	// TODO: handle error
	// 	return
	// }
	// defer resp.Body.Close()
	// err = json.NewDecoder(resp.Body).Decode(&prices)
	// if err != nil {
	// 	fmt.Println(err)
	// 	log.Errorf("GetOtherPrices: json.NewDecoder(resp.Body).Decode(&prices); returned err=%s\n", err.Error())
	// 	// TODO: handle error
	// 	return
	// }
	return
}
