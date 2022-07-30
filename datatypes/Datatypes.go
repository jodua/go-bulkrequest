package datatypes

import (
	"errors"
	"math/rand"
	"time"
)

type ProxyConfig struct {
	ProxyList        []string
	RequestsPerProxy int
}

type UserAgentConfig struct {
	UserAgentList []string
}

type DelayConfig struct {
	DelayMin int
	DelayMax int
}

func (c *UserAgentConfig) GetRandomUserAgent() string {
	return c.UserAgentList[rand.Intn(len(c.UserAgentList))]
}

func (c *DelayConfig) GetRandomDelay() time.Duration {
	delay := rand.Intn(c.DelayMax-c.DelayMin) + c.DelayMin
	return time.Duration(delay) * time.Millisecond
}

func (c *DelayConfig) Validate() error {
	if c.DelayMin > c.DelayMax {
		return errors.New("DelayMin must be less than DelayMax")
	}
	if c.DelayMin < 0 || c.DelayMax < 0 {
		return errors.New("DelayMin and DelayMax must be greater than 0")
	}
	return nil
}

func (c *ProxyConfig) Validate() error {
	if len(c.ProxyList) == 0 {
		return errors.New("ProxyList is empty")
	}
	if c.RequestsPerProxy < 1 {
		return errors.New("RequestsPerProxy must be greater than 0")
	}
	return nil
}

func (c *UserAgentConfig) Validate() error {
	if len(c.UserAgentList) == 0 {
		return errors.New("UserAgentList is empty")
	}
	return nil
}
