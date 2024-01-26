<img src="./assets/logo.png" width="40%" height="40%" align="right">

# Caddy - Discord
tl;dr: Authenticate caddy routes based on a Discord User Identity.
_<br />e.g. Accessing `/really-cool-people` requires user to have `{Role}` within `{Guild}`_

This package contains a module allowing authorization in Caddy based on a Discord Identity, by using  Discords OAuth2 flow (authorization code grant).

---

`caddy-discord` is licensed under [_GNU Affero General Public License v3.0_](https://github.com/enum-gg/caddy-discord/blob/main/LICENSE.md)
<br><i>Logo created by [@AutonomousCat](https://github.com/AutonomousCat/)</i>
---

### Caddy Modules
```
caddydiscord
http.authentication.providers.discord
http.handler.discord
```

## Docker (Container)
```sh
docker run -p 8080:8080 \
  --rm -v $PWD/Caddyfile:/etc/caddy/Caddyfile \
  enumgg/caddy-discord:v1.0.1
```

## Discord Resources
**realm** allows you to name a label and group together specific targeted Discord Users by using the directives below.

| Resource        | Description                                                 | Example                                                                                                                                                                                                                          |
|-----------------|-------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| User ID         | Discord User IDs (_optionally with guild presence_)         | <pre>realm godmode {<br />  user 314009111187026172 # Allow user regardless of which guild they are in<br />  guild 1063070451111289907 {<br />    user 314009111187026199 # Allow user if they're part of guild<br />  }<br />} |
| Guild           | Any user that exists  _within the guild_                    | <pre>realm cool_guild_users {<br />  guild 1063070451111289907 {<br />    * # Allows all users <br />  }<br />}                                                                                                                  |
| Role            | Users that assigned a specific role _within a guild_        | <pre>realm cool_role {<br />  guild 1063070451111289907 {<br />    role 106301111332755034<br />    role 106301111332755034<br />  }<br />}</pre>                                                                                |

Loosely inspired from [caddy-security's Discord OAuth2 module](https://authp.github.io/docs/authenticate/oauth/backend-oauth2-0013-discord), with a much stronger focus on coupling Discord and Caddy for authentication purposes.

<div align="center">
	<br />
	<p>
		<a href="https://discord.gg/k9tVAwws8U"><img src="https://img.shields.io/discord/1063070457047289907?color=5865F2&logo=discord&logoColor=white" alt="Discord server" /></a>
	</p>
</div>

# Install

[**Download Latest Version**](https://github.com/enum-gg/caddy-discord/releases)

1. Download caddy + caddy-discord
    - Using released binaries
    - Build yourself using `xcaddy`
      - `xcaddy build --with github.com/enum-gg/caddy-discord`
2. Create Discord Application ([Discord Developer Portal](https://discord.com/developers/applications))
    - New Application
    - OAuth2
        - Obtain your Client ID & Client secret
        - Add Redirects [Docs](https://discord.com/developers/docs/topics/oauth2#authorization-code-grant-redirect-url-example)
3. Prepare your `Caddyfile`
    - Gather your Discord App OAuth2 Client ID & Client Secret,
    - Decide your route for caddy-discords to use as the OAuth2


## Caddyfile Example
```caddyfile
{
    discord {
        client_id 1000000000000000000 # Discord app OAuth client ID 
        client_secret 8CEPZZZZZAfl_w19ZZZZW_k # Discord app OAuth secret
        redirect http://localhost:8080/discord/callback # Route you've configured with `discordauth callback`

        realm clique {
            guild 106307051119907 {
                role 10630111112755034
            }
        }
        
        realm just_for_me {
            user 31400111187026172
        }
    }
}

http://localhost:8080 {
    route /discord/callback {
         # Desigate route as OAuth callback endpoint
         discord callback 
   }

    route /discordians-only {
         # Only allow discord users that auth against 'really_cool_area' realm 
         protect using clique 
        
         respond "Hello {http.auth.user.username}!<br /><br /><img src='https://cdn.discordapp.com/avatars/{http.auth.user.id}/{http.auth.user.avatar}?size=4096.png'> "
    }

    respond "Hello, world!"
}

```

## Building
```
xcaddy build --with github.com/enum-gg/caddy-discord=./
```
