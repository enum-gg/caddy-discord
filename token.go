package caddydiscord

import (
	"errors"
	"fmt"
	"github.com/enum-gg/caddy-discord/internal/discord"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type (
	TokenSignerSignature       = func(token *jwt.Token) (string, error)
	AuthedTokenParserSignature = func(signedToken string) (*AuthenticatedClaims, error)
	FlowTokenParserSignature   = func(signedToken string) (*FlowTokenParser, error)
)

type JWTManager struct {
	key []byte
}

type AuthenticatedClaims struct {
	Realm string `json:"realm,omitempty"`

	Username   string `json:"user,omitempty"`
	Avatar     string `json:"avatar,omitempty"`
	Authorised bool   `json:"authorised,omitempty"`
	jwt.RegisteredClaims
}

func (a AuthenticatedClaims) GetAudience() string {
	return "auth"
}

type FlowTokenParser struct {
	Realm       string `json:"realm,omitempty"`
	RedirectURI string `json:"redirectURI,omitempty"`
	jwt.RegisteredClaims
}

func (f FlowTokenParser) GetAudience() string {
	return "flow"
}

func NewAuthenticatedToken(identity discord.User, realm string, exp time.Time, authorised bool) *jwt.Token {
	claims := AuthenticatedClaims{
		realm,
		identity.Username,
		identity.Avatar,
		authorised,
		jwt.RegisteredClaims{
			Audience:  []string{"auth"},
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   identity.ID,
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
}

func NewAuthFlowToken(redirectURI string, realm string, exp time.Time) *jwt.Token {
	claims := FlowTokenParser{
		realm,
		redirectURI,
		jwt.RegisteredClaims{
			Audience:  []string{"flow"},
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
}

func NewTokenSigner(key []byte) TokenSignerSignature {
	return func(token *jwt.Token) (string, error) {
		return token.SignedString(key)
	}
}

type CustomClaim interface {
	jwt.Claims
}

func NewAuthedTokenParser(key []byte) func(signedToken string) (*AuthenticatedClaims, error) {
	return func(signedToken string) (*AuthenticatedClaims, error) {
		token, err := jwt.ParseWithClaims(signedToken, &AuthenticatedClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return key, nil
		})

		if err != nil {
			return nil, err
		}

		if claims, ok := token.Claims.(*AuthenticatedClaims); ok && token.Valid {
			//if customClaims.GetAudience() != claims.GetAudience() {
			//	return nil, errors.New("failed to authenticate JWT audience")
			//}

			return claims, nil
		}

		return nil, errors.New("unknown error")
	}
}

func NewFlowTokenParser(key []byte) func(signedToken string) (*FlowTokenParser, error) {
	return func(signedToken string) (*FlowTokenParser, error) {
		token, err := jwt.ParseWithClaims(signedToken, &FlowTokenParser{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return key, nil
		})

		if err != nil {
			return nil, err
		}

		// TODO: Check audience
		if claims, ok := token.Claims.(*FlowTokenParser); ok && token.Valid {
			//if customClaims.GetAudience() != claims.GetAudience() {
			//	return nil, errors.New("failed to authenticate JWT audience")
			//}

			return claims, nil
		}

		return nil, errors.New("unknown error")
	}
}
