package api

import "net/http"

// handleHeaders helps with parsing incomming header
// for 'x-passphrase' header.
func handleHeaders(rh http.Header) http.Header {
	header := http.Header{}
	pph := rh.Get("x-passphrase")
	if pph != "" {
		header.Add("x-passphrase", pph)
	}

	return header
}
