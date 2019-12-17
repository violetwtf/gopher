package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "os/signal"
    "syscall"

    cmds "github.com/violetwtf/gopher/commands"
    "github.com/violetwtf/gopher/events"

    "github.com/bwmarrin/discordgo"
)

const prefix = "!"

var commands = make(map[string]cmds.Command)

func main() {
    commands["about"] = cmds.About
    commands["ban"]   = cmds.Ban
    commands["help"]  = cmds.Help
    commands["ping"]  = cmds.Ping

    cmds.SetCommandPrefix(prefix)
    cmds.SetCommandRegistry(commands)

    events.SetCommandPrefix(prefix)
    events.SetCommandRegistry(commands)

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
    discord.AddHandler(events.MessageCreate)

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



