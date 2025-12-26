package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PRPO-skupina-02/common/database"
	"github.com/PRPO-skupina-02/common/xtesting"
	"github.com/orgs/PRPO-skupina-02/Spored/db"
	"github.com/orgs/PRPO-skupina-02/Spored/models"
	"github.com/stretchr/testify/assert"
)

func TestMoviesList(t *testing.T) {
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
			params: "?sort=-title",
		},
		{
			name:   "ok-paginated-sort",
			status: http.StatusOK,
			params: "?limit=2&offset=1&sort=title",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/movies%s", testCase.params)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodGet, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
		})
	}
}

func TestMoviesCreate(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	r := TestingRouter(t, db)

	tests := []struct {
		name   string
		body   MovieRequest
		status int
	}{
		{
			name: "ok",
			body: MovieRequest{
				Title:         "TestMovie",
				Description:   "New Description",
				ImageURL:      "http://example.com/image.png",
				Rating:        7.6666,
				LengthMinutes: 125,
			},
			status: http.StatusCreated,
		},
		{
			name: "validation-errors",
			body: MovieRequest{
				Title:         "A",
				Description:   "B",
				ImageURL:      "randomText",
				Rating:        12,
				LengthMinutes: 5,
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

			targetURL := "/api/v1/movies"

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodPost, testCase.body)
			assert.NoError(t, err)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			ignoreResp := xtesting.ValuesCheckers{
				"id":         xtesting.ValueUUID(),
				"created_at": xtesting.ValueTimeInPastDuration(time.Second),
				"updated_at": xtesting.ValueTimeInPastDuration(time.Second),
			}

			ignoreMovies := xtesting.GenerateValueCheckersForArrays(map[string]xtesting.ValueChecker{"ID": xtesting.ValueUUID(), "CreatedAt": xtesting.ValueTime(), "UpdatedAt": xtesting.ValueTime()}, 10)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w, ignoreResp)
			xtesting.AssertGoldenDatabaseTable(t, db.Order("title"), []models.Movie{}, ignoreMovies)
		})
	}
}

func TestMoviesShow(t *testing.T) {
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
			id:     "510633ca-e23f-11f0-a626-d3b8771e2cb9",
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

			targetURL := fmt.Sprintf("/api/v1/movies/%s", testCase.id)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodGet, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
		})
	}
}

func TestMoviesUpdate(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	r := TestingRouter(t, db)

	tests := []struct {
		name   string
		body   MovieRequest
		status int
		id     string
	}{
		{
			name: "ok",
			body: MovieRequest{
				Title:         "TestMovie",
				Description:   "New Description",
				ImageURL:      "http://example.com/image.png",
				Rating:        7.6666,
				LengthMinutes: 125,
			},
			status: http.StatusOK,
			id:     "510633ca-e23f-11f0-a626-d3b8771e2cb9",
		},
		{
			name: "validation-errors",
			body: MovieRequest{
				Title:         "A",
				Description:   "B",
				ImageURL:      "randomText",
				Rating:        12,
				LengthMinutes: 5,
			},
			status: http.StatusBadRequest,
			id:     "510633ca-e23f-11f0-a626-d3b8771e2cb9",
		},
		{
			name:   "no-body",
			status: http.StatusBadRequest,
			id:     "510633ca-e23f-11f0-a626-d3b8771e2cb9",
		},
		{
			name: "invalid-id",
			body: MovieRequest{
				Title:         "TestMovie",
				Description:   "New Description",
				ImageURL:      "http://example.com/image.png",
				Rating:        7.6666,
				LengthMinutes: 125,
			},
			status: http.StatusNotFound,
			id:     "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name: "nil-id",
			body: MovieRequest{
				Title:         "TestMovie",
				Description:   "New Description",
				ImageURL:      "http://example.com/image.png",
				Rating:        7.6666,
				LengthMinutes: 125,
			},
			status: http.StatusBadRequest,
			id:     "00000000-0000-0000-0000-000000000000",
		},
		{
			name: "malformed-id",
			body: MovieRequest{
				Title:         "TestMovie",
				Description:   "New Description",
				ImageURL:      "http://example.com/image.png",
				Rating:        7.6666,
				LengthMinutes: 125,
			},
			status: http.StatusBadRequest,
			id:     "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/movies/%s", testCase.id)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodPut, testCase.body)
			assert.NoError(t, err)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			ignoreResp := xtesting.ValuesCheckers{
				"updated_at": xtesting.ValueTimeInPastDuration(time.Second),
			}

			ignoreMovies := xtesting.GenerateValueCheckersForArrays(map[string]xtesting.ValueChecker{"UpdatedAt": xtesting.ValueTime()}, 10)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w, ignoreResp)
			xtesting.AssertGoldenDatabaseTable(t, db, []models.Movie{}, ignoreMovies)
		})
	}
}

func TestMoviesDelete(t *testing.T) {
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
			id:     "510633ca-e23f-11f0-a626-d3b8771e2cb9",
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

			targetURL := fmt.Sprintf("/api/v1/movies/%s", testCase.id)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodDelete, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
			xtesting.AssertGoldenDatabaseTable(t, db, []models.Movie{}, nil)
		})
	}
}
