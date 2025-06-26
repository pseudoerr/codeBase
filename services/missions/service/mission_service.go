package service

import (
	"context"
	"log/slog"
	"sync"

	"github.com/pseudoerr/mission-service/models"
)

type MissionStore interface {
	ListMissionsByUserID(ctx context.Context, userID string) ([]models.Mission, error)
	AddMission(ctx context.Context, m models.Mission) (models.Mission, error)
	GetMissionByID(ctx context.Context, missionID string, userID string) (models.Mission, error)
	UpdateMission(ctx context.Context, m models.Mission) (models.Mission, error)
	DeleteMission(ctx context.Context, missionID string, userID string) error
	GetProfileByUserID(ctx context.Context, userID string) (models.Profile, error)
}

type InMemoryStore struct {
	mu       sync.Mutex
	missions []models.Mission
	nextID   int
}

type MissionService struct {
	Store  MissionStore
	Logger *slog.Logger
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		missions: []models.Mission{
			{ID: "1", Title: "Hello, World!", Points: 100},
			{ID: "2", Title: "FizzBuzz", Points: 200},
		},
		nextID: 3,
	}
}

func (s *InMemoryStore) ListMissions(ctx context.Context) ([]models.Mission, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]models.Mission{}, s.missions...), nil
}

func (s *InMemoryStore) AddMission(ctx context.Context, m models.Mission) (models.Mission, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m.ID = s.nextID
	s.nextID++
	s.missions = append(s.missions, m)
	return m, nil
}

func (s *MissionService) GetProfile(ctx context.Context) (models.Profile, error) {
	missions, err := s.Store.ListMissions(ctx)
	if err != nil {
		return models.Profile{}, err
	}
	total := 0
	for _, m := range missions {
		total += m.Points

	}
	var level string
	switch {
	case total >= 1000:
		level = "Expert"
	case total >= 500:
		level = "Advanced"
	case total >= 200:
		level = "Intermediate"
	default:
		level = "Beginner"
	}

	badges := []string{}
	if total >= 200 {
		badges = append(badges, "🏅200+ points")
	}

	if total >= 500 {
		badges = append(badges, "🎖 500+ points")
	}
	if total >= 1000 {
		badges = append(badges, "🏆 1000+ points")
	}

	return models.Profile{
		TotalPoints:  total,
		Level:        level,
		Achievements: badges}, nil
}
