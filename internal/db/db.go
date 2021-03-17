package db

import (
	"context"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/slawekzachcial/tournament/internal/tournament"
)

func WaitForDb(dbUrl string, tries, waitSecs int) error {
	var result error

	for i := tries; i > 0; i-- {
		conn, err := pgx.Connect(context.Background(), dbUrl)
		if err == nil {
			defer conn.Close(context.Background())
			return nil
		}
		result = err
		time.Sleep(time.Duration(waitSecs) * time.Second)
	}

	return result
}

func RunMigrations(folder, dbUrl string) error {
	m, err := migrate.New(folder, dbUrl)
	if err != nil {
		return err
	}
	if err := m.Up(); err != migrate.ErrNoChange {
		return err
	}

	return nil
}

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
