package timing

import "time"

const InfiniteDuration = 1<<63 - 1

func Since(t time.Time) float64 {
	return time.Since(t).Seconds()
}
