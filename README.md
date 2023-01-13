# Caddy Module - Discord Auth
This package contains a module to bridge Discord and Caddy together for the purpose of authentication.

Authenticate routes based on 'realms' which are a collection of your rules, corresponding with a Discord identity, both within a Guild context or globally.

| Resource    | Description                                                 | Example                                                                                                                                                 |
|-------------|-------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------|
| User ID     | Discord User IDs (_no guild context required_)       | <pre>realm godmode {<br />  user_id 314009111187026172 <br />}                                                                                          |
| Guild       | Any user that exists  _within the guild_                    | <pre>realm cool_guild_users {<br />  guild 1063070451111289907 {<br />    * # Allows all users<br />  }<br />}                                          |
| Role        | Users that assigned a specific role _within a guild_        | <pre>realm cool_role {<br />  guild 1063070451111289907 {<br />    role_id 106301111332755034<br />    role_id 106301111332755034<br />  }<br />}</pre> |
| Channel IDs | Users that have permission to view channel _within a guild_ | <pre>realm channel_whitelist {<br />  guild 1063070451111289907 {<br />    channel_id 1060211512<br />    channel_id 1060761112<br />  }<br />}</pre>  |


Loosely inspired from [caddy-security's Discord OAuth2 module](https://authp.github.io/docs/authenticate/oauth/backend-oauth2-0013-discord), with a much stronger focus on coupling Discord and Caddy for authentication purposes.

<div align="center">
	<br />
	<p>
		<a href="https://discord.gg/k9tVAwws8U"><img src="https://img.shields.io/discord/1063070457047289907?color=5865F2&logo=discord&logoColor=white" alt="Discord server" /></a>
	</p>
</div>

# Install


## Components
There are two components required for Caddy Discord-Auth

### Discord App
An app is required for the module to be able to determine Guild User information.

OAuth scopes required `identify`

### Caddy Modules
```
discordauth
http.handler.discordauth
```

## Usage
```caddyfile
{
    discordauth {
        client_id 1000000000000000000 # Discord app OAuth client ID 
        client_secret 8CEPZZZZZAfl_w19ZZZZW_k # Discord app OAuth secret
        redirect http://localhost:8080/discord/callback # Route you've configured with `discordauth callback`

        realm really_cool_area {
            guild 106307051119907 {
                role_id 10630111112755034
                user_id 31400111187026172
                channel_id 106311116556734
            }
        }
    }
}

http://localhost:8091 {
    route /discord/callback {
        discordauth callback # Desigate route as OAuth callback endpoint
    }

    route /hello {
        protect using really_cool_area # Only allow discord users that auth against 'really_cool_area' realm 
        
        respond "Only really cool people can see this!"
    }

    respond "Hello, worBASE ld!"
}

```

## Building
```
xcaddy build --with github.com/dev-this/caddy-discordauthauth=./
```