# Reviewed source and runtime evidence: pluralrules/testdata/compact-review.md.
# Original FormatJS fixtures remain unchanged; Node witnesses use their actual version.

id: pluralrules-formatjs-index-test-ts-040
source: formatjs:packages/intl-pluralrules/tests/index.test.ts
owner: pluralrules
status: accepted
reason: Reviewed 2026-09-25 at FormatJS 4bfe08a527d96a511e816a12d896e0ed0548d66e, index.test.ts:236. The fr compact test intentionally omits NumberFormat locale data, so InitializePluralRules/ResolvePlural leave e=0 and select other for 1200000. go-intl owns complete compact data and selects many from CLDR's e=6 rule, independently witnessed by Node 26.8.1. Both algorithms use the source formatted decimal; the old XFAIL description was incorrect. See compact-review.md for the source/configuration analysis.
native_witness: pluralrules-node-v26-8-1-compact-fr-1200000
review_after: 2026-12-31
removal_path: Recheck on a FormatJS reference update; resolve when the source test loads compact data and agrees, or reconsider the product contract if ECMA-402 or the pinned French CLDR rule changes. Keep the original fixture and independent witness while the missing-data test environment remains different.

id: pluralrules-formatjs-index-test-ts-041
source: formatjs:packages/intl-pluralrules/tests/index.test.ts
owner: pluralrules
status: accepted
reason: Reviewed 2026-09-25 at FormatJS 4bfe08a527d96a511e816a12d896e0ed0548d66e, index.test.ts:237. The fr compact test intentionally omits NumberFormat locale data, so e=0 selects other for 234500000. go-intl resolves the million exponent e=6 from its constructor-owned CLDR data and selects many, as does Node 26.8.1. This is the same optional-data fallback mismatch as 040, not selection from a compact display decimal. See compact-review.md.
native_witness: pluralrules-node-v26-8-1-compact-fr-234500000
review_after: 2026-12-31
removal_path: Recheck on a FormatJS reference update; resolve when the source test loads compact data and agrees, or reconsider the product contract if ECMA-402 or the pinned French CLDR rule changes. Do not alter the extracted other assertion to match the native witness.

id: pluralrules-formatjs-index-test-ts-044
source: formatjs:packages/intl-pluralrules/tests/index.test.ts
owner: pluralrules
status: accepted
reason: Reviewed 2026-09-25 at FormatJS 4bfe08a527d96a511e816a12d896e0ed0548d66e, index.test.ts:252. The compact-versus-standard test explicitly expects the no-NumberFormat-data fallback e=0 for fr compact 1200000. go-intl has complete compact data and selects many with e=6, independently verified on Node 26.8.1; standard notation still selects other. This repeats the input of 040 in a distinct original assertion. See compact-review.md.
native_witness: pluralrules-node-v26-8-1-compact-fr-1200000
review_after: 2026-12-31
removal_path: Recheck on a FormatJS reference update; resolve when this distinct source assertion reflects complete compact data, or reconsider if ECMA-402 or pinned CLDR changes. Keep the standard and compact native witnesses so the accepted difference cannot erase notation semantics.
