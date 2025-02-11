package session

import (
	logger "bedouin/bedouin/logging"
	tracing "bedouin/bedouin/tracing"
	"go.uber.org/zap"
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
	// Log request details before sending
	logger.Log.Info("Sending API request",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.Any("headers", req.Header),
	)

	startTime := time.Now()
	resp, err := c.Client.Do(req)

	if err != nil {
		logger.Log.Error("Failed to send request", zap.String("method", req.Method), zap.Error(err))
		return nil, err
	}

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

	logger.Log.Info("Received API response",
		zap.Int("status", resp.StatusCode),
		zap.Int64("response_time_ms", time.Since(startTime).Milliseconds()),
	)

	return resp, err
}

func (c *HttpSession) GetAggStats() *tracing.AggStats {
	return c.trace.GetAggStats()
}
