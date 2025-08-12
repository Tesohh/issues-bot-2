package handler

import (
	"issues/v2/db"
	"issues/v2/slash"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func GuildJoinHandler(s *discordgo.Session, event *discordgo.GuildCreate) {
	// register the guild
	guild, isNew, err := registerGuild(event)
	if err != nil {
		slog.Error("error(GuildCreate) while registering the guild ", "err", err)
		return
	}

	if isNew {
		err = registerRoles(s, event, &guild)
		if err != nil {
			slog.Error("error(GuildCreate) while registering roles for new guild", "err", err)
			return
		}
		err := db.Conn.Save(&guild).Error
		if err != nil {
			slog.Error("error(GuildCreate) while saving changes to new guild", "err", err)
			return
		}
	}

}

func registerGuild(event *discordgo.GuildCreate) (db.Guild, bool, error) {
	guild := db.Guild{
		ID: event.Guild.ID,
	}
	result := db.Conn.FirstOrCreate(&guild, guild)
	err := result.Error
	if err != nil {
		return db.Guild{}, false, err
	}

	isNew := false
	if result.RowsAffected == 1 {
		isNew = true
		slog.Info("Registered new", "guild", guild.ID)
	}

	return guild, isNew, nil
}

func registerRole(s *discordgo.Session, guildID string, role *discordgo.RoleParams, kind db.RoleKind) (db.Role, error) {
	role.Mentionable = slash.Ptr(true)

	discordRole, err := s.GuildRoleCreate(guildID, role)
	if err != nil {
		return db.Role{}, err
	}

	dbrole := db.Role{
		ID:      discordRole.ID,
		Kind:    kind,
		GuildID: guildID,
	}

	err = db.Conn.Create(&dbrole).Error
	return dbrole, err
}

func registerRoles(s *discordgo.Session, event *discordgo.GuildCreate, guild *db.Guild) error {
	// register all categories
	registeredCategoryRoles := []db.Role{}
	for i, role := range categoryRoles {
		registeredRole, err := registerRole(s, event.Guild.ID, role, db.RoleKindCategory)
		if err != nil {
			return err
		}
		registeredCategoryRoles = append(registeredCategoryRoles, registeredRole)
		if i == 0 {
			guild.DefaultCategoryRoleID = registeredRole.ID
		}
	}

	// register all priorities
	registeredPriorityRoles := []db.Role{}
	for i, role := range priorityRoles {
		registeredRole, err := registerRole(s, event.Guild.ID, role, db.RoleKindPriority)
		if err != nil {
			return err
		}
		registeredPriorityRoles = append(registeredPriorityRoles, registeredRole)
		if i == 3 { // 3 == NORMAL
			guild.DefaultPriorityRoleID = registeredRole.ID
		}
	}

	registeredNobodyRole, err := registerRole(s, event.Guild.ID, nobodyRole, db.RoleKindNobody)
	if err != nil {
		return err
	}
	guild.NobodyRoleID = registeredNobodyRole.ID

	registeredDiscussionRole, err := registerRole(s, event.Guild.ID, discussionRole, db.RoleKindDiscussion)
	if err != nil {
		return err
	}
	guild.DiscussionRoleID = registeredDiscussionRole.ID

	return nil
}
