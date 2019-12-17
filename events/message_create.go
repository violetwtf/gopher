package events

import (
    "strings"
    cmds "github.com/violetwtf/gopher/commands"
    "github.com/bwmarrin/discordgo"
)

var (
    commands map[string]cmds.Command
    prefix   string
)

// SetCommandRegistry sets the map[string]Command of commands for this package
func SetCommandRegistry(c map[string]cmds.Command) {
    commands = c
}

// SetCommandPrefix sets the command prefix for this package
func SetCommandPrefix(p string) {
    prefix = p
}

// MessageCreate is the handler for the messageCreate event
func MessageCreate(s *discordgo.Session, event *discordgo.MessageCreate) {
    selfID := s.State.User.ID

    if event.Author.ID == selfID {
        return
    }

    msg := event.Message.Content

    if strings.HasPrefix(msg, prefix) {
        // Get the command
        parts := strings.Split(msg[1:], " ")

        command, ok := commands[parts[0]]

        args := parts[1:]

        if !ok {
            return
        }

        channelID := event.ChannelID

        botPermissions, err := s.State.UserChannelPermissions(
            selfID, channelID)
        if err != nil {
            return
        }

        requiredBotPerms := command.Requirements.BotPermissions

        if (botPermissions & requiredBotPerms != requiredBotPerms) &&
            botPermissions != discordgo.PermissionAdministrator {
            // Bad permissions
            go s.ChannelMessageSend(channelID, 
                "Bot doesn't have permission to do this")
            return
        }

        userPermissions, err := s.State.UserChannelPermissions(
            selfID, channelID)
        if err != nil {
            return
        }

        requiredUserPerms := command.Requirements.UserPermissions

        if (userPermissions & requiredUserPerms != requiredUserPerms) && 
            userPermissions != discordgo.PermissionAdministrator {
            // Bad USER permissions
            go s.ChannelMessageSend(channelID,
                "User doesn't have permission to do this")
            return
        }

        requiredArgs := command.Requirements.Args
        argsCorrect := true
        argsLen := len(args)
        requiredArgsLen := len(requiredArgs)

        // argsCheck should be true when the length of args is invalid
        argsCheck := argsLen != requiredArgsLen

        if command.Requirements.LongText {
            argsCheck = argsLen < requiredArgsLen
        }

        if argsCheck {
            go s.ChannelMessageSend(
                channelID, 
                "Invalid usage! Use: ```" + 
                    prefix + parts[0] + " " + command.GetUsage() + "```")
            return
        }

        for i := range requiredArgs {
            switch requiredArgs[i].Type {
                case "mention":
                    arg := args[i]

                    if !strings.HasPrefix(arg, "<@") {
                        argsCorrect = false
                        break
                    }

                    // Check if user exists in guild
                    userID := arg[3:21]

                    _, err := s.State.Member(event.GuildID, userID)

                    if err != nil {
                        argsCorrect = false
                    }
                    break
                default:
            }
        }

        if !argsCorrect {
            go s.ChannelMessageSend(
                channelID, 
                "Invalid usage! Use: ```" + 
                    prefix + parts[0] + " " + command.GetUsage() + "```")
            return
        }

        go command.Handler(s, event, args)
    }
}
