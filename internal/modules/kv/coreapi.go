package kv

import (
	"fmt"
	"net/http"
)

func (kv *KV) CoreApiHandler(method string, parts []string, params map[string]string, body []byte) (any, int, error) {
	if method == "all" {
		data, err := kv.engine.All()
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return data, 0, nil
	}

	if len(parts) != 1 {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid request")
	}
	keyName := parts[0]

	switch method {
	case "put":
		err := kv.engine.Put(keyName, string(body))
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return nil, 0, nil
	case "upsert":
		err := kv.engine.Upsert(keyName, string(body))
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return nil, 0, nil
	case "get":
		v, err := kv.engine.Get(keyName)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return v, 0, nil
	case "delete":
		err := kv.engine.Delete(keyName)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return nil, 0, nil
	}

	return nil, http.StatusBadRequest, fmt.Errorf("unknown method %q", method)
}
