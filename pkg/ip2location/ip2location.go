package ip2location

import (
	"os"

	"github.com/ip2location/ip2location-go/v9"
)

var db *ip2location.DB

func Open() {
	binPath := os.Getenv("IP2LOCATION_BIN")
	if binPath == "" {
		return
	}

	var err error
	db, err = ip2location.OpenDB(binPath)
	if err != nil {
		panic(err)
	}
}

func GetCountry(ipAddr string) string {
	if db == nil {
		return ""
	}

	rec, err := db.Get_country_short(ipAddr)
	if err != nil {
		return ""
	}

	return rec.Country_short
}
