package discordauth_test

import (
	"errors"
	discordauth "github.com/enum-gg/caddy-discord"
	"net/url"
	"testing"
	"time"
)

func TestStoreSessionKeyRequiresUniqueness(t *testing.T) {
	want := discordauth.ErrSessionKeyExists
	s := discordauth.NewSessionStore()

	_ = s.StartAuthFlow("1234", URLMustParse("http://localhost/"), time.Now(), "realm1")
	err := s.StartAuthFlow("1234", URLMustParse("http://localhost/"), time.Now(), "realm1")

	if !errors.Is(err, want) {
		t.Fail()
	}
}

func TestStoreSessionCompletion(t *testing.T) {
	wantedURI := URLMustParse("http://eggs.com")
	want := discordauth.NewAuthInFlight(time.Date(2023, time.May, 19, 1, 2, 3, 4, time.UTC), wantedURI)
	store := discordauth.NewSessionStore()

	if err := store.StartAuthFlow("5555", URLMustParse("http://eggs.com"), time.Date(2023, time.May, 19, 1, 2, 3, 4, time.UTC), "realm1"); err != nil {
		t.Fail()
	}

	got, err := store.CompleteAuthFlow("5555")
	if err != nil {
		t.Fail()
	}

	if got.GetRedirectURI().String() != want.GetRedirectURI().String() {
		t.Errorf("redirect URI was unexpected. wanted: %s, got %s", want.GetRedirectURI(), got.GetRedirectURI())
	}

	if got.GetStartedAt().String() != want.GetStartedAt().String() {
		t.Errorf("started at timestamp was different. wanted: %s, got %s", want.GetStartedAt(), got.GetStartedAt())
	}
}

func URLMustParse(rawURL string) *url.URL {
	wantedURI, err := url.Parse("http://eggs.com")
	if err != nil {
		panic(err)
	}

	return wantedURI
}
