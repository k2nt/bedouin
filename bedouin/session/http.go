package session

import (
	tracing "bedouin/bedouin/tracing"
	"net/http"
	"sync"
	"time"
)

type HttpSession struct {
	Client *http.Client
	trace  tracing.Trace
	mutex  sync.Mutex
}

var DefaultHttpSession = &HttpSession{
	Client: http.DefaultClient,
}

func (c *HttpSession) Submit(req *http.Request) (*http.Response, error) {
	startTime := time.Now()
	resp, err := c.Client.Do(req)
	endTime := time.Now()

	c.trace.Add(
		tracing.Stat{
			StartTime:     startTime,
			EndTime:       endTime,
			Req:           req,
			Response:      resp,
			ResponseError: err,
		},
	)

	return resp, err
}

func (c *HttpSession) GetAggStats() *tracing.AggStats {
	return c.trace.GetAggStats()
}
