# Specification decisions

This register records observable choices where BCP 47, the IANA Language
Subtag Registry, RFC 4647, RFC 9110, RFC 8259, or a maintained dependency does
not by itself select one package behavior. Source pins, update monitoring,
machine bindings, and append-only history live under
[`specification/`](../specification/README.md).

## LOCALIZED-DEC-001: Canonical locale identity

Status `resolved`; owner `localized maintainers`; classification
`interoperability policy`; decision scope `application-policy`; specification
`BCP 47 language tags`; version `RFC 5646`; source authority `bcp47-rfc5646`;
authority URL https://www.rfc-editor.org/rfc/rfc5646.txt; section
`Sections 2.1, 2.2.9, and 4.5`; requirement strength `SHOULD`.

Additional authoritative source: `{"id":"iana-language-registry","version":"Registry File-Date 2026-06-14","url":"https://www.iana.org/assignments/language-subtag-registry/language-subtag-registry","specifications":["IANA Language Subtag Registry"]}`

| Field | Decision |
|---|---|
| Issue | Valid language tags can preserve caller spelling or use one registry-aware identity, and aliases can otherwise become distinct map keys. |
| Credible interpretations | Preserve every accepted spelling. Or: canonicalize every stored key through the pinned locale dependency. |
| Known peer behavior | The pinned `go-international/locale` implementation uses x/text-derived IANA tables; no independent local differential harness is claimed. |
| Selected behavior | Every stored and queried locale is canonicalized before identity comparison, ordering, hashing, persistence, or encoding. |
| Rationale | One canonical identity prevents aliases and casing variants from representing different translations. |
| Security consequences | Malformed, invalid UTF-8, and overlong locale text fails before retention. |
| Resource consequences | Canonicalization is bounded by the locale dependency and performs no network access. |
| Compatibility consequences | Preferred-value or registry-table changes require dependency, persistence, and release review. |
| Wire consequences | Encoders emit canonical locale keys rather than accepted aliases or caller casing. |
| Executable evidence | `TestStandardsCanonicalizationMatrix`; `TestStandardsRegistryProvenance` |
| Fuzz evidence | `FuzzTextProperties` |
| Public APIs | `NewText`; `TextFromMap`; `Text.Get`; `Text.Set`; `EncodeJSON` |
| Documentation | `docs/specification-decisions.md`; `docs/semantics.md`; `docs/compatibility.md` |
| Upstream status | RFC 5646 errata and the current IANA registry are monitored separately from the pinned locale dependency snapshot. |
| Reconsider when | BCP 47, its errata, the IANA registry, or the locale dependency changes canonical identity semantics. |

## LOCALIZED-DEC-002: Special and private-use locale acceptance

Status `resolved`; owner `localized maintainers`; classification `optional behavior`;
decision scope `application-policy`; specification `BCP 47 language tags`;
version `RFC 5646`; source authority `bcp47-rfc5646`; authority URL
https://www.rfc-editor.org/rfc/rfc5646.txt; section `Sections 2.2.1 and 2.2.7`;
requirement strength `MAY`.

| Field | Decision |
|---|---|
| Issue | Valid tags include `und`, `mul`, private-use tags, and registry-valid reserved values, but applications may not want each class. |
| Credible interpretations | Reject special classes globally. Or: accept valid tags by default and expose explicit construction policy. |
| Known peer behavior | Locale libraries generally parse these valid classes but do not define application acceptance policy. |
| Selected behavior | Construction accepts valid special, private-use, and registry-valid tags by default; `LocalePolicy` can reject each class explicitly. |
| Rationale | Syntax validity and application acceptance are separate decisions. |
| Security consequences | Applications can fail closed for classes their authorization or persistence model does not recognize. |
| Resource consequences | Policy evaluation is bounded and performs no registry lookup or I/O. |
| Compatibility consequences | Default acceptance remains stable; tightening requires an explicit caller policy or compatibility review. |
| Wire consequences | Accepted classes serialize as their canonical BCP 47 identity. |
| Executable evidence | `TestLocaleAcceptancePolicyIsExplicit` |
| Public APIs | `LocalePolicy`; `ConstructionOptions`; `NewTextWithOptions` |
| Documentation | `docs/specification-decisions.md`; `docs/semantics.md` |
| Upstream status | No BCP 47 rule selects one application acceptance policy for valid special classes. |
| Reconsider when | A supported application profile makes one class mandatory or forbidden. |

## LOCALIZED-DEC-003: Duplicate identity after canonicalization

Status `resolved`; owner `localized maintainers`; classification `ambiguity`;
decision scope `defensive`; specification `BCP 47 language tags`; version
`RFC 5646`; source authority `bcp47-rfc5646`; authority URL
https://www.rfc-editor.org/rfc/rfc5646.txt; section `Section 4.5`;
requirement strength `SHOULD`.

| Field | Decision |
|---|---|
| Issue | Distinct input spellings can canonicalize to the same locale and silently overwrite a translation. |
| Credible interpretations | First value wins. Last value wins. Or: reject unless an explicit duplicate policy selects a winner. |
| Known peer behavior | Generic maps commonly overwrite duplicates, but that behavior does not expose canonical identity collisions. |
| Selected behavior | Strict construction rejects canonical duplicate locales; explicit construction options may select left-wins or right-wins behavior. |
| Rationale | Silent overwrite can hide conflicting localized content or parser differentials. |
| Security consequences | Canonical duplicate smuggling fails closed by default. |
| Resource consequences | Duplicate detection is bounded by the locale-count limit. |
| Compatibility consequences | Existing strict callers retain rejection; permissive resolution requires an explicit policy. |
| Wire consequences | Strict JSON rejects duplicate raw members and distinct members that canonicalize to one locale. |
| Executable evidence | `TestDuplicatePoliciesApplyAfterCanonicalization`; `TestTextRejectsCanonicalDuplicate` |
| Public APIs | `DuplicatePolicy`; `NewTextWithOptions`; `DecodeJSON` |
| Documentation | `docs/specification-decisions.md`; `docs/semantics.md`; `docs/security.md` |
| Upstream status | BCP 47 defines tag identity but not an application map's collision policy. |
| Reconsider when | A supported wire profile defines deterministic duplicate handling that must be preserved. |

## LOCALIZED-DEC-004: Exact lookup and language matching separation

Status `resolved`; owner `localized maintainers`; classification
`interoperability policy`; decision scope `application-policy`; specification
`RFC 4647 language matching`; version `RFC 4647`; source authority
`rfc4647-source`; authority URL https://www.rfc-editor.org/rfc/rfc4647.txt;
section `Sections 3 and 3.4`; requirement strength `MAY`.

| Field | Decision |
|---|---|
| Issue | Exact locale identity, RFC lookup, and implementation-defined best-fit matching produce observably different results. |
| Credible interpretations | Match implicitly during every lookup. Or: keep exact lookup separate and expose matching as an explicit operation. |
| Known peer behavior | `golang.org/x/text/language.Matcher` provides maintained best-fit behavior, but it is the pinned implementation dependency rather than independent peer evidence. |
| Selected behavior | `Text.Get` and `Text.Has` are exact-only; `match.Best` performs explicit weighted matching against the current supported locale set. |
| Rationale | Callers can distinguish stored identity from a selected alternative and cannot receive an invented translation implicitly. |
| Security consequences | Exact authorization or persistence lookups cannot be widened by matching. |
| Resource consequences | Candidate count is bounded and matching performs no network access. |
| Compatibility consequences | Matcher-data updates require governed vectors and compatibility review even when APIs compile unchanged. |
| Wire consequences | Matching returns the actual stored locale and text; it does not rewrite the stored map. |
| Executable evidence | `TestStandardsMatchingMatrix`; `TestBestDistinguishesExactMatchedEmptyAndMissing` |
| Public APIs | `Text.Get`; `Text.Has`; `match.Best`; `match.Result` |
| Documentation | `docs/specification-decisions.md`; `docs/semantics.md`; `docs/compatibility.md` |
| Upstream status | RFC 4647 errata and pinned matcher behavior are reviewed separately; no full RFC conformance claim is made for best-fit results. |
| Reconsider when | The package adopts a versioned RFC 4647 lookup profile or replaces its matcher. |

## LOCALIZED-DEC-005: Explicit fallback ownership

Status `resolved`; owner `localized maintainers`; classification `omission`;
decision scope `application-policy`; specification `RFC 4647 language matching`;
version `RFC 4647`; source authority `rfc4647-source`; authority URL
https://www.rfc-editor.org/rfc/rfc4647.txt; section `Section 3.4.1`;
requirement strength `not specified`.

| Field | Decision |
|---|---|
| Issue | Parent traversal, application fallback chains, and a final default locale are not one universal standards-defined policy. |
| Credible interpretations | Infer parents and a global default automatically. Or: require callers to construct a bounded immutable plan. |
| Known peer behavior | Locale libraries expose parent relationships but cannot select an application's fallback graph or default. |
| Selected behavior | Fallback occurs only through an explicit bounded plan whose candidate kinds, chains, and optional default are caller-owned. |
| Rationale | Explicit plans prevent ambient global locale policy and invented translations. |
| Security consequences | Cycles, duplicates, invalid candidates, and oversized plans fail before resolution. |
| Resource consequences | Plan size and traversal are bounded and immutable. |
| Compatibility consequences | Application fallback order remains explicit and independent of exact lookup and best-fit matching. |
| Wire consequences | Resolution returns an existing stored entry and never serializes an inferred translation. |
| Executable evidence | `TestPlanUsesLocaleParentsWithoutInventingEntries`; `TestPlanRejectsCyclesDuplicatesAndBounds` |
| Fuzz evidence | `FuzzFallbackPlan` |
| Public APIs | `match.NewPlan`; `match.Plan.Resolve`; `match.NewFallbackPlan` |
| Documentation | `docs/specification-decisions.md`; `docs/semantics.md`; `docs/cookbook.md` |
| Upstream status | No upstream standard defines this package's application fallback graph. |
| Reconsider when | A supported application profile supplies a normative fallback chain. |

## LOCALIZED-DEC-006: Accept-Language quality and wildcard policy

Status `resolved`; owner `localized maintainers`; classification `ambiguity`;
decision scope `transport-specific`; specification `RFC 9110 Accept-Language`;
version `RFC 9110`; source authority `rfc9110-source`; authority URL
https://www.rfc-editor.org/rfc/rfc9110.txt; section `Section 12.5.4`;
requirement strength `not specified`.

Additional authoritative source: `{"id":"rfc4647-source","version":"RFC 4647","url":"https://www.rfc-editor.org/rfc/rfc4647.txt","specifications":["RFC 4647 language matching"]}`

| Field | Decision |
|---|---|
| Issue | Equal quality values, zero quality, wildcard selection, duplicate canonical ranges, and unsupported ranges require deterministic local policy. |
| Credible interpretations | Reorder equal weights, treat zero as acceptable, or let wildcard select any map iteration result. Or: preserve input order, exclude zero, and select wildcard deterministically. |
| Known peer behavior | HTTP frameworks differ in normalization and wildcard policy; no maintained peer is asserted as normative. |
| Selected behavior | Parsing preserves stable order within equal weights, excludes zero-weight preferences from selection, rejects canonical duplicates, and lets wildcard choose the first canonical lexical stored locale. |
| Rationale | Stable ordering and lexical wildcard selection prevent runtime map order from affecting responses. |
| Security consequences | Header bytes, candidate count, grammar, parameters, and quality precision are bounded and validated. |
| Resource consequences | Parsing and selection are bounded and perform no I/O. |
| Compatibility consequences | Tie and wildcard results are documented behavior and require compatibility review if changed. |
| Wire consequences | No default locale is added; an empty or wholly unacceptable header yields a missing result. |
| Executable evidence | `TestParseAcceptLanguagePreservesStableWeightedPreferences`; `TestAcceptLanguageSelectSupportsBoundedWildcard`; `TestAcceptLanguageBoundaryMatrix` |
| Fuzz evidence | `FuzzParseAcceptLanguage` |
| Public APIs | `http.ParseAcceptLanguage`; `http.Select`; `http.ParseOptions` |
| Documentation | `docs/specification-decisions.md`; `docs/semantics.md`; `docs/api-integrations.md` |
| Upstream status | RFC 9110 and RFC 4647 errata are monitored; wildcard tie-breaking remains package policy. |
| Reconsider when | A supported HTTP profile mandates different quality, duplicate, or wildcard behavior. |

## LOCALIZED-DEC-007: Canonical JSON object mapping

Status `resolved`; owner `localized maintainers`; classification
`implementation-defined behavior`; decision scope `application-policy`;
specification `RFC 8259 JSON`; version `RFC 8259`; source authority
`rfc8259-source`; authority URL https://www.rfc-editor.org/rfc/rfc8259.txt;
section `Sections 4 and 8.1`; requirement strength `SHOULD`.

| Field | Decision |
|---|---|
| Issue | JSON objects are unordered and duplicate member behavior is unpredictable, while localized values need deterministic bytes and collision-safe keys. |
| Credible interpretations | Use generic map encoding and last-wins decoding. Or: sort canonical keys and reject every raw or canonical duplicate. |
| Known peer behavior | Generic JSON decoders vary in duplicate handling and do not enforce BCP 47 canonical-key collisions. |
| Selected behavior | Encoding emits one compact object with canonical locale keys in lexical order; strict decoding rejects duplicates, invalid UTF-8, invalid locales, trailing values, and limit violations. |
| Rationale | Deterministic bytes support stable hashes and fixtures while fail-closed decoding prevents member shadowing. |
| Security consequences | Duplicate, invalid, oversized, and ambiguous input fails without disclosing localized content. |
| Resource consequences | Input, locale count, locale bytes, text bytes, and total retained bytes are bounded. |
| Compatibility consequences | Canonical byte order and strict rejection are public compatibility commitments. |
| Wire consequences | Missing text is an absent member and present-empty text is an empty JSON string. |
| Executable evidence | `TestTextJSONIsCanonicalAndRoundTripsEmptyValues`; `TestTextJSONRejectsHostileAndAmbiguousInput` |
| Fuzz evidence | `FuzzDecodeJSON` |
| Public APIs | `EncodeJSON`; `DecodeJSON`; `Text.MarshalJSON`; `Text.UnmarshalJSON` |
| Documentation | `docs/specification-decisions.md`; `docs/semantics.md`; `docs/security.md` |
| Upstream status | RFC 8259 errata are monitored; the canonical ordering profile is package-owned. |
| Reconsider when | A versioned canonical JSON profile replaces this package-specific byte contract. |

## LOCALIZED-DEC-008: Legacy JSON compatibility boundary

Status `resolved`; owner `localized maintainers`; classification `optional behavior`;
decision scope `application-policy`; specification `RFC 8259 JSON`; version
`RFC 8259`; source authority `rfc8259-source`; authority URL
https://www.rfc-editor.org/rfc/rfc8259.txt; section `Section 3`;
requirement strength `MAY`.

| Field | Decision |
|---|---|
| Issue | Existing payloads may use JSON null or underscore locale separators even though neither is the canonical localized representation. |
| Credible interpretations | Accept legacy forms everywhere. Reject them everywhere. Or: isolate them behind an explicit decode mode. |
| Known peer behavior | Legacy application serializers commonly emit these forms, but no maintained-peer conformance claim is made. |
| Selected behavior | Strict decoding rejects null and underscore locale separators; `PermissiveJSON` accepts them only at the explicit migration boundary and re-encodes canonical output. |
| Rationale | Compatibility input does not silently weaken the default parser or persist a second wire identity. |
| Security consequences | Permissive mode retains all ordinary duplicate, UTF-8, locale, and resource checks. |
| Resource consequences | Legacy normalization is bounded and performs no I/O. |
| Compatibility consequences | Callers must opt into legacy acceptance and can migrate stored data by canonical re-encoding. |
| Wire consequences | Output never emits null or underscore locale separators. |
| Executable evidence | `TestDecodeJSONPermissiveModeAcceptsNullAndLegacySeparators`; `TestDecodeJSONBoundaryFailures` |
| Fuzz evidence | `FuzzDecodeJSON` |
| Public APIs | `JSONMode`; `StrictJSON`; `PermissiveJSON`; `DecodeOptions` |
| Documentation | `docs/specification-decisions.md`; `docs/migration.md`; `docs/semantics.md` |
| Upstream status | This is a package migration policy, not an RFC 8259 or BCP 47 requirement. |
| Reconsider when | Legacy payloads are retired or a versioned wire format requires continued acceptance. |

## Unresolved decisions

None. New ambiguity, errata, registry drift, matching policy, or wire behavior
must be registered before it changes observable behavior.
