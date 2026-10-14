// Package localized provides immutable UTF-8 text keyed by canonical BCP 47 tags.
// Locale identity and registry policy belong to the pinned International locale
// dependency. Exact lookup, preference matching and configured fallback are
// separate operations; missing, present-empty and nullable wrapper state differ.
// Core preserves text without normalization, pluralization, interpolation, message
// formatting, language detection or translation. No CLDR conformance is claimed.
//
// Containers are snapshotted and accessors expose copies or immutable values.
// Builders, nullable wrappers and callback implementations remain caller-owned.
// CPU work and callbacks are synchronous without context deadlines or workers;
// limits count content, not all heap allocation or callback execution time.
package localized

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"iter"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/faustbrian/go-international/v4/locale"
)

// Limits bounds construction work and retained content.
type Limits struct {
	// MaxLocales caps supplied entries before duplicate reduction; default 128.
	MaxLocales int
	// MaxTagBytes caps canonical tag bytes; default 255. Raw pair/JSON boundaries also
	// have an independent 255-byte parser ceiling.
	MaxTagBytes int
	// MaxTextBytes caps UTF-8 bytes in each supplied text; default 1 MiB, not runes or
	// graphemes.
	MaxTextBytes int
	// MaxTotalBytes caps summed supplied text bytes including discarded duplicates;
	// default 8 MiB. Tags and total heap are not counted.
	MaxTotalBytes int
}

// DuplicatePolicy selects construction behavior after locale canonicalization.
type DuplicatePolicy uint8

const (
	// RejectDuplicates fails construction when canonical identities repeat.
	RejectDuplicates DuplicatePolicy = iota
	// FirstWins preserves the first canonical occurrence.
	FirstWins
	// LastWins preserves the final canonical occurrence.
	LastWins
)

// ConstructionOptions configures bounded construction behavior.
type ConstructionOptions struct {
	// Limits uses defaults only when the whole struct is zero; partially specified zero
	// fields stay zero allowances.
	Limits Limits
	// Duplicates selects handling of repeated canonical identities; zero rejects them.
	Duplicates DuplicatePolicy
	// Locales rejects selected otherwise valid locale classes without changing text.
	Locales LocalePolicy
}

// LocalePolicy explicitly rejects selected otherwise-valid BCP 47 classes.
// Its zero value accepts und, mul, and private-use tags.
type LocalePolicy struct {
	// RejectUnd rejects und and tags beginning und-.
	RejectUnd bool
	// RejectMul rejects mul and tags beginning mul-.
	RejectMul bool
	// RejectPrivateUse rejects any canonical tag containing private-use components.
	RejectPrivateUse bool
}

const (
	defaultMaxLocales    = 128
	defaultMaxTagBytes   = 255
	defaultMaxTextBytes  = 1 << 20
	defaultMaxTotalBytes = 8 << 20
)

// DefaultLimits returns explicit construction defaults, not universal process
// memory/time or application policy limits.
func DefaultLimits() Limits {
	return Limits{
		MaxLocales: defaultMaxLocales, MaxTagBytes: defaultMaxTagBytes,
		MaxTextBytes: defaultMaxTextBytes, MaxTotalBytes: defaultMaxTotalBytes,
	}
}

func constructionLimits(limits Limits) (Limits, error) {
	if limits == (Limits{}) {
		limits = DefaultLimits()
	}
	if limits.MaxLocales < 0 || limits.MaxTagBytes < 0 || limits.MaxTextBytes < 0 || limits.MaxTotalBytes < 0 {
		return Limits{}, fmt.Errorf("%w: negative limit", ErrLimitExceeded)
	}
	return limits, nil
}

// Entry associates text with a BCP 47 locale.
type Entry struct {
	// Locale is a typed locale canonicalized at construction.
	Locale locale.Tag
	// Text is UTF-8 content preserved without normalization; empty is present, not
	// missing.
	Text string
}

// Text is an immutable collection of UTF-8 strings keyed by locale. Its zero
// value is an empty collection and is safe for concurrent reads.
type Text struct {
	entries []Entry
	index   map[string]int
}

// NewText constructs Text, rejecting duplicate canonical locale identities.
func NewText(entries ...Entry) (Text, error) {
	return NewTextWithOptions(ConstructionOptions{}, entries...)
}

// NewTextWithLimits constructs Text using explicit resource limits.
func NewTextWithLimits(limits Limits, entries ...Entry) (Text, error) {
	return NewTextWithOptions(ConstructionOptions{Limits: limits}, entries...)
}

// NewTextWithOptions snapshots entries with explicit duplicate/locale/limit policy.
// All supplied entries, including discarded duplicates, consume count/text budgets.
// Caller containers must not be mutated during the call. Negative bounds reject;
// only wholly zero Limits gains defaults. Failure returns zero Text, not partial
// content. Text strings are immutable values, not Unicode-normalized copies.
func NewTextWithOptions(options ConstructionOptions, entries ...Entry) (Text, error) {
	if options.Duplicates > LastWins {
		return Text{}, ErrInvalidPolicy
	}
	limits, err := constructionLimits(options.Limits)
	if err != nil {
		return Text{}, err
	}
	if len(entries) == 0 {
		return Text{}, nil
	}
	if len(entries) > limits.MaxLocales {
		return Text{}, fmt.Errorf("%w: locale count", ErrLimitExceeded)
	}

	byLocale := make(map[string]Entry, len(entries))
	total := 0
	for _, entry := range entries {
		canonical, err := entry.Locale.Canonical()
		if err != nil {
			return Text{}, ErrInvalidLocale
		}
		entry.Locale = canonical
		key := canonical.String()
		if (options.Locales.RejectUnd && (key == "und" || strings.HasPrefix(key, "und-"))) ||
			(options.Locales.RejectMul && (key == "mul" || strings.HasPrefix(key, "mul-"))) ||
			(options.Locales.RejectPrivateUse && canonical.HasPrivateUse()) {
			return Text{}, ErrLocaleRejected
		}
		if len(key) > limits.MaxTagBytes {
			return Text{}, fmt.Errorf("%w: tag bytes", ErrLimitExceeded)
		}
		if len(entry.Text) > limits.MaxTextBytes {
			return Text{}, fmt.Errorf("%w: text bytes", ErrLimitExceeded)
		}
		if len(entry.Text) > limits.MaxTotalBytes-total {
			return Text{}, fmt.Errorf("%w: total text bytes", ErrLimitExceeded)
		}
		total += len(entry.Text)
		if !utf8.ValidString(entry.Text) {
			return Text{}, ErrInvalidUTF8
		}
		if _, exists := byLocale[key]; exists {
			switch options.Duplicates {
			case RejectDuplicates:
				return Text{}, ErrDuplicateLocale
			case FirstWins:
				continue
			}
		}
		byLocale[key] = entry
	}
	owned := make([]Entry, 0, len(byLocale))
	for _, entry := range byLocale {
		owned = append(owned, entry)
	}
	slices.SortFunc(owned, func(left, right Entry) int {
		return strings.Compare(left.Locale.String(), right.Locale.String())
	})
	return textFromCanonicalEntries(owned), nil
}

// textFromCanonicalEntries takes ownership of sorted, unique, validated entries.
// Subsets preserve those invariants and cannot exceed their admitted source, so
// they need a fresh lookup index rather than another construction-limit check.
func textFromCanonicalEntries(entries []Entry) Text {
	if len(entries) == 0 {
		return Text{}
	}
	index := make(map[string]int, len(entries))
	for i := range entries {
		index[entries[i].Locale.String()] = i
	}
	return Text{entries: entries, index: index}
}

// TextFromMap parses strict raw keys under default admission and snapshots the
// map. The caller must not mutate it during decoding; failure is content-free.
// Map iteration makes the first reported failure order unspecified.
func TextFromMap(values map[string]string) (Text, error) {
	if len(values) > defaultMaxLocales {
		return Text{}, fmt.Errorf("%w: locale count", ErrLimitExceeded)
	}
	entries := make([]Entry, 0, len(values))
	for raw, value := range values {
		if strings.ContainsAny(raw, "_ \t\r\n") {
			return Text{}, ErrInvalidLocale
		}
		if len(raw) > defaultMaxTagBytes {
			return Text{}, fmt.Errorf("%w: tag bytes", ErrLimitExceeded)
		}
		tag, err := locale.Parse(raw)
		if err != nil {
			return Text{}, ErrInvalidLocale
		}
		entries = append(entries, Entry{Locale: tag, Text: value})
	}
	return NewText(entries...)
}

// Len returns the number of present locales, including present-empty values.
func (t Text) Len() int { return len(t.entries) }

// IsEmpty reports whether no locales are present.
func (t Text) IsEmpty() bool { return len(t.entries) == 0 }

// Has reports whether locale is present. It performs no matching or fallback.
func (t Text) Has(tag locale.Tag) bool {
	_, ok := t.canonicalIndex(tag)
	return ok
}

// Get returns the exact locale value and whether it is present. It performs no
// matching or fallback, so an empty string with true is distinct from missing.
func (t Text) Get(tag locale.Tag) (string, bool) {
	i, ok := t.canonicalIndex(tag)
	if !ok {
		return "", false
	}
	return t.entries[i].Text, true
}

func (t Text) canonicalIndex(tag locale.Tag) (int, bool) {
	canonical, err := tag.Canonical()
	if err != nil {
		return 0, false
	}
	i, ok := t.index[canonical.String()]
	return i, ok
}

// Require returns the exact locale value or ErrMissingLocale.
func (t Text) Require(tag locale.Tag) (string, error) {
	value, ok := t.Get(tag)
	if !ok {
		return "", ErrMissingLocale
	}
	return value, nil
}

// Locales returns canonical locales in deterministic lexical order. The caller
// owns the returned slice.
func (t Text) Locales() []locale.Tag {
	locales := make([]locale.Tag, len(t.entries))
	for i := range t.entries {
		locales[i] = t.entries[i].Locale
	}
	return locales
}

// Entries returns a deterministic caller-owned copy of all entries.
func (t Text) Entries() []Entry { return append([]Entry(nil), t.entries...) }

// Pair is a canonical string-keyed construction and wire convenience.
type Pair struct {
	// Locale is a strict raw locale key parsed at construction; it need not already use
	// canonical casing.
	Locale string
	// Text is UTF-8 content, including meaningful present-empty values.
	Text string
}

// TextFromPairs parses strict BCP 47 pair keys and constructs an owned Text.
func TextFromPairs(pairs ...Pair) (Text, error) {
	return TextFromPairsWithOptions(ConstructionOptions{}, pairs...)
}

// TextFromPairsWithOptions parses strict keys and applies explicit construction
// policies. Raw keys keep a 255-byte parser ceiling independent of MaxTagBytes;
// the configured bound is also checked after canonicalization. Input is borrowed
// only for the call; output entries are owned and lexically ordered.
func TextFromPairsWithOptions(options ConstructionOptions, pairs ...Pair) (Text, error) {
	limits, err := constructionLimits(options.Limits)
	if err != nil {
		return Text{}, err
	}
	if len(pairs) > limits.MaxLocales {
		return Text{}, fmt.Errorf("%w: locale count", ErrLimitExceeded)
	}
	entries := make([]Entry, 0, len(pairs))
	for _, pair := range pairs {
		if strings.ContainsAny(pair.Locale, "_ \t\r\n") {
			return Text{}, ErrInvalidLocale
		}
		if len(pair.Locale) > defaultMaxTagBytes {
			return Text{}, fmt.Errorf("%w: tag bytes", ErrLimitExceeded)
		}
		tag, err := locale.Parse(pair.Locale)
		if err != nil {
			return Text{}, ErrInvalidLocale
		}
		entries = append(entries, Entry{Locale: tag, Text: pair.Text})
	}
	return NewTextWithOptions(options, entries...)
}

// Pairs returns deterministic caller-owned canonical string pairs.
func (t Text) Pairs() []Pair {
	pairs := make([]Pair, len(t.entries))
	for i, entry := range t.entries {
		pairs[i] = Pair{Locale: entry.Locale.String(), Text: entry.Text}
	}
	return pairs
}

// All returns a deterministic iterator over immutable locale and text values.
// The yield callback runs synchronously without library locks or panic recovery.
// Returning false stops iteration; its time and side effects are caller-owned.
func (t Text) All() iter.Seq2[locale.Tag, string] {
	return func(yield func(locale.Tag, string) bool) {
		for _, entry := range t.entries {
			if !yield(entry.Locale, entry.Text) {
				return
			}
		}
	}
}

// Builder accumulates mutable entries without construction admission until Build.
// It requires caller synchronization when shared; returned Text snapshots remain
// independent of future additions. Zero Builder uses default construction policy.
type Builder struct {
	options ConstructionOptions
	entries []Entry
}

// NewBuilder creates an isolated construction builder.
func NewBuilder(options ConstructionOptions) *Builder {
	return &Builder{options: options}
}

// Add appends a typed locale pair to the builder. A nil builder is a no-op.
func (b *Builder) Add(tag locale.Tag, text string) {
	if b == nil {
		return
	}
	b.entries = append(b.entries, Entry{Locale: tag, Text: text})
}

// AddString parses a strict BCP 47 key and appends it when valid.
func (b *Builder) AddString(rawLocale, text string) error {
	if b == nil {
		return ErrInvalidPolicy
	}
	if strings.ContainsAny(rawLocale, "_ \t\r\n") {
		return ErrInvalidLocale
	}
	tag, err := locale.Parse(rawLocale)
	if err != nil {
		return ErrInvalidLocale
	}
	b.Add(tag, text)
	return nil
}

// Build validates and snapshots current entries without clearing the builder.
// Nil builder produces empty Text. Feeding a builder is not bounded by Build limits.
func (b *Builder) Build() (Text, error) {
	if b == nil {
		return Text{}, nil
	}
	return NewTextWithOptions(b.options, b.entries...)
}

// Set reconstructs a copy with locale added/replaced using DefaultLimits. Custom
// construction policy is not retained; valid larger receivers can fail revalidation.
func (t Text) Set(tag locale.Tag, value string) (Text, error) {
	entries := t.Entries()
	canonical, err := tag.Canonical()
	if err != nil {
		return Text{}, ErrInvalidLocale
	}
	if i, ok := t.index[canonical.String()]; ok {
		entries[i].Text = value
	} else {
		entries = append(entries, Entry{Locale: canonical, Text: value})
	}
	return NewText(entries...)
}

// Remove removes one locale without reapplying construction limits to retained
// entries. Absent/invalid locale is a no-op. The receiver remains unchanged.
func (t Text) Remove(tag locale.Tag) Text {
	i, ok := t.canonicalIndex(tag)
	if !ok {
		return t
	}
	entries := make([]Entry, 0, len(t.entries))
	entries = append(entries, t.entries[:i]...)
	entries = append(entries, t.entries[i+1:]...)
	return textFromCanonicalEntries(entries)
}

// Filter returns entries accepted by predicate in deterministic order. A nil
// predicate returns the receiver unchanged. Callback calls are synchronous in
// canonical order, without forced cancellation or panic recovery. Retained
// entries are not rechecked against construction limits. The receiver remains
// unchanged; callback effects are not rolled back.
func (t Text) Filter(predicate func(locale.Tag, string) bool) Text {
	if predicate == nil {
		return t
	}
	entries := make([]Entry, 0, len(t.entries))
	for _, entry := range t.entries {
		if predicate(entry.Locale, entry.Text) {
			entries = append(entries, entry)
		}
	}
	return textFromCanonicalEntries(entries)
}

// Map runs synchronous transforms in canonical order, then revalidates using
// DefaultLimits, not the receiver's originating policy. Errors return zero Text
// without a partial value, but callback side effects are not rolled back. The
// callback is borrowed and owns cancellation, reentrancy and concurrent-call safety.
func (t Text) Map(transform func(locale.Tag, string) (string, error)) (Text, error) {
	if transform == nil {
		return Text{}, ErrInvalidPolicy
	}
	entries := t.Entries()
	for i := range entries {
		text, err := transform(entries[i].Locale, entries[i].Text)
		if err != nil {
			return Text{}, err
		}
		entries[i].Text = text
	}
	return NewText(entries...)
}

// Equal reports structural equality of locale and text pairs.
func (t Text) Equal(other Text) bool {
	if len(t.entries) != len(other.entries) {
		return false
	}
	for i := range t.entries {
		if t.entries[i].Locale.String() != other.entries[i].Locale.String() ||
			t.entries[i].Text != other.entries[i].Text {
			return false
		}
	}
	return true
}

// Hash returns a stable, length-framed SHA-256 digest of canonical entries.
func (t Text) Hash() [sha256.Size]byte {
	hash := sha256.New()
	var size [8]byte
	for _, entry := range t.entries {
		for _, value := range []string{entry.Locale.String(), entry.Text} {
			binary.BigEndian.PutUint64(size[:], uint64(len(value)))
			_, _ = hash.Write(size[:])
			_, _ = hash.Write([]byte(value))
		}
	}
	var sum [sha256.Size]byte
	copy(sum[:], hash.Sum(nil))
	return sum
}

// MergePolicy selects behavior for locales present in both values.
type MergePolicy uint8

const (
	// LeftWins preserves the receiver on conflicts.
	LeftWins MergePolicy = iota
	// RightWins replaces conflicts with the overlay value.
	RightWins
	// RejectConflict fails if any locale is present in both values.
	RejectConflict
	// ResolveConflict delegates overlapping values to a resolver.
	ResolveConflict
)

// EmptyPolicy determines whether present-empty participates in merging.
type EmptyPolicy uint8

const (
	// EmptyIsValue treats an empty string as a present conflicting value.
	EmptyIsValue EmptyPolicy = iota
	// EmptyIsAbsent omits empty strings before conflict resolution.
	EmptyIsAbsent
)

// MergeResolver synchronously returns text for a canonical conflict. It is
// borrowed, not made concurrency-safe or cancellable; errors propagate unchanged
// and may contain content. External side effects are not rolled back.
type MergeResolver func(locale locale.Tag, left, right string) (string, error)

// MergeOptions configures explicit conflict, empty, and output bounds.
type MergeOptions struct {
	// Conflicts selects overlapping-locale resolution; zero preserves the left value.
	Conflicts MergePolicy
	// Empty selects whether empty strings participate; zero treats them as present values.
	Empty EmptyPolicy
	// Resolver is borrowed synchronous caller code, required for ResolveConflict. Its side
	// effects are not transactional.
	Resolver MergeResolver
	// Limits bounds the final reconstructed output, not callback execution or all
	// temporary merge storage.
	Limits Limits
}

// Merge combines two values using an explicit conflict policy.
func (t Text) Merge(other Text, policy MergePolicy) (Text, error) {
	return t.MergeWithOptions(other, MergeOptions{Conflicts: policy})
}

// MergeWithOptions combines values with explicit conflict/empty/output policy.
// Resolver calls occur in the overlay's canonical order before final output
// validation; no partial Text is returned on failure, but callback effects remain.
// Limits bound final construction, not preliminary map allocation or callback time.
func (t Text) MergeWithOptions(other Text, options MergeOptions) (Text, error) {
	if options.Conflicts > ResolveConflict || options.Empty > EmptyIsAbsent {
		return Text{}, ErrInvalidPolicy
	}
	if options.Conflicts == ResolveConflict && options.Resolver == nil {
		return Text{}, ErrResolverRequired
	}
	entries := make(map[string]Entry)
	for _, entry := range t.entries {
		if options.Empty == EmptyIsAbsent && entry.Text == "" {
			continue
		}
		entries[entry.Locale.String()] = entry
	}
	for _, entry := range other.entries {
		if options.Empty == EmptyIsAbsent && entry.Text == "" {
			continue
		}
		key := entry.Locale.String()
		left, conflict := entries[key]
		if conflict {
			switch options.Conflicts {
			case LeftWins:
				continue
			case RejectConflict:
				return Text{}, ErrConflict
			case ResolveConflict:
				text, err := options.Resolver(entry.Locale, left.Text, entry.Text)
				if err != nil {
					return Text{}, err
				}
				entry.Text = text
			}
		}
		entries[key] = entry
	}
	merged := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		merged = append(merged, entry)
	}
	limits := options.Limits
	if limits == (Limits{}) {
		limits = DefaultLimits()
	}
	result, err := NewTextWithLimits(limits, merged...)
	if err != nil {
		return Text{}, err
	}
	return result, nil
}
