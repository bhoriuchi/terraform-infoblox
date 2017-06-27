package main

type Object struct {
	Ref string `json:"_ref"`
}
type Host struct {
	Object
	Name, View string
	Ttl        int
	Use_Ttl    bool
	Ipv4addrs  []Ipv4
}
type Ipv4 struct {
	Object
	Host, Ipv4addr     string
	Configure_for_dhcp bool
}
type A struct {
	Object
	Name, View, Ipv4addr, Dns_name 	string
	Ttl 				int
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