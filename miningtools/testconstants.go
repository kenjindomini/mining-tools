package miningtools

var (
	user0x01Success = &MinerGeneralInfo{
		Status: true,
		Data: MinerGeneralInfoData{
			Account:            "0x01",
			UnconfirmedBalance: "0.0",
			Balance:            "0.142",
			Workers: []MinerGeneralInfoWorker{
				{
					Rating: 2000,
				},
				{
					Rating: 8000,
				},
			},
		},
	}
	user0x02Error = &NanopoolError{
		Status: false,
		Error:  "Address does not exist",
	}
)
