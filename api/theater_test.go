package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PRPO-skupina-02/common/database"
	"github.com/PRPO-skupina-02/common/xtesting"
	"github.com/PRPO-skupina-02/spored/db"
	"github.com/PRPO-skupina-02/spored/models"
	"github.com/stretchr/testify/assert"
)

func TestTheatersList(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	r := TestingRouter(t, db)

	tests := []struct {
		name   string
		status int
		params string
	}{
		{
			name:   "ok",
			status: http.StatusOK,
		},
		{
			name:   "ok-paginated",
			status: http.StatusOK,
			params: "?limit=1&offset=1",
		},
		{
			name:   "ok-sort",
			status: http.StatusOK,
			params: "?sort=-updated_at",
		},
		{
			name:   "ok-paginated-sort",
			status: http.StatusOK,
			params: "?limit=2&offset=1&sort=updated_at",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/theaters%s", testCase.params)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodGet, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
		})
	}
}

func TestTheatersCreate(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	r := TestingRouter(t, db)

	tests := []struct {
		name   string
		body   TheaterRequest
		status int
	}{
		{
			name: "ok",
			body: TheaterRequest{
				Name: "TestTheater",
			},
			status: http.StatusCreated,
		},
		{
			name: "short-name",
			body: TheaterRequest{
				Name: "A",
			},
			status: http.StatusBadRequest,
		},
		{
			name:   "no-body",
			status: http.StatusBadRequest,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := "/api/v1/theaters"

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodPost, testCase.body)
			assert.NoError(t, err)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			ignoreResp := xtesting.ValuesCheckers{
				"id":         xtesting.ValueUUID(),
				"created_at": xtesting.ValueTimeInPastDuration(time.Second),
				"updated_at": xtesting.ValueTimeInPastDuration(time.Second),
			}

			ignoreTheaters := xtesting.GenerateValueCheckersForArrays(map[string]xtesting.ValueChecker{"ID": xtesting.ValueUUID(), "CreatedAt": xtesting.ValueTime(), "UpdatedAt": xtesting.ValueTime()}, 10)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w, ignoreResp)
			xtesting.AssertGoldenDatabaseTable(t, db.Order("name"), []models.Theater{}, ignoreTheaters)
		})
	}
}

func TestTheatersShow(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	r := TestingRouter(t, db)

	tests := []struct {
		name   string
		status int
		id     string
	}{
		{
			name:   "ok",
			status: http.StatusOK,
			id:     "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:   "invalid-id",
			status: http.StatusNotFound,
			id:     "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name:   "nil-id",
			status: http.StatusBadRequest,
			id:     "00000000-0000-0000-0000-000000000000",
		},
		{
			name:   "malformed-id",
			status: http.StatusBadRequest,
			id:     "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/theaters/%s", testCase.id)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodGet, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
		})
	}
}

func TestTheatersUpdate(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	r := TestingRouter(t, db)

	tests := []struct {
		name   string
		body   TheaterRequest
		status int
		id     string
	}{
		{
			name: "ok",
			body: TheaterRequest{
				Name: "NewTheater",
			},
			status: http.StatusOK,
			id:     "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "short-name",
			body: TheaterRequest{
				Name: "A",
			},
			status: http.StatusBadRequest,
			id:     "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:   "no-body",
			status: http.StatusBadRequest,
			id:     "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "invalid-id",
			body: TheaterRequest{
				Name: "NewTheater",
			},
			status: http.StatusNotFound,
			id:     "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name: "nil-id",
			body: TheaterRequest{
				Name: "NewTheater",
			},
			status: http.StatusBadRequest,
			id:     "00000000-0000-0000-0000-000000000000",
		},
		{
			name: "malformed-id",
			body: TheaterRequest{
				Name: "NewTheater",
			},
			status: http.StatusBadRequest,
			id:     "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/theaters/%s", testCase.id)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodPut, testCase.body)
			assert.NoError(t, err)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			ignoreResp := xtesting.ValuesCheckers{
				"updated_at": xtesting.ValueTimeInPastDuration(time.Second),
			}

			ignoreTheaters := xtesting.GenerateValueCheckersForArrays(map[string]xtesting.ValueChecker{"UpdatedAt": xtesting.ValueTime()}, 10)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w, ignoreResp)
			xtesting.AssertGoldenDatabaseTable(t, db, []models.Theater{}, ignoreTheaters)
		})
	}
}

func TestTheatersDelete(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	r := TestingRouter(t, db)

	tests := []struct {
		name   string
		status int
		id     string
	}{
		{
			name:   "ok",
			status: http.StatusNoContent,
			id:     "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:   "invalid-id",
			status: http.StatusNotFound,
			id:     "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name:   "nil-id",
			status: http.StatusBadRequest,
			id:     "00000000-0000-0000-0000-000000000000",
		},
		{
			name:   "malformed-id",
			status: http.StatusBadRequest,
			id:     "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/theaters/%s", testCase.id)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodDelete, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
			xtesting.AssertGoldenDatabaseTable(t, db, []models.Theater{}, nil)
			xtesting.AssertGoldenDatabaseTable(t, db, []models.Room{}, nil)
		})
	}
}
