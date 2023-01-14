package discordauth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"golang.org/x/oauth2"
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
	realm       string
}

type KnownSession struct {
	token  *oauth2.Token
	realms CheckedRealmMap
}

type CheckedRealmMap map[string]bool

func (a *AuthInFlight) GetStartedAt() time.Time {
	return a.startedAt
}

func (a *AuthInFlight) GetRedirectURI() *url.URL {
	return a.redirectUri
}

// TODO: Cleanup process to avoid memory leaks of stale sessions, or introduce storage compatibility.
type SessionStore struct {
	db       map[string]AuthInFlight
	known    map[string]KnownSession
	authTex  sync.Mutex
	sesshTex sync.Mutex
}

func (s *SessionStore) StartAuthFlow(sessionKey string, redirectURI *url.URL, startFrom time.Time, realm string) error {
	s.authTex.Lock()
	defer s.authTex.Unlock()

	if _, ok := s.db[sessionKey]; ok {
		return ErrSessionKeyExists
	}

	s.db[sessionKey] = AuthInFlight{
		startedAt:   startFrom,
		redirectUri: redirectURI,
		realm:       realm,
	}

	return nil
}

func (s *SessionStore) CompleteAuthFlow(sessionKey string) (AuthInFlight, error) {
	s.authTex.Lock()
	defer s.authTex.Unlock()

	sessh, ok := s.db[sessionKey]
	if !ok {
		return AuthInFlight{}, errors.New("session key does not exist")
	}

	delete(s.db, sessionKey)

	return sessh, nil
}

func (s *SessionStore) AddKnown(sessionID string, token *oauth2.Token, realm string) error {
	s.sesshTex.Lock()
	defer s.sesshTex.Unlock()

	if _, ok := s.known[sessionID]; ok {
		return errors.New("session ID is already known")
	}

	s.known[sessionID] = KnownSession{
		token: token,
		realms: CheckedRealmMap{
			realm: true,
		},
	}

	return nil
}

func (s *SessionStore) GetKnown(sessionID string) *KnownSession {
	if known, ok := s.known[sessionID]; ok {
		return &known
	}

	return nil
}

// TODO: UpdateWithCheckedRealm

func NewSessionStore() *SessionStore {
	return &SessionStore{db: map[string]AuthInFlight{}, known: map[string]KnownSession{}}
}

func NewAuthInFlight(startedAt time.Time, redirectURL *url.URL) AuthInFlight {
	return AuthInFlight{
		startedAt:   startedAt,
		redirectUri: redirectURL,
	}
}

func randomness(length uint) []byte {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return randomBytes
}

func SessionIDGenerator(length uint) string {
	return base64.RawStdEncoding.EncodeToString(randomness(length))
}
