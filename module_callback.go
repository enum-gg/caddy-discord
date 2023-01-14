package caddydiscord

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/enum-gg/caddy-discord/internal/discord"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

var (
	_ caddy.Provisioner           = (*DiscordAuthPlugin)(nil)
	_ caddyhttp.MiddlewareHandler = (*DiscordAuthPlugin)(nil)
	_ caddy.Validator             = (*DiscordAuthPlugin)(nil)
)

func init() {
	caddy.RegisterModule(DiscordAuthPlugin{})
	httpcaddyfile.RegisterHandlerDirective("discord", parseCaddyfileHandlerDirective)
}

func parseCaddyfileHandlerDirective(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var s DiscordAuthPlugin
	s.UnmarshalCaddyfile(h.Dispenser)
	return s, s.UnmarshalCaddyfile(h.Dispenser)
}

type DiscordAuthPlugin struct {
	Configuration   []string
	OAuth           *oauth2.Config
	Realms          *RealmRegistry
	Key             string
	tokenSigner     TokenSignerSignature
	flowTokenParser FlowTokenParserSignature
}

func (DiscordAuthPlugin) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.discord",
		New: func() caddy.Module { return new(DiscordAuthPlugin) },
	}
}

func (s *DiscordAuthPlugin) Provision(ctx caddy.Context) error {
	ctxApp, _ := ctx.App(moduleName)
	app := ctxApp.(*DiscordPortalApp)

	s.OAuth = app.getOAuthConfig()
	s.Realms = &app.Realms

	key, err := hex.DecodeString(app.Key)
	if err != nil {
		return err
	}

	s.tokenSigner = NewTokenSigner(key)
	s.flowTokenParser = NewFlowTokenParser(key)

	return nil
}

func (s *DiscordAuthPlugin) Validate() error {
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
func (d DiscordAuthPlugin) ServeHTTP(w http.ResponseWriter, r *http.Request, _ caddyhttp.Handler) error {
	ctx := context.Background()
	q := r.URL.Query()

	token, err := d.flowTokenParser(q.Get("state"))
	if err != nil {
		// Unable to find session. Using load balancers? Server was restarted?
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return err
	}

	realm := d.Realms.ByName(token.Realm)
	if realm == nil {
		// Unable to resolve realm
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return err
	}

	tok, err := d.OAuth.Exchange(ctx, q.Get("code"))
	if err != nil {
		return err
	}

	client := discord.NewClientWrapper(d.OAuth.Client(ctx, tok))

	allowed := false

	identity, err := client.FetchCurrentUser()
	if err != nil || len(identity.ID) == 0 {
		// Unable to resolve realm
		http.Error(w, "Failed to resolve Discord User", http.StatusInternalServerError)
		return err
	}

	for _, rule := range realm.Identifiers {
		if ResourceRequiresGuild(rule.Resource) {
			_, err := client.FetchGuildMembership(rule.GuildID)
			if err != nil {
				continue
				// TODO: check error type - probably not a member of guild...
			}

			// TODO assert guildMember has data
			allowed = true
		} else if rule.Resource == DiscordUserRule && rule.Wildcard == false && rule.Identifier == identity.ID {
			allowed = true
			break
		} else if rule.Resource == DiscordUserRule && rule.Wildcard == true {
			allowed = true
			break
		}
	}

	if !allowed {
		// User failed realm checks
		//http.Error(w, "You do not have access to this", http.StatusForbidden)
		http.Redirect(w, r, token.RedirectURI, http.StatusFound)

		return nil
	}
	// Re-validate user through OAuth2 flow every 16 hours
	expiration := time.Now().Add(time.Hour * 16)

	authedToken := NewAuthenticatedToken(*identity, realm.Ref, expiration)
	signedToken, err := d.tokenSigner(authedToken)
	if err != nil {
		// Unable to generate JWT
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return err
	}

	cookie := &http.Cookie{
		Name:     fmt.Sprintf("%s_%s", cookieName, realm.Ref),
		Value:    signedToken,
		Expires:  expiration,
		HttpOnly: true,
		// Strict mode breaks functionality - due to discord referrer.
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		//Secure // TODO: Configurable
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, token.RedirectURI, http.StatusFound)

	return nil
}
