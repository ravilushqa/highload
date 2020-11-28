create database if not exists app;
use app;

create table if not exists users
(
    id         int AUTO_INCREMENT primary key,
    email      varchar(255)                     not null,
    password   char(60)                         not null,
    firstname  varchar(255)                     not null,
    lastname   varchar(255)                     not null,
    birthday   DATE,
    sex        enum ('male', 'female', 'other') not null,
    interests  varchar(1024),
    city       varchar(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP                        null
) ENGINE = INNODB
  DEFAULT CHARSET = utf8
  AUTO_INCREMENT = 1;

create unique index users_email_uindex
    on users (email);

create table if not exists friends
(
    user_id   int not null,
    friend_id int not null,
    approved  bool default false,
    constraint friends_pk
        primary key (user_id, friend_id),
    constraint friends_users_id_fk
        foreign key (user_id) references users (id),
    constraint friends_users_id_fk_2
        foreign key (friend_id) references users (id)
);


create table messages
(
    id         int AUTO_INCREMENT primary key,
    user_id    int                                 null,
    chat_id    int                                 not null,
    text       text                                not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    updated_at timestamp default null              null,
    deleted_at timestamp default null              null
);

