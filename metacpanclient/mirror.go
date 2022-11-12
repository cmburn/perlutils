package metacpanclient

import (
	"github.com/relvacode/iso8601"
)

type Mirror struct {
	AKAName  string `json:"aka_name"`
	AOrCName string `json:"A_or_CNAME"`
	CCode    string `json:"ccode"`
	City     string `json:"city"`
	Contact  []struct {
		ContactSite string `json:"contact_site"`
		ContactUser string `json:"contact_user"`
	} `json:"contact"`
	Country      string       `json:"country"`
	DNSRR        string       `json:"dnsrr"`
	FTP          string       `json:"ftp"`
	Freq         string       `json:"freq"`
	HTTP         string       `json:"http"`
	InceptDate   iso8601.Time `json:"inceptdate"`
	Location     [2]float64   `json:"location"`
	Name         string       `json:"name"`
	Note         string       `json:"note"`
	Organization string       `json:"org"`
	RSync        string       `json:"rsync"`
	Region       string       `json:"region"`
	ReitreDate   iso8601.Time `json:"reitredate"`
	Src          string       `json:"src"`
	Tz           string       `json:"tz"`
}
