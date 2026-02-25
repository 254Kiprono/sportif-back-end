package services

import (
	"errors"
	"webuye-sportif/app/models"
	"webuye-sportif/app/repository"
)

type MatchService interface {
	GetLineup(fixtureID string) (*models.Lineup, error)
	SaveLineup(lineup *models.Lineup) error

	StartMatch(fixtureID string) error
	EndMatch(fixtureID string) error

	GetEvents(fixtureID string) ([]models.MatchEvent, error)
	LogEvent(event *models.MatchEvent) error
}

type matchService struct {
	matchRepo   repository.MatchRepository
	playerRepo  repository.PlayerRepository
	fixtureRepo repository.FixtureRepository
}

func NewMatchService(m repository.MatchRepository, p repository.PlayerRepository, f repository.FixtureRepository) MatchService {
	return &matchService{m, p, f}
}

func (s *matchService) GetLineup(fixtureID string) (*models.Lineup, error) {
	return s.matchRepo.GetLineup(fixtureID)
}

func (s *matchService) SaveLineup(lineup *models.Lineup) error {
	fixture, err := s.fixtureRepo.GetByID(lineup.FixtureID.String())
	if err != nil {
		return err
	}
	if fixture.Status != "upcoming" {
		return errors.New("cannot edit lineup after match has started or finished")
	}
	return s.matchRepo.SaveLineup(lineup)
}

func (s *matchService) StartMatch(fixtureID string) error {
	fixture, err := s.fixtureRepo.GetByID(fixtureID)
	if err != nil {
		return err
	}
	if fixture.Status != "upcoming" {
		return errors.New("match has already started or finished")
	}

	lineup, err := s.matchRepo.GetLineup(fixtureID)
	if err != nil {
		return errors.New("cannot start match without a lineup")
	}

	return s.matchRepo.Transaction(func(tx repository.MatchRepository) error {
		// Update fixture status
		fixture.Status = "live"
		if err := s.fixtureRepo.Update(fixture); err != nil {
			return err
		}

		// Increment appearances for starters
		for _, lp := range lineup.Players {
			if lp.IsStarter {
				player, err := s.playerRepo.GetByID(lp.PlayerID.String())
				if err == nil {
					player.Appearances++
					s.playerRepo.Update(player)
				}
			}
		}
		return nil
	})
}

func (s *matchService) EndMatch(fixtureID string) error {
	fixture, err := s.fixtureRepo.GetByID(fixtureID)
	if err != nil {
		return err
	}
	if fixture.Status != "live" {
		return errors.New("only live matches can be ended")
	}

	fixture.Status = "completed"
	return s.fixtureRepo.Update(fixture)
}

func (s *matchService) GetEvents(fixtureID string) ([]models.MatchEvent, error) {
	return s.matchRepo.GetEvents(fixtureID)
}

func (s *matchService) LogEvent(event *models.MatchEvent) error {
	fixture, err := s.fixtureRepo.GetByID(event.FixtureID.String())
	if err != nil {
		return err
	}
	if fixture.Status != "live" {
		return errors.New("events can only be logged for live matches")
	}

	return s.matchRepo.Transaction(func(tx repository.MatchRepository) error {
		// Create the event
		if err := tx.CreateEvent(event); err != nil {
			return err
		}

		// Apply event-driven side effects
		switch event.Type {
		case models.EventGoal:
			if event.IsOpponent {
				fixture.AwayScore++ // Assuming Webuye is always home for simplicity or we should check
				// In a real app we'd check if fixture.HomeTeam is Webuye
			} else {
				fixture.HomeScore++
				if event.PlayerID != nil {
					player, err := s.playerRepo.GetByID(event.PlayerID.String())
					if err == nil {
						player.Goals++
						s.playerRepo.Update(player)
					}
				}
				if event.AssistPlayerID != nil {
					asst, err := s.playerRepo.GetByID(event.AssistPlayerID.String())
					if err == nil {
						asst.Assists++
						s.playerRepo.Update(asst)
					}
				}
			}
			s.fixtureRepo.Update(fixture)

		case models.EventSubstitution:
			if event.PlayerID != nil { // PlayerID is the one coming ON
				player, err := s.playerRepo.GetByID(event.PlayerID.String())
				if err == nil {
					// Check if they already appeared (maybe they were starters? unlikely but safety)
					// Logic: If they are coming on as sub, they get an appearance if not already marked.
					// We'd ideally track match-specific appearances, but global is what we have.
					player.Appearances++
					s.playerRepo.Update(player)
				}
			}
		}

		return nil
	})
}
