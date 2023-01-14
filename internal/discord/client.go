package discord

import (
	"fmt"
	"net/http"
	"net/url"
)

type APIClient struct {
	client *http.Client
}

func (d *APIClient) getRequest(url string) (*http.Response, error) {
	return d.client.Get(url)
}

func NewClientWrapper(client *http.Client) *APIClient {
	return &APIClient{
		client: client,
	}
}

func (d *APIClient) FetchCurrentUser() (*User, error) {
	return fetch[User](d, "https://discord.com/api/users/@me")
}

func (d *APIClient) FetchGuildMembership(guildID string) (*GuildMemberResponse, error) {
	return fetch[GuildMemberResponse](d, fmt.Sprintf("https://discord.com/api/users/@me/guilds/%s/member", url.QueryEscape(guildID)))
}

func fetch[T any](client *APIClient, url string) (*T, error) {
	response, err := client.getRequest(url)
	if err != nil {
		return nil, err
	}

	if response.Body != nil {
		defer response.Body.Close()
	}

	if response.StatusCode == http.StatusOK {
		normalised, err := getBody[T](response)
		if err != nil {
			return nil, err
		}

		return normalised, nil
	}

	normalisedError, err := getBody[ErrorResponse](response)
	if err != nil {
		//failed to parse response body into error
		return nil, err
	}

	// Invalid requests
	// https://discord.com/developers/docs/topics/rate-limits#invalid-request-limit-aka-cloudflare-bans
	// https://discord.com/developers/docs/topics/opcodes-and-status-codes#http
	if response.StatusCode == http.StatusUnauthorized {
		// Token expired?
		// log .Message, .Code?
		return nil, ErrInsufficientScope
	}

	if response.StatusCode == http.StatusForbidden {
		// Scopes insufficient?
		// log .Message, .Code?
		return nil, ErrInsufficientScope
	}

	// Rate limited
	// https://discord.com/developers/docs/topics/rate-limits#rate-limits
	if response.StatusCode == http.StatusTooManyRequests {
		// TODO: http.client transport for retrying
		// log .Message, .Code?
		return nil, ErrRateLimited
	}

	return nil, resolveError(normalisedError.Code)
}
