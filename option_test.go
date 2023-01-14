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
	testcases := []struct {
		name      string
		dispenser *caddyfile.Dispenser
		want      string
	}{
		{
			name: "all discord users",
			dispenser: caddyfile.NewTestDispenser(`{
				discordauth {
					client_id 1000000000000005
					client_secret 7SEWAAAA1AP_k
					redirect http://localhost:8080/discord/callback
			
					realm really_cool_area {
						*
					}
				}
			}`),
			want: `{
				"clientID":"1000000000000005",
				"clientSecret":"7SEWAAAA1AP_k",
				"redirectURL":"http://localhost:8080/discord/callback",
				"realms":[
					{
						"Ref":"really_cool_area",
						"Identifiers": [
							{"Resource":3,"Wildcard":true}
						]
					}
				],
				"inFlightState": null
			}`,
		},
		{
			name: "guild members only",
			dispenser: caddyfile.NewTestDispenser(`{
				discordauth {
					client_id 1000000000000005
					client_secret 7SEWAAAA1AP_k
					redirect http://localhost:8080/discord/callback
			
					realm nice_guild {
						guild 12354 {
							*
						}
					}
				}
			}`),
			want: `{
				"clientID":"1000000000000005",
				"clientSecret":"7SEWAAAA1AP_k",
				"redirectURL":"http://localhost:8080/discord/callback",
				"realms":[
					{
						"Ref":"nice_guild",
						"Identifiers": [
							{"Resource":1,"GuildID":"12354","Wildcard":true}
						]
					}
				],
				"inFlightState": null
			}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			app, err := parseCaddyfileGlobalOption(tc.dispenser, nil)
			if err != nil {
				t.Fail()
			}

			got := string(app.(httpcaddyfile.App).Value)
			if diff := cmp.Diff(tc.want, got, WithoutSpaces); diff != "" {
				t.Error(diff)
			}
		})
	}

}
