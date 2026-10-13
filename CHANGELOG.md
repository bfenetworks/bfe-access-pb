# Changelog

All notable changes to `bfe-access-pb` will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.3.16]

### Added

- Add zero-cost reason marker field to `RequestLog`:
  - `ai_zero_cost_reason` (field `765`)

This optional field records why a 2xx request that holds an RMB quota plan was charged 0: `unrecognized` (upstream usage neither recognized nor estimatable), `aborted` (the client disconnected before the final usage), `cache_hit` (AI cache hit, the upstream was never called), `batch_release` (batch contract released without billing), `price_miss` (no price entry and no fallback), `price_zero` (price entry present but the billed dimensions are priced 0), `model_table_missing` / `config_missing` (pricing configuration unavailable), `no_usage_tokens` (no billable token count although no usage problem was detected) or `cost_zero` (fallback: 0 for an unclassified reason). It is absent when the request was billed, when no RMB plan was bound, or when the upstream call never succeeded. Set by `mod_ai_token_auth` together with the `BillingZeroCost` module counter (issue #1414), it makes every zero-cost 2xx request countable and attributable in access-log reconciliation instead of relying on cross-field heuristics, closing the structural observability gap that the single-branch issue #1410 fix left open: billing bypasses outside the unrecognized-usage branch (batch releases, price tables whose billed dimensions are priced 0, deduction failures, client aborts) used to produce no signal at all.

### Changed

- `ai_usage_source` (field `764`): document the additional `aborted` value, set by `mod_ai_token_auth` when the client disconnects before the final usage is seen (issue #1352 yardstick: no deduction, the token fields are auth-time seeding residue or the initial usage of the stream). The wire format is unchanged, no renumbering.
- `bfe_access_pb/bfe_access.proto`: add `ai_zero_cost_reason` field (765) and document the `aborted` value of `ai_usage_source`.
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document field 765 and the `aborted` value of field 764.

## [v0.3.15]

### Added

- Add usage source marker field to `RequestLog`:
  - `ai_usage_source` (field `764`)

This optional field records where the reported token usage came from: `parsed` (final usage parsed from the upstream response, including genuine zero-usage requests), `estimated` (EstimateToken estimate approved for billing, mirrored into the shared usage by the issue #1398 fix), or `unrecognized` (no final usage parsed and the request is not estimatable: billing fields are zeroed, no RMB deduction happens, and the reported token fields are auth-time seeding residue). It is absent when the upstream call never succeeded (non-2xx, where the usage fields are suppressed per the issue #1409 fix). Set by `mod_ai_token_auth` (issue #1410), it closes the observability gap where a billing bypass was indistinguishable from a legitimate zero-usage row (`ai_total_tokens=0`) and where an estimate-billed row was indistinguishable from a parsed-usage row, so billing bypasses become directly queryable in access-log reconciliation and alertable via the module counter instead of relying on out-of-band manual reconciliation.

### Changed

- `bfe_access_pb/bfe_access.proto`: add `ai_usage_source` field (764).
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document field 764.

## [v0.3.14]

### Added

- Add price lookup miss marker field to `RequestLog`:
  - `ai_price_miss` (field `763`)

This optional field marks billing responses whose `(model, mode)` price entry was absent from the cluster price table (`AIConf.ModelTable`). It is set by `mod_ai_token_auth` when the cost is computed by falling back to the model's chat-mode price entry (the issue #1408 fix), and also when no fallback price exists and the request is billed 0. It makes under-priced / zero-cost rows caused by missing price entries distinguishable from legitimate zero-cost rows (e.g. `/count_tokens`, client aborts, pre-pricing traffic) in access-log reconciliation, closing the observability gap where `ai_cost_value` is only written when greater than 0. Empty (unset) for normally priced requests.

### Changed

- `bfe_access_pb/bfe_access.proto`: add `ai_price_miss` field (763).
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document field 763.

## [v0.3.13]

### Added

- Add batch task (mod_ai_batch) fields to `RequestLog`, in the new AI Observability sub-range 901-920 (range extended from 701-900):
  - `ai_batch_id` (field `901`)
  - `ai_file_id` (field `902`)
  - `ai_batch_op` (field `903`)
  - `ai_file_lines` (field `904`)
  - `ai_file_bytes` (field `905`)
  - `ai_batch_status` (field `906`)
  - `ai_batch_settle` (field `907`)

These optional fields support OpenAI Batch API (`/v1/files`, `/v1/batches`) passthrough observability: `ai_batch_id` / `ai_file_id` identify the batch task and file reported by `mod_ai_batch` (file operations belonging to a known batch backfill their batch ID); `ai_batch_op` records the operation type (`upload` / `create` / `get` / `list` / `cancel` / `download`); `ai_file_lines` / `ai_file_bytes` record the jsonl line count (counted at upload, or parsed from a batch output file) and the file size in bytes; `ai_batch_status` snapshots the provider batch state machine (`validating` / `queued` / `in_progress` / `finalizing` / `completed` / `expired` / `failed` / `cancelled` / `cancelling`); `ai_batch_settle` marks the quota settlement outcome (`settle` = settled by result-file usage, `release` = reserve released without settlement, `none`). All are empty for non-batch traffic.

### Changed

- `bfe_access_pb/bfe_access.proto`: extend the AI Observability range header to 701-920; add batch task fields (901-907).
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: add sub-range 901-920 (batch tasks, mod_ai_batch) to the range plan; document fields 901-907; update the reserved-range table.

## [v0.3.12]

### Added

- Add upstream error normalization (`AIConf.NormalizeUpstreamError`) fields to `RequestLog`:
  - `ai_upstream_status` (field `810`)
  - `ai_upstream_err_code` (field `811`)
  - `ai_err_normalized` (field `812`)
  - `ai_err_normalize_miss` (field `813`)
  - `ai_stream_error_rewritten` (field `814`)
  - `ai_stream_truncated` (field `815`)

These optional fields record the upstream error normalization outcome of an AI request: `ai_upstream_status` / `ai_upstream_err_code` keep the original upstream HTTP status and error code (OpenAI `error.code` / Anthropic `error.type` / Gemini `error.status`, redacted when `RedactSecrets=true`); `ai_err_normalized` marks whether the upstream error was rewritten into the unified error body (`true`) or passed through (`false`, including unrecognized errors and disabled clusters); `ai_err_normalize_miss` marks normalization misses (parser returned nil) to complete the protocol-to-catalog mapping tables; `ai_stream_error_rewritten` marks that at least one SSE error event of the stream was rewritten (`StreamEnabled`); `ai_stream_truncated` marks a truncated stream (EOF without the protocol's terminal event, `StreamEnabled`; Gemini streams end at HTTP EOF and never set this). All are empty when the cluster does not enable `NormalizeUpstreamError`.

### Changed

- `bfe_access_pb/bfe_access.proto`: add upstream error normalization fields (810-815).
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document fields 810-815.

## [v0.3.11]

### Added

- Add AI context compression (mod_ai_context) fields to `RequestLog`:
  - `ai_context_compress_status` (field `793`)
  - `ai_context_tokens_before` (field `794`)
  - `ai_context_tokens_after` (field `795`)
  - `ai_context_compress_mode` (field `796`)

These optional fields record the `mod_ai_context` context compression and pruning outcome of an AI request: `ai_context_compress_status` is `trim` / `rewrite` on success, `skip_no_rule` / `skip_protocol` / `skip_body_incomplete` / `skip_parse_err` / `skip_under_threshold` when skipped, and `repair_rollback` when protocol repair failed and the original request was forwarded unchanged; `ai_context_tokens_before` / `ai_context_tokens_after` record the heuristic prompt-token estimates before and after compression (empty when not compressed); `ai_context_compress_mode` records the effective mode (`conservative` / `balanced` / `aggressive`). All are empty when the request was not processed by the module.

### Changed

- `bfe_access_pb/bfe_access.proto`: add AI context compression fields (793-796).
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document fields 793-796; extend the 761-800 range description to cover context compression.

## [v0.3.10]

### Added

- Add AI cache (mod_ai_cache) semantic caching fields to `RequestLog`:
  - `ai_cache_semantic` (field `791`)
  - `ai_cache_similarity` (field `792`)

These optional fields support the second-phase semantic cache of `mod_ai_cache` (embedding + vector similarity lookup): `ai_cache_semantic` marks that the hit was served by the semantic cache (`hit_semantic` status), and `ai_cache_similarity` records the normalized similarity score (larger = more similar) for threshold calibration. Both are empty when no semantic lookup was performed.

### Changed

- `bfe_access_pb/bfe_access.proto`: add AI cache semantic fields (791-792).
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document fields 791-792.

## [v0.3.9]

### Added

- Add AI intent (semantic routing) fields to `RequestLog`:
  - `ai_intent_question` (field `803`)
  - `ai_intent_answer` (field `804`)
  - `ai_intent_confidence` (field `805`)
  - `ai_intent_source` (field `806`)
  - `ai_intent_latency_us` (field `807`)
  - `ai_intent_cache_hit` (field `808`)
  - `ai_intent_questions_version` (field `809`)

These optional fields record the `mod_ai_intent` semantic routing intent outcome consumed by routing: the question name, classified answer (`unknown` when classification failed/timed out/below the confidence gate), post-gate confidence, answer source (`explicit_header` / `classifier` / `cache`), decision service latency, intent LRU cache hit, and the version of the questions config used. They enable routing audit and confidence-threshold calibration.

### Changed

- `bfe_access_pb/bfe_access.proto`: add AI intent fields (803-809).
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document fields 803-809.

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
