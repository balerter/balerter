package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type vaultResponse struct {
	Data struct {
		Data map[string]string `json:"data"`
	} `json:"data"`
}

func DecodeSecret(key []byte) ([]byte, error) {
	pair := strings.Split(string(key), "@")
	if len(pair) != 2 {
		return nil, fmt.Errorf("wrong format vault secret")
	}

	var result string

	vaultURL := os.Getenv("BALERTER_VAULT_URL")
	vaultToken := os.Getenv("BALERTER_VAULT_TOKEN")
	vaultNS := os.Getenv("BALERTER_VAULT_NAMESPACE")

	vaultURL += "/v1/" + pair[0]

	req, errRequest := http.NewRequest(http.MethodGet, vaultURL, http.NoBody)
	if errRequest != nil {
		return nil, fmt.Errorf("error create a request, %w", errRequest)
	}
	if vaultToken != "" {
		req.Header.Add("X-Vault-Token", vaultToken)
	}
	if vaultNS != "" {
		req.Header.Add("X-Vault-Namespace", vaultNS)
	}

	resp, errResponse := http.DefaultClient.Do(req)
	if errResponse != nil {
		return nil, fmt.Errorf("error send a request to the vault, %w", errResponse)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("vault secret '%s' not found", key)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status code = %d", resp.StatusCode)
	}

	var r vaultResponse

	errUnmarshal := json.NewDecoder(resp.Body).Decode(&r)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}

	result, ok := r.Data.Data[pair[1]]
	if !ok {
		return nil, fmt.Errorf("vault secret '%s' not found", key)
	}

	return []byte(result), nil
}
