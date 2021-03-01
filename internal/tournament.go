package tournament

import (
	"errors"
	"sort"
)

type Tournament struct {
	games Games
}

type Game struct {
	TeamA  string
	ScoreA int
	TeamB  string
	ScoreB int
}

type Stats struct {
	Team   string
	Played int
	Won    int
	Drawn  int
	Lost   int
	Points int
}

var ErrTeamNotFound = errors.New("Team not found")

type Games interface {
	Save(game *Game) error
	FindByTeam(team string) ([]Game, error)
	FindAll() ([]Game, error)
}

func NewTournament(games Games) *Tournament {
	return &Tournament{
		games: games,
	}
}

func (t *Tournament) GetStats(team string) (Stats, error) {
	teamGames, err := t.games.FindByTeam(team)
	if err != nil {
		return Stats{}, err
	}

	stats := []*Stats{}
	for _, game := range teamGames {
		stats = updateStats(stats, &game)
	}

	for _, s := range stats {
		if s.Team == team {
			return *s, nil
		}
	}

	return Stats{}, nil
}

func (t *Tournament) GetAllStats() ([]Stats, error) {
	allGames, err := t.games.FindAll()
	if err != nil {
		return nil, err
	}

	allStats := []*Stats{}
	for _, game := range allGames {
		allStats = updateStats(allStats, &game)
	}

	sort.Slice(allStats, func(i, j int) bool {
		return allStats[i].Points > allStats[j].Points || allStats[i].Points == allStats[i].Points && allStats[i].Team < allStats[i].Team
	})

	result := make([]Stats, 0, len(allStats))
	for _, stats := range allStats {
		result = append(result, *stats)
	}

	return result, nil
}

func (t *Tournament) Play(game Game) error {
	return t.games.Save(&game)
}

func updateStats(stats []*Stats, game *Game) []*Stats {
	var teamAStats, teamBStats *Stats

	for _, s := range stats {
		switch s.Team {
		case game.TeamA:
			teamAStats = s
		case game.TeamB:
			teamBStats = s
		}
	}

	if teamAStats == nil {
		teamAStats = &Stats{Team: game.TeamA}
		stats = append(stats, teamAStats)
	}
	if teamBStats == nil {
		teamBStats = &Stats{Team: game.TeamB}
		stats = append(stats, teamBStats)
	}

	teamAStats.Played++
	teamBStats.Played++

	if game.ScoreA > game.ScoreB {
		teamAStats.Won++
		teamBStats.Lost++
		teamAStats.Points += 3
	} else if game.ScoreA < game.ScoreB {
		teamAStats.Lost++
		teamBStats.Won++
		teamBStats.Points += 3
	} else {
		teamAStats.Points++
		teamBStats.Points++
		teamAStats.Drawn++
		teamBStats.Drawn++
	}

	return stats
}
