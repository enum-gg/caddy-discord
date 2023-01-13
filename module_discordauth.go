package discordauth

import (
	"context"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
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
	fmt.Println(app)

	s.OAuth = app.getOAuthConfig()

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
		s.Configuration = append(s.Configuration, d.Val())
		if d.NextArg() {
			if d.Val() == "callback" {
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
	//d.w.Write([]byte(r.RemoteAddr))
	ctx := context.Background()

	tok, err := d.OAuth.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		// TODO: ERROR
	}

	client := d.OAuth.Client(ctx, tok)
	res, _ := client.Get("https://discord.com/api/users/@me")

	fmt.Println(res)

	return next.ServeHTTP(w, r)
}
