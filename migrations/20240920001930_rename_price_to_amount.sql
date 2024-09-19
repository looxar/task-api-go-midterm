-- +goose Up
ALTER TABLE items RENAME COLUMN price TO amount;
-- +goose StatementBegin
SELECT
    'Renamed column price to amount in items table';
-- +goose StatementEnd

-- +goose Down
ALTER TABLE items RENAME COLUMN amount TO price;
-- +goose StatementBegin
SELECT
    'Reverted column amount to price in items table';
-- +goose StatementEnd
