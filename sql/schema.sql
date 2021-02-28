DROP DATABASE tournament;

CREATE DATABASE tournament;

\c tournament

CREATE TABLE games (
    team_a varchar(40) NOT NULL,
    score_a int NOT NULL,
    team_b varchar(40) NOT NULL,
    score_b int NOT NULL
);
