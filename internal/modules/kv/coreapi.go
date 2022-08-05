package kv

import (
	"fmt"
	"net/http"
)

func (kv *KV) CoreApiHandler(req []string, body []byte) (any, int, error) {
	if len(req) == 0 {
		return nil, http.StatusBadRequest, fmt.Errorf("empty path")
	}

	method := req[0]

	switch method {
	case "all":
		data, err := kv.engine.All()
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return data, 0, nil
	case "put":
		if len(req) != 2 {
			return nil, http.StatusBadRequest, fmt.Errorf("wrong number of params")
		}
		err := kv.engine.Put(req[1], string(body))
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return nil, 0, nil
	case "upsert":
		if len(req) != 2 {
			return nil, http.StatusBadRequest, fmt.Errorf("wrong number of params")
		}
		err := kv.engine.Upsert(req[1], string(body))
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return nil, 0, nil
	case "get":
		if len(req) != 2 {
			return nil, http.StatusBadRequest, fmt.Errorf("wrong number of params")
		}
		v, err := kv.engine.Get(req[1])
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return v, 0, nil
	case "delete":
		if len(req) != 2 {
			return nil, http.StatusBadRequest, fmt.Errorf("wrong number of params")
		}
		err := kv.engine.Delete(req[1])
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return nil, 0, nil
	}

	return nil, http.StatusBadRequest, fmt.Errorf("unknown method %q", method)
}
