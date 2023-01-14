package discordauth

import (
	"context"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/enum-gg/caddy-discord/internal/discord"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

func init() {
	caddy.RegisterModule(DiscordAuthPlugin{})
	httpcaddyfile.RegisterHandlerDirective(moduleName, parseCaddyfileHandlerDirective)
}

func parseCaddyfileHandlerDirective(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var s DiscordAuthPlugin
	s.UnmarshalCaddyfile(h.Dispenser)
	return s, s.UnmarshalCaddyfile(h.Dispenser)
}

type DiscordAuthPlugin struct {
	Configuration []string
	OAuth         *oauth2.Config
	SessionStore  *SessionStore
	Realms        *RealmRegistry
}

func (DiscordAuthPlugin) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.discordauth",
		New: func() caddy.Module { return new(DiscordAuthPlugin) },
	}
}

func (s *DiscordAuthPlugin) Provision(ctx caddy.Context) error {
	ctxApp, _ := ctx.App(moduleName)
	app := ctxApp.(*DiscordPortalApp)

	s.OAuth = app.getOAuthConfig()
	s.SessionStore = app.InFlightState
	s.Realms = &app.Realms

	return nil
}

// UnmarshalCaddyfile will extract discordauth directives on a server-level
//
//	route /some/path/callback {
//	    discordauth callback
//	}
func (s *DiscordAuthPlugin) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	s.Configuration = []string{}

	for d.Next() {
		if d.NextArg() {
			if d.Val() == "callback" {
				s.Configuration = append(s.Configuration, d.Val())

				if d.NextArg() {
					return d.ArgErr()
				}
			}
		}
	}

	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (d DiscordAuthPlugin) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	ctx := context.Background()
	q := r.URL.Query()

	session, err := d.SessionStore.CompleteAuthFlow(q.Get("state"))
	if err != nil {
		// Unable to find session. Using load balancers? Server was restarted?
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return nil
	}

	realm := d.Realms.ByName(session.realm)
	if realm == nil {
		// Unable to resolve realm
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return nil
	}

	tok, err := d.OAuth.Exchange(ctx, q.Get("code"))
	if err != nil {
		log.Fatal(err)
	}

	client := discord.NewClientWrapper(d.OAuth.Client(ctx, tok))
	// REALM CHECKS HERE...

	allowed := false

	_, err = client.FetchCurrentUser()
	if err != nil {
		// Unable to resolve realm
		http.Error(w, "Failed to resolve Discord User", http.StatusInternalServerError)
		return nil
	}

	/*
		for _, rule := range realm.Identifiers {
			if ResourceRequiresGuild(rule.Resource) {
				_, err := client.FetchGuildMembership(rule.GuildID)
				if err != nil {
					// check TYPE - probably not a member of guild...
				}
			} else if rule.Resource == DiscordUserRule {
				if rule.Wildcard == true {
					allowed = true
					break
				}
			}
		}*/

	if allowed {

		cookie := &http.Cookie{
			Name:     cookieName,
			Value:    SessionIDGenerator(64),
			MaxAge:   0,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			//Secure // TODO: Configurable
		}

		http.SetCookie(w, cookie)

		return next.ServeHTTP(w, r)
	}

	// User failed realm checks
	http.Error(w, "You do not have access to this", http.StatusForbidden)
	return nil
}
