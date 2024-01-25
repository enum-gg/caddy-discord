package caddydiscord

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"golang.org/x/oauth2"
)

var (
	_ caddy.App         = (*DiscordPortalApp)(nil)
	_ caddy.Module      = (*DiscordPortalApp)(nil)
	_ caddy.Provisioner = (*DiscordPortalApp)(nil)
)

const (
	moduleName = "discord"
	cookieName = "_DISCORDCADDY"
)

func init() {
	caddy.RegisterModule(DiscordPortalApp{})
	httpcaddyfile.RegisterGlobalOption("discord", parseCaddyfileGlobalOption)
}

type DiscordPortalApp struct {
	ClientID     string        `json:"clientID"`
	ClientSecret string        `json:"clientSecret"`
	RedirectURL  string        `json:"redirectURL"`
	Realms       RealmRegistry `json:"realms"`
	oauthConfig  *oauth2.Config
	Key          string `json:"key,omitempty"`
	Signature    string `json:"signature,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (DiscordPortalApp) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  moduleName,
		New: func() caddy.Module { return new(DiscordPortalApp) },
	}
}

func (d *DiscordPortalApp) Provision(_ caddy.Context) error {
	// Discord App ID is used as entropy for JWT signing keys.
	d.Key = hashString512(d.ClientID)

	// TODO: Signature will be used for cookie integrity checks, to ensure checks are inline with most recent Caddyfile.
	// TODOTODO: Use parsed caddyfile signature for checks, instead of just Discord App Client ID.
	d.Signature = hashString256(d.ClientID, 16)
	return nil
}

func (d DiscordPortalApp) Start() error {
	return nil
}

// Stop stops the App.
func (d DiscordPortalApp) Stop() error {
	return nil
}

func (d DiscordPortalApp) Validate() error {
	if d.ClientID == "" {
		return fmt.Errorf("client ID is missing")
	}

	if d.ClientSecret == "" {
		return fmt.Errorf("discord OAuth client secret has not been set")
	}

	if d.RedirectURL == "" {
		return fmt.Errorf("redirect URL has not bee configured")
	}

	return nil
}

// getOAuthConfig singleton
func (d DiscordPortalApp) getOAuthConfig() *oauth2.Config {
	if d.oauthConfig == nil {
		d.oauthConfig = &oauth2.Config{
			ClientID:     d.ClientID,
			ClientSecret: d.ClientSecret,
			Scopes:       []string{"identify", "guilds.members.read"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://discord.com/oauth2/authorize",
				TokenURL: "https://discord.com/api/oauth2/token",
			},
			RedirectURL: d.RedirectURL,
		}
	}

	return d.oauthConfig
}
