package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mikerybka/brass"
	"github.com/mikerybka/util"
	"golang.org/x/net/publicsuffix"
)

func main() {
	b, _ := os.ReadFile("/etc/brassd/email")
	email := strings.TrimSpace(string(b))

	err := util.ServeHTTPS(util.NewLiveObject("/var/brassd/data.json", &brass.Server{}), email, "/var/brassd/certs", func(host string) bool {
		domain, err := publicsuffix.EffectiveTLDPlusOne(host)
		if err != nil {
			return false
		}
		minusSLD := strings.TrimSuffix(host, domain)
		if minusSLD == "" {
			return true
		}
		allowedSubdomains := []string{
			"www",
			"api",
			"admin",
			"docs",
		}
		for _, sd := range allowedSubdomains {
			if minusSLD == sd+"." {
				return true
			}
		}
		return false
	})
	fmt.Println(err)
}
