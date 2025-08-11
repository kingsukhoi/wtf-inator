-- migrate:up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "requests"
(
    id           uuid primary key,
    method       text        not null,
    content      bytea,
    source_ip    text        not null,
    timestamp    timestamptz not null default now(),
    request_path text not null
);

create table if not exists "responses"
(
    requestId        uuid primary key,
    response_content bytea,
    response_code    int         not null,
    timestamp        timestamptz not null default now()
);

create table if not exists response_headers (
    requestId    uuid not null,
    name         text not null,
    value        text,
    foreign key (requestId) references responses (requestId)
);

create table if not exists request_headers
(
    request_id uuid not null,
    name       text not null,
    value      text,
    foreign key (request_id) references requests (id)
);
create table if not exists request_query_parameters
(
    request_id uuid not null,
    name       text not null,
    value      text,
    foreign key (request_id) references requests (id)
);

-- migrate:down

drop table if exists request_query_parameters;
drop table if exists request_headers;
drop table if exists requests;
drop table if exists response_headers;
drop table if exists responses;