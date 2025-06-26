package service_test

import (
	"context"
	"github.com/pseudoerr/mission-service/models"
	"github.com/pseudoerr/mission-service/service"
	"testing"
)

func TestAddMission(t *testing.T) {

	store := service.NewInMemoryStore()

	newM := models.Mission{
		Title:  "Test Mission",
		Points: 50,
	}

	added, err := store.AddMission(context.Background(), newM)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if added.ID == 0 {
		t.Errorf("expected non-zero ID, got 0")
	}

	if added.Title != newM.Title {
		t.Errorf("expected title %s, got %s", newM.Title, added.Title)
	}

	missions, err := store.ListMissions(context.Background())
	if err != nil {
		t.Fatalf("could not list missions: %v", err)
	}

	found := false

	for _, m := range missions {
		if m.ID == added.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("added mission not found in store")
	}
}
