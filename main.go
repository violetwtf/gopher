package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "os/signal"
    "strings"
    "syscall"

    "github.com/bwmarrin/discordgo"
)

const prefix = "!"

var commands = make(map[string]Command)

func main() {
    commands["ping"] = CommandPing
    commands["help"] = CommandHelp
    commands["ban"]  = CommandBan

    token, err := ioutil.ReadFile("./SUPER_SECRET_TOKEN.txt")
    if err != nil {
        log.Fatal("error loading token:", err.Error())
    }

    discord, err := discordgo.New("Bot " + string(token))
    if err != nil {
        log.Fatal("error creating client:", err.Error())
    }

    if err = discord.Open(); err != nil {
        log.Fatal("error connecting:", err.Error())
    }

    defer discord.Close()

    discord.AddHandler(ready)
    discord.AddHandler(messageCreate)

    // Wait here until CTRL-C or other term signal is received.
    fmt.Println("Gopher is alive! Exit with CTRL-C.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    fmt.Println("Closing client")
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
    s.UpdateStatus(0, "!help")
}

func messageCreate(s *discordgo.Session, event *discordgo.MessageCreate) {
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


