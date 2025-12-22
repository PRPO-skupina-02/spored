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

func TestRoomsList(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	r := TestingRouter(t, db)

	tests := []struct {
		name      string
		status    int
		params    string
		theaterID string
	}{
		{
			name:      "ok",
			status:    http.StatusOK,
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:      "ok-paginated",
			status:    http.StatusOK,
			params:    "?limit=1&offset=1",
			theaterID: "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
		},
		{
			name:      "ok-sort",
			status:    http.StatusOK,
			params:    "?sort=-updated_at",
			theaterID: "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
		},
		{
			name:      "ok-paginated-sort",
			status:    http.StatusOK,
			params:    "?limit=2&offset=1&sort=updated_at",
			theaterID: "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
		},
		{
			name:      "ok-no-rooms",
			status:    http.StatusOK,
			theaterID: "ea0b7f96-ddc9-11f0-9635-23efd36396bd",
		},
		{
			name:      "invalid-theater-id",
			status:    http.StatusNotFound,
			theaterID: "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name:      "nil-theater-id",
			status:    http.StatusBadRequest,
			theaterID: "00000000-0000-0000-0000-000000000000",
		},
		{
			name:      "malformed-theater-id",
			status:    http.StatusBadRequest,
			theaterID: "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/theaters/%s/rooms%s", testCase.theaterID, testCase.params)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodGet, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
		})
	}
}

func TestRoomsCreate(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	r := TestingRouter(t, db)

	tests := []struct {
		name      string
		body      RoomRequest
		status    int
		theaterID string
	}{
		{
			name: "ok",
			body: RoomRequest{
				Name:    "TestRoom",
				Rows:    10,
				Columns: 20,
			},
			status:    http.StatusCreated,
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "validation-errors",
			body: RoomRequest{
				Name:    "A",
				Rows:    -1,
				Columns: 1000,
			},
			status:    http.StatusBadRequest,
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:      "no-body",
			status:    http.StatusBadRequest,
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "invalid-theater-id",
			body: RoomRequest{
				Name:    "TestRoom",
				Rows:    10,
				Columns: 20,
			},
			status:    http.StatusNotFound,
			theaterID: "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name: "nil-theater-id",
			body: RoomRequest{
				Name:    "TestRoom",
				Rows:    10,
				Columns: 20,
			},
			status:    http.StatusBadRequest,
			theaterID: "00000000-0000-0000-0000-000000000000",
		},
		{
			name: "malformed-theater-id",
			body: RoomRequest{
				Name:    "TestRoom",
				Rows:    10,
				Columns: 20,
			},
			status:    http.StatusBadRequest,
			theaterID: "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/theaters/%s/rooms", testCase.theaterID)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodPost, testCase.body)
			assert.NoError(t, err)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			ignoreResp := xtesting.ValuesCheckers{
				"id":         xtesting.ValueUUID(),
				"created_at": xtesting.ValueTimeInPastDuration(time.Second),
				"updated_at": xtesting.ValueTimeInPastDuration(time.Second),
			}

			ignoreRooms := xtesting.GenerateValueCheckersForArrays(map[string]xtesting.ValueChecker{"ID": xtesting.ValueUUID(), "CreatedAt": xtesting.ValueTime(), "UpdatedAt": xtesting.ValueTime()}, 10)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w, ignoreResp)
			xtesting.AssertGoldenDatabaseTable(t, db.Order("name"), []models.Room{}, ignoreRooms)
		})
	}
}

func TestRoomsShow(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	r := TestingRouter(t, db)

	tests := []struct {
		name      string
		status    int
		roomID    string
		theaterID string
	}{
		{
			name:      "ok",
			status:    http.StatusOK,
			roomID:    "ec19b8aa-df42-11f0-9018-53ba2f5e5e7c",
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:      "room-from-different-theater",
			status:    http.StatusNotFound,
			roomID:    "e0722c3a-df42-11f0-9579-3734395be62a",
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:      "invalid-room-id",
			status:    http.StatusNotFound,
			roomID:    "01234567-0123-0123-0123-0123456789ab",
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:      "nil-room-id",
			status:    http.StatusBadRequest,
			roomID:    "00000000-0000-0000-0000-000000000000",
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:      "malformed-room-id",
			status:    http.StatusBadRequest,
			roomID:    "000",
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:      "invalid-theater-id",
			status:    http.StatusNotFound,
			roomID:    "ec19b8aa-df42-11f0-9018-53ba2f5e5e7c",
			theaterID: "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name:      "nil-theater-id",
			status:    http.StatusBadRequest,
			roomID:    "ec19b8aa-df42-11f0-9018-53ba2f5e5e7c",
			theaterID: "00000000-0000-0000-0000-000000000000",
		},
		{
			name:      "malformed-theater-id",
			status:    http.StatusBadRequest,
			roomID:    "ec19b8aa-df42-11f0-9018-53ba2f5e5e7c",
			theaterID: "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/theaters/%s/rooms/%s", testCase.theaterID, testCase.roomID)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodGet, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
		})
	}
}

func TestRoomsUpdate(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	r := TestingRouter(t, db)

	tests := []struct {
		name      string
		body      RoomRequest
		status    int
		roomID    string
		theaterID string
	}{
		{
			name: "ok",
			body: RoomRequest{
				Name:    "UpdatedRoom",
				Rows:    12,
				Columns: 24,
			},
			status:    http.StatusOK,
			roomID:    "ec19b8aa-df42-11f0-9018-53ba2f5e5e7c",
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "validation-errors",
			body: RoomRequest{
				Name:    "A",
				Rows:    -1,
				Columns: 1000,
			},
			status:    http.StatusBadRequest,
			roomID:    "ec19b8aa-df42-11f0-9018-53ba2f5e5e7c",
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:      "no-body",
			status:    http.StatusBadRequest,
			roomID:    "ec19b8aa-df42-11f0-9018-53ba2f5e5e7c",
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "room-from-different-theater",
			body: RoomRequest{
				Name:    "UpdatedRoom",
				Rows:    12,
				Columns: 24,
			},
			status:    http.StatusNotFound,
			roomID:    "e0722c3a-df42-11f0-9579-3734395be62a",
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "invalid-room-id",
			body: RoomRequest{
				Name:    "UpdatedRoom",
				Rows:    12,
				Columns: 24,
			},
			status:    http.StatusNotFound,
			roomID:    "01234567-0123-0123-0123-0123456789ab",
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "nil-room-id",
			body: RoomRequest{
				Name:    "UpdatedRoom",
				Rows:    12,
				Columns: 24,
			},
			status:    http.StatusBadRequest,
			roomID:    "00000000-0000-0000-0000-000000000000",
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "malformed-room-id",
			body: RoomRequest{
				Name:    "UpdatedRoom",
				Rows:    12,
				Columns: 24,
			},
			status:    http.StatusBadRequest,
			roomID:    "000",
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "invalid-theater-id",
			body: RoomRequest{
				Name:    "UpdatedRoom",
				Rows:    12,
				Columns: 24,
			},
			status:    http.StatusNotFound,
			roomID:    "ec19b8aa-df42-11f0-9018-53ba2f5e5e7c",
			theaterID: "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name: "nil-theater-id",
			body: RoomRequest{
				Name:    "UpdatedRoom",
				Rows:    12,
				Columns: 24,
			},
			status:    http.StatusBadRequest,
			roomID:    "ec19b8aa-df42-11f0-9018-53ba2f5e5e7c",
			theaterID: "00000000-0000-0000-0000-000000000000",
		},
		{
			name: "malformed-theater-id",
			body: RoomRequest{
				Name:    "UpdatedRoom",
				Rows:    12,
				Columns: 24,
			},
			status:    http.StatusBadRequest,
			roomID:    "ec19b8aa-df42-11f0-9018-53ba2f5e5e7c",
			theaterID: "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/theaters/%s/rooms/%s", testCase.theaterID, testCase.roomID)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodPut, testCase.body)
			assert.NoError(t, err)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			ignoreResp := xtesting.ValuesCheckers{
				"updated_at": xtesting.ValueTimeInPastDuration(time.Second),
			}

			ignoreRooms := xtesting.GenerateValueCheckersForArrays(map[string]xtesting.ValueChecker{"UpdatedAt": xtesting.ValueTime()}, 10)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w, ignoreResp)
			xtesting.AssertGoldenDatabaseTable(t, db, []models.Room{}, ignoreRooms)
		})
	}
}

func TestRoomsDelete(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	r := TestingRouter(t, db)

	tests := []struct {
		name      string
		status    int
		roomID    string
		theaterID string
	}{
		{
			name:      "ok",
			status:    http.StatusNoContent,
			roomID:    "ec19b8aa-df42-11f0-9018-53ba2f5e5e7c",
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:      "room-from-different-theater",
			status:    http.StatusNotFound,
			roomID:    "e0722c3a-df42-11f0-9579-3734395be62a",
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:      "invalid-room-id",
			status:    http.StatusNotFound,
			roomID:    "01234567-0123-0123-0123-0123456789ab",
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:      "nil-room-id",
			status:    http.StatusBadRequest,
			roomID:    "00000000-0000-0000-0000-000000000000",
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:      "malformed-room-id",
			status:    http.StatusBadRequest,
			roomID:    "000",
			theaterID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:      "invalid-theater-id",
			status:    http.StatusNotFound,
			roomID:    "ec19b8aa-df42-11f0-9018-53ba2f5e5e7c",
			theaterID: "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name:      "nil-theater-id",
			status:    http.StatusBadRequest,
			roomID:    "ec19b8aa-df42-11f0-9018-53ba2f5e5e7c",
			theaterID: "00000000-0000-0000-0000-000000000000",
		},
		{
			name:      "malformed-theater-id",
			status:    http.StatusBadRequest,
			roomID:    "ec19b8aa-df42-11f0-9018-53ba2f5e5e7c",
			theaterID: "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/theaters/%s/rooms/%s", testCase.theaterID, testCase.roomID)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodDelete, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
			xtesting.AssertGoldenDatabaseTable(t, db, []models.Room{}, nil)
		})
	}
}
