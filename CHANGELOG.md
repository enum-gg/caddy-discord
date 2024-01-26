# Changelog

## [1.1.2](https://github.com/enum-gg/caddy-discord/compare/v1.1.1...v1.1.2) (2024-01-26)


### Bug Fixes

* Authorisation when using roles would cease after first role check ([43fbe4e](https://github.com/enum-gg/caddy-discord/commit/43fbe4ecf894aa40d1d2d878a4bbff62198097c7))
* JWT Signing key will use hashed Discord App Client ID to prevent breaking on server reboot ([1b3219f](https://github.com/enum-gg/caddy-discord/commit/1b3219f164f7157cbcf90a9bc084c4b694e28dbe))
* Remove stale token from goreleaser workflow ([b8ea22c](https://github.com/enum-gg/caddy-discord/commit/b8ea22c16f444c3bf25768f5af7503b9a2e751ef))
* return a 401 HTTP Response on auth fail instead of looping OAuth forever ([8886d2c](https://github.com/enum-gg/caddy-discord/commit/8886d2c635c3ee779ad44006977b277b8ccaeb5b))

## [1.1.1](https://github.com/enum-gg/caddy-discord/compare/v1.1.0...v1.1.1) (2023-08-11)


### Bug Fixes

* Role checking will actually check the role now ([e5949d9](https://github.com/enum-gg/caddy-discord/commit/e5949d943805e10bdbaefbdffa2fbb273b52e11a))
