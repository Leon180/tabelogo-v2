# Changelog

## [0.7.0](https://github.com/Leon180/tabelogo-v2/compare/tabelogo-v2/pkg-v0.6.0...tabelogo-v2/pkg-v0.7.0) (2025-12-30)


### Features

* Week 2 - Prometheus metrics enhancement + Grafana dashboards ([9874e1f](https://github.com/Leon180/tabelogo-v2/commit/9874e1f5cb48f6affb6142a3cddc7df5a73c046e))

## [0.6.0](https://github.com/Leon180/tabelogo-v2/compare/tabelogo-v2/pkg-v0.5.0...tabelogo-v2/pkg-v0.6.0) (2025-12-24)


### Features

* **middleware:** update auth middleware with session validation (Phase 2.1) ([2c32769](https://github.com/Leon180/tabelogo-v2/commit/2c3276967393d18694aa04fbda97533e4b2e91ce))


### Bug Fixes

* **auth:** update tests for new JWT signature ([b8498cf](https://github.com/Leon180/tabelogo-v2/commit/b8498cfe4c3f53911eac4f5372a390f74a8355d0))
* **middleware:** handle JSON string session format from Redis ([553ad08](https://github.com/Leon180/tabelogo-v2/commit/553ad08329c3bca215cec586db163022d5c6428b))

## [0.5.0](https://github.com/Leon180/tabelogo-v2/compare/tabelogo-v2/pkg-v0.4.0...tabelogo-v2/pkg-v0.5.0) (2025-12-14)


### Features

* Add Japanese name support to Restaurant model ([64694ce](https://github.com/Leon180/tabelogo-v2/commit/64694cec4fa9312567efde43ac8616ba07b8e674))
* Add Prometheus metrics for Restaurant Service cache performance ([66f50b7](https://github.com/Leon180/tabelogo-v2/commit/66f50b7c0e1dd6eb96333397a8a993b1e936dfad))
* Complete Map Service Phase 1 + Start Phase 2 Integration ([7739085](https://github.com/Leon180/tabelogo-v2/commit/773908597ee87969db5889f983e14fff3abcd90d))

## [0.4.0](https://github.com/Leon180/tabelogo-v2/compare/tabelogo-v2/pkg-v0.3.0...tabelogo-v2/pkg-v0.4.0) (2025-12-02)


### Features

* **map-service:** implement complete Map Service with Phase 1-4 ([c4670ee](https://github.com/Leon180/tabelogo-v2/commit/c4670ee0675de81ae24e8a0bfec132ff5489079a))

## [0.3.0](https://github.com/Leon180/tabelogo-v2/compare/tabelogo-v2/pkg-v0.2.0...tabelogo-v2/pkg-v0.3.0) (2025-11-24)


### Features

* **auth:** finalize docker and makefile configuration ([d84778a](https://github.com/Leon180/tabelogo-v2/commit/d84778a50322b99f1e1ab5bcfe7488d3eb2515ae))
* **auth:** implement domain, infra, app, and grpc layers ([ce6abd7](https://github.com/Leon180/tabelogo-v2/commit/ce6abd7ae9dd67fab16e19c0b96cfd677e7455a6))

## [0.2.0](https://github.com/Leon180/tabelogo-v2/compare/tabelogo-v2/pkg-v0.1.0...tabelogo-v2/pkg-v0.2.0) (2025-11-22)


### Features

* allow prefix env key ([9a8f6dd](https://github.com/Leon180/tabelogo-v2/commit/9a8f6dd07fb9538f9476badf07d7e3df7572ca68))
* complete Phase 1 - migrations, shared packages, and middleware ([6e4933c](https://github.com/Leon180/tabelogo-v2/commit/6e4933c132eb48d7f431926f4d34ad8084c907e5))
* init migrations pkg and test case ([7b624f7](https://github.com/Leon180/tabelogo-v2/commit/7b624f76113d03ce02aadb12395f38a435f09be8))


### Bug Fixes

* modify mod name, and apply fx to pkg ([960e177](https://github.com/Leon180/tabelogo-v2/commit/960e17712c30c9b162551f94aed0c38b6942f1af))
* standardize Go version to 1.23 and update gitignore for workspace files ([05f10d9](https://github.com/Leon180/tabelogo-v2/commit/05f10d9f5985e1de842886b66bef0c6a5b5aabc6))
