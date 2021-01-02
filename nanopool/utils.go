package nanopool

import (
	"math/rand"
	"time"
)

// GenerateShareRate is used to generate a MinerShareRate struct for testing
func GenerateShareRate(status bool, shares int64, hoursRequired int64, addExtraElements bool) (shareRate MinerShareRate) {
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
