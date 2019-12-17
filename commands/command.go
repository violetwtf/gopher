package commands

import (
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

var (
    commands map[string]Command
    prefix   string
)

// SetCommandRegistry sets the map[string]Command of commands for this package
func SetCommandRegistry(c map[string]Command) {
    commands = c
}

// SetCommandPrefix sets the command prefix for this package
func SetCommandPrefix(p string) {
    prefix = p
}
