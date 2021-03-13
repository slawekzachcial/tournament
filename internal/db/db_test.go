package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	tournament "github.com/slawekzachcial/tournament/internal"
)

const DB_NAME = "tournament_test"

var dbServerUrl = getEnv("TEST_DATABASE_URL", "postgres://postgres:secret@localhost:5432")
var testDbUrl = fmt.Sprintf("%s/%s?sslmode=disable", dbServerUrl, DB_NAME)

var dbPool *pgxpool.Pool

func TestFindByTeam(t *testing.T) {
	defer deleteAllGames()

	g1, g2, g3 := tournament.Game{"A", 1, "B", 0}, tournament.Game{"B", 2, "C", 2}, tournament.Game{"C", 3, "A", 4}

	gd := GamesData{dbPool}
	gd.Save(&g1)
	gd.Save(&g2)
	gd.Save(&g3)

	aGamesExpected := []tournament.Game{g1, g3}
	aGamesGot, err := gd.FindByTeam("A")
	if err != nil {
		t.Fatalf("Error getting team A games: %v", err)
	}

	if !reflect.DeepEqual(asMap(aGamesExpected), asMap(aGamesGot)) {
		t.Errorf("Expected A games %v but got %v", aGamesExpected, aGamesGot)
	}
}

func TestFindByTeamNotFound(t *testing.T) {
	gd := GamesData{dbPool}

	_, err := gd.FindByTeam("YOU_SHOULD_NOT_FIND_ME")
	if err != tournament.ErrTeamNotFound {
		t.Fatalf("Expecting ErrTeamNotFound error but got %v", err)
	}
}

func TestFindAll(t *testing.T) {
	defer deleteAllGames()

	g1, g2, g3 := tournament.Game{"A", 1, "B", 0}, tournament.Game{"B", 2, "C", 2}, tournament.Game{"C", 3, "A", 4}

	gd := GamesData{dbPool}
	gd.Save(&g1)
	gd.Save(&g2)
	gd.Save(&g3)

	expected := []tournament.Game{g1, g2, g3}
	got, err := gd.FindAll()
	if err != nil {
		t.Fatalf("Error getting all games: %v", err)
	}

	if !reflect.DeepEqual(asMap(expected), asMap(got)) {
		t.Errorf("Expected games %v but got %v", expected, got)
	}
}

func TestMain(m *testing.M) {
	testExitCode := 0
	defer func() { os.Exit(testExitCode) }()
	defer dropTestDatabase()

	err := createTestDatabase()
	if err != nil {
		log.Panicf("Unable to create the database: %v", err)
	}

	dbPool, err = pgxpool.Connect(context.Background(), testDbUrl)
	if err != nil {
		log.Panicf("Unable to connect to the database: %v", err)
	}
	defer dbPool.Close()

	testExitCode = m.Run()
}

func asMap(games []tournament.Game) map[tournament.Game]bool {
	m := make(map[tournament.Game]bool)
	for _, g := range games {
		m[g] = true
	}
	return m
}

func createTestDatabase() error {
	conn, err := pgx.Connect(context.Background(), dbServerUrl)
	if err != nil {
		return fmt.Errorf("Unable to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %v;", DB_NAME))
	if err != nil {
		return err
	}

	if err := RunMigrations("file://../../sql", testDbUrl); err != nil {
		return err
	}

	return nil
}

func dropTestDatabase() {
	conn, err := pgx.Connect(context.Background(), dbServerUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to the database: %v", err)
		return
	}
	defer conn.Close(context.Background())

	// TODO: I don't want to "WITH (FORCE)" - what is still connected when we get here?
	_, err = conn.Exec(context.Background(), fmt.Sprintf("DROP DATABASE %v WITH (FORCE);", DB_NAME))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error dropping the database: %v", err)
	}
}

func deleteAllGames() {
	_, err := dbPool.Exec(context.Background(), "TRUNCATE games;")
	if err != nil {
		log.Panicf("Unable to delete all games: %v", err)
	}
}

func getEnv(name, defaultValue string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return defaultValue
}
