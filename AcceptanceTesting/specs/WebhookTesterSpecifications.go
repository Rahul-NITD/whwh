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

type Client interface {
}

type WebhookTesterSubject interface {
	CreateChannel() (chanID string, err error)
	ConnectClientAndServer(chanID string) (client Client, err error)
	Dispatch(req *http.Request, hook_url string) (res *http.Response, err error)
}

func TestWebhookTester(t *testing.T, subject WebhookTesterSubject) {

	reqTemplate := postNilBody

	hook_url, close := SpinHelperHook()
	defer close()

	chanID, err := subject.CreateChannel()
	assert.NoError(t, err, "Error while creating Channel")

	_, err = subject.ConnectClientAndServer(chanID)
	assert.NoError(t, err, "Couldn't Connect client to server")

	req, err := buildRequest(hook_url, reqTemplate)
	assert.NoError(t, err, "Could not build request using given template")

	hres, err := request(req.Clone(context.Background()), reqTemplate)
	assert.NoError(t, err, "Error while executing request")

	whwhres, err := subject.Dispatch(req.Clone(context.Background()), hook_url)
	assert.NoError(t, err, "Error while subject dispatch")

	wres, err := ExtractBodyFromRes(whwhres)
	assert.NoError(t, err, "Error while extracting body from res")

	assert.Equal(t, hres, wres)
}

func postNilBody(url string) (*http.Request, error) {
	return http.NewRequest(http.MethodPost, url, nil)
}

func buildRequest(url string, builder func(string) (*http.Request, error)) (*http.Request, error) {
	return builder(url)
}

func request(req *http.Request, builder func(string) (*http.Request, error)) ([]byte, error) {
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
