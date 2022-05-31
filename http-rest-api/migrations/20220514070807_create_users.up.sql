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
    caption text not null,
    likes bigserial not null
);

CREATE TABLE users_liked (
    post_id bigserial not null,
    user_id bigserial not null
);

CREATE TABLE comment (
    id bigserial not null primary key,
    user_id bigserial not null,
    username varchar not null,
    post_id bigserial not null,
    created_date DATE,
    text text not null
);