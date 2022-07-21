create table Users (
    id text primary key,
    name text not null,
    email text not null,
    login boolean not null,
    passwordHash integer not null,
    role text not null,
    registeredFrom text not null,
    refreshToken text not null,
    refreshTokenExpired boolean not null
)

create table News (
    id text primary key,
    title text not null,
    author text references Users(name),
    active boolean not null,
    activeFrom integer not null,
    text text not null,
    textJSON text not null,
    userId text references Users(id),
    isImportant boolean not null
)

create table News2Tags (
    news_id text references News (id),
    tag_id text references Tag (id)
)

create table Tags (
    id text primary key,
    name text,
    unique(name)
)

create table Files (
    id text primary key,
    name text not null,
    ext text not null,
    base64 text not null,
    dateCreate integer not null,
    userId string references Users(id)
)