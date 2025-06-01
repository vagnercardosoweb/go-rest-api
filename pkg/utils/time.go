package utils

import "time"

var (
	LocationBrl, _ = time.LoadLocation(TimezoneBrl)
	TimezoneBrl    = "America/Sao_Paulo"
)

func NowBrl() time.Time {
	return time.Now().In(LocationBrl)
}

func NowUtc() time.Time {
	return time.Now().UTC()
}
