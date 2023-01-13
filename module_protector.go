package discordauth

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"golang.org/x/oauth2"
	"net/http"
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
	OAuthConfig *oauth2.Config
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (e *ProtectorPlugin) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	//if NOT AUTHED...
	// TODO: Proper state checking
	url := e.OAuthConfig.AuthCodeURL("state", oauth2.ApprovalForce)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return nil
	//return next.ServeHTTP(w, r)
}

func (e *ProtectorPlugin) Provision(ctx caddy.Context) error {
	ctxApp, _ := ctx.App(moduleName)
	app := ctxApp.(*DiscordPortalApp)
	fmt.Println(app)
	e.OAuthConfig = app.getOAuthConfig()

	return nil
}

func (ProtectorPlugin) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.protect",
		New: func() caddy.Module { return new(ProtectorPlugin) },
	}
}

func (ProtectorPlugin) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			if d.Val() != "with" {
				return d.ArgErr()
			}

			if !d.NextArg() {
				return d.ArgErr()
			}

			accessGroup := d.Val()
			fmt.Println(accessGroup)
		}
	}

	return nil
}
