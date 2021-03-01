package db

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/slawekzachcial/tournament"
	trn "github.com/slawekzachcial/tournament"
)

const DB_NAME = "tournament_test"

// TEST_DATABASE_URL="postgres://postgres:secret@localhost:5432"
var dbServerUrl = getEnv("TEST_DATABASE_URL", "postgres://postgres:secret@localhost:5432")
var testDbUrl = fmt.Sprintf("%s/%s", dbServerUrl, DB_NAME)

var dbPool *pgxpool.Pool

func TestFindByTeam(t *testing.T) {
	defer deleteAllGames()

	g1, g2, g3 := trn.Game{"A", 1, "B", 0}, trn.Game{"B", 2, "C", 2}, trn.Game{"C", 3, "A", 4}

	gd := GamesData{dbPool}
	gd.Save(&g1)
	gd.Save(&g2)
	gd.Save(&g3)

	aGamesExpected := []trn.Game{g1, g3}
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

	g1, g2, g3 := trn.Game{"A", 1, "B", 0}, trn.Game{"B", 2, "C", 2}, trn.Game{"C", 3, "A", 4}

	gd := GamesData{dbPool}
	gd.Save(&g1)
	gd.Save(&g2)
	gd.Save(&g3)

	expected := []trn.Game{g1, g2, g3}
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
		panic(fmt.Sprintf("Unable to create the database: %v", err))
	}

	dbPool, err = pgxpool.Connect(context.Background(), testDbUrl)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to the database: %v", err))
	}
	defer dbPool.Close()

	testExitCode = m.Run()
}

func asMap(games []trn.Game) map[trn.Game]bool {
	m := make(map[trn.Game]bool)
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

	conn2, err := pgx.Connect(context.Background(), testDbUrl)
	if err != nil {
		return fmt.Errorf("Unable to connect to %v: %v", DB_NAME, err)
	}

	_, err = conn2.Exec(context.Background(), "CREATE TABLE games (team_a varchar(40) NOT NULL, score_a int NOT NULL, team_b varchar(40) NOT NULL, score_b int NOT NULL);")
	if err != nil {
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
		panic(fmt.Sprintf("Unable to delete all games: %v", err))
	}
}

func getEnv(name, defaultValue string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return defaultValue
}
