package main

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

type Named struct {
	Name string
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
