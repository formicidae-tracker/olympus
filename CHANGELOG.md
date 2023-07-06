# Change Log

All notable changes to this project will be documented in this file.


## v0.5.0 / 2023-06-12

Major rewrite of the olympus application.

### Added

  - Docker based deployment.

### Changed

 - Moves to gRPC.
 - Uses Angular Material
 - Makes the webapp a PWA.
 - Refactores all webapi.
 - Added common libraries for olympus gRPC
   `github.com/formicidae-tracker/olympus/pkg/api` and telemetry
   `github.com/formicidae-tracker/olympus/pkg/tm`
 - Alarm Push notifications

## Fixed

 - Various UI bugs
