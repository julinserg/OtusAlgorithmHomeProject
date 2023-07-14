-- +goose Up
CREATE table document_source (
    id               serial primary key,
    url              text unique not null, 
    title            text not null,
    data             text not null
);

-- +goose Down
drop table document_source;
