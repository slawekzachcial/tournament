package tournament

import (
	"reflect"
	"testing"
)

type GamesArray struct {
	games []Game
}

func (ga *GamesArray) Save(game *Game) error {
	if ga.games == nil {
		ga.games = []Game{}
	}
	ga.games = append(ga.games, *game)
	return nil
}

func (ga *GamesArray) FindByTeam(team string) ([]Game, error) {
	teamGames := []Game{}
	for _, game := range ga.games {
		if team == game.TeamA || team == game.TeamB {
			teamGames = append(teamGames, game)
		}
	}
	if len(teamGames) == 0 {
		return nil, ErrTeamNotFound
	}
	return teamGames, nil
}

func (ga *GamesArray) FindAll() ([]Game, error) {
	return ga.games, nil
}

var gameTestData = []struct {
	testName string
	games    []Game
	stats    []Stats
}{
	{
		testName: "single game",
		games: []Game{
			{"a", 2, "b", 1},
		},
		stats: []Stats{
			Stats{Team: "a", Played: 1, Won: 1, Points: 3},
			Stats{Team: "b", Played: 1, Lost: 1},
		},
	},
	{
		testName: "drawn game",
		games: []Game{
			{"a", 1, "b", 1},
		},
		stats: []Stats{
			Stats{Team: "a", Played: 1, Drawn: 1, Points: 1},
			Stats{Team: "b", Played: 1, Drawn: 1, Points: 1},
		},
	},
	{
		testName: "multiple games",
		games: []Game{
			{"a", 1, "b", 0},
			{"a", 3, "c", 3},
			{"b", 0, "c", 1},
		},
		stats: []Stats{
			Stats{Team: "a", Played: 2, Won: 1, Drawn: 1, Points: 4},
			Stats{Team: "b", Played: 2, Lost: 2, Points: 0},
			Stats{Team: "c", Played: 2, Won: 1, Drawn: 1, Points: 4},
		},
	},
}

func TestGetStatsError(t *testing.T) {
	tournament := NewTournament(&GamesArray{})
	_, err := tournament.GetStats("unknown")
	if err != ErrTeamNotFound {
		t.Errorf("Expected error when getting stats for team that has not played yet")
	}
}

func TestGames(t *testing.T) {
	for _, testData := range gameTestData {
		tournament := NewTournament(&GamesArray{})
		for _, game := range testData.games {
			tournament.Play(game)
		}
		for _, expectedStats := range testData.stats {
			gotStats, _ := tournament.GetStats(expectedStats.Team)
			if !reflect.DeepEqual(gotStats, expectedStats) {
				t.Errorf("%v: team '%v' stats - expected: %v, got: %v", testData.testName, expectedStats.Team, expectedStats, gotStats)
			}
		}
	}
}

func TestGetAllStats(t *testing.T) {
	tournament := NewTournament(&GamesArray{})

	tournament.Play(Game{"a", 1, "b", 0})
	tournament.Play(Game{"a", 3, "c", 3})
	tournament.Play(Game{"b", 0, "c", 1})

	allStats, _ := tournament.GetAllStats()
	expectedStats := []Stats{
		Stats{Team: "a", Played: 2, Won: 1, Drawn: 1, Points: 4},
		Stats{Team: "c", Played: 2, Won: 1, Drawn: 1, Points: 4},
		Stats{Team: "b", Played: 2, Lost: 2, Points: 0},
	}

	if !reflect.DeepEqual(allStats, expectedStats) {
		t.Errorf("All stats - expected: %v, got: %v", expectedStats, allStats)
	}
}
