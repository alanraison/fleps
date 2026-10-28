-- Maps a team to its slug in the football-logos.cc catalog (https://football-logos.cc/england/{logo_slug}/).
CREATE TABLE team_logos (
  team_key VARCHAR(3) PRIMARY KEY REFERENCES teams(key) ON DELETE CASCADE,
  logo_slug VARCHAR(80) NOT NULL
);
