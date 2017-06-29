package main

import (
	"github.com/go-resty/resty"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"fmt"
	"time"
)

func dataSourceInfobloxAnameRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInfobloxReadAnameRecords,

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
				Optional: true,
				Description: "Infoblox query string",
			},
			"free_name_prefix": {
				Type: schema.TypeString,
				Optional: true,
				Description: "Prefix to use for name generation",
			},
			"free_name_domain": {
				Type: schema.TypeString,
				Optional: true,
				Description: "Domain to use for name generation",
			},
			"free_name_pad": {
				Type: schema.TypeString,
				Optional: true,
				Description: "Pad to use for name generation",
			},
			"free_name_start": {
				Type: schema.TypeInt,
				Optional: true,
				Description: "Start to use for name generation",
				Default: 1,
			},
			"free_name_limit": {
				Type: schema.TypeInt,
				Optional: true,
				Description: "Limit number of names generated",
				Default: 20,
			},
		},
	}
}

func dataSourceInfobloxReadAnameRecords(d *schema.ResourceData, meta interface{}) error {
	queryString	:= d.Get("query_string").(string)
	prefix		:= d.Get("free_name_prefix").(string)
	domain		:= d.Get("free_name_domain").(string)
	pad		:= d.Get("free_name_pad").(string)
	start		:= d.Get("free_name_start").(int)
	limit		:= d.Get("free_name_limit").(int)
	wapiErr 	:= WapiError{}
	response 	:= []A{}
	records 	:= make([]map[string]string, 0)
	names		:= make([]string, 0)
	ips		:= make([]string, 0)

	log.Printf("[infoblox-provider] QueryString: %s", queryString)

	d.SetId(time.Now().UTC().String())

	if prefix != "" && pad != "" && domain != "" {
		namez, err := getAvailableHostnames(prefix, pad, domain, "a", start, limit)
		if err != nil {
			return err
		}
		names = namez
	} else {
		resp, err := resty.R().
			SetResult(&response).
			SetError(&wapiErr).
			SetQueryString(queryString).
			Get("/record:a")

		if handler := handleError(err, resp, wapiErr); handler != nil {
			return handler
		}

		log.Printf("[infoblox-provider] Response: %v", response)

		for _, v := range response {
			record := make(map[string]string)
			record["ref"] 		= v.Object.Ref
			record["name"] 		= v.Name
			record["ipv4addr"] 	= v.Ipv4addr
			record["dns_name"] 	= v.Dns_name
			record["comment"] 	= v.Comment
			record["view"] 		= v.View
			record["zone"] 		= v.Zone

			records = append(records, record)
			names 	= append(names, v.Name)
			ips 	= append(ips, v.Ipv4addr)
		}
	}

	if err := d.Set("records", records); err != nil {
		return fmt.Errorf("[infoblox-provider] Error setting records")
	}
	if err := d.Set("names", names); err != nil {
		return fmt.Errorf("[infoblox-provider] Error setting names")
	}
	if err := d.Set("ips", ips); err != nil {
		return fmt.Errorf("[infoblox-provider] Error setting ips")
	}
	return nil
}