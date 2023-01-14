package discordauth

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

func init() {
	caddy.RegisterModule(ProtectorPlugin{})
	httpcaddyfile.RegisterHandlerDirective("protect", parseCaddyfileHandlerDirective2)
}

func parseCaddyfileHandlerDirective2(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var s ProtectorPlugin
	s.UnmarshalCaddyfile(h.Dispenser)
	return &s, s.UnmarshalCaddyfile(h.Dispenser)
}

type ProtectCfg struct {
	ClientID     string
	ClientSecret string
}

type ProtectorPlugin struct {
	OAuthConfig  *oauth2.Config
	SessionStore *SessionStore
	Realm        string
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (e ProtectorPlugin) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	sID := SessionIDGenerator(32)

	existingSession, _ := r.Cookie(cookieName)
	if existingSession != nil {
		// Check if real...
	}

	if err := e.SessionStore.StartAuthFlow(sID, r.URL, time.Now(), e.Realm); err != nil {
		// TODO: Configurable redirect to error/borkage page
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return nil
	}

	// TODO: Configurable entropy for state
	url := e.OAuthConfig.AuthCodeURL(sID, oauth2.ApprovalForce)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return nil
}

func (p *ProtectorPlugin) Provision(ctx caddy.Context) error {
	ctxApp, _ := ctx.App(moduleName)
	app := ctxApp.(*DiscordPortalApp)
	p.OAuthConfig = app.getOAuthConfig()
	p.SessionStore = app.InFlightState

	return nil
}

func (ProtectorPlugin) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.protect",
		New: func() caddy.Module { return new(ProtectorPlugin) },
	}
}

func (p *ProtectorPlugin) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			// allow "with" or "using"
			if d.Val() != "with" && d.Val() != "using" {
				return d.ArgErr()
			}

			if !d.NextArg() {
				return d.ArgErr()
			}

			p.Realm = d.Val()
		}
	}

	return nil
}
