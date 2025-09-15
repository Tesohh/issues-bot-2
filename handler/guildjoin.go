package handler

import (
	"issues/v2/data"
	"issues/v2/db"
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
	result := db.Conn.FirstOrCreate(&guild)
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

func registerRole(s *discordgo.Session, guildID string, role *data.RoleDef, kind db.RoleKind) (db.Role, error) {
	discordRole, err := s.GuildRoleCreate(guildID, role.ToDiscordRoleParams())
	if err != nil {
		return db.Role{}, err
	}

	dbrole := db.Role{
		ID:      discordRole.ID,
		Key:     role.Key,
		Emoji:   role.Emoji,
		Kind:    kind,
		GuildID: guildID,
	}

	err = db.Roles.Create(db.Ctx, &dbrole)
	return dbrole, err
}

func registerRoles(s *discordgo.Session, event *discordgo.GuildCreate, guild *db.Guild) error {
	// register all categories
	registeredCategoryRoles := []db.Role{}
	rolePtrs := []*string{&guild.GenericCategoryRoleID, &guild.FeatCategoryRoleID, &guild.FixCategoryRoleID, &guild.ChoreCategoryRoleID}
	for i, role := range data.CategoryRoles {
		registeredRole, err := registerRole(s, event.Guild.ID, &role, db.RoleKindCategory)
		if err != nil {
			return err
		}
		registeredCategoryRoles = append(registeredCategoryRoles, registeredRole)
		*rolePtrs[i] = registeredRole.ID
	}

	// register all priorities
	registeredPriorityRoles := []db.Role{}
	rolePtrs = []*string{&guild.LowPriorityRoleID, &guild.NormalPriorityRoleID, &guild.ImportantPriorityRoleID, &guild.CriticalPriorityRoleID}
	for i, role := range data.PriorityRoles {
		registeredRole, err := registerRole(s, event.Guild.ID, &role, db.RoleKindPriority)
		if err != nil {
			return err
		}
		registeredPriorityRoles = append(registeredPriorityRoles, registeredRole)
		*rolePtrs[i] = registeredRole.ID
	}

	registeredNobodyRole, err := registerRole(s, event.Guild.ID, &data.NobodyRole, db.RoleKindNobody)
	if err != nil {
		return err
	}
	guild.NobodyRoleID = registeredNobodyRole.ID

	registeredDiscussionRole, err := registerRole(s, event.Guild.ID, &data.DiscussionRole, db.RoleKindDiscussion)
	if err != nil {
		return err
	}
	guild.DiscussionRoleID = registeredDiscussionRole.ID

	return nil
}
