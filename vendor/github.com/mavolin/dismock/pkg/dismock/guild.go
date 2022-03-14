package dismock

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/dismock/internal/mockutil"
	"github.com/mavolin/dismock/internal/sanitize"
)

// CreateGuild mocks a CreateGuild request.
//
// This method will sanitize Guild.ID, Guild.OwnerID, Guild.Emojis.ID and
// Guild.Roles.ID.
func (m *Mocker) CreateGuild(d api.CreateGuildData, g discord.Guild) {
	g = sanitize.Guild(g, 1, 1, 1, 1)

	m.MockAPI("CreateGuild", http.MethodPost, "/guilds",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.CreateGuildData), &d)
			mockutil.WriteJSON(t, w, g)
		})
}

// Guild mocks a Guild request.
//
// The ID field of the passed discord.Guild must be set.
//
// This method will sanitize Guild.ID, Guild.OwnerID, Guild.Emojis.ID and
// Guild.Roles.ID.
func (m *Mocker) Guild(g discord.Guild) {
	g = sanitize.Guild(g, 1, 1, 1, 1)

	m.MockAPI("Guild", http.MethodGet, "/guilds/"+g.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, g)
		})
}

// GuildWithCount mocks a GuildWithCount request.
//
// The ID field of the passed discord.Guild must be set.
//
// This method will sanitize Guild.ID, Guild.OwnerID, Guild.Emojis.ID and
// Guild.Roles.ID.
func (m *Mocker) GuildWithCount(g discord.Guild) {
	g = sanitize.Guild(g, 1, 1, 1, 1)

	m.MockAPI("GuildWithCount", http.MethodGet, "/guilds/"+g.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckQuery(t, r.URL.Query(), url.Values{
				"with_counts": {"true"},
			})
			mockutil.WriteJSON(t, w, g)
		})
}

// GuildPreview mocks a GuildPreview request.
//
// This method will sanitize GuildPreview.ID, GuildPreview.Emojis.ID and
// GuildPreview.Emojis.User.ID.
func (m *Mocker) GuildPreview(p discord.GuildPreview) {
	p = sanitize.GuildPreview(p, 1)

	m.MockAPI("GuildPreview", http.MethodGet, "/guilds/"+p.ID.String()+"/preview",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, p)
		})
}

// Guilds mocks a Guilds request.
//
// This method will sanitize Guilds.ID, Guilds.OwnerID, Guilds.Emojis.ID and
// Guilds.Roles.ID.
func (m *Mocker) Guilds(limit uint, g []discord.Guild) {
	if g == nil {
		g = []discord.Guild{}
	}

	if len(g) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent guilds (%d vs. %d)", len(g), limit))
	}

	const hardLimit uint = 100

	var after discord.GuildID

	for i := 0; i <= len(g)/int(hardLimit); i++ {
		var (
			from = uint(i) * hardLimit
			to   = uint(math.Min(float64(from+hardLimit), float64(len(g))))

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

		m.guildsRange(0, after, fmt.Sprintf("Guilds #%d", i+1), fetch, g[from:to])

		if to-from < hardLimit {
			break
		}

		after = g[to-1].ID
	}
}

// GuildsBefore mocks a GuildsBefore request.
//
// This method will sanitize Guilds.ID, Guilds.OwnerID, Guilds.Emojis.ID and
// Guilds.Roles.ID.
func (m *Mocker) GuildsBefore(before discord.GuildID, limit uint, g []discord.Guild) {
	if g == nil {
		g = []discord.Guild{}
	}

	if len(g) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent guilds (%d vs. %d)", len(g), limit))
	}

	const hardLimit = 100

	req := len(g)/hardLimit + 1

	from := uint(math.Min(float64(uint(req)*hardLimit), float64(len(g))))

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

		m.guildsRange(before, 0, fmt.Sprintf("GuildsBefore #%d", no), fetch, g[from:to])

		if to-from < hardLimit {
			break
		}

		before = g[from].ID
	}
}

// GuildsAfter mocks a GuildsAfter request.
//
// This method will sanitize Guilds.ID, Guilds.OwnerID, Guilds.Emojis.ID and
// Guilds.Roles.ID.WithToken
func (m *Mocker) GuildsAfter(after discord.GuildID, limit uint, g []discord.Guild) {
	if g == nil {
		g = []discord.Guild{}
	}

	if len(g) > int(limit) && limit != 0 {
		panic(fmt.Sprintf("limit may not be less than the number of sent guilds (%d vs. %d)", len(g), limit))
	}

	const hardLimit uint = 100

	for i := 0; i <= len(g)/int(hardLimit); i++ {
		var (
			from = uint(i) * hardLimit
			to   = uint(math.Min(float64(from+hardLimit), float64(len(g))))

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

		m.guildsRange(0, after, fmt.Sprintf("GuildsAfter #%d", i+1), fetch, g[from:to])

		if to-from < hardLimit {
			break
		}

		after = g[to-1].ID
	}
}

// guildsRange mocks a single request to the GET /guilds endpoint.
//
// This method will sanitize Guilds.ID, Guilds.OwnerID, Guilds.Emojis.ID and
// Guilds.Roles.ID.
func (m *Mocker) guildsRange(before, after discord.GuildID, name string, limit uint, g []discord.Guild) {
	for i, guild := range g {
		g[i] = sanitize.Guild(guild, 1, 1, 1, 1)
	}

	m.MockAPI(name, http.MethodGet, "/users/@me/guilds",
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
			mockutil.WriteJSON(t, w, g)
		})
}

// LeaveGuild mocks a LeaveGuild request.
func (m *Mocker) LeaveGuild(id discord.GuildID) {
	m.MockAPI("LeaveGuild", http.MethodDelete, "/users/@me/guilds/"+id.String(), nil)
}

// ModifyGuild mocks a ModifyGuild request.
//
// The ID field of the passed discord.Guild must be set.
//
// This method will sanitize Guild.ID, Guild.OwnerID, Guild.Emojis.ID and
// Guild.Roles.ID.
func (m *Mocker) ModifyGuild(d api.ModifyGuildData, g discord.Guild) {
	g = sanitize.Guild(g, 1, 1, 1, 1)

	m.MockAPI("ModifyGuild", http.MethodPatch, "/guilds/"+g.ID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.ModifyGuildData), &d)
			mockutil.WriteJSON(t, w, g)
		})
}

// DeleteGuild mocks a DeleteGuild request.
func (m *Mocker) DeleteGuild(id discord.GuildID) {
	m.MockAPI("DeleteGuild", http.MethodDelete, "/guilds/"+id.String(), nil)
}

// VoiceRegionsGuild mocks a VoiceRegionsGuild request.
func (m *Mocker) VoiceRegionsGuild(guildID discord.GuildID, vr []discord.VoiceRegion) {
	if vr == nil {
		vr = []discord.VoiceRegion{}
	}

	m.MockAPI("VoiceRegionsGuild", http.MethodGet, "/guilds/"+guildID.String()+"/regions",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, vr)
		})
}

// AuditLog mocks a AuditLog request.
//
// This method will sanitize AuditLog.Webhook.ID, AuditLog.Webhook.User.ID,
// AuditLog.Users.ID, AuditLog.Entries.ID, AuditLog.Entries.UserID,
// AuditLog.Integrations.ID, AuditLog.Integrations.RoleID and
// AuditLog.Integrations.User.ID.
func (m *Mocker) AuditLog(guildID discord.GuildID, d api.AuditLogData, al discord.AuditLog) {
	switch {
	case d.Limit == 0:
		d.Limit = 50
	case d.Limit > 100:
		d.Limit = 100
	}

	al = sanitize.AuditLog(al)

	m.MockAPI("AuditLog", http.MethodGet, "/guilds/"+guildID.String()+"/audit-logs",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := url.Values{
				"limit": {strconv.Itoa(int(d.Limit))},
			}

			if d.UserID != 0 {
				expect["user_id"] = []string{d.UserID.String()}
			}

			if d.ActionType != 0 {
				expect["action_type"] = []string{strconv.FormatUint(uint64(d.ActionType), 10)}
			}

			if d.Before != 0 {
				expect["before"] = []string{d.Before.String()}
			}

			mockutil.CheckQuery(t, r.URL.Query(), expect)
			mockutil.WriteJSON(t, w, al)
		})
}

// Integrations mocks a Integrations request.
//
// This method will sanitize Integration.ID, Integration.RoleID and
// Integration.User.ID.
func (m *Mocker) Integrations(guildID discord.GuildID, integrations []discord.Integration) {
	if integrations == nil {
		integrations = []discord.Integration{}
	}

	for i, integration := range integrations {
		integrations[i] = sanitize.Integration(integration, 1, 1, 1)
	}

	m.MockAPI("Integrations", http.MethodGet, "/guilds/"+guildID.String()+"/integrations",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, integrations)
		})
}

type attachIntegrationPayload struct {
	Type discord.Service       `json:"type"`
	ID   discord.IntegrationID `json:"id"`
}

// AttachIntegration mocks a AttachIntegration request.
func (m *Mocker) AttachIntegration(
	guildID discord.GuildID, integrationID discord.IntegrationID, integrationType discord.Service,
) {
	m.MockAPI("AttachIntegration", http.MethodPost, "/guilds/"+guildID.String()+"/integrations",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			expect := &attachIntegrationPayload{
				Type: integrationType,
				ID:   integrationID,
			}

			mockutil.CheckJSON(t, r.Body, new(attachIntegrationPayload), expect)
			w.WriteHeader(http.StatusNoContent)
		})
}

// ModifyIntegration mocks a ModifyIntegration request.
func (m *Mocker) ModifyIntegration(
	guildID discord.GuildID, integrationID discord.IntegrationID, d api.ModifyIntegrationData,
) {
	m.MockAPI("ModifyIntegration", http.MethodPatch,
		"/guilds/"+guildID.String()+"/integrations/"+integrationID.String(),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.ModifyIntegrationData), &d)
			w.WriteHeader(http.StatusNoContent)
		})
}

// SyncIntegration mocks a SyncIntegration request.
func (m *Mocker) SyncIntegration(guildID discord.GuildID, integrationID discord.IntegrationID) {
	m.MockAPI("SyncIntegration", http.MethodPost,
		"/guilds/"+guildID.String()+"/integrations/"+integrationID.String()+"/sync", nil)
}

// GuildWidget mocks a GuildWidget request.
func (m *Mocker) GuildWidget(guildID discord.GuildID, e discord.GuildWidget) {
	m.MockAPI("GuildWidget", http.MethodGet, "/guilds/"+guildID.String()+"/widget",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, e)
		})
}

// ModifyGuildWidget mocks a ModifyGuildWidget request.
func (m *Mocker) ModifyGuildWidget(guildID discord.GuildID, d api.ModifyGuildWidgetData, e discord.GuildWidget) {
	m.MockAPI("ModifyGuild", http.MethodPatch, "/guilds/"+guildID.String()+"/widget",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckJSON(t, r.Body, new(api.ModifyGuildWidgetData), &d)
			mockutil.WriteJSON(t, w, e)
		})
}

// GuildVanityURL mocks a GuildVanityURL request.
//
// Although those fields are normally not sent, this method will sanitize
// Invite.Guild.ID, Invite.Guild.OwnerID, Invite.Guild.Emojis.ID,
// Invite.Guild.Roles.ID, Invite.Channel.ID, Invite.Inviter.ID.
func (m *Mocker) GuildVanityURL(guildID discord.GuildID, i discord.Invite) {
	i = sanitize.Invite(i, 1, 1, 1, 1, 1, 1, 1)

	m.MockAPI("GuildVanityURL", http.MethodGet, "/guilds/"+guildID.String()+"/vanity-url",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.WriteJSON(t, w, i)
		})
}

// GuildImage mocks a GuildImage request.
func (m *Mocker) GuildImage(guildID discord.GuildID, style api.GuildImageStyle, img io.Reader) {
	m.MockAPI("GuildImage", http.MethodGet, "/guilds/"+guildID.String()+"/widget.png",
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			mockutil.CheckQuery(t, r.URL.Query(), url.Values{
				"style": {string(style)},
			})

			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}
