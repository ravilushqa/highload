create database if not exists app;
use app;

create table users
(
    id         int auto_increment
        primary key,
    email      varchar(255)                        not null,
    password   char(60)                            not null,
    firstname  varchar(255)                        not null,
    lastname   varchar(255)                        not null,
    birthday   date                                null,
    sex        enum ('male', 'female', 'other')    not null,
    interests  varchar(1024)                       null,
    city       varchar(255)                        null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    deleted_at timestamp                           null,
    constraint users_email_uindex
        unique (email)
)
    charset = utf8;


-- auto-generated definition
create table if not exists friends
(
    user_id   int                  not null,
    friend_id int                  not null,
    approved  tinyint(1) default 0 null,
    primary key (user_id, friend_id),
    constraint friends_users_id_fk
        foreign key (user_id) references users (id),
    constraint friends_users_id_fk_2
        foreign key (friend_id) references users (id)
);



-- auto-generated definition
create table if not exists chat_users
(
    id         int auto_increment
        primary key,
    user_id    int       not null,
    chat_id    int       not null,
    deleted_at timestamp null
);

-- auto-generated definition
create table if not exists chats
(
    id         int auto_increment
        primary key,
    name       text                     null,
    type       enum ('dialog', 'group') null,
    deleted_at timestamp                null
);

-- auto-generated definition
create table if not exists posts
(
    id         int auto_increment
        primary key,
    user_id    int                                 not null,
    text       text                                not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    deleted_at timestamp                           null,
    constraint posts_users_id_fk
        foreign key (user_id) references users (id)
);


############### sharding instances bellow


-- auto-generated definition
create table if not exists messages
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

