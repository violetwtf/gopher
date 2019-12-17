package commands

import (
    "github.com/bwmarrin/discordgo"
)

// Ping is the command for !ping
var Ping = Command{
    "Do the ping thing",
    CommandRequirements{},
    func (s *discordgo.Session, e *discordgo.MessageCreate, _ []string) {
        s.ChannelMessageSend(e.ChannelID, "Pong!")        
    },
}
