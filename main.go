package main

import (
	"os"
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"github.com/fatih/color"
	"net"
	"regexp"
)

var httpClient = &http.Client{Timeout: 5 * time.Second}

type IPInfo struct {
	IP           string `json:"ip"`
	Location     string `json:"location"`
	City         string `json:"city"`
	Region       string `json:"region"`
	Country      string `json:"country"`
	Organization string `json:"org"`
	Hostname     string `json:"hostname"`
}

func main() {

	currentInfo := new(IPInfo)
	if len(os.Args) < 2 {
		fmt.Print("Command usage: ipinfo {IP ADDRESS OR DOMAIN}")
	} else {

		if checkIP(os.Args[1]) {
			currentInfo.IP = os.Args[1]
		} else {
			currentInfo.IP = getHostnameFromIP(os.Args[1])
		}

		getIPInfo(currentInfo.IP, currentInfo)
		color.Green("IP Address " + currentInfo.IP + " was scanned successfully. \n")
		color.Red("		Hostname: " + currentInfo.Hostname)
		color.Red("		ARIN Organization: " + currentInfo.Organization)
		color.Green("Location Data:")
		color.White("		Country: " + currentInfo.Country)
		color.White("		Region: " + currentInfo.Region)
		color.White("		City: " + currentInfo.City)
	}

}

func checkIP(address string) bool {
	ipreg, _ := regexp.Compile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)
	if ipreg.MatchString(os.Args[1]) {
		return true
	}
	return false
}

func getHostnameFromIP(hostname string) string {
	addr, err := net.LookupIP(hostname)
	if err != nil {
		color.Red("Unknown Hostname")
		os.Exit(1)
	}
	if len(addr) > 1 {
		color.Red("More than one IP address discovered. Using the first.")
	}
	return addr[0].String()

}

func getIPInfo(address string, target interface{}) error {
	r, err := httpClient.Get("https://ipinfo.io/" + address)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
