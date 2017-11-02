package checkers

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/valyala/fasthttp"
	"os"
)

type HTTPChecker struct {
	HealthChecker HealthChecker
	URL           string
}

// Create new HttpChecker by URL
func NewHTTPChecker(name, url string) HTTPChecker {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	logger.Log("msg", "Start HttpChecker")
	return HTTPChecker{
		HealthChecker: NewHealthChecker(name),
		URL:           url,
	}
}

func (c *HTTPChecker) Name() string {
	return c.HealthChecker.Name
}

func (c *HTTPChecker) Check() (Health, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(c.URL)

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}

	err := client.Do(req, resp)

	c.HealthChecker.Up()
	if err != nil {
		logger.Log("err", err)
		c.HealthChecker.Down()
		c.HealthChecker.Msg = err.Error()
		return c.HealthChecker.Health, err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		logger.Log("httpStatus", resp.StatusCode())
		c.HealthChecker.Down()
		c.HealthChecker.Msg = fmt.Sprintf("%s Status is %d", c.Name(), resp.StatusCode())
		return c.HealthChecker.Health, err
	}

	return c.HealthChecker.Health, nil
}
