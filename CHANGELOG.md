# Changelog

## [0.6.1](https://github.com/kallangerard/pantalon/compare/v0.6.0...v0.6.1) (2026-11-23)


### Miscellaneous

* **deps:** update actions/setup-go action to v7.1.0 ([#126](https://github.com/kallangerard/pantalon/issues/126)) ([c4e61e9](https://github.com/kallangerard/pantalon/commit/c4e61e92d2898cfdb6e22813ffe157a1b4a6f71e))

## [0.6.0](https://github.com/kallangerard/pantalon/compare/v0.5.1...v0.6.0) (2026-11-21)


### Features

* add --version flag ([#124](https://github.com/kallangerard/pantalon/issues/124)) ([d45cab0](https://github.com/kallangerard/pantalon/commit/d45cab0eba7ccec9d4ac7eaad93c89b543d6c9aa))


### Bug Fixes

* require label and context values to be YAML strings ([#121](https://github.com/kallangerard/pantalon/issues/121)) ([7818242](https://github.com/kallangerard/pantalon/commit/7818242e4b019598bef9fcbf4fc467ca0d51951f))
* validate metadata.name as an RFC 1123 label (max 63 chars) ([#122](https://github.com/kallangerard/pantalon/issues/122)) ([ed5f1e5](https://github.com/kallangerard/pantalon/commit/ed5f1e5b41830d4d4718ff2f85419dd8f7a0a410))


### Performance Improvements

* skip .git and .terraform directories during discovery ([#125](https://github.com/kallangerard/pantalon/issues/125)) ([b7fd7e5](https://github.com/kallangerard/pantalon/commit/b7fd7e5fe00c72c614a82fb00a09ce1b0c60c14c))


### Miscellaneous

* **deps:** update goreleaser/goreleaser to v2.19.0 ([#111](https://github.com/kallangerard/pantalon/issues/111)) ([cef0069](https://github.com/kallangerard/pantalon/commit/cef0069ee055174c0ea1f1a81f820319c0c5b6b9))
* **deps:** update goreleaser/goreleaser to v2.19.0 ([#118](https://github.com/kallangerard/pantalon/issues/118)) ([b409fc5](https://github.com/kallangerard/pantalon/commit/b409fc5e516904533757c5e6d7fcd57bee15080b))
* **deps:** update goreleaser/goreleaser to v2.20.0 ([#119](https://github.com/kallangerard/pantalon/issues/119)) ([30b0deb](https://github.com/kallangerard/pantalon/commit/30b0deb7b09dfba0164d2211097e15758049e180))
* **deps:** update goreleaser/goreleaser to v2.20.2 ([#120](https://github.com/kallangerard/pantalon/issues/120)) ([3731d37](https://github.com/kallangerard/pantalon/commit/3731d37ce609f02cd2058cd6d5a5dc14b521dafe))

## [0.5.1](https://github.com/kallangerard/pantalon/compare/v0.5.0...v0.5.1) (2026-11-04)


### Miscellaneous

* **deps:** update module github.com/bmatcuk/doublestar/v4 to v4.10.4 ([#112](https://github.com/kallangerard/pantalon/issues/112)) ([18498ba](https://github.com/kallangerard/pantalon/commit/18498bafa61f4194c07b1f48f913efba5a6142ff))

## [0.5.0](https://github.com/kallangerard/pantalon/compare/v0.4.0...v0.5.0) (2026-10-30)


### Features

* reject YAML anchors, aliases and merge keys in pantalon.yaml ([#110](https://github.com/kallangerard/pantalon/issues/110)) ([77adb93](https://github.com/kallangerard/pantalon/commit/77adb93d453cecda12369a0ed60301b496c1bfea)), closes [#107](https://github.com/kallangerard/pantalon/issues/107)


### Bug Fixes

* name the file in validation errors, warn on unknown fields, reject duplicate names ([#106](https://github.com/kallangerard/pantalon/issues/106)) ([bb312a9](https://github.com/kallangerard/pantalon/commit/bb312a997bc7e23b7ad3ca48cbb7332bff0e5fef)), closes [#103](https://github.com/kallangerard/pantalon/issues/103)


### Miscellaneous

* add infra and QA reviewer agents ([#109](https://github.com/kallangerard/pantalon/issues/109)) ([9f660dd](https://github.com/kallangerard/pantalon/commit/9f660dd1025a22098fb893494bf09512fae2e046))

## [0.4.0](https://github.com/kallangerard/pantalon/compare/v0.3.1...v0.4.0) (2026-10-27)


### Features

* select configurations by dependencies.paths ([#100](https://github.com/kallangerard/pantalon/issues/100)) ([9365878](https://github.com/kallangerard/pantalon/commit/93658786fb7e878fdc82954e48c98d1e7385d5d8))

## [0.3.1](https://github.com/kallangerard/pantalon/compare/v0.3.0...v0.3.1) (2026-10-22)


### Bug Fixes

* match changed dirs on whole path segments and normalise input ([#97](https://github.com/kallangerard/pantalon/issues/97)) ([60cd6e6](https://github.com/kallangerard/pantalon/commit/60cd6e6423084d91ec0bc403b6399343e3e90dcb))

## [0.3.0](https://github.com/kallangerard/pantalon/compare/v0.2.2...v0.3.0) (2026-10-21)


### Features

* filter configurations by kubectl-style label selectors ([#95](https://github.com/kallangerard/pantalon/issues/95)) ([1a970e7](https://github.com/kallangerard/pantalon/commit/1a970e7a01b54045b9f795d4b689c08aa66ef792))

## [0.2.2](https://github.com/kallangerard/pantalon/compare/v0.2.1...v0.2.2) (2026-06-10)


### Bug Fixes

* dedupe ChangedFiles output to avoid duplicate matrix jobs ([#78](https://github.com/kallangerard/pantalon/issues/78)) ([cc3614e](https://github.com/kallangerard/pantalon/commit/cc3614ee80a6651bed92bc847237f439c158efb3))

## [0.2.1](https://github.com/kallangerard/pantalon/compare/v0.2.0...v0.2.1) (2026-06-06)


### Bug Fixes

* allow not passing changed-dirs ([#75](https://github.com/kallangerard/pantalon/issues/75)) ([0437179](https://github.com/kallangerard/pantalon/commit/0437179ffa4881e10f01fdc99535c52c5988bf89))

## [0.2.0](https://github.com/kallangerard/pantalon/compare/v0.1.0...v0.2.0) (2026-06-05)


### Features

* **cmd:** enhance help output ([2effedb](https://github.com/kallangerard/pantalon/commit/2effedb8bae545c7f4276478faf25faa164b533a))


### Bug Fixes

* **release:** add explicit packages block to resolve v5 component mismatch ([#73](https://github.com/kallangerard/pantalon/issues/73)) ([7ac233e](https://github.com/kallangerard/pantalon/commit/7ac233ea17c94b1ecb9f5842adcb0ea4c210b907))

## 0.1.0 (2026-03-05)


### Features

* **ci:** create release workflow ([#44](https://github.com/kallangerard/pantalon/issues/44)) ([9d8b29f](https://github.com/kallangerard/pantalon/commit/9d8b29fe664046ddce12270d9f037495c5233729))


### Bug Fixes

* **ci:** write permission for issues in workflow ([8b1cd0d](https://github.com/kallangerard/pantalon/commit/8b1cd0dbcde3855c2e559315bd4d5fd20ab3cfe3))
