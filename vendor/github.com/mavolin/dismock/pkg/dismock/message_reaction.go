package dismock

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/dismock/internal/mockutil"
	"github.com/mavolin/dismock/internal/sanitize"
)

// React mocks a React request.
func (m *Mocker) React(channelID, messageID discord.Snowflake, e api.Emoji) {
	m.MockAPI("React", http.MethodPut,
		"/channels/"+channelID.String()+"/messages/"+messageID.String()+"/reactions/"+url.PathEscape(e)+"/@me", nil)
}

// Unreact mocks a Unreact request.
func (m *Mocker) Unreact(channelID, messageID discord.Snowflake, e api.Emoji) {
	m.MockAPI("Unreact", http.MethodDelete,
		"/channels/"+channelID.String()+"/messages/"+messageID.String()+"/reactions/"+url.PathEscape(e)+"/@me", nil)
}

// Reactions mocks a Reactions request.
//
// This method will sanitize Users.ID.
func (m *Mocker) Reactions(channelID, messageID discord.Snowflake, limit uint, e api.Emoji, u []discord.User) {
	if u == nil {
		u = []discord.User{}
	}

	if len(u) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent users (%d vs. %d)", len(u), limit))
	}

	const hardLimit uint = 100

	var after discord.Snowflake

	for i := 0; i <= len(u)/int(hardLimit); i++ {
		var (
			from = uint(i) * hardLimit
			to   = uint(math.Min(float64(from+hardLimit), float64(len(u))))

			fetch = to - from // we expect this as the sent limit
		)

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// expect either limit or hardlimit, depending on which is smaller, instead.
			if fetch < hardLimit {
				fetch = uint(math.Min(float64(limit), float64(hardLimit)))
			}

			limit -= fetch
		} else { // this means there is no limit, hence we should expect hardlimit
			fetch = hardLimit
		}

		m.reactionsRange(channelID, messageID, 0, after, fmt.Sprintf("Reactions #%d", i+1), fetch, e, u[from:to])

		if to-from < hardLimit {
			break
		}

		after = u[to-1].ID
	}
}

// ReactionsBefore mocks a ReactionsBefore request.
//
// This method will sanitize Users.ID.
func (m *Mocker) ReactionsBefore(
	channelID, messageID discord.Snowflake, before discord.Snowflake, limit uint, e api.Emoji, u []discord.User,
) {
	if u == nil {
		u = []discord.User{}
	}

	if len(u) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent users (%d vs. %d)", len(u), limit))
	}

	const hardLimit = 100

	req := len(u)/hardLimit + 1

	from := uint(math.Min(float64(uint(req)*hardLimit), float64(len(u))))

	for i := req; i > 0; i-- {
		no := req - i + 1

		to := from
		from = uint(math.Max(float64(0), float64(int(to-hardLimit))))

		fetch := to - from // we expect this as the sent limit

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// expect either limit or hardlimit, depending on which is smaller, instead.
			if fetch < hardLimit {
				fetch = uint(math.Min(float64(limit), float64(hardLimit)))
			}

			limit -= fetch
		} else { // this means there is no limit, hence we should expect hardlimit
			fetch = hardLimit
		}

		m.reactionsRange(channelID, messageID, before, 0, fmt.Sprintf("ReactionsBefore #%d", no), fetch, e, u[from:to])

		if to-from < hardLimit {
			break
		}

		before = u[from].ID
	}
}

// ReactionsAfter mocks a ReactionsAfter request.
//
// This method will sanitize Users.ID.
func (m *Mocker) ReactionsAfter(
	channelID, messageID, after discord.Snowflake, limit uint, e api.Emoji, u []discord.User,
) {
	if u == nil {
		u = []discord.User{}
	}

	if len(u) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent users (%d vs. %d)", len(u), limit))
	}

	const hardLimit uint = 100

	for i := 0; i <= len(u)/int(hardLimit); i++ {
		var (
			from = uint(i) * hardLimit
			to   = uint(math.Min(float64(from+hardLimit), float64(len(u))))

			fetch = to - from // we expect this as the sent limit
		)

		// but if limit != unlimited
		if limit > 0 {
			// and the max data we can send (fetch) is smaller than what could be requested max, we
			// expect either limit or hardlimit, depending on which is smaller, instead.
			if fetch < hardLimit {
				fetch = uint(math.Min(float64(limit), float64(hardLimit)))
			}

			limit -= fetch
		} else { // this means there is no limit, hence we should expect hardlimit
			fetch = hardLimit
		}

		m.reactionsRange(channelID, messageID, 0, after, fmt.Sprintf("ReactionsAfter #%d", i+1), fetch, e, u[from:to])

		if to-from < hardLimit {
			break
		}

		after = u[to-1].ID
	}
}

// reactionsRange mocks a single request to the GET /reactions endpoint.
//
// This method will sanitize Users.ID.
func (m *Mocker) reactionsRange(
	channelID, messageID, before, after discord.Snowflake, name string, limit uint, e api.Emoji, u []discord.User,
) {
	for i, user := range u {
		u[i] = sanitize.User(user, 1)
	}

	m.MockAPI(name, http.MethodGet,
		"/channels/"+channelID.String()+"/messages/"+messageID.String()+"/reactions/"+url.PathEscape(e),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := url.Values{
				"limit": {strconv.FormatUint(uint64(limit), 10)},
			}

			if after != 0 {
				expect["after"] = []string{after.String()}
			}

			if before != 0 {
				expect["before"] = []string{before.String()}
			}

			mockutil.CheckQuery(t, r.URL.Query(), expect)
			mockutil.WriteJSON(t, w, u)
		})
}

// DeleteUserReaction mocks a DeleteUserReaction request.
func (m *Mocker) DeleteUserReaction(channelID, messageID, userID discord.Snowflake, e api.Emoji) {
	user := "@me"
	if userID > 0 {
		user = userID.String()
	}

	m.MockAPI("DeleteUserReaction", http.MethodDelete,
		"/channels/"+channelID.String()+"/messages/"+messageID.String()+"/reactions/"+url.PathEscape(e)+"/"+user, nil)
}

// DeleteReactions mocks a DeleteReactions request.
func (m *Mocker) DeleteReactions(channelID, messageID discord.Snowflake, e api.Emoji) {
	m.MockAPI("DeleteReactions", http.MethodDelete,
		"/channels/"+channelID.String()+"/messages/"+messageID.String()+"/reactions/"+url.PathEscape(e), nil)
}

// DeleteAllReactions mocks a DeleteAllReactions request.
func (m *Mocker) DeleteAllReactions(channelID, messageID discord.Snowflake) {
	m.MockAPI("DeleteAllReactions", http.MethodDelete,
		"/channels/"+channelID.String()+"/messages/"+messageID.String()+"/reactions", nil)
}
