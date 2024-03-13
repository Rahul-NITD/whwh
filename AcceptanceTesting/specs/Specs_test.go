package specs_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aargeee/whwh/AcceptanceTesting/specs"
	"github.com/stretchr/testify/assert"
)

func TestEchoHook(t *testing.T) {
	url, cancel := specs.SpinHelperHook()
	defer cancel()
	req, err := http.NewRequest(http.MethodPost, url, http.NoBody)
	assert.NoError(t, err, "Error in building request")

	MakeRequestAssertString(t, req.Clone(context.Background()), "REGISTERED")
	MakeRequestAssertString(t, req.Clone(context.Background()), "SUCCESS")
}

func MakeRequestAssertString(t *testing.T, req *http.Request, assertion string) {
	res, err := http.DefaultClient.Do(req.Clone(context.Background()))
	assert.NoError(t, err, "Error in executing request")
	defer res.Body.Close()
	body, err := specs.ExtractBodyFromRes(res)
	assert.NoError(t, err, "Error extracting body")
	assert.Equal(t, assertion, string(body))
}
