CREATE TABLE users (
  email TEXT Primary Key,
  first_name TEXT,
  last_name TEXT,
  password TEXT not null
);

CREATE TABLE board (
  id Text PRIMARY KEY,
  start_date Date not null,
  end_date Date not null,
  is_active Boolean not null,
  user_id TEXT not null,
  FOREIGN KEY (user_id) REFERENCES users (email)
);

CREATE TABLE habit (
  id Text PRIMARY KEY,
  title TEXT not null,
  description TEXT not null,
  content Text[] not null,
  board_id TEXT not null REFERENCES board on delete cascade,
  FOREIGN KEY (board_id) REFERENCES board (id)
);

CREATE TABLE note (
  id Text PRIMARY KEY,
  body TEXT not null,
  board_id TEXT not null REFERENCES board on delete cascade,
  FOREIGN KEY (board_id) REFERENCES board (id)
);
