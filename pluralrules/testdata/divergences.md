# Reviewed source and runtime evidence: pluralrules/testdata/compact-review.md.
# Historical FormatJS observations below were removed by the 2026-10-03 source refresh; exact-version Node witnesses are retained.

id: pluralrules-formatjs-index-test-ts-040
source: formatjs:packages/intl-pluralrules/tests/index.test.ts
owner: pluralrules
status: resolved
reason: Resolved 2026-10-03. The original assertion reviewed at FormatJS 4bfe08a527d96a511e816a12d896e0ed0548d66e is absent from the current source. The replacement compact tests expect many without ambient NumberFormat data. This historical fixture ID no longer identifies that original observation, so its exemption must not be transferred by numeric suffix. See compact-review.md for the original analysis; exact-version native witnesses remain active.
native_witness: pluralrules-node-v26-8-1-compact-fr-1200000
review_after: 2026-12-31
removal_path: Recheck on a FormatJS reference update; resolve when the source test loads compact data and agrees, or reconsider the product contract if ECMA-402 or the pinned French CLDR rule changes. Keep the original fixture and independent witness while the missing-data test environment remains different.

id: pluralrules-formatjs-index-test-ts-041
source: formatjs:packages/intl-pluralrules/tests/index.test.ts
owner: pluralrules
status: resolved
reason: Resolved 2026-10-03. The original assertion reviewed at FormatJS 4bfe08a527d96a511e816a12d896e0ed0548d66e is absent from the current source. The replacement compact tests expect many without ambient NumberFormat data. This historical fixture ID no longer identifies that original observation, so its exemption must not be transferred by numeric suffix. See compact-review.md for the original analysis; exact-version native witnesses remain active.
native_witness: pluralrules-node-v26-8-1-compact-fr-234500000
review_after: 2026-12-31
removal_path: Recheck on a FormatJS reference update; resolve when the source test loads compact data and agrees, or reconsider the product contract if ECMA-402 or the pinned French CLDR rule changes. Do not alter the extracted other assertion to match the native witness.

id: pluralrules-formatjs-index-test-ts-044
source: formatjs:packages/intl-pluralrules/tests/index.test.ts
owner: pluralrules
status: resolved
reason: Resolved 2026-10-03. The original assertion reviewed at FormatJS 4bfe08a527d96a511e816a12d896e0ed0548d66e is absent from the current source. The replacement compact tests expect many without ambient NumberFormat data. This historical fixture ID no longer identifies that original observation, so its exemption must not be transferred by numeric suffix. See compact-review.md for the original analysis; exact-version native witnesses remain active.
native_witness: pluralrules-node-v26-8-1-compact-fr-1200000
review_after: 2026-12-31
removal_path: Recheck on a FormatJS reference update; resolve when this distinct source assertion reflects complete compact data, or reconsider if ECMA-402 or pinned CLDR changes. Keep the standard and compact native witnesses so the accepted difference cannot erase notation semantics.
