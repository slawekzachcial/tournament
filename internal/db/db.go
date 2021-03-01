package db

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	tournament "github.com/slawekzachcial/tournament/internal"
)

type GamesData struct {
	pool *pgxpool.Pool
}

func NewGameData(p *pgxpool.Pool) *GamesData {
	return &GamesData{p}
}

func (g *GamesData) Save(game *tournament.Game) error {
	_, err := g.pool.Exec(context.Background(), "INSERT INTO games(team_a, score_a, team_b, score_b) VALUES ($1, $2, $3, $4)",
		game.TeamA, game.ScoreA, game.TeamB, game.ScoreB)
	if err != nil {
		return err
	}
	return nil
}

func (g *GamesData) FindByTeam(team string) ([]tournament.Game, error) {
	rows, err := g.pool.Query(context.Background(),
		"SELECT team_a, score_a, team_b, score_b FROM games WHERE team_a=$1 OR team_b=$1",
		team)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games, err := rowsToGames(rows)
	if len(games) == 0 {
		return nil, tournament.ErrTeamNotFound
	}
	return games, err
}

func (g *GamesData) FindAll() ([]tournament.Game, error) {
	rows, err := g.pool.Query(context.Background(),
		"SELECT team_a, score_a, team_b, score_b FROM games")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsToGames(rows)
}

func rowsToGames(rows pgx.Rows) ([]tournament.Game, error) {
	games := []tournament.Game{}
	for rows.Next() {
		var teamA, teamB string
		var scoreA, scoreB int
		err := rows.Scan(&teamA, &scoreA, &teamB, &scoreB)
		if err != nil {
			return nil, err
		}

		game := tournament.Game{teamA, scoreA, teamB, scoreB}
		games = append(games, game)
	}
	return games, nil
}
