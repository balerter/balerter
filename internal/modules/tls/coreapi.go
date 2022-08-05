package tls

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type coreapiRequest struct {
	Hostname string `json:"hostname"`
}

type coreapiResponse struct {
	Issuer         string   `json:"issuer"`
	Expiry         int64    `json:"expiry"`
	DNSNames       []string `json:"dns_names"`
	EmailAddresses []string `json:"email_addresses"`
}

func (a *TLS) CoreApiHandler(req []string, body []byte) (any, int, error) {
	if len(req) != 1 {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid request, expected 1 argument, got %d", len(req))
	}
	if req[0] != "get" {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid request, unknown method %q", req[0])
	}

	var r coreapiRequest
	errUnmarshalBody := json.Unmarshal(body, &r)
	if errUnmarshalBody != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("error unmarshal body: %s", errUnmarshalBody)
	}

	conn, errDial := a.dialFunc("tcp", r.Hostname+":443", &tls.Config{InsecureSkipVerify: true})
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
