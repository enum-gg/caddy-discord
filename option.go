package discordauth

import (
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
)

// parseCaddyfileGlobalOption implements caddyfile.Unmarshaler.
func parseCaddyfileGlobalOption(d *caddyfile.Dispenser, _ any) (any, error) {
	dpApp := new(DiscordPortalApp)

	for d.Next() {
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "client_id":
				if d.NextArg() {
					dpApp.ClientID = d.Val()
				}
				if d.NextArg() {
					return nil, d.ArgErr()
				}
			case "redirect":
				if d.NextArg() {
					dpApp.RedirectURL = d.Val()
				}
				if d.NextArg() {
					return nil, d.ArgErr()
				}
			case "client_secret":
				if d.NextArg() {
					dpApp.ClientSecret = d.Val()
				}
				if d.NextArg() {
					return nil, d.ArgErr()
				}
			case "realm":
				ag := &Realm{
					Identifiers: []*AccessIdentifier{},
				}

				dpApp.Realms = append(dpApp.Realms, ag)

				if d.NextArg() {
					ag.Ref = d.Val()
				}

				for nesting := d.Nesting(); d.NextBlock(nesting); {
					switch d.Val() {
					case "guild":
						if !d.NextArg() {
							return nil, d.Errf("unrecognized subdirective '%s'", d.Val())
						}

						guildID := d.Val()

						for nesting := d.Nesting(); d.NextBlock(nesting); {
							switch d.Val() {

							case "role":
								if d.NextArg() {
									ag.Identifiers = append(ag.Identifiers, &AccessIdentifier{
										Resource:   "role",
										Identifier: d.Val(),
										GuildID:    guildID,
									})
								}
								if d.NextArg() {
									return nil, d.ArgErr()
								}
							case "user":
								if d.NextArg() {
									ag.Identifiers = append(ag.Identifiers, &AccessIdentifier{
										Resource:   "user",
										Identifier: d.Val(),
										GuildID:    guildID,
									})
								}
								if d.NextArg() {
									return nil, d.ArgErr()
								}

							case "*":
								ag.Identifiers = append(ag.Identifiers, &AccessIdentifier{
									Resource:   "*",
									Identifier: "",
									GuildID:    guildID,
								})

								if d.NextArg() {
									return nil, d.ArgErr()
								}
								break
							default:
								return nil, d.Errf("unrecognized subdirective '%s'", d.Val())

							}
						}

					case "user":
						if d.NextArg() {
							ag.Identifiers = append(ag.Identifiers, &AccessIdentifier{
								Resource:   "user",
								Identifier: d.Val(),
							})
						}
						if d.NextArg() {
							return nil, d.ArgErr()
						}

						break

					default:
						return nil, d.Errf("unrecognized subdirective '%s'", d.Val())
					}
				}

			default:
				return nil, d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}

	return httpcaddyfile.App{
		Name:  "discordauth",
		Value: caddyconfig.JSON(dpApp, nil),
	}, nil
}
