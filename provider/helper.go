package main

import (
	"github.com/go-resty/resty"
	"log"
	"fmt"
	"regexp"
	"strconv"
	"sort"
	"strings"
)

func containsInt(list []int, value int) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
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

func getAvailableHostnames(prefix string, pad string, domain string, recordType string, start int, limit int) ([]string, error) {
	wapiErr 	:= WapiError{}
	used 		:= make([]int, 0)
	response        := []Named{}
	regex 		:= regexp.MustCompile("^(" + prefix + ")(\\d{" + pad + "}).(" + domain + ")$")
	max, err	:= strconv.Atoi(strings.Replace(fmt.Sprintf("%0" + pad + "d", 0), "0", "9", -1))
	list		:= make([]string, 0)

	if err != nil {
		return list, fmt.Errorf("Cannot compute max index based on pad")
	}
	if start >= max {
		return list, fmt.Errorf("Start index %d is too high for max index %d", start, max)
	}

	resp, err := resty.R().
		SetResult(&response).
		SetError(&wapiErr).
		SetQueryString("name~=^" + prefix + "\\d{" + pad + "}." + domain + "$").
		Get("/record:" + recordType)

	if handler := handleError(err, resp, wapiErr); handler != nil {
		return list, handler
	}

	for _, v := range response {
		matches := regex.FindAllStringSubmatch(v.Name, -1)

		if len(matches) > 0 && len(matches[0]) > 3 {
			idx, err := strconv.Atoi(matches[0][2])
			if err != nil {
				return list, err
			}
			used = append(used, idx)
		}
	}

	sort.Ints(used)

	for i := start; i < max; i++ {
		if len(list) >= limit {
			return list, nil
		}
		if !containsInt(used, i) {
			list = append(list, fmt.Sprintf("%s%0" + pad + "d", prefix, i))
		}
	}

	return list, nil
}
