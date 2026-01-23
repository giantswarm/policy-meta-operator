# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Resolve updated code linter findings.
- Use AppVersion for image tag defaulting.

## [0.0.3] - 2025-02-04

### Changed

- Change destination namespace to `policy-system`.

## [0.0.2] - 2025-02-03

### Changed

- Fix `AutomatedExceptions` handling and implement cleanup.

## [0.0.1] - 2025-01-23

### Added

- First release of the `policy-meta-operator` app.
- disabled logger development mode to avoid panicking
- changed: `app.giantswarm.io` label group was changed to `application.giantswarm.io`

[Unreleased]: https://github.com/giantswarm/policy-meta-operator/compare/v0.0.3...HEAD
[0.0.3]: https://github.com/giantswarm/policy-meta-operator/compare/v0.0.2...v0.0.3
[0.0.2]: https://github.com/giantswarm/policy-meta-operator/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/giantswarm/policy-meta-operator/releases/tag/v0.0.1
