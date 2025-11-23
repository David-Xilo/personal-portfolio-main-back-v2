package middleware

import (
	"net"
	"net/http"
)

// getRemoteIP returns the client IP based solely on the TCP remote address.
// This ignores any forwarded headers to avoid spoofing when running directly on the internet.
func getRemoteIP(r *http.Request) string {
	if r == nil {
		return ""
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}

	return r.RemoteAddr
}
