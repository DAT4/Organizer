package main

import (
	"errors"
	"fmt"
	discordgo "github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	token := "NzUwMDQ0MjkwNDE5NTIzNzg2.X00zLQ.NpS_3aAKusBY2nU-imhmeVz89jY"

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Register ready as a callback for the ready events.
	dg.AddHandler(ready)
	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// We need information about guilds (which includes their channels),
	// messages and voice states.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("BOT is ON")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	s.UpdateStatus(0, "")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	c, err := getChannel(s, m, "hemmelig")
	if err != nil {
		fmt.Println("Looking for hemmelig channel:", err)
	}

	if m.ChannelID == c.ID {
		if strings.HasPrefix(m.Content, "do stuff") {
			//deleteChannels(s,m,"opgaverum")
			createChannels(s, m, "opgaverum")
		}
	}
}
func getChannel(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) (*discordgo.Channel, error) {
	channels, err := s.GuildChannels(m.GuildID)
	if err != nil {
		return nil, err
	}
	for _, e := range channels {
		if strings.HasPrefix(e.Name, prefix) {
			return e, nil
		}
	}
	return nil, errors.New("No channel with that name.")
}

func createChannels(s *discordgo.Session, m *discordgo.MessageCreate, title string) {
	c, err := getChannel(s, m, "random")
	if err != nil {
		fmt.Println("Find general:", err)
	}
	parrent := c.ParentID
	for i := 1; i <= 10; i++ {
		c, err := s.GuildChannelCreateComplex(m.GuildID, discordgo.GuildChannelCreateData{
			Name:                 title + "-" + strconv.Itoa(i),
			Type:                 2,
			Topic:                "OpgaveArbejde",
			Bitrate:              0,
			UserLimit:            10,
			RateLimitPerUser:     0,
			Position:             0,
			PermissionOverwrites: nil,
			ParentID:             parrent,
			NSFW:                 false,
		})
		if err != nil {
			fmt.Println("Creating channel:", err)
		}
		fmt.Println("Created channel", c.Name)
	}
}

func getRole(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) (roleID string, err error) {
	roles, err := s.GuildRoles(m.GuildID)
	if err != nil {
		return "", err
	}
	for _, e := range roles {
		if strings.HasPrefix(e.Name, prefix) {
			fmt.Println("Role permission", e.Permissions)
			return e.ID, nil
		}
	}
	return "", errors.New("Something went wrong...")
}

func deleteChannels(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {
	channels, err := s.GuildChannels(m.GuildID)
	if err != nil {
		fmt.Println("Get all channels of:", m.GuildID+":", err)
	}
	for _, e := range channels {
		c, err := s.Channel(e.ID)
		if err != nil {
			fmt.Println("Get channel with ID", e.ID, ":", err)
		}
		if strings.HasPrefix(c.Name, prefix) {
			fmt.Println(c.Name, c.ID, "slettes")
			_, err := s.ChannelDelete(c.ID)
			if err != nil {
				fmt.Println("Deleting channel:", err)
			}

		}

	}
}
