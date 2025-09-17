-- +goose Up
create table public.go_template_json (
	id serial not null,
	data jsonb
);

-- +goose Down
drop table public.go_template_json;