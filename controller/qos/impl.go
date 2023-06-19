package qos

import "cheatppt/controller/charge"

// 3 / per mins
var FreeLimiter = NewIDRateLimiter(0.05, 3)

var PaidLimiter = NewIDRateLimiter(0.2, 12)

type Meta struct {
	charge.Consumer
	Model    string
	Provider string
}

func Allow(meta *Meta) bool {

	return true
}
