package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PRPO-skupina-02/common/database"
	"github.com/PRPO-skupina-02/common/xtesting"
	"github.com/PRPO-skupina-02/spored/db"
	"github.com/stretchr/testify/assert"
)

func TestTimeSlotsList(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	r := TestingRouter(t, db)

	tests := []struct {
		name      string
		status    int
		params    string
		theaterID string
		roomID    string
	}{
		{
			name:      "ok",
			status:    http.StatusOK,
			theaterID: "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:    "925c2358-df46-11f0-a38e-abe580bde3d1",
		},
		{
			name:      "ok-paginated",
			status:    http.StatusOK,
			params:    "?limit=1&offset=1",
			theaterID: "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:    "925c2358-df46-11f0-a38e-abe580bde3d1",
		},
		{
			name:      "ok-sort",
			status:    http.StatusOK,
			params:    "?sort=-updated_at",
			theaterID: "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:    "925c2358-df46-11f0-a38e-abe580bde3d1",
		},
		{
			name:      "ok-paginated-sort",
			status:    http.StatusOK,
			params:    "?limit=2&offset=1&sort=updated_at",
			theaterID: "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:    "925c2358-df46-11f0-a38e-abe580bde3d1",
		},
		{
			name:      "ok-filter-date",
			status:    http.StatusOK,
			params:    "?date=2025-12-30",
			theaterID: "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:    "925c2358-df46-11f0-a38e-abe580bde3d1",
		},
		{
			name:      "ok-filter-date-no-results",
			status:    http.StatusOK,
			params:    "?date=2025-01-01",
			theaterID: "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:    "925c2358-df46-11f0-a38e-abe580bde3d1",
		},
		{
			name:      "ok-filter-date-with-pagination",
			status:    http.StatusOK,
			params:    "?date=2025-12-30&limit=2&offset=1",
			theaterID: "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:    "925c2358-df46-11f0-a38e-abe580bde3d1",
		},
		{
			name:      "ok-filter-date-with-sort",
			status:    http.StatusOK,
			params:    "?date=2025-12-30&sort=-start_time",
			theaterID: "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:    "925c2358-df46-11f0-a38e-abe580bde3d1",
		},
		{
			name:      "ok-invalid-date-format",
			status:    http.StatusOK,
			params:    "?date=invalid",
			theaterID: "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:    "925c2358-df46-11f0-a38e-abe580bde3d1",
		},
		{
			name:      "invalid-theater-id",
			status:    http.StatusNotFound,
			theaterID: "01234567-0123-0123-0123-0123456789ab",
			roomID:    "925c2358-df46-11f0-a38e-abe580bde3d1",
		},
		{
			name:      "invalid-room-id",
			status:    http.StatusNotFound,
			theaterID: "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:    "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name:      "nil-theater-id",
			status:    http.StatusBadRequest,
			theaterID: "00000000-0000-0000-0000-000000000000",
			roomID:    "925c2358-df46-11f0-a38e-abe580bde3d1",
		},
		{
			name:      "nil-room-id",
			status:    http.StatusBadRequest,
			theaterID: "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:    "00000000-0000-0000-0000-000000000000",
		},
		{
			name:      "malformed-theater-id",
			status:    http.StatusBadRequest,
			theaterID: "000",
			roomID:    "925c2358-df46-11f0-a38e-abe580bde3d1",
		},
		{
			name:      "malformed-room-id",
			status:    http.StatusBadRequest,
			theaterID: "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:    "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/spored/theaters/%s/rooms/%s/timeslots%s", testCase.theaterID, testCase.roomID, testCase.params)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodGet, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
		})
	}
}

func TestTimeSlotsShow(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	r := TestingRouter(t, db)

	tests := []struct {
		name       string
		status     int
		theaterID  string
		roomID     string
		timeSlotID string
	}{
		{
			name:       "ok",
			status:     http.StatusOK,
			theaterID:  "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:     "925c2358-df46-11f0-a38e-abe580bde3d1",
			timeSlotID: "9d71d7fd-d88e-41a1-86dc-21b7f2550295",
		},
		{
			name:       "invalid-theater-id",
			status:     http.StatusNotFound,
			theaterID:  "01234567-0123-0123-0123-0123456789ab",
			roomID:     "925c2358-df46-11f0-a38e-abe580bde3d1",
			timeSlotID: "9d71d7fd-d88e-41a1-86dc-21b7f2550295",
		},
		{
			name:       "invalid-room-id",
			status:     http.StatusNotFound,
			theaterID:  "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:     "01234567-0123-0123-0123-0123456789ab",
			timeSlotID: "9d71d7fd-d88e-41a1-86dc-21b7f2550295",
		},
		{
			name:       "invalid-timeslot-id",
			status:     http.StatusNotFound,
			theaterID:  "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:     "925c2358-df46-11f0-a38e-abe580bde3d1",
			timeSlotID: "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name:       "nil-theater-id",
			status:     http.StatusBadRequest,
			theaterID:  "00000000-0000-0000-0000-000000000000",
			roomID:     "925c2358-df46-11f0-a38e-abe580bde3d1",
			timeSlotID: "9d71d7fd-d88e-41a1-86dc-21b7f2550295",
		},
		{
			name:       "nil-room-id",
			status:     http.StatusBadRequest,
			theaterID:  "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:     "00000000-0000-0000-0000-000000000000",
			timeSlotID: "9d71d7fd-d88e-41a1-86dc-21b7f2550295",
		},
		{
			name:       "nil-timeslot-id",
			status:     http.StatusBadRequest,
			theaterID:  "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:     "925c2358-df46-11f0-a38e-abe580bde3d1",
			timeSlotID: "00000000-0000-0000-0000-000000000000",
		},
		{
			name:       "malformed-theater-id",
			status:     http.StatusBadRequest,
			theaterID:  "000",
			roomID:     "925c2358-df46-11f0-a38e-abe580bde3d1",
			timeSlotID: "9d71d7fd-d88e-41a1-86dc-21b7f2550295",
		},
		{
			name:       "malformed-room-id",
			status:     http.StatusBadRequest,
			theaterID:  "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:     "000",
			timeSlotID: "9d71d7fd-d88e-41a1-86dc-21b7f2550295",
		},
		{
			name:       "malformed-timeslot-id",
			status:     http.StatusBadRequest,
			theaterID:  "bae209f6-d059-11f0-b2a4-cbf992c2eb6d",
			roomID:     "925c2358-df46-11f0-a38e-abe580bde3d1",
			timeSlotID: "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/spored/theaters/%s/rooms/%s/timeslots/%s", testCase.theaterID, testCase.roomID, testCase.timeSlotID)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodGet, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
		})
	}
}
