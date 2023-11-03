-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE items(
  id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
  name text NOT NULL,
  warehouse_location text NOT NULL,
  weight double NOT NULL,
  
  description text,
  attributes jsonb,
  tags jsonb,
  delivery_time interval,

  last_updated timestamp without timezeone NOT NULL DEFAULT current_timestamp NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;

DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd
