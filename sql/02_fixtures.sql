CREATE TABLE fixtures (
  id INTEGER PRIMARY KEY,
  home_team VARCHAR(3) REFERENCES teams(abbr) NOT NULL,
  away_team VARCHAR(3) REFERENCES teams(abbr) NOT NULL,
  date_time TIMESTAMP NOT NULL
);
CREATE TABLE results (
  fixture_id INTEGER PRIMARY KEY REFERENCES fixtures(id),
  home_goals INTEGER NOT NULL,
  away_goals INTEGER NOT NULL
);
