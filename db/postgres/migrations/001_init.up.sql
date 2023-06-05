BEGIN;

create table if not exists chats
(
    id          serial,
    name        varchar   not null
    );

create table if not exists users
(
    id          serial,
    name        varchar   not null
);

create table if not exists users_chats
(
    user_id          integer not null,
    chat_id          integer not null
);


create table if not exists messages
(
    id          serial,
    chat_id integer not null,
    sender_id        integer   not null,
    text varchar,
    created_at  timestamp not null

);


insert into chats (name) values ('global');
insert into chats (name) values ('dev');
insert into chats (name) values ('updates');

COMMIT;