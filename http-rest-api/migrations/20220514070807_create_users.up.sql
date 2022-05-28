CREATE TABLE users (
    id bigserial not null primary key,
    username varchar not null unique,
    email varchar not null unique,
    password varchar not null
);

CREATE TABLE posts (
    id bigserial not null primary key,
    user_id bigserial not null,
    username varchar not null,
    CreatedDate DATE,
    caption varchar not null
);