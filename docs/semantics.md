# Normative semantics

[Documentation index](README.md) · [Repository entry](../README.md)

## Lookup and resolution

| Operation | Input | May choose another tag | Uses app default | Mutates value |
|---|---|---:|---:|---:|
| `Get`, `Has`, `Require` | one tag | no | no | no |
| `match.Best` | weighted preferences | yes, BCP 47 matcher | no | no |
| `FallbackPlan.Resolve` | fixed exact chain | yes, configured only | optional | no |
| `Plan.Resolve` | requested tag and graph | yes, exact/parent graph | optional | no |
| `http.Select` | `Accept-Language` | yes, including wildcard | no | no |

Result kinds are `Exact`, `Matched`, `Fallback`, `Default`, and `Missing`.
`Present` is true for all successful kinds. `Empty` is true only when the
selected entry is present with an empty string.

## Special locale classes

| Locale class | Construction default | Matching |
|---|---|---|
| `und` | accepted; optionally rejected | exact when present, otherwise locale matcher policy |
| `mul` | accepted; optionally rejected | exact when present, otherwise locale matcher policy |
| private use | accepted; optionally rejected | exact identity; no invented relationship |
| zero/absent tag | rejected | invalid preference or missing plan request |
| unknown but registry-valid | accepted | pinned locale matcher policy |

HTTP `*` is represented only inside the HTTP adapter and selects the first
canonical stored locale after quality ordering. It is not stored as `und` and
does not change `match.Best` semantics.

## Missing, empty, and invalid

| State | Stored entry | `Get` text | `Get` present | Encoding |
|---|---:|---|---:|---|
| Missing | no | `""` | false | key absent |
| Present-empty | yes | `""` | true | `"tag":""` |
| Whitespace-only | yes | unchanged | true | unchanged string |
| Invalid UTF-8 | never | n/a | n/a | typed rejection |
| SQL NULL | nullable wrapper invalid | n/a | n/a | database NULL |
| Zero `Text` | no entries | `""` | false | `{}` |

Whitespace-only is valid core text. Applications opt into
`validation.RequireNonWhitespace` from `adapters/validation`.

## Merge conflicts

| Left | Right | `LeftWins` | `RightWins` | `RejectConflict` | Resolver |
|---|---|---|---|---|---|
| absent | absent | absent | absent | absent | not called |
| value | absent | left | left | left | not called |
| absent | value | right | right | right | not called |
| value | value | left | right | error | callback |
| empty | value, `EmptyIsValue` | empty | right | error | callback |
| empty | value, `EmptyIsAbsent` | right | right | right | not called |
| value | empty, `EmptyIsAbsent` | left | left | left | not called |
| empty | empty, `EmptyIsAbsent` | absent | absent | absent | not called |

Conflict identity is computed after tag canonicalization. Merge output is
deterministically ordered and revalidated against output limits.

## Duplicate construction

| Policy | Canonical duplicate result |
|---|---|
| `RejectDuplicates` | `ErrDuplicateLocale` |
| `FirstWins` | first input occurrence retained |
| `LastWins` | last input occurrence retained |

Input cardinality and bytes are checked before duplicate reduction so duplicate
floods cannot bypass work limits.

## Encoding

| Form | Order | Duplicate visibility | Null | Empty |
|---|---|---|---|---|
| Canonical JSON object | canonical lexical keys | detected while tokenizing | strict reject | `{}` |
| Permissive JSON object | canonical lexical keys | detected | accepted as zero | `{}` |
| Entry-array JSON | canonical lexical entries | detected after parsing | reject | `[]` |
| SQL/pgx JSONB | semantic object | PostgreSQL/object rules | explicit wrapper | `{}` |
| wire adapters | format deterministic rules | adapter and core rules | format-specific | map/object |

Canonicalization changes tag spelling and preferred aliases only. It never
creates, copies, translates, or removes a localized string.

## Existing operation and decoder limitations

These describe current implementation limitations, not additional normative
protocol interpretations or conformance success.

Only wholly zero `Limits` selects defaults. A partially specified struct keeps
individual zero allowances. Count and total-text-byte admission includes all
input entries before duplicate reduction; tag bytes and total heap are not part
of `MaxTotalBytes`. Raw pair and core JSON keys also keep an independent
255-byte parser ceiling when a caller raises the canonical tag limit.

`Set`, `Map` and explicit normalization reconstruct under default limits,
not the originating custom construction policy. `Remove` and non-nil `Filter`
also reconstruct under defaults but discard failure: a valid custom-limit value
can silently become empty. Use explicit construction with the intended limits
when that distinction matters. Builder feeding, merge temporaries and callback
execution occur before final value admission. Callback side effects are not
rolled back; observers isolate panics, not blocking or concurrent access.

Raw malformed UTF-8 is rejected, but encoding/json can repair unpaired surrogate
escapes in decoded text to U+FFFD. Entry-array decoding accepts case aliases and
repeated object fields with the last occurrence winning; omitted/null `text`
becomes a present-empty string. Canonical duplicate locales across entries still
reject. The null column above refers to the top-level representation, not strict
required-member or nested scalar-type enforcement. Encoders have no output-byte
ceiling: escaping can exceed default decoder admission even for a value admitted
under default construction limits. Callers own subsequent wire-boundary limits.

`Best` tries exact then approximate matching for each preference in descending
weight/stable-tie order; a higher-weight approximate match can precede a
lower-weight exact one. Matching delegates to pinned x/text data, not this
package's own CLDR implementation. `ParentRange` starts at a candidate's strict
parent, excluding that candidate. HTTP wildcard picks the first lexical stored
locale and zero-weight ranges are not a global exclusion set for wildcard or
approximate selection. These helpers do not perform network I/O or response
header policy.

Choose the next guide in the [documentation index](README.md)
or return to the [repository entry](../README.md) for installation and ownership.
