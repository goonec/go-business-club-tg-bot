create type role as enum ('user','admin');

create table if not exists "user"(
    id bigint unique,
    tg_username text unique,
    create_at timestamp,
    role role default 'user',
    primary key (id)
);

create table if not exists resident(
    id int generated always as identity,
    tg_username text unique not null,
    firstname varchar(50) not null,
    lastname varchar(50) not null,
    patronymic varchar(50) not null,
    resident_data text null,
    photo_file_id varchar(150) null,
    primary key (id)
);

create table if not exists business_cluster
(
    id int generated always as identity,
    name varchar(200) unique ,
    primary key (id)
);

insert into business_cluster (name) values ('1');
insert into business_cluster (name) values ('2');
insert into business_cluster (name) values ('3');
insert into business_cluster (name) values ('4');
insert into business_cluster (name) values ('5');

create unique index id on business_cluster (name);

create table if not exists business_cluster_resident
(
    id_business_cluster int,
    id_resident int,
    primary key (id_business_cluster,id_resident),
    foreign key (id_business_cluster)
        references business_cluster (id) on delete cascade,
    foreign key (id_resident)
        references resident (id) on delete cascade
);


--select * from "user" where tg_username IN (select tg_username from resident);
