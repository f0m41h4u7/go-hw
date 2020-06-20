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

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
