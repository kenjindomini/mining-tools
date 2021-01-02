package miningtools

import (
	"mining-tools/nanopool"
	"testing"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func Test_calcSharesPerHour(t *testing.T) {
	simple01 := nanopool.GenerateShareRate(true, 10, 24, true)
	lackingHistory01 := nanopool.GenerateShareRate(true, 10, 10, false)
	type args struct {
		shareRate []nanopool.MinerShareRateData
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
