package servertest

import (
	"net"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/cenkalti/backoff/v3"
)

func Setup(t testing.TB, opts ...Opt) string {
	t.Helper()

	c := config{}

	for _, opt := range opts {
		opt(&c)
	}

	endpoint := c.Endpoint()
	if endpoint != "" {
		waitReady(t, endpoint, c.Backoff())
		return endpoint
	}

	server := httptest.NewServer(c.makeHandler())
	addr := server.Listener.Addr().String()

	t.Cleanup(server.Close)

	return "http://" + addr
}

func waitReady(t testing.TB, endpoint string, bo backoff.BackOff) {
	addr := parseEndpoint(t, endpoint)

	err := backoff.Retry(func() error {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return err
		}
		conn.Close()

		return nil
	}, bo)
	if err != nil {
		t.Fatalf("failed to connect to server: %v", err)
	}
}

func parseEndpoint(t testing.TB, endpoint string) string {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		t.Fatal(err)
	}

	port := parsed.Port()

	switch {
	case port != "":
		return parsed.Host
	case parsed.Scheme == "http":
		return parsed.Host + ":80"
	case parsed.Scheme == "https":
		return parsed.Host + ":443"
	default:
		t.Fatalf("unsupported endpoint %q", endpoint)
		return ""
	}
}
