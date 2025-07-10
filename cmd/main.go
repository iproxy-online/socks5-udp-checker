package main

import (
	"log"
	"net"

	"github.com/beevik/ntp"
	"github.com/txthinking/socks5"
)

type config struct {
	socks5Address  string
	socks5Username string
	socks5Password string
	ntpAddress     string
}

const timeout = 1000 // NTP uses UDP, so this value is not applicable for TCP.

func do(cfg *config) {
	socks5Cl, err := socks5.NewClient(
		cfg.socks5Address,
		cfg.socks5Username,
		cfg.socks5Password,
		timeout,
		timeout)
	if err != nil {
		log.Panicf("failed to create SOCKS5 client: %s", err.Error())
	}

	resp, err := ntp.QueryWithOptions(cfg.ntpAddress, ntp.QueryOptions{
		Dialer: func(_, remoteAddress string) (net.Conn, error) {
			return socks5Cl.Dial("udp", remoteAddress)
		},
	})
	if err != nil {
		log.Panicf("NTP request failed: %s", err)
	}

	log.Println("=== NTP Response ===")
	log.Printf("Server Information:")
	log.Printf("  Stratum:         %d", resp.Stratum)
	log.Printf("  Version:         %d", resp.Version)
	log.Printf("  Leap Indicator:  %v", resp.Leap)
	log.Printf("  Poll Interval:   %v", resp.Poll)
	log.Printf("  Precision:       %v", resp.Precision)
	log.Printf("  Reference ID:    %v (%s)", resp.ReferenceID, resp.ReferenceString())

	log.Printf("\nTiming Information:")
	log.Printf("  Reference Time:  %v", resp.ReferenceTime.Format("2006-01-02 15:04:05.000 MST"))
	log.Printf("  Server Time:     %v", resp.Time.Format("2006-01-02 15:04:05.000 MST"))

	log.Printf("\nNetwork Metrics:")
	log.Printf("  Root Delay:      %v", resp.RootDelay)
	log.Printf("  Root Dispersion: %v", resp.RootDispersion)
	log.Printf("  Root Distance:   %v", resp.RootDistance)
	log.Printf("  Round Trip Time: %v", resp.RTT)
	log.Printf("  Clock Offset:    %v", resp.ClockOffset)
	log.Printf("  Min Error:       %v", resp.MinError)

	if resp.IsKissOfDeath() {
		log.Printf("\n⚠️  Kiss of Death Response:")
		log.Printf("  Kiss Code:       %s", resp.KissCode)
	}
	log.Println("===================")
}

func main() {
	cfg := &config{
		socks5Address:  "localhost:1080",
		socks5Username: "",
		socks5Password: "",
		ntpAddress:     "time.google.com:123",
	}

	do(cfg)
}
