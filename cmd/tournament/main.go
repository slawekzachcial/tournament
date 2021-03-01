package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/slawekzachcial/tournament"
	"github.com/slawekzachcial/tournament/db"
	"github.com/slawekzachcial/tournament/gen/models"
	"github.com/slawekzachcial/tournament/gen/restapi"
	"github.com/slawekzachcial/tournament/gen/restapi/operations"
)

var portFlag = flag.Int("port", 3000, "Port to run this service on")

func main() {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		fmt.Fprintf(os.Stderr, "DB_URL environment variable not set")
		os.Exit(1)
	}

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec: %v", err)
		os.Exit(1)
	}

	api := operations.NewTournamentAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	flag.Parse()
	server.Port = *portFlag

	dbPool, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	games := db.NewGameData(dbPool)
	theTournament := tournament.NewTournament(games)

	api.PlayHandler = playHandler(theTournament)
	api.GetAllStatsHandler = getAllStatsHandler(theTournament)
	api.GetTeamStatsHandler = getTeamStatsHandler(theTournament)

	if err := server.Serve(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v", err)
		os.Exit(1)
	}
}

func playHandler(theTournament *tournament.Tournament) operations.PlayHandlerFunc {
	return func(params operations.PlayParams) middleware.Responder {
		game := tournament.Game{*params.Body.TeamA, int(*params.Body.ScoreA), *params.Body.TeamB, int(*params.Body.ScoreB)}
		err := theTournament.Play(game)
		if err != nil {
			msg := err.Error()
			return operations.NewPlayDefault(400).WithPayload(&models.Error{400, &msg})
		}

		return operations.NewPlayCreated()
	}
}

func getAllStatsHandler(theTournament *tournament.Tournament) operations.GetAllStatsHandlerFunc {
	return func(params operations.GetAllStatsParams) middleware.Responder {
		stats, err := theTournament.GetAllStats()
		if err != nil {
			msg := err.Error()
			return operations.NewGetAllStatsDefault(500).WithPayload(&models.Error{500, &msg})
		}

		payload := make([]*models.Stats, 0, len(stats))
		for _, s := range stats {
			played, won, drawn, lost, points := int64(s.Played), int64(s.Won), int64(s.Drawn), int64(s.Lost), int64(s.Points)
			ms := models.Stats{
				Team:   &s.Team,
				Played: &played,
				Won:    &won,
				Drawn:  &drawn,
				Lost:   &lost,
				Points: &points,
			}
			payload = append(payload, &ms)
		}
		return operations.NewGetAllStatsOK().WithPayload(payload)
	}
}

func getTeamStatsHandler(theTournament *tournament.Tournament) operations.GetTeamStatsHandlerFunc {
	return func(params operations.GetTeamStatsParams) middleware.Responder {
		s, err := theTournament.GetStats(params.Team)
		if err != nil {
			if err == tournament.ErrTeamNotFound {
				msg := fmt.Sprintf("Team '%s' not found", params.Team)
				return operations.NewGetTeamStatsDefault(404).WithPayload(&models.Error{404, &msg})
			} else {
				msg := err.Error()
				return operations.NewGetTeamStatsDefault(500).WithPayload(&models.Error{500, &msg})
			}
		}
		played, won, drawn, lost, points := int64(s.Played), int64(s.Won), int64(s.Drawn), int64(s.Lost), int64(s.Points)
		ms := models.Stats{
			Team:   &s.Team,
			Played: &played,
			Won:    &won,
			Drawn:  &drawn,
			Lost:   &lost,
			Points: &points,
		}
		return operations.NewGetTeamStatsOK().WithPayload(&ms)
	}
}
