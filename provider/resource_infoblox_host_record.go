package main

import (
	"log"
	"strings"

	"github.com/go-resty/resty"
	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
)

func resourceInfobloxHostRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceInfobloxHostRecordCreate,
		Read:   resourceInfobloxHostRecordRead,
		Update: resourceInfobloxHostRecordUpdate,
		Delete: resourceInfobloxHostRecordDelete,

		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Description: "The domain name to create these records in",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": &schema.Schema{
				Description: "The subdomain of the record",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"name_prefix": &schema.Schema{
				Description: "Name generation prefix",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"name_index_pad": &schema.Schema{
				Description: "Name generation index padding length. Pads with leading 0s",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"name_index_start": &schema.Schema{
				Description: "Name generation index start",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     1,
			},
			"ipv4": &schema.Schema{
				Description: "The ip-address or function used to generate one",
				Type:        schema.TypeString,
				Required:    true,
			},
			"ttl": &schema.Schema{
				Description: "The TTL of the DNS record in seconds, used for client-cache invalidation",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     600,
			},
		},
	}
}

func resourceInfobloxHostRecordCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("\n[infoblox-provider] %s", "----------------- host record create")

	var name string
	prefix := d.Get("name_prefix").(string)
	pad    := d.Get("name_index_pad").(string)
	domain := d.Get("domain").(string)
	start  := d.Get("name_index_start").(int)
	ipv4   := d.Get("ipv4").(string)
	ttl    := d.Get("ttl").(int)

	if prefix != "" && pad != "" {
		names, err := getAvailableHostnames(prefix, pad, domain, "a", start, 1)
		if err != nil {
			return err
		}
		if len(names) == 0 {
			return fmt.Errorf("Could not find a free hostname")
		}
		name = names[0]
	} else {
		name = d.Get("name").(string)
	}

	fqdn    := name + "." + domain
	wapiErr := WapiError{}

	resp, err := resty.R().
		SetError(&wapiErr).
		SetBody(map[string]interface{}{
			"name": fqdn,
			"ipv4addrs": []map[string]interface{}{
				map[string]interface{}{
					"ipv4addr": ipv4,
				},
			},
			"ttl":     ttl,
			"use_ttl": true,
		}).
		Post("/record:host")
	if handler := handleError(err, resp, wapiErr); handler != nil {
		return handler
	}
	d.SetId(strings.Replace(resp.String(), "\"", "", 2))
	return resourceInfobloxHostRecordRead(d, meta)
}

func resourceInfobloxHostRecordRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("\n[infoblox-provider] %s", "----------------- host record read")
	host := Host{}
	wapiErr := WapiError{}
	resp, err := resty.R().
		SetResult(&host).
		SetError(&wapiErr).
		SetQueryParams(map[string]string{
			"_return_fields+": "ttl,use_ttl",
		}).
		Get("/" + d.Id())
	log.Printf("\n[infoblox-provider] Wapi Object: %+v", host)
	if handler := handleError(err, resp, wapiErr); handler != nil {
		return handler
	}
	splitFqdn := strings.Split(host.Name, ".")
	d.Set("fqdn", host.Name)
	d.Set("name", splitFqdn[0])
	d.Set("domain", strings.Join(splitFqdn[1:], "."))
	d.Set("ipv4", host.Ipv4addrs[0].Ipv4addr)
	d.Set("ttl", host.Ttl)
	d.Set("view", host.View)
	return nil
}

func resourceInfobloxHostRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("\n[infoblox-provider] %s", "----------------- host record update")
	ipv4 := d.Get("ipv4").(string)

	wapiErr := WapiError{}
	resp, err := resty.R().
		SetError(&wapiErr).
		SetBody(map[string]interface{}{
			"ipv4addrs": []map[string]interface{}{
				map[string]interface{}{
					"ipv4addr": ipv4,
				},
			},
		}).
		Put("/" + d.Id())
	if handler := handleError(err, resp, wapiErr); handler != nil {
		return handler
	}
	return resourceInfobloxHostRecordRead(d, meta)
}

func resourceInfobloxHostRecordDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("\n[infoblox-provider] %s", "----------------- host record delete")
	wapiErr := WapiError{}
	resp, err := resty.R().
		SetError(&wapiErr).
		Delete("/" + d.Id())
	if handler := handleError(err, resp, wapiErr); handler != nil {
		return handler
	}
	return nil
}
