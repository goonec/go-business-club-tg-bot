create type role as enum ('user','admin');

create table if not exists "user"(
    id bigint unique,
    tg_username text unique,
    created_at timestamp,
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


select r.firstname, substring(r.lastname,1,1), substring(r.patronymic,1,1) from resident r
JOIN business_cluster_resident bcr on bcr.id_resident = r.id
JOIN business_cluster bc on bc.id = bcr.id_business_cluster
WHERE bc.id = 2;

