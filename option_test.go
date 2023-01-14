package discordauth

import (
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

var (
	spaceReplacer = strings.NewReplacer(" ", "", "\r", "", "\n", "", "\t", "")
	WithoutSpaces = cmp.Transformer("SpacesIgnored", func(in string) string {
		return spaceReplacer.Replace(in)
	})
)

func TestParsingGlobalOptions(t *testing.T) {
	want := `{
		"clientID":"1000000000000005",
		"clientSecret":"7SEWAAAA1AP_k",
		"redirectURL":"http://localhost:8080/discord/callback",
		"realms":[
			{
				"Ref":"really_cool_area",
				"Identifiers": [
					{"Resource":"role","Identifier":"10681122442122222","GuildID":"106010111111"},
					{"Resource":"user","Identifier":"30681122442122222","GuildID":""}
				]
			}
		],
		"inFlightState": {}
	}`

	d := caddyfile.NewTestDispenser(`{
		discordauth {
			client_id 1000000000000005
			client_secret 7SEWAAAA1AP_k
			redirect http://localhost:8080/discord/callback
	
			realm really_cool_area {
				guild 106010111111 {
					role 10681122442122222
				}
	
				user 30681122442122222
			}
		}
	}`)

	app, err := parseCaddyfileGlobalOption(d, nil)
	if err != nil {
		t.Fail()
	}

	if diff := cmp.Diff(want, string(app.(httpcaddyfile.App).Value), WithoutSpaces); diff != "" {
		t.Error(diff)
	}
}
