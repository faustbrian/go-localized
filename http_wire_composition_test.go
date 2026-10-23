package localized_test

import (
	"context"
	"net/http"
	"testing"

	client "github.com/faustbrian/go-http-client/v3"
	"github.com/faustbrian/go-international/v5/locale"
	localized "github.com/faustbrian/go-localized/v6"
	httpadapter "github.com/faustbrian/go-localized/v6/adapters/httpclient"
	wireadapter "github.com/faustbrian/go-localized/v6/adapters/wire"
	localizedhttp "github.com/faustbrian/go-localized/v6/http"
	localizedmatch "github.com/faustbrian/go-localized/v6/match"
	"github.com/faustbrian/go-wire/v3/jsonwire"
)

func TestPublicHTTPWireComposition(t *testing.T) {
	t.Parallel()
	english, err := locale.Parse("en-US")
	if err != nil {
		t.Fatal(err)
	}
	value, err := localized.NewText(localized.Entry{Locale: english, Text: "Hello"})
	if err != nil {
		t.Fatal(err)
	}
	spec, err := client.NewRequestSpec("https://example.test", "/content")
	if err != nil {
		t.Fatal(err)
	}
	for _, preferences := range []func(client.RequestSpec, client.RequestLayer, ...localizedmatch.Preference) (client.RequestSpec, error){
		httpadapter.WithPreferences,
	} {
		selected, err := preferences(spec, client.LayerRequest, localizedmatch.Preference{Locale: english, Weight: 1})
		if err != nil {
			t.Fatal(err)
		}
		request, err := selected.Build(context.Background(), http.MethodGet)
		if err != nil {
			t.Fatal(err)
		}
		if got := request.Header.Get("Accept-Language"); got != "en-US" {
			t.Fatalf("Accept-Language = %q", got)
		}
		original, err := spec.Build(context.Background(), http.MethodGet)
		if err != nil || original.Header.Get("Accept-Language") != "" {
			t.Fatalf("original request changed: %v", err)
		}
		for _, selectResponse := range []func(localized.Text, *http.Response, localizedhttp.ParseOptions) (localizedmatch.Result, error){
			httpadapter.SelectResponse,
		} {
			result, err := selectResponse(value, &http.Response{Request: request}, localizedhttp.ParseOptions{})
			if err != nil || result.Text != "Hello" || result.Locale.String() != "en-US" {
				t.Fatalf("selected response = %+v, %v", result, err)
			}
		}
	}
	for _, encode := range []func(localized.Text, jsonwire.EncodeOptions) ([]byte, error){
		wireadapter.EncodeJSON,
	} {
		encoded, err := encode(value, jsonwire.EncodeOptions{})
		if err != nil || string(encoded) != `{"en-US":"Hello"}` {
			t.Fatalf("canonical JSON = %q, %v", encoded, err)
		}
		for _, decode := range []func([]byte, jsonwire.DecodeOptions) (localized.Text, error){
			wireadapter.DecodeJSON,
		} {
			decoded, err := decode(encoded, jsonwire.DecodeOptions{})
			if err != nil || !decoded.Equal(value) {
				t.Fatalf("decoded value = %v, %v", decoded.Entries(), err)
			}
		}
	}
}
