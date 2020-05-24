package mock

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const stateCreateErrMessage = "unable to create new state"

// NewState provides a *discordgo.State instance to be used in unit testing.
func NewState() (*discordgo.State, error) {
	state := discordgo.NewState()

	state.User = &discordgo.User{
		ID:       TestSession,
		Username: TestSession,
		Bot:      true,
	}

	err := state.GuildAdd(mockGuild(TestGuild))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", stateCreateErrMessage, err)
	}

	return state, nil
}
