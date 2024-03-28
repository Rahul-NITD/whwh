package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aargeee/whwh/handlers"
	"github.com/stretchr/testify/assert"
)

func TestCreateChannel(t *testing.T) {
	t.Run("/create POST request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/create", http.NoBody)
		assert.NoError(t, err, "Could not make request to /create")
		res := httptest.NewRecorder()
		handlers.NewServer().ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
	})

}
