package caddydiscord

import (
	"encoding/hex"
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

// DiscordPortalApp allows you to authenticate caddy routes based
// on a Discord User Identity.
//
// e.g. Accessing /really-cool-people requires user to have {Role}
// within {Guild}
//
// Discord's OAuth flow is used for identity using your
// own Discord developer application.
//
// See an example Caddyfile https://github.com/enum-gg/caddy-discord#caddyfile-example
type DiscordPortalApp struct {
	// ClientID is the "Client ID" from your Discord Application OAuth information
	ClientID string `json:"clientID"`

	// ClientSecret is the "Client Secret" from your Discord Application
	// OAuth information.
	//
	// Treat this is sensitive. Do not share or expose it to anyone.
	ClientSecret string `json:"clientSecret"`

	// RedirectURL is the destination for clients for the OAuth flow
	// Your Discord Application's OAuth "Redirects" needs to be aware
	// of this endpoint.
	//
	// Within your Caddyfile this URL should be configured with "discord callback".
	RedirectURL string `json:"redirectURL"`

	// Realms group together explicit rules about whom to authorise.
	Realms      RealmRegistry `json:"realms"`
	oauthConfig *oauth2.Config

	// Key is the signing key used for the JWT stored as the client's cookie
	// it is ephemeral alongside the caddy server process.
	Key string `json:"key,omitempty"`

	// ExecutionKey is an ephemeral identifier for the client's cookie which contains
	// the JWT payload proving Discord identity. It is the 'public' version of the signing Key.
	//
	// End users will have to perform the OAuth flow once uniquely per ExecutionKey,
	//  which will be a touchless experience barely noticeably from their end.
	//
	// ExecutionKey exists as an alternative to the server operator providing their own
	// JWT signing key; which should eventually become an optional configuration.
	ExecutionKey string `json:"signature,omitempty"`
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
	d.Key = hex.EncodeToString(randomness(64))
	d.ExecutionKey = hashString512(d.Key)

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
		return fmt.Errorf("redirect URL is not configured")
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
