-- +goose Up
-- +goose StatementBegin
INSERT INTO users(guid, email)
VALUES  ('acd5146e-c98d-4ca8-a5c8-5139324b1e18', 'test1@mail.com'),
        ('66dc497d-39fc-4c80-ac74-504701365774', 'test2@mail.com'),
        ('61269951-e406-40f2-a2ac-1243b64c6310', 'test3@mail.com'),
        ('39710db9-1efb-42bd-b239-28adf0683bbf', 'test4@mail.com'),
        ('fca18db6-b094-41c9-ac23-6fc77bab41e3', 'test5@mail.com');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE guid IN ('acd5146e-c98d-4ca8-a5c8-5139324b1e18', '66dc497d-39fc-4c80-ac74-504701365774', '61269951-e406-40f2-a2ac-1243b64c6310', '39710db9-1efb-42bd-b239-28adf0683bbf', 'fca18db6-b094-41c9-ac23-6fc77bab41e3')
-- +goose StatementEnd
