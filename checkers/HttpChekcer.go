package checkers

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/valyala/fasthttp"
	"os"
)

var (
	logger log.Logger
)

type HttpChecker struct {
	HealthChecker HealthChecker
	Url           string
}

func NewHttpChecker(name, url string) HttpChecker {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	return HttpChecker{
		HealthChecker: NewHealthChecker(name),
		Url:           url,
	}
}

func (c *HttpChecker) Name() string {
	return c.HealthChecker.Name
}

func (c *HttpChecker) Check() (Health, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(c.Url)

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

func (c *HttpChecker) AddChecker(adding Checker) {
	c.HealthChecker.Checkers[adding.Name()] = adding
}
