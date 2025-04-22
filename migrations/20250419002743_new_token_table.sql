-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tokens (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    refresh_hash TEXT UNIQUE NOT NULL,
    guid UUID UNIQUE NOT NULL,
    generation INT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tokens;
-- +goose StatementEnd
