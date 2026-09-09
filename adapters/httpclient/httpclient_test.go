package httpclient_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	client "github.com/faustbrian/go-http-client"
	"github.com/faustbrian/go-international/locale"
	localized "github.com/faustbrian/go-localized"
	adapter "github.com/faustbrian/go-localized/adapters/httpclient"
	localizedhttp "github.com/faustbrian/go-localized/http"
	legacy "github.com/faustbrian/go-localized/localizedhttpclient"
	localizedmatch "github.com/faustbrian/go-localized/match"
)

func TestCanonicalHTTPClientAdapterPreservesNegotiationAndIdentity(t *testing.T) {
	t.Parallel()

	spec, err := client.NewRequestSpec("https://example.test", "/content")
	if err != nil {
		t.Fatal(err)
	}
	english, _ := locale.Parse("EN-us")
	spec, err = adapter.WithPreferences(spec, client.LayerRequest,
		localizedmatch.Preference{Locale: english, Weight: 1},
	)
	if err != nil {
		t.Fatal(err)
	}
	request, err := spec.Build(context.Background(), http.MethodGet)
	if err != nil {
		t.Fatal(err)
	}
	value, _ := localized.TextFromMap(map[string]string{"en-US": "Hello"})
	result, err := adapter.SelectResponse(value, &http.Response{Request: request}, localizedhttp.ParseOptions{})
	if err != nil || result.Text != "Hello" {
		t.Fatalf("SelectResponse() = %+v, %v", result, err)
	}
	if reflect.TypeOf(adapter.ErrInvalidResponse) != reflect.TypeOf(legacy.ErrInvalidResponse) {
		t.Fatal("compatibility Error identity changed")
	}
}
