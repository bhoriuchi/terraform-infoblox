package main

import (
	"github.com/go-resty/resty"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"fmt"
	"time"
)

func dataSourceInfobloxHostRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInfobloxReadHostRecords,

		Schema: map[string]*schema.Schema{
			"records": {
				Type: schema.TypeList,
				Computed: true,
				Description: "Records list",
				Elem: &schema.Schema{Type: schema.TypeMap},
			},
			"names": {
				Type: schema.TypeList,
				Computed: true,
				Description: "Name list",
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"ips": {
				Type: schema.TypeList,
				Computed: true,
				Description: "IP list",
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"query_string": {
				Type: schema.TypeString,
				Required: true,
				Description: "Infoblox query string",
			},
		},
	}
}


func dataSourceInfobloxReadHostRecords(d *schema.ResourceData, meta interface{}) error {
	queryString	:= d.Get("query_string").(string)
	wapiErr 	:= WapiError{}
	response 	:= []Host{}

	d.SetId(time.Now().UTC().String())

	resp, err := resty.R().
		SetResult(&response).
		SetError(&wapiErr).
		SetQueryString(queryString).
		Get("/record:host")

	if handler := handleError(err, resp, wapiErr); handler != nil {
		return handler
	}

	log.Printf("[infoblox-provider] Response: %v", response)

	records := make([]map[string]interface{}, len(response))
	names	:= make([]string, len(response))
	ips	:= make([]string, len(response))

	for i, v := range response {
		record := make(map[string]interface{})
		record["ref"] 		= v.Object.Ref
		record["name"] 		= v.Name
		record["ipv4addr"] 	= v.Ipv4addrs[0].Ipv4addr
		record["view"] 		= v.View
		record["ttl"] 		= v.Ttl
		record["use_ttl"] 	= v.Use_Ttl

		records[i] 	= record
		names[i]	= v.Name
		ips[i]		= v.Ipv4addrs[0].Ipv4addr
	}

	if err := d.Set("records", records); err != nil {
		return fmt.Errorf("[infoblox-provider] Error setting records")
	}
	return nil
}