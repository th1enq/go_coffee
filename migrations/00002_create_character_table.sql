-- +goose Up
-- +goose StatementBegin
CREATE TABLE characters (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    name VARCHAR(255) NOT NULL UNIQUE,
    rarity INT NOT NULL,
    region VARCHAR(255) NOT NULL,
    vision VARCHAR(255) NOT NULL,
    weapon_type VARCHAR(255) NOT NULL,
    constellation VARCHAR(255) NOT NULL,
    birthday TIMESTAMP WITH TIME ZONE NOT NULL,
    affilliation VARCHAR(255) NOT NULL,
    release_date TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_characters_deleted_at ON characters(deleted_at);
CREATE INDEX idx_characters_name ON characters(name);
CREATE INDEX idx_characters_region ON characters(region);
CREATE INDEX idx_characters_vision ON characters(vision);
CREATE INDEX idx_characters_rarity ON characters(rarity);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS characters;
-- +goose StatementEnd
