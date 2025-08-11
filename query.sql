-- name: CreateRequest :exec
INSERT INTO requests (id, method, content, source_ip, timestamp, request_path)
VALUES ($1, $2, $3,
        $4, $5, $6);

-- name: CreateResponse :exec
insert into responses (requestid, response_content, response_code, timestamp)
values ($1, $2, $3, $4);

-- name: CreateRequestHeaders :copyfrom
insert into request_headers (request_id, name, value)
values ($1, $2, $3);

-- name: CreateRequestQueryParameters :copyfrom
insert into request_query_parameters (request_id, name, value)
values ($1, $2, $3);

-- name: CreateResponseHeaders :copyfrom
insert into response_headers (requestid, name, value)
values ($1, $2, $3);

