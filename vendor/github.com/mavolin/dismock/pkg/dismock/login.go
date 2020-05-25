package dismock

import (
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/api"

	"github.com/mavolin/dismock/internal/mockutil"
)

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login mocks a Login request.
func (m *Mocker) Login(email, password string, response api.LoginResponse) {
	m.MockAPI("Login", http.MethodPost, "/auth/login",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := loginPayload{
				Email:    email,
				Password: password,
			}

			mockutil.CheckJSON(t, r.Body, new(loginPayload), &expect)
			mockutil.WriteJSON(t, w, response)
		})
}

type totpPayload struct {
	Code   string `json:"code"`
	Ticket string `json:"ticket"`
}

// TOTP mocks a TOTP request.
func (m *Mocker) TOTP(code, ticket string, response api.LoginResponse) {
	m.MockAPI("TOTP", http.MethodPost, "/auth/mfa/totp",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := totpPayload{
				Code:   code,
				Ticket: ticket,
			}

			mockutil.CheckJSON(t, r.Body, new(totpPayload), &expect)
			mockutil.WriteJSON(t, w, response)
		})
}
