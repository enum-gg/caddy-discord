package discord

import "time"

type (
	// https://discord.com/developers/docs/topics/opcodes-and-status-codes#json
	APIErrorCode int

	// https://discord.com/developers/docs/resources/user#user-object-premium-types
	PremiumTypes int

	// https://discord.com/developers/docs/resources/user#user-object-user-flags
	UserFlags int
)

const (
	ErrorCodeGeneral      APIErrorCode = 0
	ErrorCodeUnknownGuild APIErrorCode = 10004

	PremiumTypeNone         PremiumTypes = 0
	PremiumTypeNitroClassic PremiumTypes = 1
	PremiumTypeNitro        PremiumTypes = 2
	PremiumTypeNitroBasic   PremiumTypes = 3

	UserFlagStaff                 UserFlags = 1 << 0
	UserFlagPartner               UserFlags = 1 << 1
	UserFlagHypesquad             UserFlags = 1 << 2
	UserFlagBugHunterLevel1       UserFlags = 1 << 3
	UserFlagHypesquadOnlineHouse1 UserFlags = 1 << 6
	UserFlagHypesquadOnlineHouse2 UserFlags = 1 << 7
	UserFlagHypesquadOnlineHouse3 UserFlags = 1 << 8
	UserFlagPremiumEarlySupporter UserFlags = 1 << 9
	UserFlagTeamPseudoUser        UserFlags = 1 << 10
	UserFlagBugHunterLevel2       UserFlags = 1 << 14
	UserFlagVerifiedBot           UserFlags = 1 << 16
	UserFlagVerifiedDeveloper     UserFlags = 1 << 17
	UserFlagCertifiedModerator    UserFlags = 1 << 18
	UserFlagBotHttpInteractions   UserFlags = 1 << 19
	UserFlagActiveDeveloper       UserFlags = 1 << 22
)

type User struct {
	ID               string        `json:"id"`
	Username         string        `json:"username"`
	Discriminator    string        `json:"discriminator"`
	Avatar           string        `json:"avatar"`
	AvatarDecoration string        `json:"avatar_decoration"`
	Verified         *bool         `json:"verified"`
	Email            *string       `json:"email"`
	Flags            *UserFlags    `json:"flags"`
	Bot              *bool         `json:"bot"`
	System           *bool         `json:"system"`
	Banner           *string       `json:"banner"`
	MFAEnabled       *bool         `json:"mfa_enabled"`
	Locale           *string       `json:"locale"`
	AccentColor      *int          `json:"accent_color"`
	PremiumType      *PremiumTypes `json:"premium_type"`
	PublicFlags      *UserFlags    `json:"public_flags"`
}

type GuildMemberResponse struct {
	User                       User        `json:"user"`
	Nick                       string      `json:"nick"`
	Avatar                     interface{} `json:"avatar"`
	Roles                      []string    `json:"roles"`
	JoinedAt                   time.Time   `json:"joined_at"`
	PremiumSince               *time.Time  `json:"premiumSince"`
	Deaf                       bool        `json:"deaf"`
	Mute                       bool        `json:"mute"`
	Pending                    *bool       `json:"pending"`
	Permissions                *string     `json:"permissions"`
	CommunicationDisabledUntil *time.Time  `json:"communication_disabled_until"`
}

type ErrorResponse struct {
	Message string       `json:"message"`
	Code    APIErrorCode `json:"code"`
}
