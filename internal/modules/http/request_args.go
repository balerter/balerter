package http

import (
	"fmt"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type requestArgs struct {
	Method             string            `json:"method"`
	URI                string            `json:"uri"`
	Body               []byte            `json:"body"`
	Headers            map[string]string `json:"headers"`
	Timeout            time.Duration     `json:"timeout"`
	InsecureSkipVerify bool              `json:"insecure_skip_verify"`
}

func newRequestArgs() *requestArgs {
	return &requestArgs{
		Headers: make(map[string]string),
	}
}

func (r *requestArgs) parseFromTable(tbl *lua.LTable) error {
	methodValue := tbl.RawGetString("method")
	if methodValue.Type() != lua.LTNil {
		r.Method = parseMethod(methodValue.String())
		if r.Method == "" {
			return fmt.Errorf("bad http method %s", methodValue.String())
		}
	}

	uriValue := tbl.RawGetString("uri")
	if uriValue.Type() != lua.LTNil {
		if uriValue.Type() != lua.LTString {
			return fmt.Errorf("uri must be a string")
		}
		r.URI = uriValue.String()
	}

	bodyValue := tbl.RawGetString("body")
	if bodyValue.Type() != lua.LTNil {
		r.Body = []byte(bodyValue.String())
	}

	skipVerifyValue := tbl.RawGetString("insecureSkipVerify")
	if skipVerifyValue.Type() != lua.LTNil {
		if skipVerifyValue.Type() != lua.LTBool {
			return fmt.Errorf("insecureSkipVerify must be a bool")
		}
		r.InsecureSkipVerify = skipVerifyValue.String() == "true"
	}

	timeoutValue := tbl.RawGetString("timeout")
	if timeoutValue.Type() != lua.LTNil {
		timeoutDuration, errParseTimeout := time.ParseDuration(timeoutValue.String())
		if errParseTimeout != nil {
			return fmt.Errorf("timeout must be a time.Duration")
		}
		r.Timeout = timeoutDuration
	}

	headersValue := tbl.RawGetString("headers")
	if headersValue.Type() != lua.LTNil {
		if headersValue.Type() != lua.LTTable {
			return fmt.Errorf("headers must be a table")
		}

		headersValue.(*lua.LTable).ForEach(func(value lua.LValue, value2 lua.LValue) {
			r.Headers[value.String()] = value2.String()
		})
	}

	return nil
}

func parseMethod(method string) string {
	for _, m := range []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE"} {
		if strings.ToUpper(method) == m {
			return m
		}
	}
	return ""
}

// parseRequestArgs parses arguments for methods http.get, http.post etc
// 1 - uri
// 2 - body
// 3 - headers
func (h *HTTP) parseRequestArgs(luaState *lua.LState) (*requestArgs, error) {
	args := newRequestArgs()

	argURI := luaState.Get(1)
	if argURI.Type() != lua.LTString {
		return nil, fmt.Errorf("uri argument must be a string")
	}
	args.URI = argURI.(lua.LString).String()

	argBody := luaState.Get(2) //nolint:gomnd // param position
	switch argBody.Type() {
	case lua.LTString:
		args.Body = []byte(argBody.(lua.LString).String())
	case lua.LTNil:
	default:
		return nil, fmt.Errorf("body argument must be a string or nil")
	}

	argHeaders := luaState.Get(3) //nolint:gomnd // param position
	switch argHeaders.Type() {
	case lua.LTTable:
		argHeaders.(*lua.LTable).ForEach(func(value lua.LValue, value2 lua.LValue) {
			args.Headers[value.String()] = value2.String()
		})
	case lua.LTNil:
	default:
		return nil, fmt.Errorf("headers argument must be a table or nil")
	}

	return args, nil
}
