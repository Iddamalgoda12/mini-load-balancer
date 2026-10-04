package config

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

// Config contains the process-level settings for the load balancer.
type Config struct {
	ListenAddr        string
	Backends          []*url.URL
	ReadHeaderTimeout time.Duration
}

// Parse reads flags and environment defaults, then validates the result.
func Parse(args []string, lookupEnv func(string) (string, bool)) (Config, error) {
	if lookupEnv == nil {
		lookupEnv = os.LookupEnv
	}

	listenDefault := envOr(lookupEnv, "LB_LISTEN_ADDR", ":8080")
	backendsDefault := envOr(lookupEnv, "LB_BACKENDS", "")

	flags := flag.NewFlagSet("loadbalancer", flag.ContinueOnError)
	listenAddr := flags.String("listen", listenDefault, "address for the HTTP listener")
	backends := flags.String("backends", backendsDefault, "comma-separated backend URLs")
	if err := flags.Parse(args); err != nil {
		return Config{}, err
	}
	if flags.NArg() != 0 {
		return Config{}, fmt.Errorf("unexpected positional arguments: %s", strings.Join(flags.Args(), " "))
	}
	if strings.TrimSpace(*listenAddr) == "" {
		return Config{}, errors.New("listen address cannot be empty")
	}

	parsedBackends, err := parseBackends(*backends)
	if err != nil {
		return Config{}, err
	}
	return Config{
		ListenAddr:        *listenAddr,
		Backends:          parsedBackends,
		ReadHeaderTimeout: 5 * time.Second,
	}, nil
}

func envOr(lookupEnv func(string) (string, bool), key, fallback string) string {
	if value, ok := lookupEnv(key); ok {
		return value
	}
	return fallback
}

func parseBackends(raw string) ([]*url.URL, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, errors.New("at least one backend is required; set -backends or LB_BACKENDS")
	}

	parts := strings.Split(raw, ",")
	backends := make([]*url.URL, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		parsed, err := url.Parse(value)
		if err != nil {
			return nil, fmt.Errorf("invalid backend URL %q: %w", value, err)
		}
		if parsed.Scheme != "http" && parsed.Scheme != "https" {
			return nil, fmt.Errorf("backend %q must use http or https", value)
		}
		if parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
			return nil, fmt.Errorf("backend %q must have a host and cannot contain credentials, a query, or a fragment", value)
		}
		backends = append(backends, parsed)
	}
	return backends, nil
}
