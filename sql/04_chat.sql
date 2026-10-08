CREATE TABLE chat_admins (
  email VARCHAR(255) PRIMARY KEY REFERENCES players(email)
);
