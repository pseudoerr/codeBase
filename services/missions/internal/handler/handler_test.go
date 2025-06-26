package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pseudoerr/mission-service/internal/handler"
	"github.com/pseudoerr/mission-service/models"
	"github.com/pseudoerr/mission-service/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeStore struct {
	missions []models.Mission
	nextID   int
}

func (f *fakeStore) ListMissions(ctx context.Context) ([]models.Mission, error) {
	return f.missions, nil
}

func (f *fakeStore) AddMission(ctx context.Context, m models.Mission) (models.Mission, error) {
	f.nextID++
	m.ID = f.nextID
	f.missions = append(f.missions, m)
	return m, nil
}

func (f *fakeStore) GetByID(ctx context.Context, id int) (models.Mission, error) {
	for _, m := range f.missions {
		if m.ID == id {
			return m, nil
		}
	}
	return models.Mission{}, fmt.Errorf("not found")
}

func (f *fakeStore) UpdateMission(ctx context.Context, m models.Mission) (models.Mission, error) {
	return m, nil
}

func (f *fakeStore) DeleteMission(ctx context.Context, id int) error {
	return nil
}

func TestGetMissions(t *testing.T) {

	store := &fakeStore{
		missions: []models.Mission{
			{ID: 1, Title: "Test", Points: 100},
		},
		nextID: 1,
	}
	svc := &service.MissionService{Store: store}
	newHandler := &handler.Handler{Service: svc}
	router := handler.NewRouter(newHandler)

	req := httptest.NewRequest(http.MethodGet, "/missions", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var data []models.Mission
	if err := json.NewDecoder(rec.Body).Decode(&data); err != nil {
		t.Fatalf("bad json: %v", err)
	}

	if len(data) != 1 || data[0].Title != "Test" {
		t.Fatalf("unexpected data: %+v", data)
	}
}
