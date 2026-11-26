# Upstream review history

Earlier reviews are recorded in the authority review sections of the
[specification decision register](../docs/specification-decisions.md).

## 2026-11-26: errata and registry monitoring review

The two authorities whose digests no longer matched their monitoring pins were
retrieved twice over HTTPS with the same request profile as the online
specification check. Each pair of retrievals produced the same digest. The
RFC 5646, RFC 4647, RFC 9110 and RFC 8259 texts, their other errata pages and
the IANA XML release index still match their pins. The thirty-day cadence and
all authority URLs remain unchanged.

| Authority | Previous SHA-256 | Current SHA-256 |
| --- | --- | --- |
| `rfc9110-errata` | `cec32fd170146656d933f627b512f2e027ae5c3592f5ec7760c3627493b30505` | `bdd312bd1c1dc5b04ac86286bf6f1eca35fb70371c3900d8638fe98ba16c963f` |
| `iana-language-registry` | `755fad43283be7b41ebe3c89ad054b6eaf928f404f9c0edb74799e0eab74beb1` | `8ac4e58567593c96e9e6d66e2017dbb8d3953dbcb42fe32ecd3b3fccdfdaba7b` |

The RFC 9110 errata page now lists six verified and eight rejected reports and
no open reports. Errata 9162 (field-line combination, Section 5.2) and 9164
(collected ABNF rendering, Appendix A), both reviewed as reported and
behavior-neutral on 2026-10-04, were rejected on 2026-11-13. No new erratum
was filed. Their rejection confirms the existing reading:
`LOCALIZED-DEC-006` keeps its quality, ordering, duplicate, wildcard and
resource-bound contracts, and the Accept-Language adapter still parses one
caller-supplied value.

The IANA Language Subtag Registry now has File-Date 2026-11-18. A public
archive copy from 2026-10-04 reproduces the previous pinned digest exactly
(File-Date 2026-09-17). Compared with it, the current registry changes only
the File-Date line and adds two variant subtags, both on 2026-11-18 with
prefixes `uz` and `uz-Latn`: `1995` (Uzbek Latin orthography of 1995) and
`2026` (Uzbek Latin orthography of 2026). No existing record changed, so no
subtag was deprecated or given a new preferred value. Monitoring a newer registry does not adopt it into runtime
data: `LOCALIZED-DEC-001` and `LOCALIZED-DEC-002` follow the pinned
International and x/text locale dataset, whose provenance remains registry
version 2026-06-14. Canonicalization, tag acceptance and matching are
unchanged; recognising the new Uzbek variants would require a reviewed locale
dependency update.

The `iana-language-registry` monitoring version and its entry in the decision
register now read "Registry File-Date 2026-11-18" so the label matches the
pinned content. Neither change alters a decision, conformance binding,
generated table or runtime behavior.
