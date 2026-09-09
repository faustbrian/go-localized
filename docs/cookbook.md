# Cookbook

## Regional fallback

Create a `match.Plan` from `fr-CA` to exact `fr` and then an application
default. Keep the default in `PlanOptions`, not a global.

## Script-sensitive Chinese

Use `ParentRange` from `zh-Hant-TW`; locale-layer parents can select `zh-Hant`.
Do not strip `-TW` or `-Hant` manually.

## Required languages

Check `Has` for each application-required tag after construction. Return a
domain validation error rather than embedding required-language policy in
`Text`.

## Non-empty validated text

```go
err := validation.Validate(value,
    validation.RequireNonWhitespace(),
    validation.MaxBytes(4<<10),
    validation.MaxLines(4),
    validation.NoControlCharacters(),
)
```

## Explicit NFC

Call `validation.Normalize(value, validation.NFC)` and retain
the returned value. The original remains unchanged.

## Resolver merge

Use `ResolveConflict` with a deterministic callback. The callback MUST NOT
perform remote translation or mutate external state if callers require
repeatable merge results.

## SQL nullable values

Use `postgres.NewText(value)` for non-null values. A zero `postgres.Text`
encodes SQL NULL; `localized.Text{}` encodes `{}`.

## Wire formats

Use `adapters/wire` functions with explicit wire-format options. YAML, TOML, and
MessagePack decode into a string map and then revalidate canonical locale keys.

## Tests

```go
value := localizedtest.New(t).Add("en", "Hello").Add("fi", "Hei").Build()
english, _ := locale.Parse("en")
localizedtest.AssertExact(t, value, english, "Hello")
```
