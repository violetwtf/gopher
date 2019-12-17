package commands

import (
    "strings"
    "github.com/bwmarrin/discordgo"
)

// Ban is the command for !ban
var Ban = Command{
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
