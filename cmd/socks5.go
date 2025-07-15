package main

import (
	"fmt"
	"net"
	"net/url"
	"regexp"

	"github.com/beevik/ntp"
	"github.com/txthinking/socks5"
)

type socks5Config struct {
	address  string
	username string
	password string
}

type config struct {
	socks5Config
	ntpAddress string
}

func performNTPTest(cfg *config) (*ntp.Response, error) {
	const timeout = 1000

	socks5Cl, err := socks5.NewClient(
		cfg.address,
		cfg.username,
		cfg.password,
		timeout,
		timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create SOCKS5 client: %w", err)
	}

	resp, err := ntp.QueryWithOptions(cfg.ntpAddress, ntp.QueryOptions{
		Dialer: func(_, remoteAddress string) (net.Conn, error) {
			return socks5Cl.Dial("udp", remoteAddress)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("NTP request failed: %w", err)
	}

	return resp, nil
}

func parseSocks5String(socks5Str string) (socks5Config, error) {
	var cfg socks5Config
	if socks5Str == "" {
		return cfg, fmt.Errorf("empty input string")
	}

	// Try URL parsing first for socks5:// format
	if parsed, err := url.Parse(socks5Str); err == nil && parsed.Scheme == "socks5" && parsed.Host != "" {
		cfg.address = parsed.Host
		if parsed.User != nil {
			cfg.username = parsed.User.Username()
			cfg.password, _ = parsed.User.Password()
		}
		return cfg, nil
	}

	// Regex patterns with named capture groups
	patterns := []*regexp.Regexp{
		// socks5://username:password@host:port
		regexp.MustCompile(`^socks5://(?P<username>[^:]+):(?P<password>[^@]+)@(?P<host>[^:]+):(?P<port>\d+)$`),
		// socks5://host:port:username:password
		regexp.MustCompile(`^socks5://(?P<host>[^:]+):(?P<port>\d+):(?P<username>[^:]+):(?P<password>.+)$`),
		// socks5:host:port:username:password
		regexp.MustCompile(`^socks5:(?P<host>[^:]+):(?P<port>\d+):(?P<username>[^:]+):(?P<password>.+)$`),
		// socks5://host:port
		regexp.MustCompile(`^socks5://(?P<host>[^:]+):(?P<port>\d+)$`),
		// socks5:host:port
		regexp.MustCompile(`^socks5:(?P<host>[^:]+):(?P<port>\d+)$`),
	}

	// Try each pattern
	for _, pattern := range patterns {
		matches := pattern.FindStringSubmatch(socks5Str)
		if matches != nil {
			names := pattern.SubexpNames()
			result := make(map[string]string)

			// Create a map of named captures
			for i, name := range names {
				if i != 0 && name != "" {
					result[name] = matches[i]
				}
			}

			// Build the config from named captures
			cfg.address = fmt.Sprintf("%s:%s", result["host"], result["port"])
			cfg.username = result["username"]
			cfg.password = result["password"]

			return cfg, nil
		}
	}

	return cfg, fmt.Errorf("invalid format")
}
