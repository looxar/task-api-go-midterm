-- +goose Up
ALTER TABLE items ALTER COLUMN amount TYPE int USING amount::integer;
-- +goose StatementBegin
SELECT
    'Changed amount column type to int in items table';
-- +goose StatementEnd

-- +goose Down
ALTER TABLE items ALTER COLUMN amount TYPE real USING amount::real;
-- +goose StatementBegin
SELECT
    'Reverted amount column type to real in items table';
-- +goose StatementEnd