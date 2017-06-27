package main

import (
	"github.com/go-resty/resty"
	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
	"net/url"
	"net"
	"strings"
)

const (
	RECORD_A = "a"
	RECORD_AAAA = "aaaa"
	RECORD_CNAME = "cname"
	RECORD_HOST = "host"
)

func dataSourceInfobloxRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInfobloxRead,

		Schema: map[string]*schema.Schema{
			"record_type": {
				Type: schema.TypeString,
				Required: true,
				ValidateFunc: validateRecordType,
			},
			"filter": {
				Type: schema.TypeString,
				Required: true,
			},
		},
	}
}


func dataSourceInfobloxRead(d *schema.ResourceData, meta interface{}) error {
	recordType	:= d.Get("record_type").(string)
	filter		:= url.QueryEscape(d.Get("filter").(string))
	wapiErr 	:= WapiError{}
	var rec interface{}

	if recordType == RECORD_A {
		rec = A{}
	} else if recordType == RECORD_AAAA {
		rec = Aaaa{}
	} else if recordType == RECORD_CNAME {
		rec = Cname{}
	} else if recordType == RECORD_HOST {
		rec = Host{}
	}

	resp, err := resty.R().
		SetResult(&rec).
		SetError(&wapiErr).
		Get("/record:" + recordType + "?" + filter)

	if handler := handleError(err, resp, wapiErr); handler != nil {
		return handler
	}
	d.SetId(strings.Replace(resp.String(), "\"", "", 2))
	return nil
}

func validateRecordType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	var validType map[string]bool
	validType[RECORD_A] = true
	validType[RECORD_AAAA] = true
	validType[RECORD_CNAME] = true
	validType[RECORD_HOST] = true

	if !validType[value] {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid record type. Valid states are: %q, %q, %q, and %q",
		k, RECORD_A, RECORD_AAAA, RECORD_CNAME, RECORD_HOST))
	}
	return
}