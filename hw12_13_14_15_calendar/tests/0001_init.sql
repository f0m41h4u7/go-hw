-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
  uuid TEXT, 
  title TEXT,
  start TEXT,
  end TEXT,
  description TEXT,
  ownerid TEXT,
  notifyin TEXT
);
-- +goose StatementEnd
-- +goose StatementBegin
INSERT INTO events (uuid, title, start, end, description, ownerid, notifyin) VALUES ("9bed7c53-c3bd-4f7e-92d1-5d98c04fb83a", "event", "2020-07-29T20:00:00", "2020-07-29T22:00:00", "test", "9bed7c53-c3bd-4f7e-92d1-5d98c04fb83a", "70h")
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
