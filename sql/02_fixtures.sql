CREATE TABLE seasons (
  id TEXT PRIMARY KEY
);
CREATE TABLE rounds (
  id TEXT PRIMARY KEY,
  season_id TEXT REFERENCES seasons(id) NOT NULL
);
CREATE TABLE fixtures (
  id INTEGER PRIMARY KEY,
  round_id TEXT REFERENCES rounds(id) NOT NULL,
  home_team VARCHAR(3) REFERENCES teams(key) NOT NULL,
  away_team VARCHAR(3) REFERENCES teams(key) NOT NULL,
  date_time TIMESTAMP NOT NULL,
  UNIQUE (round_id, home_team, away_team)
);
CREATE TABLE results (
  fixture_id INTEGER PRIMARY KEY REFERENCES fixtures(id),
  home_goals INTEGER NOT NULL,
  away_goals INTEGER NOT NULL
);
