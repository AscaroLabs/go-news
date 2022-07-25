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
    registeredFrom integer not null,
    unique(email)
);

create table RefreshSessions (
    id serial primary key,
    userId text references Users(id),
    refreshToken text not null,
    expiresIn integer not null,
    createdAt integer not null,
    unique(refreshToken)
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
    isImportant boolean not null,
    tags text[],
    files text[]
);

create table Tags (
    id text primary key,
    name text,
    unique(name)
);

-- create table News2Tags (
--     news_id text references News (id),
--     tag_id text references Tags (id)
-- );

create table Files (
    id text primary key,
    name text not null,
    ext text not null,
    base64 text not null,
    dateCreate integer not null,
    userId text references Users(id)
);

-- create table News2Files (
--     news_id text references News (id),
--     files_id text references Files (id)
-- );

insert into Users (id,name,email,passwordHash,role,registeredFrom) values 
('1','ilya','qwe@qwe.ru','qqq111','dealer',123),
('2','vova','email@qq.q','pas22','member',251),
('3','pipka','qwe,,v@a','qwe2222','moder',5555),
('4','pupa','rr@email.z','xhvi323-s','guest',222);

insert into News (id,title,author,active,activeFrom,text,textJSON,userId,isImportant, tags, files) values 
('322','Hello','pupa',true,220,'Hello i am here','{"text":"Hello i am here"}',4,false,'{"7001"}','{"6903"}'),
('12','Alive!','ilya',false,8088,'I am alive!','{"text":"I am alive!"}',1,true,'{"7003","7002"}','{"6902"}'),
('17','Taxes','vova',true,2000,'PAY TAXES!!!11!','{"text":"PAY TAXES!!!11!"}',2,true,'{"7002","7003","7004"}','{}'),
('18','Arch','ilya',true,2022,'Clean Arch is garbage','{"text":"Clean Arch is garbage"}',1,false,'{}','{}');

insert into Tags (id, name) values
('7001', 'Fun'),
('7002', 'Hard'),
('7003', 'Truth'),
('7004', 'No choice');

-- insert into News2Tags (news_id,tag_id) values 
-- ('322','7001'),
-- ('12','7003'),
-- ('12','7002'),
-- ('17','7002'),
-- ('17','7003'),
-- ('17','7004');

insert into Files (id,name,ext,base64,dateCreate,userId) values
('6900','Doc42','txt','base01',126,'1'),
('6901','Raschetnaya_vypiska','docx','base02',1,'4'),
('6902','Annot_','pptx','base03',256,'2'),
('6903','Screen125','png','base03',188,'1'),
('6904','main','cpp','base01',1024,'1');

-- insert into News2Files (news_id,files_id) values 
-- ('322','6903'),
-- ('12','6902');