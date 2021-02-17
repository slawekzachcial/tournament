package tournament

import (
	"fmt"
	"sort"
)

type Tournament struct {
	teamStats map[string]*Stats
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

func NewTournament() *Tournament {
	return &Tournament{
		teamStats: make(map[string]*Stats),
	}
}

func (t *Tournament) GetStats(team string) (Stats, error) {
	if stats, ok := t.teamStats[team]; ok {
		return *stats, nil
	} else {
		return Stats{}, fmt.Errorf("Team not found or it has not played yet: %v", team)
	}
}

func (t *Tournament) GetAllStats() []Stats {
	allStats := make([]Stats, 0, len(t.teamStats))

	for _, stats := range t.teamStats {
		allStats = append(allStats, *stats)
	}

	sort.Slice(allStats, func(i, j int) bool {
		return allStats[i].Points > allStats[j].Points || allStats[i].Points == allStats[i].Points && allStats[i].Team < allStats[i].Team
	})

	return allStats
}

func (t *Tournament) Play(game Game) {
	if _, ok := t.teamStats[game.TeamA]; !ok {
		t.teamStats[game.TeamA] = &Stats{Team: game.TeamA}
	}
	if _, ok := t.teamStats[game.TeamB]; !ok {
		t.teamStats[game.TeamB] = &Stats{Team: game.TeamB}
	}

	t.teamStats[game.TeamA].Played++
	t.teamStats[game.TeamB].Played++

	if game.ScoreA > game.ScoreB {
		t.teamStats[game.TeamA].Won++
		t.teamStats[game.TeamB].Lost++
		t.teamStats[game.TeamA].Points += 3
	} else if game.ScoreA < game.ScoreB {
		t.teamStats[game.TeamA].Lost++
		t.teamStats[game.TeamB].Won++
		t.teamStats[game.TeamB].Points += 3
	} else {
		t.teamStats[game.TeamA].Points++
		t.teamStats[game.TeamB].Points++
		t.teamStats[game.TeamA].Drawn++
		t.teamStats[game.TeamB].Drawn++
	}
}
