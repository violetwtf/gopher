package internal

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/violetwtf/gopher/commands"
)

// GetPunishUserCommand gets the command for a specific punishment action
func GetPunishUserCommand(punishment string, permission int) commands.Command {
	return commands.Command{
		Definition: 
			"Gives a user a " + punishment + ", and sends them the reason.",
		Requirements: commands.CommandRequirements{
			Args: []commands.Arg{
				commands.Arg{Name: "user",   Type: "mention"},
				commands.Arg{Name: "reason", Type: "string" }},
			LongText:        true,
			BotPermissions:  permission,
			UserPermissions: permission,
		},
		Handler: func(
			s *discordgo.Session, e *discordgo.MessageCreate, args []string,
		) {
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

			// Send before the punishment, so that we can still MSG if they have
			// mutuals DMs on
			s.ChannelMessageSend(
				dm.ID,
				"You have recieved a **"+punishment+"** from **"+guild.Name+
					"** for **"+reason+"**.")

			switch punishment {
				case "ban":
					go s.GuildBanCreateWithReason(guildID, target, reason, 0)
					return
				case "kick":
					go s.GuildMemberDeleteWithReason(guildID, target, reason)
					return
				default:
			}
		},
	}
}
