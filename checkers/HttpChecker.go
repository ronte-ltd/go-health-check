package checkers

import (
	"fmt"
	"net/http"
)

//HTTPChecker is truct provide HealthChecker and URL for checking
type HTTPChecker struct {
	HealthChecker HealthChecker
	URL           string
}

// NewHTTPChecker return new HttpChecker by URL
func NewHTTPChecker(name, url string) *HTTPChecker {
	return &HTTPChecker{
		HealthChecker: NewHealthChecker(name),
		URL:           url,
	}
}

//Name return name current cehcker
func (c *HTTPChecker) Name() string {
	return c.HealthChecker.Health.Name
}

//Check send GET request on URL and return Healt or err
func (c *HTTPChecker) Check() (Health, error) {
	resp, err := http.Get(c.URL)
	checkError(err)

	c.HealthChecker.Up()
	if err != nil {
		checkError(err)
		c.HealthChecker.Down()
		c.HealthChecker.Msg = err.Error()
		return c.HealthChecker.Health, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("httpStatus", resp.StatusCode)
		c.HealthChecker.Down()
		c.HealthChecker.Msg = fmt.Sprintf("%s Status is %d", c.Name(), resp.StatusCode)
		return c.HealthChecker.Health, nil
	}

	return c.HealthChecker.Health, nil
}
