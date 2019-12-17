package commands

import (
    "github.com/bwmarrin/discordgo"
)

// About is the command for !about
var About = Command{
    "Info about the Gopher! bot.",
    CommandRequirements{},
    func (s *discordgo.Session, e *discordgo.MessageCreate, _ []string) {
        aboutMessage := `
Gopher! is an open-source Discord bot by **Violet#6096**.
Gopher! is currently in development, but we're happy to accept feedback.
Visit our GitHub repository: https://github.com/violetwtf/gopher for more info.`
        s.ChannelMessageSend(e.ChannelID, aboutMessage) 
    },
}
