package mock

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

const (
	TestSession        = "testSession"
	TestGuild          = "testGuild"
	TestGuildLarge     = "testGuildLarge"
	TestChannel        = "testChannel"
	TestChannel2       = "testChannel2"
	TestPrivateChannel = "testPrivateChannel"
	TestRole           = "testRole"
	TestUser           = "testUser"
)

const (
	memberCountLimit        = 1000
	sessionCreateErrMessage = "unable to create new session"
)

// TestingInstance is an interface intended for testing.T and testing.B
// instances.
type TestingInstance interface {
	Error(args ...interface{})
}

// RoundTripperFunc allows functions to satisfy the http.RoundTripper
// interface.
type RoundTripperFunc func(req *http.Request) (*http.Response, error)

// RoundTrip implements the http.RoundTripper interface.
func (rt RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

// NewMirrorRoundTripper returns an http.RoundTripper that mirrors the request
// body in the response body.
func NewMirrorRoundTripper() http.RoundTripper {
	return RoundTripperFunc(
		func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				Status:     http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Request:    req,
			}

			if req.Body == nil {
				resp.ContentLength = 0
				resp.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))

				return resp, nil
			}

			reqBody, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}

			err = req.Body.Close()
			if err != nil {
				return nil, err
			}

			resp.ContentLength = int64(len(reqBody))
			resp.Body = ioutil.NopCloser(bytes.NewReader(reqBody))

			return resp, nil
		},
	)
}

func mockGuild(guildID string) *discordgo.Guild {
	guild := &discordgo.Guild{
		ID:   guildID,
		Name: guildID,
	}

	guild.Members = mockMembers()
	guild.MemberCount = len(guild.Members)
	guild.Roles = mockRoles()
	guild.Channels = mockChannels()

	return guild
}

func mockUser(userID string) *discordgo.User {
	return &discordgo.User{
		ID:       userID,
		Username: userID,
	}
}

func mockLargeMembers() []*discordgo.Member {
	members := make([]*discordgo.Member, memberCountLimit)

	for i := 0; i < memberCountLimit; i++ {
		members[i] = mockMember(fmt.Sprintf("%s-%d", TestUser, i))
	}

	return members
}

func mockMembers() []*discordgo.Member {
	return []*discordgo.Member{
		mockMember(TestSession),
		mockMember(TestUser),
	}
}

func mockMember(userID string) *discordgo.Member {
	return &discordgo.Member{
		GuildID: TestGuild,
		User:    mockUser(userID),
		Roles:   mockRoleIDs(mockRoles()),
	}
}

func mockRoleIDs(roles discordgo.Roles) []string {
	roleIDs := make([]string, len(roles))

	for i, role := range roles {
		roleIDs[i] = role.ID
	}

	return roleIDs
}

func mockRoles() discordgo.Roles {
	return discordgo.Roles{
		mockRole(TestRole),
		mockRole("{eph} " + TestChannel),
	}
}

func mockRole(roleID string) *discordgo.Role {
	return &discordgo.Role{
		ID:   roleID,
		Name: roleID,
	}
}

func mockChannels() []*discordgo.Channel {
	return []*discordgo.Channel{
		mockChannel(TestChannel),
		mockChannel(TestChannel2),
		mockChannel(TestPrivateChannel),
	}
}

func mockChannel(channelID string) *discordgo.Channel {
	var permissionOverwrites []*discordgo.PermissionOverwrite

	if channelID == TestPrivateChannel {
		permissionOverwrites = []*discordgo.PermissionOverwrite{
			{
				ID:    TestGuild,
				Type:  "role",
				Deny:  discordgo.PermissionSendMessages,
				Allow: 0,
			},
		}
	} else {
		permissionOverwrites = []*discordgo.PermissionOverwrite{
			{
				ID:    TestGuild,
				Type:  "role",
				Deny:  0,
				Allow: discordgo.PermissionSendMessages,
			},
		}
	}

	return &discordgo.Channel{
		ID:                   channelID,
		Name:                 channelID,
		GuildID:              TestGuild,
		Type:                 discordgo.ChannelTypeGuildVoice,
		PermissionOverwrites: permissionOverwrites,
	}
}
