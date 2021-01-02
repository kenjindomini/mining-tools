package nanopool

// ErrorResponse is a struct for marshaling json of any error coming from nanopool's API
type ErrorResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

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

// MinerBalance is for decoding json from a successful response of the nanopool miner balance api endpoint
type MinerBalance struct {
	Status bool    `json:"status"`
	Data   float64 `json:"data"`
}

// OtherPrices is for decoding json from a successful response of the nanopool other prices api endpoint
type OtherPrices struct {
	Status bool            `json:"status"`
	Data   OtherPricesData `json:"data"`
}

// OtherPricesData is for decoding json from a successful response of the nanopool other prices api endpoint
type OtherPricesData struct {
	PriceUSD float64 `json:"price_usd"`
	PriceEUR float64 `json:"price_eur"`
	PriceRUR float64 `json:"price_rur"`
	PriceCNY float64 `json:"price_cny"`
	PriceBTC float64 `json:"price_btc"`
}
