package main

import (
	"crypto/tls"
	"time"

	"github.com/go-resty/resty"
	"log"
	"fmt"
	"regexp"
	"strconv"
	"sort"
)

type Semver struct {
	Major int
	Minor int
	Patch int
}

type Config struct {
	User             string
	Password         string
	InfobloxEndpoint string
	InsecureFlag     bool
	InfobloxVersion  Semver
	HTTPTimeout      int
}

type WapiError struct {
	Error, Code, Text string
}

type Object struct {
	Ref string `json:"_ref"`
}

type Ipv4 struct {
	Object
	Host, Ipv4addr     string
	Configure_for_dhcp bool
}

type Host struct {
	Object
	Name, View string
	Ttl        int
	Use_Ttl    bool
	Ipv4addrs  []Ipv4
}

type A struct {
	Object
	Comment, Name, View, Ipv4addr, Dns_name, Zone 	string
	Disable, Use_ttl				bool
	Ttl 						int
}

type Aaaa struct {
	Object
	Name, View, Ipv6addr, Dns_name 	string
	Ttl				int
}

type Cname struct {
	Object
	Name, View, Canonical, Dns_name string
	Ttl				int
}

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

func handleError(err error, resp *resty.Response, wapiErr WapiError) error {
	log.Printf("\n[infoblox-provider] HTTP Code: (%v) Response Body: %v", resp.StatusCode(), resp)
	if err != nil {
		return fmt.Errorf("[infoblox-provider] Resty Error: %+v", err)
	} else if resp.StatusCode() >= 300 && resp.StatusCode() < 400 {
		return fmt.Errorf("[infoblox-provider] HTTP Redirect: (%v)", resp.StatusCode())
	} else if resp.StatusCode() >= 400 && wapiErr.Error != "" {
		return fmt.Errorf("[infoblox-provider] WAPI Error: (%v) %+v", resp.StatusCode(), wapiErr)
	} else if resp.StatusCode() >= 400 {
		return fmt.Errorf("[infoblox-provider] Unknown HTTP Error: (%v) %+v", resp.StatusCode(), resp.String())
	}
	return nil
}

func getNextHostname(name string, prefix string, pad string, domain string, recordType string) (string, error) {
	wapiErr 	:= WapiError{}
	response 	:= []A{}
	used 		:= make([]int, 0)
	regex 		:= regexp.MustCompile("^(" + prefix + ")(\\d{" + pad + "}).(" + domain + ")$")

	resp, err := resty.R().
		SetResult(&response).
		SetError(&wapiErr).
		SetQueryString("name~=^" + prefix + "\\d{" + pad + "}." + domain + "$").
		Get("/record:" + recordType)

	if handler := handleError(err, resp, wapiErr); handler != nil {
		return nil, handler
	}

	for _, v := range response {
		matches := regex.FindAllStringSubmatch(v.Name, -1)

		if len(matches) > 0 && len(matches[0]) > 3 {
			idx, err := strconv.Atoi(matches[0][2])
			if err != nil {
				return nil, err
			}
			used = append(used, idx)
		}
	}

	sort.Ints(used)



	return nil
}