id: datetimeformat-formatjs-format-range-test-ts-014
source: formatjs:packages/intl-datetimeformat/tests/format-range.test.ts
owner: datetimeformat
status: resolved
reason: Resolved 2026-10-03. Regenerated formatRangeToParts from the current FormatJS source now marks the common month, following literal, and range separator as shared, agreeing with the retained native witness and go-intl.
native_witness: datetimeformat-node-v26-range-shared-month-prefix
review_after: 2026-12-31
removal_path: Resolve if the FormatJS fixture aligns with native range-source semantics, or keep accepted while the FormatJS source remains the only disagreeing reference.

id: datetimeformat-formatjs-index-test-ts-000
source: formatjs:packages/intl-datetimeformat/tests/index.test.ts
owner: datetimeformat
status: resolved
reason: Resolved 2026-10-03. Strict options/declaration extraction no longer generates the Amsterdam observation associated with this historical ID. The regenerated index-test-ts-000 is a distinct UTC date observation and must execute without this spacing exemption; the native spacing witness remains unchanged.
native_witness: datetimeformat-node-v26-day-period-time-zone-name-spacing
review_after: 2026-12-31
removal_path: Remove this divergence if FormatJS aligns with native Intl spacing, or if a future Node witness changes the native output.

id: datetimeformat-formatjs-format-range-test-ts-013
source: formatjs:packages/intl-datetimeformat/tests/format-range.test.ts
owner: datetimeformat
status: accepted
reason: Reviewed 2026-10-03. The refreshed FormatJS assertion uses U+202F before AM/PM, while Node 26.10.0 with the same locale, component options, and endpoints uses ASCII space. go-intl retains its native-aligned day-period spacing; this is a localized literal difference, not a missing range operation. The independent witness records both complete text and parts.
native_witness: datetimeformat-node-v26-10-0-range-spacing-cross-date
review_after: 2026-12-31
removal_path: Resolve if the source assertion or native/pinned locale pattern aligns. Recheck the exact input and spacing on data/reference refresh; do not transfer this exemption to a different observation after fixture ID changes.

id: datetimeformat-formatjs-index-test-ts-003
source: formatjs:packages/intl-datetimeformat/tests/index.test.ts
owner: datetimeformat
status: accepted
reason: Reviewed 2026-10-03. The refreshed FormatJS assertion uses U+202F before AM/PM, while Node 26.10.0 with the same locale, component options, and endpoints uses ASCII space. go-intl retains its native-aligned day-period spacing; this is a localized literal difference, not a missing range operation. The independent witness records both complete text and parts.
native_witness: datetimeformat-node-v26-10-0-range-spacing-new-york
review_after: 2026-12-31
removal_path: Resolve if the source assertion or native/pinned locale pattern aligns. Recheck the exact input and spacing on data/reference refresh; do not transfer this exemption to a different observation after fixture ID changes.
