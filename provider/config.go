package main

import (
	"crypto/tls"
	"time"

	"github.com/go-resty/resty"
)

func (c *Config) Client() (*Config, error) {
	resty.
		SetHostURL(c.InfobloxEndpoint).
		SetBasicAuth(c.User, c.Password).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetTimeout(time.Duration(c.HTTPTimeout) * time.Second)
	if c.InsecureFlag == true {
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	return c, nil
}