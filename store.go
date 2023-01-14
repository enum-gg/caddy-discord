package discordauth

import (
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"
)

var (
	ErrSessionKeyExists   = errors.New("session key already exists")
	ErrInvalidRedirectURI = errors.New("redirect URI is invalid")
)

type AuthInFlight struct {
	startedAt   time.Time
	redirectUri *url.URL
}

func (a *AuthInFlight) GetStartedAt() time.Time {
	return a.startedAt
}

func (a *AuthInFlight) GetRedirectURI() *url.URL {
	return a.redirectUri
}

type SessionStore struct {
	db map[string]AuthInFlight
	m  sync.Mutex
}

func (s *SessionStore) Start(sessionKey string, redirectURI string, startFrom time.Time) error {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.db[sessionKey]; ok {
		return ErrSessionKeyExists
	}

	redirect, err := url.ParseRequestURI(redirectURI)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidRedirectURI, err.Error())
	}

	s.db[sessionKey] = AuthInFlight{
		startedAt:   startFrom,
		redirectUri: redirect,
	}

	return nil
}

func (s *SessionStore) Complete(sessionKey string) (AuthInFlight, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sessh, ok := s.db[sessionKey]
	if !ok {
		return AuthInFlight{}, errors.New("session key does not exist")
	}

	delete(s.db, sessionKey)

	return sessh, nil
}

func NewSessionStore() *SessionStore {
	return &SessionStore{db: map[string]AuthInFlight{}}
}

func NewAuthInFlight(startedAt time.Time, redirectURL *url.URL) AuthInFlight {
	return AuthInFlight{
		startedAt:   startedAt,
		redirectUri: redirectURL,
	}
}
