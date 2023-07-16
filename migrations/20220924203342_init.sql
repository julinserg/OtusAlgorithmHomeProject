-- +goose Up
CREATE table document_source (
    id               serial primary key,
    url              text unique not null, 
    title            text not null,
    data             text not null
);

CREATE table document_invert_index (
    word             text primary key,
    documents_list   bytea not null     
);

-- +goose Down
drop table document_source;
drop table document_invert_index;