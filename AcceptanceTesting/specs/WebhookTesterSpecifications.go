package specs

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Client interface{}

type WebhookTesterSubject interface {
	CreateChannel(t *testing.T) (chanID string)
	ConnectClientAndServer(t *testing.T, chanID string) (client Client)
	Dispatch(t *testing.T, req *http.Request, hook_url string) (res *http.Response)
}

func TestWebhookTester(t *testing.T, subject WebhookTesterSubject) {
	hook_url, close := SpinHelperHook()
	defer close()

	chanID := subject.CreateChannel(t)
	t.Skip()
	subject.ConnectClientAndServer(t, chanID)
	req := mustBuildRequest(t, hook_url, postNilBody)
	assert.Equal(t, "REGISTERED", RequestAndExtract(t, req))
	assert.Equal(t, "SUCCESS", DispatchAndExtract(t, req, subject, hook_url))
}

func mustBuildRequest(t *testing.T, url string, builder func(string) (*http.Request, error)) *http.Request {
	req, err := buildRequest(url, builder)
	assert.NoError(t, err, "Could not build request")
	return req
}

func RequestAndExtract(t *testing.T, req *http.Request) string {
	hres, err := request(req.Clone(context.Background()))
	assert.NoError(t, err, "Error while executing request")
	return string(hres)
}

func DispatchAndExtract(t *testing.T, req *http.Request, subject WebhookTesterSubject, hook_url string) string {
	whwhres := subject.Dispatch(t, req.Clone(context.Background()), hook_url)
	wres, err := ExtractBodyFromRes(whwhres)
	assert.NoError(t, err, "Error while extracting body from res")
	return string(wres)
}

func postNilBody(url string) (*http.Request, error) {
	return http.NewRequest(http.MethodPost, url, nil)
}

func buildRequest(url string, builder func(string) (*http.Request, error)) (*http.Request, error) {
	return builder(url)
}

func request(req *http.Request) ([]byte, error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while executing request, %w", err)
	}
	return ExtractBodyFromRes(res)
}

func ExtractBodyFromRes(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response body %w", err)
	}
	return resBody, nil
}

func SpinHelperHook() (url string, close func()) {
	var rgot *http.Request
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if rgot != nil {

			if reflect.DeepEqual(rgot, r.Clone(context.Background())) {
				fmt.Fprint(w, "SUCCESS")
			} else {
				fmt.Fprintln(w, "FAILURE")
				fmt.Fprintln(w, *rgot)
				fmt.Fprintln(w, *r)
			}

		} else {
			rgot = r.Clone(context.Background())
			fmt.Fprint(w, "REGISTERED")
		}
	}))

	return svr.URL, svr.Close

}
