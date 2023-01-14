package caddydiscord

const (
	Unknown          DiscordResource = 0
	DiscordGuildRule DiscordResource = 1
	DiscordRoleRule  DiscordResource = 2
	// DiscordMemberRule represents a specific Discord User within a specific guild
	DiscordMemberRule DiscordResource = 3
	// DiscordUserRule represents a Discord User Snowflake ID
	DiscordUserRule DiscordResource = 4
)

type DiscordResource int

func ResourceRequiresGuild(resource DiscordResource) bool {
	switch resource {
	case DiscordGuildRule, DiscordRoleRule, DiscordMemberRule:
		return true
	}

	return false
}

type Realm struct {
	Ref         string              `json:"Ref"`
	Identifiers []*AccessIdentifier `json:"Identifiers"`
}

func (r Realm) GetAllGuilds() []string {
	// Use map to avoid doubling up values
	guildMap := make(map[string]bool)

	for _, resource := range r.Identifiers {
		switch resource.Resource {
		case DiscordGuildRule, DiscordRoleRule, DiscordMemberRule:
			guildMap[resource.GuildID] = true
		}
	}

	guildIDs := make([]string, 0, len(guildMap))

	for guildID, _ := range guildMap {
		guildIDs = append(guildIDs, guildID)
	}

	return guildIDs
}

type AccessIdentifier struct {
	Resource   DiscordResource `json:"Resource"` // role, user
	Identifier string          `json:"Identifier,omitempty"`
	GuildID    string          `json:"GuildID,omitempty"`
	Wildcard   bool            `json:"Wildcard,omitempty"`
}
