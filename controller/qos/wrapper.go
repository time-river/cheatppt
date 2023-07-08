package qos

import "cheatppt/controller/billing"

// 3 / per mins
var FreeLimiter = NewIDRateLimiter(0.05, 3)

// 12 / per mins
var PaidLimiter = NewIDRateLimiter(0.2, 12)

type Meta struct {
	*billing.Consumer
	Model    string
	Provider string
}

func Allow(meta *Meta) bool {
	if meta.Free {
		return FreeLimiter.Allow(uint(meta.UserId))
	} else {
		return PaidLimiter.Allow(uint(meta.UserId))
	}
}
