package main

import (
    "strings"
    "github.com/bwmarrin/discordgo"
)

// Command defines a command
type Command struct {
    Definition   string
    Requirements CommandRequirements
    Handler      func(*discordgo.Session, *discordgo.MessageCreate, []string)
}

// CommandRequirements represent what a command requires to execute
type CommandRequirements struct {
    Args                            []Arg
    LongText                        bool
    UserPermissions, BotPermissions int
}

// Arg is an argument
type Arg struct {
    Name, Type string
}

// GetUsage returns the args in a readable format
func (c Command) GetUsage() (usage string) {
    args := c.Requirements.Args

    usage = ""

    if len(args) == 0 {
        return 
    }

    for i := range args {
        arg := args[i]
        t := arg.Type

        // Clarify names of types for users
        switch t {
            case "mention":
                t = "@mention"
            default:
                t = "text"
        }

        usage += "<" + arg.Name + "(" + t + ")>"

        // Append space
        if i < len(args) {
            usage += " "
        }
    }

    return
}

// CommandPing is the command for !ping
var CommandPing = Command{
    "Do the ping thing",
    CommandRequirements{},
    func (s *discordgo.Session, e *discordgo.MessageCreate, _ []string) {
        s.ChannelMessageSend(e.ChannelID, "Pong!")        
    },
}

// CommandHelp is the command for !help
var CommandHelp = Command{
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

// CommandBan is the command for !ban
var CommandBan = Command{
    "Ban a user from the server, send them the reason",
    CommandRequirements{
        LongText: true,
        Args: []Arg{
            Arg{"user", "mention"}, Arg{"reason", "string"}},

        BotPermissions:  discordgo.PermissionBanMembers,
        UserPermissions: discordgo.PermissionBanMembers},
    func (s *discordgo.Session, e *discordgo.MessageCreate, args []string) {
        mentions := e.Mentions

        target := mentions[0].ID

        dm, err := s.UserChannelCreate(target)

        if err != nil {
            return
        }

        guildID := e.GuildID
        reason := strings.Join(args[1:], " ")

        guild, err := s.State.Guild(guildID)

        if err != nil {
            return
        }

        // Send before the ban, so that we can still MSG if they have mutual
        // DMs on
        s.ChannelMessageSend(
            dm.ID, 
            "You have been banned from **" + guild.Name +
                "** for **" + reason + "**.")

        s.GuildBanCreateWithReason(guildID, target, reason, 0)
    },
}
