create table requests
(
    "Host"        varchar   not null,
    "requestData" varchar   not null,
    id            bigserial not null
        constraint requests_pk
            primary key
);