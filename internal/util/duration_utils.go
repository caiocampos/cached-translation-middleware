package util

import (
	"math/rand/v2"
	"time"
)

const TenSeconds = 10 * time.Second
const TwentySeconds = 20 * time.Second

const OneDay = 24 * time.Hour
const SevenDays = 7 * OneDay

const TenSecondsString = "10s"
const TwentySecondsString = "20s"

const OneDayString = "24h"
const SevenDaysString = "168h"

func GenerateRandomCacheTime(minDays int, maxDays int) time.Duration {
	days := rand.IntN(maxDays-minDays+1) + minDays
	return time.Duration(days) * OneDay
}
