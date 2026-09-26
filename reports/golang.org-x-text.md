# golang.org/x/text v0.41.0: newer CLDR locale subtags

The pinned CLDR 48.1.0 likely-subtag source contains valid language/script
identifiers that `language.Parse` does not recognize. `locale.Parse`, which uses
that dependency as its BCP 47 boundary, consequently returns ErrInvalidValue.

Reproduction: call `locale.Parse("und-Gara")` or `locale.Parse("sjc")` on
Go 1.27.0 with x/text v0.41.0. The error is `locale: invalid value languageTag
"und-Gara": expected a well-formed BCP 47 language tag; got "und-Gara"`;
the wrapped dependency error reports an unrecognized subtag.

Expected: structurally valid identifiers present in the pinned CLDR source can
cross the locale parser. Actual rejected source keys are `bap-Krai`, `hnm`,
`luh`, `rrm`, `sjc`, `ynb`, `und-Berf`, `und-Gara`, `und-Gukh`, `und-Krai`,
`und-Onao`, `und-Sidt`, `und-Tayo`, `und-Todr`, `und-Tols`, and `und-Tutg`.
Source: `tools/gen-cldr/.cldr-json/node_modules/cldr-core/supplemental/likelySubtags.json`.

The generator validates and retains the full source through the project's
Unicode subtag grammar; no runtime replacement parser or dependency fork is
introduced. Upgrade x/text when its recognition data includes these subtags.
Public minimize tests cover the supported locale profile and explicit boundary
cases; full source maximize coverage remains in the generator round-trip gate.
This limitation precedes the minimize algorithm change and is not evidence of
an incorrect candidate-equivalence algorithm.
