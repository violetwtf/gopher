package commands

import (
    "github.com/bwmarrin/discordgo"
)

// Help is the command for !help
var Help = Command{
    "List all commands",
    CommandRequirements{},
    func (s *discordgo.Session, e *discordgo.MessageCreate, _ []string) {
        output := "command - description - usage\n"

        for i := range commands {
            command := commands[i]
            displayCommand := prefix + i
            def := "**" + displayCommand + "** - " + 
                command.Definition + "```" + 
                displayCommand + " " + command.GetUsage() + "```"

            output += "\n" + def
        }

        s.ChannelMessageSend(e.ChannelID, output)
    },
}
