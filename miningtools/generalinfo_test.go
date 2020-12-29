package miningtools

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

type mockAPIClient struct {
	mock.Mock
}

func (mac *mockAPIClient) Get(url string) (resp *http.Response, err error) {
	args := mac.Called(url)
	body, _ := json.Marshal(args.Get(0))
	resp = &http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer(body)),
	}
	return resp, args.Error(1)
}

func Test_getMinerShareRate(t *testing.T) {
	// Set up
	mockClient := &mockAPIClient{}
	apiClient = mockClient
	mockClient.On("Get", "http://test.com/shareratehistory/0x01").Return(shareratehistory0x01Success, nil)
	mockClient.On("Get", "http://test.com/shareratehistory/").Return(*new(MinerShareRate), errors.New("Timedout"))
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
			wantShareRate: *shareratehistory0x01Success,
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
				shareRate: shareratehistory0x01Success.Data,
				hours:     24,
			},
			wantSharesPerHour: 60,
			wantHours:         24,
		},
		{
			name: "AllHours01",
			args: args{
				shareRate: shareratehistory0x01Success.Data,
				hours:     -1,
			},
			wantSharesPerHour: 12048,
			wantHours:         30,
		},
		{
			name: "NotEnoutHistory01",
			args: args{
				shareRate: shareratehistory0x01Success.Data,
				hours:     45,
			},
			wantSharesPerHour: 12048,
			wantHours:         30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hours := tt.args.hours
			if gotSharesPerHour := calcSharesPerHour(tt.args.shareRate, &hours); gotSharesPerHour != tt.wantSharesPerHour || hours != tt.wantHours {
				if gotSharesPerHour != tt.wantSharesPerHour {
					t.Errorf("calcSharesPerHour(tt.args.shareRate, hours) = %v, want %v", gotSharesPerHour, tt.wantSharesPerHour)
				}
				if tt.args.hours != tt.wantHours {
					t.Errorf("calcSharesPerHour(tt.args.shareRate, tt.args.hours) hours= %v, want %v", hours, tt.wantHours)
				}
			}
		})
	}
}
