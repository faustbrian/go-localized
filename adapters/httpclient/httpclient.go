package httpclient

import (
	"net/http"

	client "github.com/faustbrian/go-http-client"
	localized "github.com/faustbrian/go-localized"
	localizedhttp "github.com/faustbrian/go-localized/http"
	legacy "github.com/faustbrian/go-localized/localizedhttpclient"
	localizedmatch "github.com/faustbrian/go-localized/match"
)

// Error preserves the released localizedhttpclient error identity.
type Error = legacy.Error

// ErrInvalidResponse reports a nil response or originating request.
const ErrInvalidResponse = legacy.ErrInvalidResponse

// WithPreferences returns a request spec with a canonical Accept-Language
// header at layer. An empty preference list removes the header at that layer.
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
