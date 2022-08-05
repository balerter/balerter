package tls

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

type coreapiResponse struct {
	Issuer         string   `json:"issuer"`
	Expiry         int64    `json:"expiry"`
	DNSNames       []string `json:"dns_names"`
	EmailAddresses []string `json:"email_addresses"`
}

func (a *TLS) CoreApiHandler(method string, parts []string, params map[string]string, body []byte) (any, int, error) {
	if method != "get" {
		return nil, http.StatusBadRequest, fmt.Errorf("unknown method %q", method)
	}

	conn, errDial := a.dialFunc("tcp", string(body)+":443", &tls.Config{InsecureSkipVerify: true})
	if errDial != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("error dial to host: %s", errDial)
	}
	defer conn.Close()

	var resp []coreapiResponse

	for _, cert := range conn.ConnectionState().PeerCertificates {
		resp = append(resp, coreapiResponse{
			Issuer:         cert.Issuer.String(),
			Expiry:         cert.NotAfter.UTC().Unix(),
			DNSNames:       cert.DNSNames,
			EmailAddresses: cert.EmailAddresses,
		})
	}

	return &resp, 0, nil
}
