package models

import "time"

// ScraperConfig is a value object for scraper configuration
type ScraperConfig struct {
	maxConcurrent        int
	timeout              time.Duration
	retryCount           int
	userAgents           []string
	requestsPerSecond    int
	delayBetweenRequests time.Duration
	maxLinksToCollect    int
}

// NewScraperConfig creates a new ScraperConfig with defaults
func NewScraperConfig() *ScraperConfig {
	return &ScraperConfig{
		maxConcurrent:        4,
		timeout:              10 * time.Second,
		retryCount:           3,
		requestsPerSecond:    2,
		delayBetweenRequests: 500 * time.Millisecond,
		maxLinksToCollect:    4,
		userAgents: []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		},
	}
}

// MaxConcurrent returns the maximum concurrent requests
func (c *ScraperConfig) MaxConcurrent() int {
	return c.maxConcurrent
}

// Timeout returns the request timeout
func (c *ScraperConfig) Timeout() time.Duration {
	return c.timeout
}

// RetryCount returns the retry count
func (c *ScraperConfig) RetryCount() int {
	return c.retryCount
}

// UserAgents returns the list of user agents
func (c *ScraperConfig) UserAgents() []string {
	return c.userAgents
}

// RequestsPerSecond returns the rate limit
func (c *ScraperConfig) RequestsPerSecond() int {
	return c.requestsPerSecond
}

// DelayBetweenRequests returns the delay between requests
func (c *ScraperConfig) DelayBetweenRequests() time.Duration {
	return c.delayBetweenRequests
}

// MaxLinksToCollect returns the maximum links to collect
func (c *ScraperConfig) MaxLinksToCollect() int {
	return c.maxLinksToCollect
}

// WithMaxConcurrent sets the maximum concurrent requests
func (c *ScraperConfig) WithMaxConcurrent(max int) *ScraperConfig {
	c.maxConcurrent = max
	return c
}

// WithTimeout sets the request timeout
func (c *ScraperConfig) WithTimeout(timeout time.Duration) *ScraperConfig {
	c.timeout = timeout
	return c
}

// WithRetryCount sets the retry count
func (c *ScraperConfig) WithRetryCount(count int) *ScraperConfig {
	c.retryCount = count
	return c
}

// WithRequestsPerSecond sets the rate limit
func (c *ScraperConfig) WithRequestsPerSecond(rps int) *ScraperConfig {
	c.requestsPerSecond = rps
	return c
}
