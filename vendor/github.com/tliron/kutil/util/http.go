package util

import (
	"net/http"
)

//
// ForceHTTPSRoundTripper
//

type ForceHTTPSRoundTripper struct {
	roundTripper http.RoundTripper
}

func NewForceHTTPSRoundTripper(roundTripper http.RoundTripper) *ForceHTTPSRoundTripper {
	return &ForceHTTPSRoundTripper{roundTripper}
}

// http.RoundTripper interface
func (self *ForceHTTPSRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	if request.URL.Scheme != "https" {
		// Rewrite URL
		url := *request.URL
		url.Scheme = "https"
		request.URL = &url
	}

	return self.roundTripper.RoundTrip(request)
}
