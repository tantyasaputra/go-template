-- +goose Up
CREATE TABLE go_template_examples(
    id serial not null primary key,
    "name" varchar(255) not null,
    created_at timestamptz not null default current_timestamp,
    deleted_at timestamptz 
);
CREATE TABLE go_template_samples(
    id serial not null primary key,
    "name" varchar(255) not null,
    created_at timestamptz not null default current_timestamp,
    deleted_at timestamptz 
);

-- +goose Down
DROP TABLE go_template_examples;
DROP TABLE go_template_samples;