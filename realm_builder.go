package caddydiscord

type RealmBuilder struct {
	realm *Realm
}

func (r RealmBuilder) AllowAllGuildMembers(guildID string) {
	r.realm.Identifiers = append(r.realm.Identifiers, &AccessIdentifier{
		Resource: DiscordGuildRule,
		GuildID:  guildID,
		Wildcard: true,
	})
}

func (r RealmBuilder) AllowGuildMember(guildID string, userID string) {
	r.realm.Identifiers = append(r.realm.Identifiers, &AccessIdentifier{
		Resource:   DiscordMemberRule,
		GuildID:    guildID,
		Identifier: userID,
	})
}

func (r RealmBuilder) AllowDiscordUser(userID string) {
	r.realm.Identifiers = append(r.realm.Identifiers, &AccessIdentifier{
		Resource:   DiscordUserRule,
		Identifier: userID,
	})
}

func (r RealmBuilder) AllowAllDiscordUsers() {
	r.realm.Identifiers = append(r.realm.Identifiers, &AccessIdentifier{
		Resource: DiscordUserRule,
		Wildcard: true,
	})
}

func (r RealmBuilder) AllowGuildRole(guildID string, roleID string) {
	r.realm.Identifiers = append(r.realm.Identifiers, &AccessIdentifier{
		Resource:   DiscordRoleRule,
		Identifier: roleID,
		GuildID:    guildID,
	})
}

func (r RealmBuilder) Build() *Realm {
	builtRealm := r.realm
	r.realm = &Realm{}

	return builtRealm
}

func (r RealmBuilder) Name(name string) {
	r.realm.Ref = name
}

func NewRealmBuilder() RealmBuilder {
	return RealmBuilder{
		realm: &Realm{},
	}
}
