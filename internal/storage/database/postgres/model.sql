drop table if exists Users CASCADE;
drop table if exists News CASCADE;
drop table if exists News2Tags CASCADE;
drop table if exists Tags CASCADE;
drop table if exists Files CASCADE;

create table Users (
    id text primary key,
    name text not null,
    email text not null,
    passwordHash text not null,
    role text not null,
    registeredFrom text not null,
    refreshToken text,
    refreshTokenExpired int,
    unique(email)
);

create table News (
    id text primary key,
    title text not null,
    author text not null,
    active boolean not null,
    activeFrom integer not null,
    text text not null,
    textJSON text not null,
    userId text references Users(id),
    isImportant boolean not null
);

create table Tags (
    id text primary key,
    name text,
    unique(name)
);

create table News2Tags (
    news_id text references News (id),
    tag_id text references Tags (id)
);

create table Files (
    id text primary key,
    name text not null,
    ext text not null,
    base64 text not null,
    dateCreate integer not null,
    userId text references Users(id)
);