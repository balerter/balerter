package mock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const unsupportedMockRequest = "unsupported mock request"

func restClient() *http.Client {
	return &http.Client{Transport: RoundTripperFunc(discordAPIResponse)}
}

func discordAPIResponse(r *http.Request) (*http.Response, error) {
	switch {
	case strings.Contains(r.URL.Path, "users"):
		return usersResponse(r), nil
	case strings.Contains(r.URL.Path, "members"):
		return membersResponse(r), nil
	case strings.Contains(r.URL.Path, "roles"):
		return rolesResponse(r), nil
	case strings.Contains(r.URL.Path, "channels"):
		return channelsResponse(r), nil
	case strings.Contains(r.URL.Path, "guilds"):
		return guildsResponse(r), nil
	}

	return nil, fmt.Errorf(unsupportedMockRequest)
}

func usersResponse(r *http.Request) *http.Response {
	pathTokens := strings.Split(r.URL.Path, "/")
	userID := pathTokens[len(pathTokens)-1]

	respBody, err := json.Marshal(mockUser(userID))
	if err != nil {
		return newResponse(http.StatusInternalServerError, []byte(err.Error()))
	}

	return newResponse(http.StatusOK, respBody)
}

func membersResponse(r *http.Request) *http.Response {
	pathTokens := strings.Split(r.URL.Path, "/")
	userID := pathTokens[len(pathTokens)-1]
	guildID := pathTokens[len(pathTokens)-2]

	var (
		respBody []byte
		err      error
	)

	if userID == "members" {
		if guildID == TestGuild {
			respBody, err = json.Marshal(mockMembers())
		}

		if guildID == TestGuildLarge {
			queryParamters := r.URL.Query()

			if len(queryParamters["after"]) == 0 {
				respBody, err = json.Marshal(mockLargeMembers())
			} else {
				respBody, err = json.Marshal(mockMembers())
			}
		}
	} else {
		respBody, err = json.Marshal(mockMember(userID))
	}

	if err != nil {
		return newResponse(http.StatusInternalServerError, []byte(err.Error()))
	}

	return newResponse(http.StatusOK, respBody)
}

func rolesResponse(r *http.Request) *http.Response {
	switch r.Method {
	case http.MethodGet:
		respBody, err := json.Marshal(mockRoles())
		if err != nil {
			return newResponse(http.StatusInternalServerError, []byte(err.Error()))
		}

		return newResponse(http.StatusOK, respBody)
	case http.MethodPost:
		respBody, err := json.Marshal(mockRole(TestRole))
		if err != nil {
			return newResponse(http.StatusInternalServerError, []byte(err.Error()))
		}

		return newResponse(http.StatusOK, respBody)
	case http.MethodPatch:
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return newResponse(http.StatusInternalServerError, []byte(err.Error()))
		}

		err = r.Body.Close()
		if err != nil {
			return newResponse(http.StatusInternalServerError, []byte(err.Error()))
		}

		return newResponse(http.StatusOK, reqBody)
	case http.MethodDelete:
		return newResponse(http.StatusOK, nil)
	}

	return newResponse(http.StatusMethodNotAllowed, []byte{})
}

func channelsResponse(r *http.Request) *http.Response {
	var (
		respBody []byte
		err      error
	)

	if strings.Contains(r.URL.Path, "guilds") {
		respBody, err = json.Marshal(mockChannels())
	} else {
		respBody, err = json.Marshal(mockChannel(TestChannel))
	}

	if err != nil {
		return newResponse(http.StatusInternalServerError, []byte(err.Error()))
	}

	return newResponse(http.StatusOK, respBody)
}

func guildsResponse(r *http.Request) *http.Response {
	pathTokens := strings.Split(r.URL.Path, "/")
	guildID := pathTokens[len(pathTokens)-1]

	respBody, err := json.Marshal(mockGuild(guildID))
	if err != nil {
		return newResponse(http.StatusInternalServerError, []byte(err.Error()))
	}

	return newResponse(http.StatusOK, respBody)
}

func newResponse(status int, respBody []byte) *http.Response {
	return &http.Response{
		Status:     http.StatusText(status),
		StatusCode: status,
		Header:     make(http.Header),
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}
}
