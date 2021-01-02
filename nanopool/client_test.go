package nanopool

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type mockAPIClient struct {
	mock.Mock
}

func init() {
	log.SetLevel(log.DebugLevel)
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
	success01 := GenerateShareRate(true, 10, 24, true)
	success02 := GenerateShareRate(true, 10, 24, true)
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
			gotShareRate, err := GetMinerShareRate(tt.args.apiRoot, tt.args.address)
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
			gotInfo, err := GetMinerGeneralInfo(tt.args.apiRoot, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMinerGeneralInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotInfo, tt.wantInfo) {
				t.Errorf("GetMinerGeneralInfo() = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}
