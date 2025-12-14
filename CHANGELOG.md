# Changelog

## [0.7.0](https://github.com/mikemrm/masscan-exporter/compare/v0.6.4...v0.7.0) (2025-12-14)


### Features

* **release:** build multi-arch images ([#28](https://github.com/mikemrm/masscan-exporter/issues/28)) ([97cf6a4](https://github.com/mikemrm/masscan-exporter/commit/97cf6a466755ceddac4978b9250e60c93ba0ca0d))

## [0.6.4](https://github.com/mikemrm/masscan-exporter/compare/v0.6.3...v0.6.4) (2025-12-07)


### Bug Fixes

* panic when unhealthy_failed_scrapes is configured ([44d18a9](https://github.com/mikemrm/masscan-exporter/commit/44d18a96b519e91007cc9f222509dedf3a9f896e))

## [0.6.3](https://github.com/mikemrm/masscan-exporter/compare/v0.6.2...v0.6.3) (2025-12-07)


### Bug Fixes

* **release:** correct environment variables ([8776030](https://github.com/mikemrm/masscan-exporter/commit/8776030bd7acd9809f164dd738dca12c56e59ba5))

## [0.6.2](https://github.com/mikemrm/masscan-exporter/compare/v0.6.1...v0.6.2) (2025-12-07)


### Bug Fixes

* **release:** correct chart repo ([e0df252](https://github.com/mikemrm/masscan-exporter/commit/e0df25280d6ed1f6b40d54ef565cb657572eb90e))

## [0.6.1](https://github.com/mikemrm/masscan-exporter/compare/v0.6.0...v0.6.1) (2025-12-07)


### Miscellaneous Chores

* re-release v0.6.0 ([eaae551](https://github.com/mikemrm/masscan-exporter/commit/eaae551aea36d224479fe8cd615a1fcca5652a84))

## [0.6.0](https://github.com/mikemrm/masscan-exporter/compare/v0.5.1...v0.6.0) (2025-12-07)


### Features

* add health metrics and endpoints ([#21](https://github.com/mikemrm/masscan-exporter/issues/21)) ([172fef8](https://github.com/mikemrm/masscan-exporter/commit/172fef850a4b91b064d65d4252443bfef935a0c2))
* add helm chart ([#22](https://github.com/mikemrm/masscan-exporter/issues/22)) ([081e898](https://github.com/mikemrm/masscan-exporter/commit/081e89816779b0f390a33bd75b479f777fe2e501))


### Bug Fixes

* **deps:** update module github.com/adhocore/gronx to v1.19.6 ([#9](https://github.com/mikemrm/masscan-exporter/issues/9)) ([827e265](https://github.com/mikemrm/masscan-exporter/commit/827e2653a7afab101f723272db36bf184a821570))
* **deps:** update module github.com/prometheus/client_golang to v1.23.2 ([#11](https://github.com/mikemrm/masscan-exporter/issues/11)) ([db3e9b0](https://github.com/mikemrm/masscan-exporter/commit/db3e9b06c47d78884a506fdef7cdc72b6318f911))
* **deps:** update module github.com/spf13/cobra to v1.10.2 ([#14](https://github.com/mikemrm/masscan-exporter/issues/14)) ([caab0e5](https://github.com/mikemrm/masscan-exporter/commit/caab0e50c5c2ecdb18c32b7ec9746849d8172af6))
* **deps:** update module github.com/spf13/viper to v1.21.0 ([#16](https://github.com/mikemrm/masscan-exporter/issues/16)) ([f25f815](https://github.com/mikemrm/masscan-exporter/commit/f25f81569e6e253276eff3f14516454da09e5fdb))
* gracefully handle no ports found ([#20](https://github.com/mikemrm/masscan-exporter/issues/20)) ([fcfd365](https://github.com/mikemrm/masscan-exporter/commit/fcfd365570cf9533456c9e8373608b685d3cfcb8))

## [0.5.1](https://github.com/mikemrm/masscan-exporter/compare/v0.5.0...v0.5.1) (2025-04-26)


### Bug Fixes

* add docs referencing the grafana dashboard id to import ([bfa97b2](https://github.com/mikemrm/masscan-exporter/commit/bfa97b2722c2db8f768c6e22794d3faf98f2c099))
* add license ([62046a3](https://github.com/mikemrm/masscan-exporter/commit/62046a3575b335ecf7ac3a6f678e586b36b6081b))
* cleanup some unused code, add grafana dashboard ([2eb1533](https://github.com/mikemrm/masscan-exporter/commit/2eb15333fdb5f9b8264caf6e5e982e5dc1d220b1))

## [0.5.0](https://github.com/mikemrm/masscan-exporter/compare/v0.4.0...v0.5.0) (2025-04-26)


### Features

* changes to supporting multiple collectors ([935f935](https://github.com/mikemrm/masscan-exporter/commit/935f935c867a1b7935423410d3569c16cbf8e5ae))

## [0.4.0](https://github.com/mikemrm/masscan-exporter/compare/v0.3.0...v0.4.0) (2025-04-25)


### Features

* add caching ([71bbfc2](https://github.com/mikemrm/masscan-exporter/commit/71bbfc2b91a12ddd486dd5fcaaa2707ce71ed9b2))
* support providing a masscan config or config path ([99b3c96](https://github.com/mikemrm/masscan-exporter/commit/99b3c96afd092d678278d307853d868d82b9727c))


### Bug Fixes

* remove old main ([859f969](https://github.com/mikemrm/masscan-exporter/commit/859f9698bad003b3b3c5e42d015a2f38d2075a54))

## [0.3.0](https://github.com/mikemrm/masscan-exporter/compare/v0.2.0...v0.3.0) (2025-04-20)


### Features

* include scrape seconds as a metric ([aa73177](https://github.com/mikemrm/masscan-exporter/commit/aa73177a218911899137667f2ed00e8dc835d10d))

## [0.2.0](https://github.com/mikemrm/masscan-exporter/compare/v0.1.1...v0.2.0) (2025-04-20)


### Features

* ensure only a single scrape happens at a time ([1de8749](https://github.com/mikemrm/masscan-exporter/commit/1de8749216d2a14c243468d58c27c100f412efc6))

## [0.1.1](https://github.com/mikemrm/masscan-exporter/compare/v0.1.0...v0.1.1) (2025-04-19)


### Bug Fixes

* add missing libpcap ([9b935a4](https://github.com/mikemrm/masscan-exporter/commit/9b935a4b5309507ed8e6eb758344567ac8cced5f))

## 0.1.0 (2025-04-19)


### Features

* initial commit ([1c66fee](https://github.com/mikemrm/masscan-exporter/commit/1c66fee7e9ad4c75cc3da4efe5972bf5a4145702))


### Bug Fixes

* update example to show more results ([4f8009c](https://github.com/mikemrm/masscan-exporter/commit/4f8009c66db823514c04106d773f6ed32097a582))


### Miscellaneous Chores

* release 0.1.0 ([265fa97](https://github.com/mikemrm/masscan-exporter/commit/265fa971a6a22f2c875e474172d2dfcca0f94c61))
