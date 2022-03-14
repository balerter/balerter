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

// Member mocks a Member request.
//
// The User.ID field of the passed member must be set.
func (m *Mocker) Member(guildID discord.GuildID, member discord.Member) {
	m.MockAPI("Member", http.MethodGet, "/guilds/"+guildID.String()+"/members/"+member.User.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, member)
		})
}

// Members mocks a Members request.
//
// This method will sanitize Member.User.ID.
func (m *Mocker) Members(guildID discord.GuildID, limit uint, members []discord.Member) {
	if members == nil {
		members = []discord.Member{}
	}

	if len(members) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent Members (%d vs. %d)", len(members), limit))
	}

	const hardLimit uint = 1000

	var after discord.UserID

	for i := 0; i <= len(members)/int(hardLimit); i++ {
		var (
			from = uint(i) * hardLimit
			to   = uint(math.Min(float64(from+hardLimit), float64(len(members))))

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

		m.membersAfter(guildID, after, fmt.Sprintf("Members #%d", i+1), fetch, members[from:to])

		if to-from < hardLimit {
			break
		}

		after = members[to-1].User.ID
	}
}

// MembersAfter mocks a MembersAfter request.
//
// This method will sanitize Member.User.ID.
func (m *Mocker) MembersAfter(guildID discord.GuildID, after discord.UserID, limit uint, members []discord.Member) {
	if members == nil {
		members = []discord.Member{}
	}

	if len(members) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent Members (%d vs. %d)", len(members), limit))
	}

	const hardLimit uint = 1000

	for i := 0; i <= len(members)/int(hardLimit); i++ {
		var (
			from = uint(i) * hardLimit
			to   = uint(math.Min(float64(from+hardLimit), float64(len(members))))

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

		m.membersAfter(guildID, after, fmt.Sprintf("MembersAfter #%d", i+1), fetch, members[from:to])

		if to-from < hardLimit {
			break
		}

		after = members[to-1].User.ID
	}
}

// membersAfter mocks a single request to the GET /Members endpoint.
//
// This method will sanitize Member.User.ID.
func (m *Mocker) membersAfter(
	guildID discord.GuildID, after discord.UserID, name string, limit uint, g []discord.Member,
) {
	for i, Member := range g {
		g[i] = sanitize.Member(Member, 1)
	}

	m.MockAPI(name, http.MethodGet, "/guilds/"+guildID.String()+"/members",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := url.Values{
				"limit": {strconv.FormatUint(uint64(limit), 10)},
			}

			if after != 0 {
				expect["after"] = []string{after.String()}
			}

			mockutil.CheckQuery(t, r.URL.Query(), expect)
			mockutil.WriteJSON(t, w, g)
		})
}

// AddMember mocks a AddMember request.
//
// The User.ID field of the passed discord.Member must be set.
func (m *Mocker) AddMember(guildID discord.GuildID, d api.AddMemberData, member discord.Member) {
	m.MockAPI("AddMember", http.MethodPut, "/guilds/"+guildID.String()+"/members/"+member.User.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.AddMemberData), &d)
			mockutil.WriteJSON(t, w, member)
		})
}

// ModifyMember mocks a ModifyMember request.
func (m *Mocker) ModifyMember(guildID discord.GuildID, userID discord.UserID, d api.ModifyMemberData) {
	m.MockAPI("ModifyMember", http.MethodPatch, "/guilds/"+guildID.String()+"/members/"+userID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.ModifyMemberData), &d)
			w.WriteHeader(http.StatusNoContent)
		})
}

type pruneBody struct {
	Pruned uint `json:"pruned"`
}

// PruneCount mocks a PruneCount request.
func (m *Mocker) PruneCount(guildID discord.GuildID, d api.PruneCountData, pruned uint) {
	if d.Days == 0 {
		d.Days = 7
	}

	m.MockAPI("PruneCount", http.MethodGet, "/guilds/"+guildID.String()+"/prune",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := url.Values{
				"days": {strconv.Itoa(int(d.Days))},
			}

			if len(d.IncludedRoles) > 0 {
				expect["include_roles"] = make([]string, len(d.IncludedRoles))

				for i, r := range d.IncludedRoles {
					expect["include_roles"][i] = r.String()
				}
			}

			resp := pruneBody{
				Pruned: pruned,
			}

			mockutil.CheckQuery(t, r.URL.Query(), expect)
			mockutil.WriteJSON(t, w, resp)
		})
}

// Prune mocks a Prune request.
func (m *Mocker) Prune(guildID discord.GuildID, d api.PruneData, pruned uint) {
	if d.Days == 0 {
		d.Days = 7
	}

	m.MockAPI("Prune", http.MethodPost, "/guilds/"+guildID.String()+"/prune",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := url.Values{
				"days":                {strconv.Itoa(int(d.Days))},
				"compute_prune_count": {strconv.FormatBool(d.ReturnCount)},
			}

			if len(d.IncludedRoles) > 0 {
				expect["include_roles"] = make([]string, len(d.IncludedRoles))

				for i, r := range d.IncludedRoles {
					expect["include_roles"][i] = r.String()
				}
			}

			resp := pruneBody{
				Pruned: pruned,
			}

			mockutil.CheckQuery(t, r.URL.Query(), expect)
			mockutil.WriteJSON(t, w, resp)
		})
}

// Kick mocks a Kick request.
func (m *Mocker) Kick(guildID discord.GuildID, userID discord.UserID) {
	m.MockAPI("Kick", http.MethodDelete, "/guilds/"+guildID.String()+"/members/"+userID.String(), nil)
}

// Bans mocks a Bans request.
//
// This method will sanitize Bans.User.ID.
func (m *Mocker) Bans(guildID discord.GuildID, b []discord.Ban) {
	for i, ban := range b {
		b[i] = sanitize.Ban(ban, 1)
	}

	m.MockAPI("Bans", http.MethodGet, "/guilds/"+guildID.String()+"/bans",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, b)
		})
}

// GetBan mocks a GetBan request.
//
// The User.ID field of the passed discord.Ban must be set.
func (m *Mocker) GetBan(guildID discord.GuildID, b discord.Ban) {
	m.MockAPI("GetBan", http.MethodGet, "/guilds/"+guildID.String()+"/bans/"+b.User.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, b)
		})
}

// Ban mocks a Ban request.
func (m *Mocker) Ban(guildID discord.GuildID, userID discord.UserID, d api.BanData) {
	if *d.DeleteDays > 7 {
		*d.DeleteDays = 7
	}

	m.MockAPI("Ban", http.MethodPut, "/guilds/"+guildID.String()+"/bans/"+userID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := make(url.Values)

			if d.DeleteDays != nil {
				expect["delete_message_days"] = []string{strconv.Itoa(int(*d.DeleteDays))}
			}

			if d.Reason != nil {
				expect["reason"] = []string{*d.Reason}
			}

			mockutil.CheckQuery(t, r.URL.Query(), expect)
			w.WriteHeader(http.StatusNoContent)
		})
}

// Unban mocks a Unban request.
func (m *Mocker) Unban(guildID discord.GuildID, userID discord.UserID) {
	m.MockAPI("Unban", http.MethodDelete, "/guilds/"+guildID.String()+"/bans/"+userID.String(), nil)
}
