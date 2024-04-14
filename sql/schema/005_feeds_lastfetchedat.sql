-- +goose Up
ALTER TABLE feeds ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose Down
ALTER TABLE feeds DROP COLOUMN last_fetched_at;