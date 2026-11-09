package railway

import (
	"net/http"
	"strings"
)

// projectTokenPrefix marks Railway project tokens, which authenticate with the
// Project-Access-Token header. Every other token (workspace, account, legacy UUID)
// uses Authorization: Bearer.
const projectTokenPrefix = "rlwy_proj_"

// setAuthHeader sets the Railway auth header for token on h. Never log h afterwards.
func setAuthHeader(h http.Header, token string) {
	if strings.HasPrefix(token, projectTokenPrefix) {
		h.Set("Project-Access-Token", token)
		return
	}

	h.Set("Authorization", "Bearer "+token)
}

type authedTransport struct {
	token   string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	setAuthHeader(req.Header, t.token)
	req.Header.Set("Content-Type", "application/json")

	return t.wrapped.RoundTrip(req)
}
