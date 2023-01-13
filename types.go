package discordauth

type AccessIdentifier struct {
	Resource   string `json:"Resource"` // role, user, channel... `
	Identifier string `json:"Identifier"`
	GuildID    string `json:"GuildID"`
}

type Realm struct {
	Ref         string              `json:"Ref"`
	Identifiers []*AccessIdentifier `json:"Identifiers"`
}
