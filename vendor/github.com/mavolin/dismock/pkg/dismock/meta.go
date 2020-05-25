package dismock

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/require"
)

// ================================ Channel ================================

// ChannelIconURL mocks a ChannelIconURL request.
func (m *Mocker) ChannelIconURL(channelID discord.Snowflake, icon discord.Hash, img io.Reader) {
	m.Mock("ChannelIconURL", http.MethodGet,
		"/channel-icons/"+channelID.String()+"/"+formatImageType(icon, discord.PNGImage),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// ChannelIconURLWithType mocks a ChannelIconURLWithType request.
func (m *Mocker) ChannelIconURLWithType(
	channelID discord.Snowflake, icon discord.Hash, t discord.ImageType, img io.Reader,
) {
	m.Mock("ChannelIconURLWithType", http.MethodGet, "/channel-icons/"+channelID.String()+"/"+formatImageType(icon, t),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// ================================ Emoji ================================

// EmojiURL mocks a EmojiURL request.
func (m *Mocker) EmojiURL(emojiID discord.Snowflake, animated bool, img io.Reader) {
	var url string
	if animated {
		url = "/emojis/" + formatImageType(emojiID.String(), discord.GIFImage)
	} else {
		url = "/emojis/" + formatImageType(emojiID.String(), discord.PNGImage)
	}

	m.Mock("EmojiURLWithType", http.MethodGet, url, func(w http.ResponseWriter, r *http.Request, t *testing.T) {
		_, err := io.Copy(w, img)
		require.NoError(t, err)
	})
}

// EmojiURLWithType mocks a EmojiURLWithType request.
func (m *Mocker) EmojiURLWithType(emojiID discord.Snowflake, animated bool, t discord.ImageType, img io.Reader) {
	if t == discord.AutoImage {
		m.EmojiURL(emojiID, animated, img)
		return
	}

	m.Mock("EmojiURLWithType", http.MethodGet, "/emojis/"+formatImageType(emojiID.String(), t),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// ================================ Guild ================================

// GuildIconURL mocks a GuildIconURL request.
//
// This method can be used for both Guilds and GuildPreviews.
func (m *Mocker) GuildIconURL(guildID discord.Snowflake, icon discord.Hash, img io.Reader) {
	m.Mock("GuildIconURL", http.MethodGet, "/icons/"+guildID.String()+"/"+formatImageType(icon, discord.AutoImage),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// GuildIconURLWithType mocks a GuildIconURLWithType request.
//
// This method can be used for both Guilds and GuildPreviews.
func (m *Mocker) GuildIconURLWithType(
	guildID discord.Snowflake, icon discord.Hash, t discord.ImageType, img io.Reader,
) {
	m.Mock("GuildIconURLWithType", http.MethodGet, "/icons/"+guildID.String()+"/"+formatImageType(icon, t),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// BannerURL mocks a BannerURL request.
//
// This method can be used for both Guilds and GuildPreviews.
func (m *Mocker) BannerURL(guildID discord.Snowflake, banner discord.Hash, img io.Reader) {
	m.Mock("BannerURL", http.MethodGet, "/banners/"+guildID.String()+"/"+formatImageType(banner, discord.PNGImage),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// BannerURLWithType mocks a BannerURLWithType request.
//
// This method can be used for both Guilds and GuildPreviews.
func (m *Mocker) BannerURLWithType(
	guildID discord.Snowflake, banner discord.Hash, t discord.ImageType, img io.Reader,
) {
	m.Mock("BannerURLWithType", http.MethodGet, "/banners/"+guildID.String()+"/"+formatImageType(banner, t),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// SplashURL mocks a SplashURL request.
//
// This method can be used for both Guilds and GuildPreviews.
func (m *Mocker) SplashURL(guildID discord.Snowflake, splash discord.Hash, img io.Reader) {
	m.Mock("SplashURL", http.MethodGet, "/splashes/"+guildID.String()+"/"+formatImageType(splash, discord.PNGImage),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// SplashURLWithType mocks a SplashURLWithType request.
//
// This method can be used for both Guilds and GuildPreviews.
func (m *Mocker) SplashURLWithType(
	guildID discord.Snowflake, splash discord.Hash, t discord.ImageType, img io.Reader,
) {
	m.Mock("SplashURLWithType", http.MethodGet, "/splashes/"+guildID.String()+"/"+formatImageType(splash, t),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// DiscoverySplashURL mocks a DiscoverySplashURL request.
//
// This method can be used for both Guilds and GuildPreviews.
func (m *Mocker) DiscoverySplashURL(guildID discord.Snowflake, splash discord.Hash, img io.Reader) {
	m.Mock("DiscoverySplashURL", http.MethodGet,
		"/splashes/"+guildID.String()+"/"+formatImageType(splash, discord.PNGImage),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// DiscoverySplashURLWithType mocks a DiscoverySplashURLWithType request.
//
// This method can be used for both Guilds and GuildPreviews.
func (m *Mocker) DiscoverySplashURLWithType(
	guildID discord.Snowflake, splash discord.Hash, t discord.ImageType, img io.Reader,
) {
	m.Mock("DiscoverySplashURLWithType", http.MethodGet,
		"/splashes/"+guildID.String()+"/"+formatImageType(splash, t),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// =============================================================================
// Utils
// =====================================================================================

func formatImageType(name string, t discord.ImageType) string {
	if t == discord.AutoImage {
		if strings.HasPrefix(name, "a_") {
			return name + ".gif"
		}

		return name + ".png"
	}

	return name + string(t)
}
