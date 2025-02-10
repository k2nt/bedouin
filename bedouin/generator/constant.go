package generator

import (
	"bedouin/bedouin/timing"
	"time"
)

type ConstantGenerator struct {
	Field
	ReqPerSec       int32
	DurationSeconds float64
}

func NewConstantGenerator(
	handleFunc func(),
	isAsync bool,
	reqPerSec int32,
	durationSeconds float64,
) *ConstantGenerator {
	if durationSeconds == 0 {
		durationSeconds = timing.InfiniteDuration
	}

	return &ConstantGenerator{
		Field: Field{
			HandleFunc: handleFunc,
			IsAsync:    isAsync,
		},
		ReqPerSec:       reqPerSec,
		DurationSeconds: durationSeconds,
	}
}

func (g *ConstantGenerator) Run() {
	ticker := time.NewTicker(time.Second / time.Duration(g.ReqPerSec))
	defer ticker.Stop()

	startTime := time.Now()
	for time.Since(startTime).Seconds() <= g.DurationSeconds {
		select {
		case <-ticker.C:
			if g.IsAsync {
				go g.HandleFunc()
			} else {
				g.HandleFunc()
			}
		}
	}
}
