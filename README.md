# Tournament Golang Microservice Vertical Slice

This repository contains a vertical slice of a golang sample microservice that
tracks tournament games statistics.

## Code Packages

The slice has the following components:

* `tournament` package - contains the service business logic and data access
  interface `Games`
* `db` package - contains PostgreSQL database persitence layer that implements
  `tournament.Games` interface
* `gen` package - contains [go-swagger](https://github.com/go-swagger/go-swagger)
  generated code from the specification located in `api/swagger.yml`
  * `/games` API endpoint requires authentication, `/stats` and `/stats/{team}`
    can be accessed anonymously
* `cmd/tournament/main.go` - the main microservice file that stiches all the
  elements together

## Implementation Principles

The following implementation principles were used:

* `tournament` and `db` packages were implemented using test-driven development
* [Go Watch](https://github.com/mitranim/gow) has been used to run tests each
  time the source code was updated, using command `gow test .`
* Database schema is managed using [golang-migrate](https://github.com/golang-migrate/migrate)
  package with migrations present in `sql` folder
* Microservice has been packaged as Docker Compose stack with definitions in
  `deploy` folder
* The Docker image is built based on `scratch` and so the application must be
  compiled completely statically
* Project layout follows https://github.com/golang-standards/project-layout

## Running the Microservice

To start the microservice (API and database):

```shell
docker compose up -f deploy/docker-compose.yml
```

It is also possible to start API and database separately:
```shell
tools/run-postgres.sh`
```
```shell
DB_URL=postgres://postgres:secret@localhost:5432/tournament?sslmode=disable \
  go run ./cmd/tournament/main.go --port 3000
```

To record a game score:

```shell
curl -X POST http://localhost:3000/games \
  -H 'Content-Type: application/slawekzachcial.tournament.v1+json' \
  -H 'x-token: qwerty' \
  -d '{"TeamA": "C", "ScoreA": 1, "TeamB": "A", "ScoreB": 0}'
```

To get team statistics:

```shell
curl -s http://localhost:3000/stats/A \
  -H 'Content-Type: application/slawekzachcial.tournament.v1+json
```

To get all statistics:

```shell
curl -s http://localhost:3000/stats \
  -H 'Content-Type: application/slawekzachcial.tournament.v1+json'
```

