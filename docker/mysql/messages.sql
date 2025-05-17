CREATE DATABASE IF NOT EXISTS messages;
USE messages;

CREATE TABLE IF NOT EXISTS messages
(
    uuid       char(36)                            not null
        primary key,
    user_id    int                                 null,
    chat_id    int                                 not null,
    text       text                                not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    updated_at timestamp                           null,
    deleted_at timestamp                           null
);
