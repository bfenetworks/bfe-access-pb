# Changelog

All notable changes to `bfe-access-pb` will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.3.8]

### Added

- Add traffic mirroring fields to `RequestLog`:
  - `mirror_hit` (field `842`)
  - `mirror_cluster` (field `843`)
  - `mirror_status` (field `844`)
  - `mirror_latency_us` (field `845`)
  - `mirror_ttfb_us` (field `846`)
  - `mirror_prompt_tokens` (field `847`)
  - `mirror_completion_tokens` (field `848`)
  - `mirror_finish_reason` (field `849`)
  - `mirror_error` (field `850`)

These optional fields record the `mod_traffic_mirror` traffic mirroring (shadow traffic) outcome of a request: `mirror_hit` and `mirror_cluster` are set synchronously when a request is selected for mirroring; `mirror_status`, latency/TTFB, and the parsed `usage` / `finish_reason` / OpenAI error fields are reserved for async mirror results.

### Changed

- `bfe_access_pb/bfe_access.proto`: add traffic mirroring fields (842-850).
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document fields 842-850; split the 841-880 range into quota-plan hit (841), traffic mirroring (842-860), and security/privacy reservation (861-880).

## [v0.3.7]

### Added

- Add AI cache status and cache key fields to `RequestLog`:
  - `ai_cache_status` (field `789`)
  - `ai_cache_key` (field `790`)

These optional fields record the `mod_ai_cache` exact-match cache outcome of an AI request: `ai_cache_status` is one of `hit` / `miss` / `skip` (empty when the cache is not enabled), and `ai_cache_key` logs the cache key for debugging only when explicitly enabled.

### Changed

- `bfe_access_pb/bfe_access.proto`: add `ai_cache_status` field (789) and `ai_cache_key` field (790).
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document fields 789 and 790.

## [v0.3.6]

### Added

- Add 1h-TTL cache write token metering field to `RequestLog`:
  - `ai_cache_write_1h_tokens` (field `788`)

This optional field records the number of cache write tokens with a 1-hour TTL, included in `ai_cache_write_tokens`. It supports cache billing where 1h-TTL cache writes are priced separately from 5-minute-TTL writes (e.g., Anthropic extended cache TTL).

### Changed

- `bfe_access_pb/bfe_access.proto`: add `ai_cache_write_1h_tokens` field (788).
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document field 788.

## [v0.3.5]

### Added

- Add image input token and video metering fields to `RequestLog`:
  - `ai_image_input_tokens` (field `786`)
  - `ai_video_count` (field `787`)

These optional fields support AI image-aware billing where image input tokens are priced separately from text tokens, and video generation billing where cost is calculated per generated video.

### Changed

- `bfe_access_pb/bfe_access.proto`: add image input token and video count fields.
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document fields 786 and 787.

## [v0.3.4]

### Added

- Add AI protocol style field to `RequestLog`:
  - `ai_protocol` (field `717`)

This optional field records the AI request protocol style (e.g., `openai`, `anthropic`) and is intended to support protocol-aware log analysis, routing reconciliation, and provider-specific billing in the AI Gateway.

### Changed

- `bfe_access_pb/bfe_access.proto`: add `ai_protocol` field (717).
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document field 717.

## [v0.3.3]

### Added

- Add AI request mode field to `RequestLog`:
  - `ai_mode` (field `716`)

- Add image metering field to `RequestLog`:
  - `ai_image_count` (field `785`)

These optional fields support AI image generation billing where cost is calculated per generated image, and enable mode-based log analysis and cost reconciliation.

### Changed

- `bfe_access_pb/bfe_access.proto`: add `ai_mode` and `ai_image_count` fields.
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document fields 716 and 785.
- `b2log/test_data/`: add missing test fixtures and `gen_test_data.go`.

## [v0.3.2]

### Added

- Add audio token metering fields to `RequestLog`:
  - `ai_audio_input_tokens` (field `783`)
  - `ai_audio_output_tokens` (field `784`)

These optional fields support AI audio billing where audio tokens are priced separately from text tokens.

### Changed

- `bfe_access_pb/bfe_access.proto`: add audio token fields.
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document fields 783 and 784.
