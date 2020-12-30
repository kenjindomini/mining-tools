package miningtools

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type mockAPIClient struct {
	mock.Mock
}

func init() {
	logLevel = 5
}

func (mac *mockAPIClient) Get(url string) (resp *http.Response, err error) {
	args := mac.Called(url)
	body, _ := json.Marshal(args.Get(0))
	resp = &http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer(body)),
	}
	return resp, args.Error(1)
}

func generateShareRate(status bool, shares int64, hoursRequired int64, addExtraElements bool) (shareRate MinerShareRate) {
	shareRate.Status = status
	now := time.Now().UTC()
	d := time.Duration(10 * time.Minute)
	epoch := now.Truncate(d).Unix()
	hoursAgo := now.Add(time.Duration(hoursRequired) * -1 * time.Hour)
	thePast := hoursAgo.Truncate(d).Unix()
	elements := hoursRequired * 6
	if addExtraElements {
		elements += rand.Int63n(100)
	}
	for i := int64(0); i < elements; i++ {
		if epoch > thePast {
			msrd := MinerShareRateData{
				Date:   epoch,
				Shares: shares,
			}
			shareRate.Data = append(shareRate.Data, msrd)
		} else {
			msrd := MinerShareRateData{
				Date:   epoch,
				Shares: rand.Int63n(10000),
			}
			shareRate.Data = append(shareRate.Data, msrd)
		}
		epoch -= 600
	}
	return
}

func Test_getMinerShareRate(t *testing.T) {
	// Set up
	mockClient := &mockAPIClient{}
	apiClient = mockClient
	success01 := generateShareRate(true, 10, 24, true)
	success02 := generateShareRate(true, 10, 24, true)
	mockClient.On("Get", "http://test.com/shareratehistory/0x01").Return(success01, nil)
	mockClient.On("Get", "http://test.com/shareratehistory/").Return(*new(MinerShareRate), errors.New("Timedout"))
	mockClient.On("Get", "http://test.com/shareratehistory/0x02").Return(success02, nil)
	type args struct {
		apiRoot string
		address string
	}
	tests := []struct {
		name          string
		args          args
		wantShareRate MinerShareRate
		wantErr       bool
	}{
		{
			name: "Success01",
			args: args{
				apiRoot: "http://test.com/",
				address: "0x01",
			},
			wantShareRate: success01,
			wantErr:       false,
		},
		{
			name: "Error01",
			args: args{
				apiRoot: "http://test.com/",
				address: "",
			},
			wantShareRate: *new(MinerShareRate),
			wantErr:       true,
		},
		{
			name: "Success02",
			args: args{
				apiRoot: "http://test.com/",
				address: "0x02",
			},
			wantShareRate: success02,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotShareRate, err := getMinerShareRate(tt.args.apiRoot, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("getMinerShareRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotShareRate, tt.wantShareRate) {
				t.Errorf("getMinerShareRate() = %v, want %v", gotShareRate, tt.wantShareRate)
			}
		})
	}
}

func Test_getMinerGeneralInfo(t *testing.T) {
	// Set up
	mockClient := &mockAPIClient{}
	apiClient = mockClient
	mockClient.On("Get", "http://test.com/user/0x01").Return(user0x01Success, nil)
	type args struct {
		apiRoot string
		address string
	}
	tests := []struct {
		name     string
		args     args
		wantInfo MinerGeneralInfo
		wantErr  bool
	}{
		{
			name: "Success01",
			args: args{
				apiRoot: "http://test.com/",
				address: "0x01",
			},
			wantInfo: *user0x01Success,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfo, err := getMinerGeneralInfo(tt.args.apiRoot, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("getMinerGeneralInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotInfo, tt.wantInfo) {
				t.Errorf("getMinerGeneralInfo() = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}

func Test_calcSharesPerHour(t *testing.T) {
	simple01 := generateShareRate(true, 10, 24, true)
	lackingHistory01 := generateShareRate(true, 10, 10, false)
	type args struct {
		shareRate []MinerShareRateData
		hours     int64
	}
	tests := []struct {
		name              string
		args              args
		wantSharesPerHour int64
		wantHours         int64
	}{
		{
			name: "Simple01",
			args: args{
				shareRate: simple01.Data,
				hours:     24,
			},
			wantSharesPerHour: 60,
			wantHours:         24,
		},
		{
			name: "LackingHistory01",
			args: args{
				shareRate: lackingHistory01.Data,
				hours:     24,
			},
			wantSharesPerHour: 60,
			wantHours:         10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hours := tt.args.hours
			if gotSharesPerHour := calcSharesPerHour(tt.args.shareRate, &hours); gotSharesPerHour != tt.wantSharesPerHour || hours != tt.wantHours {
				if gotSharesPerHour != tt.wantSharesPerHour {
					t.Errorf("calcSharesPerHour(tt.args.shareRate, hours) = %v, want %v", gotSharesPerHour, tt.wantSharesPerHour)
				}
				if hours != tt.wantHours {
					t.Errorf("calcSharesPerHour(tt.args.shareRate, hours) hours= %v, want %v", hours, tt.wantHours)
				}
			}
		})
	}
}
