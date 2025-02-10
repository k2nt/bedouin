package stats

import (
	"net/http"
	"slices"
	"sort"
	"time"
)

type Stat struct {
	StartTime     time.Time
	EndTime       time.Time
	Req           *http.Request
	Response      *http.Response
	ResponseError error
}

type LatencyStats struct {
	Average time.Duration
	Min     time.Duration
	Max     time.Duration
	Q1      time.Duration
	Q2      time.Duration
	Q3      time.Duration
	P99     time.Duration
}

type AggStats struct {
	LatencyStats
}

type Trace struct {
	trace []Stat
}

func (t *Trace) Add(s Stat) {
	t.trace = append(t.trace, s)
}

func (t *Trace) getLatencyStats() LatencyStats {
	var latencies []time.Duration
	for _, stat := range t.trace {
		latency := stat.EndTime.Sub(stat.StartTime)
		latencies = append(latencies, latency)
	}

	sort.Slice(
		latencies, func(i, j int) bool {
			return latencies[i] < latencies[j]
		},
	)

	total := len(latencies)
	average := time.Duration(0)
	for _, latency := range latencies {
		average += latency
	}
	average /= time.Duration(total)

	q1 := latencies[int(0.25*float64(total))]
	q2 := latencies[int(0.5*float64(total))]
	q3 := latencies[int(0.75*float64(total))]
	p99 := latencies[int(0.99*float64(total))]

	return LatencyStats{
		Min:     slices.Min(latencies),
		Max:     slices.Max(latencies),
		Average: average,
		Q1:      q1,
		Q2:      q2,
		Q3:      q3,
		P99:     p99,
	}
}

func (t *Trace) GetAggStats() *AggStats {
	return &AggStats{
		LatencyStats: t.getLatencyStats(),
	}
}
