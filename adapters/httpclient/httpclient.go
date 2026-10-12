package httpclient

import (
	"net/http"

	client "github.com/faustbrian/go-http-client/v2"
	localized "github.com/faustbrian/go-localized/v5"
	localizedhttp "github.com/faustbrian/go-localized/v5/http"
	legacy "github.com/faustbrian/go-localized/v5/localizedhttpclient"
	localizedmatch "github.com/faustbrian/go-localized/v5/match"
)

// Error preserves the released localizedhttpclient error identity.
type Error = legacy.Error

// ErrInvalidResponse reports a nil response or originating request.
const ErrInvalidResponse = legacy.ErrInvalidResponse

// WithPreferences returns a request spec with a canonical Accept-Language
// header at layer. Empty preferences remove that layer's header. The full header
// is built before admission; derived size bounds are not inbound defaults. Bound
// untrusted preferences before calling; this helper sends no request.
func WithPreferences(
	spec client.RequestSpec,
	layer client.RequestLayer,
	preferences ...localizedmatch.Preference,
) (client.RequestSpec, error) {
	return legacy.WithPreferences(spec, layer, preferences...)
}

// SelectResponse selects a localized value from the Accept-Language header on
// the response's originating request. It applies matching but no fallback.
func SelectResponse(
	value localized.Text,
	response *http.Response,
	options localizedhttp.ParseOptions,
) (localizedmatch.Result, error) {
	return legacy.SelectResponse(value, response, options)
}
