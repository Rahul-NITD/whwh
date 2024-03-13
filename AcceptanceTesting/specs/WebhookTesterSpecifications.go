package specs

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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

	hook_url, close := spinEchoHook()
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

	wres, err := extractBodyFromRes(whwhres)
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
		return nil, fmt.Errorf("Error while executing request, %w", err)
	}
	return extractBodyFromRes(res)
}

func extractBodyFromRes(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while reading response body %w", err)
	}
	return resBody, nil
}

func spinEchoHook() (url string, close func()) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		fmt.Fprint(w, r)
	}))

	return svr.URL, svr.Close

}
