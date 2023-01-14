package discordauth_test

import (
	"errors"
	discordauth "github.com/dev-this/caddy-discordauth"
	"net/url"
	"testing"
	"time"
)

func TestStoreSessionKeyRequiresUniqueness(t *testing.T) {
	want := discordauth.ErrSessionKeyExists
	s := discordauth.NewSessionStore()

	_ = s.Start("1234", "http://localhost/", time.Now())
	err := s.Start("1234", "http://localhost/", time.Now())

	if !errors.Is(err, want) {
		t.Fail()
	}
}

func TestStoreSessionChecksURIValidity(t *testing.T) {
	want := discordauth.ErrInvalidRedirectURI
	s := discordauth.NewSessionStore()

	err := s.Start("1234", "notarealurl", time.Now())

	if !errors.Is(err, want) {
		t.Error(err)
		t.Fail()
	}
}

func TestStoreSessionCompletion(t *testing.T) {
	wantedURI, _ := url.Parse("http://eggs.com")
	want := discordauth.NewAuthInFlight(time.Date(2023, time.May, 19, 1, 2, 3, 4, time.UTC), wantedURI)
	store := discordauth.NewSessionStore()

	if err := store.Start("5555", "http://eggs.com", time.Date(2023, time.May, 19, 1, 2, 3, 4, time.UTC)); err != nil {
		t.Fail()
	}

	got, err := store.Complete("5555")
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
