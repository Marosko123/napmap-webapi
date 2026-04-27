package napmap

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Marosko123/napmap-webapi/internal/db_service"
)

// mockDbService is a thread-safe in-memory implementation of DbService[Station] for tests.
type mockDbService struct {
	mu    sync.Mutex
	store map[string]*Station
	pingErr error
}

func newMockDbService() *mockDbService {
	return &mockDbService{store: map[string]*Station{}}
}

func (m *mockDbService) CreateDocument(_ context.Context, id string, doc *Station) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.store[id]; ok {
		return db_service.ErrConflict
	}
	clone := *doc
	m.store[id] = &clone
	return nil
}

func (m *mockDbService) FindDocument(_ context.Context, id string) (*Station, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	doc, ok := m.store[id]
	if !ok {
		return nil, db_service.ErrNotFound
	}
	clone := *doc
	return &clone, nil
}

func (m *mockDbService) FindDocuments(_ context.Context, filter interface{}) ([]*Station, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]*Station, 0, len(m.store))
	for _, s := range m.store {
		clone := *s
		out = append(out, &clone)
	}
	return out, nil
}

func (m *mockDbService) UpdateDocument(_ context.Context, id string, doc *Station) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.store[id]; !ok {
		return db_service.ErrNotFound
	}
	clone := *doc
	m.store[id] = &clone
	return nil
}

func (m *mockDbService) DeleteDocument(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.store[id]; !ok {
		return db_service.ErrNotFound
	}
	delete(m.store, id)
	return nil
}

func (m *mockDbService) Ping(_ context.Context) error {
	return m.pingErr
}

func (m *mockDbService) Disconnect(_ context.Context) error { return nil }

func setupTestContext(method, path, body string, params gin.Params, db db_service.DbService[Station]) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("db_service", db)
	c.Params = params
	if body != "" {
		c.Request = httptest.NewRequest(method, path, bytes.NewBufferString(body))
		c.Request.Header.Set("Content-Type", "application/json")
	} else {
		c.Request = httptest.NewRequest(method, path, nil)
	}
	return c, w
}

func validStationJSON() string {
	return `{
		"name": "NAPMap Test",
		"stationType": "CHARGING",
		"fuels": ["ELECTRIC"],
		"operatorName": "NAPMap Energy s.r.o.",
		"address": "Mlynské nivy 1",
		"city": "Bratislava",
		"lat": 48.1486,
		"lng": 17.1077,
		"maxPowerKw": 50
	}`
}

func TestCreateStation_assignsUUIDWhenIdMissing(t *testing.T) {
	mock := newMockDbService()
	c, w := setupTestContext(http.MethodPost, "/stations", validStationJSON(), nil, mock)

	NewStationsApi().CreateStation(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
	var resp Station
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Id == "" || resp.Id == "@new" {
		t.Errorf("expected generated UUID, got %q", resp.Id)
	}
	if resp.Status != "ACTIVE" {
		t.Errorf("expected default status ACTIVE, got %q", resp.Status)
	}
	if resp.Country != "SK" {
		t.Errorf("expected default country SK, got %q", resp.Country)
	}
}

func TestCreateStation_rejectsMissingRequiredFields(t *testing.T) {
	cases := []struct {
		name string
		body string
	}{
		{"no name", `{"stationType":"CHARGING","fuels":["ELECTRIC"],"operatorName":"X","address":"Y","city":"Z"}`},
		{"no stationType", `{"name":"A","fuels":["ELECTRIC"],"operatorName":"X","address":"Y","city":"Z"}`},
		{"empty fuels", `{"name":"A","stationType":"CHARGING","fuels":[],"operatorName":"X","address":"Y","city":"Z"}`},
		{"no operator", `{"name":"A","stationType":"CHARGING","fuels":["ELECTRIC"],"address":"Y","city":"Z"}`},
		{"no city", `{"name":"A","stationType":"CHARGING","fuels":["ELECTRIC"],"operatorName":"X","address":"Y"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, w := setupTestContext(http.MethodPost, "/stations", tc.body, nil, newMockDbService())
			NewStationsApi().CreateStation(c)
			if w.Code != http.StatusBadRequest {
				t.Errorf("%s: expected 400, got %d body=%s", tc.name, w.Code, w.Body.String())
			}
		})
	}
}

func TestCreateStation_returns409OnDuplicate(t *testing.T) {
	mock := newMockDbService()
	mock.store["st-001"] = &Station{Id: "st-001", Name: "Existing"}
	body := strings.Replace(validStationJSON(), `"name": "NAPMap Test"`, `"id":"st-001","name":"NAPMap Test"`, 1)

	c, w := setupTestContext(http.MethodPost, "/stations", body, nil, mock)
	NewStationsApi().CreateStation(c)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetStation_returns404WhenMissing(t *testing.T) {
	c, w := setupTestContext(http.MethodGet, "/stations/missing", "", gin.Params{{Key: "stationId", Value: "missing"}}, newMockDbService())
	NewStationsApi().GetStation(c)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestDeleteStation_softDeletesByMarkingInactive(t *testing.T) {
	mock := newMockDbService()
	mock.store["st-001"] = &Station{Id: "st-001", Name: "Foo", Status: "ACTIVE"}

	c, w := setupTestContext(http.MethodDelete, "/stations/st-001", "", gin.Params{{Key: "stationId", Value: "st-001"}}, mock)
	NewStationsApi().DeleteStation(c)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d body=%s", w.Code, w.Body.String())
	}
	if got := mock.store["st-001"]; got == nil || got.Status != "INACTIVE" {
		t.Errorf("expected status INACTIVE, got %+v", got)
	}
}

func TestUpdateStation_replacesDocumentAndPreservesId(t *testing.T) {
	mock := newMockDbService()
	mock.store["st-001"] = &Station{Id: "st-001", Name: "Old", City: "Bratislava"}

	body := `{"id":"will-be-overwritten","name":"New","stationType":"CHARGING","fuels":["ELECTRIC"],"operatorName":"X","address":"A","city":"Žilina","lat":1,"lng":2}`
	c, w := setupTestContext(http.MethodPut, "/stations/st-001", body, gin.Params{{Key: "stationId", Value: "st-001"}}, mock)
	NewStationsApi().UpdateStation(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	got := mock.store["st-001"]
	if got.Id != "st-001" {
		t.Errorf("expected id preserved as st-001, got %q", got.Id)
	}
	if got.City != "Žilina" {
		t.Errorf("expected city updated, got %q", got.City)
	}
}

func TestGetStations_returnsListFromDb(t *testing.T) {
	mock := newMockDbService()
	mock.store["a"] = &Station{Id: "a", Name: "A"}
	mock.store["b"] = &Station{Id: "b", Name: "B"}

	c, w := setupTestContext(http.MethodGet, "/stations", "", nil, mock)
	NewStationsApi().GetStations(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var stations []*Station
	if err := json.Unmarshal(w.Body.Bytes(), &stations); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(stations) != 2 {
		t.Errorf("expected 2 stations, got %d", len(stations))
	}
}
