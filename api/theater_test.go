package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/orgs/PRPO-skupina-02/Spored/xtesting"
	"github.com/stretchr/testify/assert"
)

func TestTheatersList(t *testing.T) {
	db, fixtures := xtesting.PrepareTestDatabase(t)
	r := TestingRouter(t, db)

	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "ok",
			status: http.StatusOK,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := "/api/v1/theaters"

			req, err := http.NewRequest(http.MethodGet, targetURL, nil)
			assert.NoError(t, err)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
		})
	}
}
