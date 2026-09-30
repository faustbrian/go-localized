# Security policy

## Supported versions

The latest stable v1 release is supported. Security fixes are developed on
`main` and published as compatible v1 patches when the public contract permits.
Older v1 versions should upgrade to the latest patch; no separate maintenance
branch is promised.

## Reporting

Report suspected vulnerabilities privately through GitHub Security Advisories.
Do not include production localized content, private-use locale tags, access
tokens, database URLs, or customer identifiers in a public issue.

Include the affected version, operation, bounded reproduction, and expected
versus observed behavior. Maintainers will acknowledge a report, assess impact,
prepare tests and a fix, and coordinate disclosure. The acknowledgement target
is three business days. Severity considers confidentiality, integrity,
availability, attacker control, and required deployment configuration.

Confirmed critical and high issues take priority over ordinary work. A medium
issue must be fixed or explicitly accepted with an owner and mitigation.
Maintainers coordinate affected versions, a regression test, release notes,
upgrade guidance, and a GitHub advisory with the reporter. Exploit details
remain private until a fix is available or an agreed embargo ends. A release
or advisory is not implied by an unconfirmed scanner result.

The detailed threat model is in [docs/security.md](docs/security.md).
