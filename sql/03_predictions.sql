CREATE TABLE players (
  email VARCHAR(100) PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  active BOOLEAN DEFAULT TRUE
);
CREATE TABLE predictions (
  player VARCHAR(100) NOT NULL REFERENCES players(email),
  fixture_id INTEGER NOT NULL REFERENCES fixtures(id),
  home_goals INTEGER NOT NULL,
  away_goals INTEGER NOT NULL,
  PRIMARY KEY (player, fixture_id)
);
