package discord

import "errors"

var (
	ErrDiscordGeneric = errors.New("discord response indicated failed")
	ErrUnknownGuild   = errors.New("user does not exist within guild")

	//ErrBadRequest        = errors.New("failed")
	ErrRateLimited       = errors.New("request has failed due to exceeding Discord rate limits")
	ErrInsufficientScope = errors.New("request failed due to lack of permission")
	ErrTokenExpired      = errors.New("discord authorization failed")
)

type ErrJSONError struct {
	message string
	code    APIErrorCode
}

func (e *ErrJSONError) Error() string {
	return e.message
}

func NewErrorDiscordResponse(message string, code APIErrorCode) *ErrJSONError {
	return &ErrJSONError{
		message: message,
		code:    code,
	}
}

func resolveError(discordErrorCode APIErrorCode) error {
	switch discordErrorCode {
	case ErrorCodeUnknownGuild:
		return ErrUnknownGuild
	}

	return ErrDiscordGeneric
}
