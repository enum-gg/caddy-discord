package discordauth

import (
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"testing"
)

func TestParsingGlobalOptions(t *testing.T) {
	d := caddyfile.NewTestDispenser(`{
    discordauth {
        client_id 1000000000000005
        client_secret 7SEWyAANTo1AP_k
        redirect http://localhost:8080/discord/callback

        realm really_cool_area {
            guild 106010111111 {
                role_id 10681122442122222
                user_id 30681122442122222
                channel_id 7106811224421222
            }
        }
    }
}`)

	_, _ = parseCaddyfileGlobalOption(d, nil)
}
