package caddydiscord

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
				realmBuilder := NewRealmBuilder()

				if d.NextArg() {
					realmBuilder.Name(d.Val())
					//ag.Ref = d.Val()
				}

				for subNesting := d.Nesting(); d.NextBlock(subNesting); {
					switch d.Val() {
					case "guild":
						if !d.NextArg() {
							return nil, d.Errf("unrecognized subdirective '%s'", d.Val())
						}

						guildID := d.Val()

						for subSubNesting := d.Nesting(); d.NextBlock(subSubNesting); {
							switch d.Val() {

							case "role":
								if d.NextArg() {
									realmBuilder.AllowGuildRole(guildID, d.Val())
								}
								if d.NextArg() {
									return nil, d.ArgErr()
								}
							case "user":
								if d.NextArg() {
									realmBuilder.AllowGuildMember(guildID, d.Val())
								}
								if d.NextArg() {
									return nil, d.ArgErr()
								}

							case "*":
								realmBuilder.AllowAllGuildMembers(guildID)

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
							realmBuilder.AllowDiscordUser(d.Val())
						}
						if d.NextArg() {
							return nil, d.ArgErr()
						}

					case "*":
						// Anyone with a Discord Account...
						realmBuilder.AllowAllDiscordUsers()

						if d.NextArg() {
							return nil, d.ArgErr()
						}

						break

					default:
						return nil, d.Errf("unrecognized subdirective '%s'", d.Val())
					}
				}

				dpApp.Realms = append(dpApp.Realms, realmBuilder.Build())

			default:
				return nil, d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}

	return httpcaddyfile.App{
		Name:  "discord",
		Value: caddyconfig.JSON(dpApp, nil),
	}, nil
}
