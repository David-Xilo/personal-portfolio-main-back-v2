package middleware

import (
	"net"
	"net/http"
	"strings"
)

func getRemoteIP(r *http.Request) string {
	if r == nil {
		return ""
	}

	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		for _, part := range strings.Split(forwardedFor, ",") {
			if ip := strings.TrimSpace(part); ip != "" {
				return ip
			}
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}

	return r.RemoteAddr
}
