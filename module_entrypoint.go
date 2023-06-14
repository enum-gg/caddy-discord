package caddydiscord

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/caddyauth"
	"golang.org/x/oauth2"
)

var (
	_ caddyfile.Unmarshaler   = (*ProtectorPlugin)(nil)
	_ caddy.Validator         = (*ProtectorPlugin)(nil)
	_ caddyauth.Authenticator = (*ProtectorPlugin)(nil)
)

func init() {
	caddy.RegisterModule(ProtectorPlugin{})
	httpcaddyfile.RegisterHandlerDirective("protect", parseCaddyfileHandlerDirective2)
}

func parseCaddyfileHandlerDirective2(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var s ProtectorPlugin
	s.UnmarshalCaddyfile(h.Dispenser)
	return caddyauth.Authentication{
		ProvidersRaw: caddy.ModuleMap{
			"discord": caddyconfig.JSON(s, nil),
		},
	}, nil
}

type ProtectCfg struct {
	ClientID     string
	ClientSecret string
}

type ProtectorPlugin struct {
	OAuthConfig       *oauth2.Config
	tokenSigner       TokenSignerSignature
	authedTokenParser AuthedTokenParserSignature
	flowTokenParser   FlowTokenParserSignature
	CookieName        string
	Realm             string
}

// Authenticate implements caddyhttp.MiddlewareHandler.
func (e ProtectorPlugin) Authenticate(w http.ResponseWriter, r *http.Request) (caddyauth.User, bool, error) {
	existingSession, _ := r.Cookie(fmt.Sprintf("%s_%s", defaultCookieName, e.Realm))

	// Handle passing through signed token over to support multiple domains.
	if existingSession == nil && r.URL.Query().Has("DISCO_PASSTHROUGH") && r.URL.Query().Has("DISCO_REALM") {
		q := r.URL.Query()
		signedToken := q.Get("DISCO_PASSTHROUGH")
		realm := q.Get("DISCO_REALM")
		q.Del("DISCO_PASSTHROUGH")
		q.Del("DISCO_REALM")
		r.URL.RawQuery = q.Encode()

		cookie := &http.Cookie{
			Name:     fmt.Sprintf("%s_%s", defaultCookieName, realm),
			Value:    signedToken,
			Expires:  time.Now().Add(time.Hour * 16),
			HttpOnly: true,
			// Strict mode breaks functionality - due to discord referrer.
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
			// Secure // TODO: Configurable
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, r.URL.String(), http.StatusFound)
		return caddyauth.User{}, false, nil
	}

	if existingSession != nil {
		claims, err := e.authedTokenParser(existingSession.Value)
		if err != nil {
			return caddyauth.User{}, false, err
		}

		return caddyauth.User{
			ID: claims.Subject,
			Metadata: map[string]string{
				"username": claims.Username,
				"avatar":   claims.Avatar,
			},
		}, true, nil
	}

	// 15 minutes to make it through Discord consent.
	exp := time.Now().Add(time.Minute * 15)
	backToURL := *r.URL
	if !backToURL.IsAbs() {
		backToURL.Scheme = "http"
		if r.TLS != nil {
			backToURL.Scheme = "https"
		}

		backToURL.Host = r.Host
	}
	token := NewAuthFlowToken(backToURL.String(), e.Realm, exp)
	signedToken, err := e.tokenSigner(token)
	if err != nil {
		// Unable to generate JWT
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return caddyauth.User{}, false, err
	}

	url := e.OAuthConfig.AuthCodeURL(signedToken, oauth2.SetAuthURLParam("prompt", "none"))

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return caddyauth.User{}, false, nil
}

func (p *ProtectorPlugin) Provision(ctx caddy.Context) error {
	ctxApp, _ := ctx.App(moduleName)
	app := ctxApp.(*DiscordPortalApp)
	p.OAuthConfig = app.getOAuthConfig()

	p.CookieName = app.CookieName
	if p.CookieName == "" {
		// Use default cookie name if none provided
		p.CookieName = defaultCookieName
	}

	key, err := hex.DecodeString(app.Key)
	if err != nil {
		return err
	}

	p.tokenSigner = NewTokenSigner(key)
	p.authedTokenParser = NewAuthedTokenParser(key)
	p.flowTokenParser = NewFlowTokenParser(key)

	return nil
}

func (ProtectorPlugin) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.authentication.providers.discord",
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

func (p *ProtectorPlugin) Validate() error {
	return nil
}
