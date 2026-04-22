package ip2location

import (
	"net"
	"strings"

	"github.com/ip2location/ip2location-go/v9"
)

var dbV4 *ip2location.DB
var dbV6 *ip2location.DB

func Open() {
	var err error

	dbV4, err = ip2location.OpenDB("./geoip_ipv4.bin")
	if err != nil {
		panic("erro abrindo db: " + err.Error())
	}

	dbV6, err = ip2location.OpenDB("./geoip_ipv6.bin")
	if err != nil {
		panic("erro abrindo db: " + err.Error())
	}
}

func GetCountry(ipAddr string) string {
	ipAddr = strings.TrimSpace(ipAddr)

	if host, _, err := net.SplitHostPort(ipAddr); err == nil {
		ipAddr = host
	}

	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return ""
	}

	if ip.IsLoopback() || ip.IsPrivate() {
		return ""
	}

	var db *ip2location.DB
	if ip.To4() != nil {
		db = dbV4
	} else {
		db = dbV6
	}

	if db == nil {
		return ""
	}

	rec, err := db.Get_country_short(ipAddr)
	if err != nil {
		return ""
	}

	return rec.Country_short
}
