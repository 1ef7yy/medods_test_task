-- +goose Up
-- +goose StatementBegin
ALTER TABLE tokens
ADD CONSTRAINT guid_fk
FOREIGN KEY (guid)
REFERENCES users(guid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tokens
DROP CONSTRAINT guid_fk;
-- +goose StatementEnd
