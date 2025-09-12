# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

### [Unreleased]

### Added

### Fixed

### Changed

- The backend was updated to work with Go v18.

### Removed

## [v0.4.0] - 01-15-2019

### Added

### Fixed

### Changed

### Removed

- Remove `/api/address` endpoint. Use `/api/transactions` instead.
- Remove precompiled `dist/` folder. Users should compile the frontend with `make build-ng` or use a docker image.
- Remove `-use-unversioned-api` option

[Unreleased]: https://github.com/skycoin/skycoin/compare/master...develop
[0.25.0]: https://github.com/skycoin/skycoin/compare/v0.4.0...v0.3.0
