package session

import (
	stats "bedouin/bedouin/tracing"
	"net/http"
)

type Session interface {
	Submit(req *http.Request) (*http.Response, error)
	GetAggregateStats() *stats.AggStats
}
