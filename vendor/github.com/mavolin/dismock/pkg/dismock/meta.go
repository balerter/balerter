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

// ChannelIcon mocks a ChannelIcon request.
func (m *Mocker) ChannelIcon(channelID discord.ChannelID, icon discord.Hash, img io.Reader) {
	m.Mock("ChannelIcon", http.MethodGet,
		"/channel-icons/"+channelID.String()+"/"+formatImageType(icon, discord.PNGImage),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// ChannelIconWithType mocks a ChannelIconWithType request.
func (m *Mocker) ChannelIconWithType(
	channelID discord.ChannelID, icon discord.Hash, t discord.ImageType, img io.Reader,
) {
	m.Mock("ChannelIconWithType", http.MethodGet, "/channel-icons/"+channelID.String()+"/"+formatImageType(icon, t),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// ================================ Emoji ================================

// EmojiPicture mocks a EmojiPicture request.
func (m *Mocker) EmojiPicture(emojiID discord.EmojiID, animated bool, img io.Reader) {
	var url string
	if animated {
		url = "/emojis/" + formatImageType(emojiID.String(), discord.GIFImage)
	} else {
		url = "/emojis/" + formatImageType(emojiID.String(), discord.PNGImage)
	}

	m.Mock("EmojiPictureWithType", http.MethodGet, url, func(w http.ResponseWriter, r *http.Request, t *testing.T) {
		_, err := io.Copy(w, img)
		require.NoError(t, err)
	})
}

// EmojiPictureWithType mocks a EmojiPictureWithType request.
func (m *Mocker) EmojiPictureWithType(emojiID discord.EmojiID, animated bool, t discord.ImageType, img io.Reader) {
	if t == discord.AutoImage {
		m.EmojiPicture(emojiID, animated, img)
		return
	}

	m.Mock("EmojiPictureWithType", http.MethodGet, "/emojis/"+formatImageType(emojiID.String(), t),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// ================================ Guild ================================

// GuildIcon mocks a GuildIcon request.
//
// This method can be used for both discord.Guild and discord.GuildPreview.
func (m *Mocker) GuildIcon(guildID discord.GuildID, icon discord.Hash, img io.Reader) {
	m.Mock("GuildIcon", http.MethodGet, "/icons/"+guildID.String()+"/"+formatImageType(icon, discord.AutoImage),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// GuildIconWithType mocks a GuildIconWithType request.
//
// This method can be used for both discord.Guild and discord.GuildPreview.
func (m *Mocker) GuildIconWithType(guildID discord.GuildID, icon discord.Hash, t discord.ImageType, img io.Reader) {
	m.Mock("GuildIconWithType", http.MethodGet, "/icons/"+guildID.String()+"/"+formatImageType(icon, t),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// Banner mocks a Banner request.
//
// This method can be used for both discord.Guild and discord.GuildPreview.
func (m *Mocker) Banner(guildID discord.GuildID, banner discord.Hash, img io.Reader) {
	m.Mock("Banner", http.MethodGet, "/banners/"+guildID.String()+"/"+formatImageType(banner, discord.PNGImage),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// BannerWithType mocks a BannerWithType request.
//
// This method can be used for both discord.Guild and discord.GuildPreview.
func (m *Mocker) BannerWithType(guildID discord.GuildID, banner discord.Hash, t discord.ImageType, img io.Reader) {
	m.Mock("BannerWithType", http.MethodGet, "/banners/"+guildID.String()+"/"+formatImageType(banner, t),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// Splash mocks a Splash request.
//
// This method can be used for both discord.Guild and discord.GuildPreview.
func (m *Mocker) Splash(guildID discord.GuildID, splash discord.Hash, img io.Reader) {
	m.Mock("Splash", http.MethodGet, "/splashes/"+guildID.String()+"/"+formatImageType(splash, discord.PNGImage),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// SplashWithType mocks a SplashWithType request.
//
// This method can be used for both discord.Guild and discord.GuildPreview.
func (m *Mocker) SplashWithType(guildID discord.GuildID, splash discord.Hash, t discord.ImageType, img io.Reader) {
	m.Mock("SplashWithType", http.MethodGet, "/splashes/"+guildID.String()+"/"+formatImageType(splash, t),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// DiscoverySplash mocks a DiscoverySplash request.
//
// This method can be used for both discord.Guild and discord.GuildPreview.
func (m *Mocker) DiscoverySplash(guildID discord.GuildID, splash discord.Hash, img io.Reader) {
	m.Mock("DiscoverySplash", http.MethodGet,
		"/splashes/"+guildID.String()+"/"+formatImageType(splash, discord.PNGImage),
		func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			_, err := io.Copy(w, img)
			require.NoError(t, err)
		})
}

// DiscoverySplashWithType mocks a DiscoverySplashWithType request.
//
// This method can be used for both discord.Guild and discord.GuildPreview.
func (m *Mocker) DiscoverySplashWithType(
	guildID discord.GuildID, splash discord.Hash, t discord.ImageType, img io.Reader,
) {
	m.Mock("DiscoverySplashWithType", http.MethodGet,
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
