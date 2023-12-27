DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role') THEN
            CREATE TYPE role AS ENUM ('user', 'admin');
        END IF;
    END $$;


create table if not exists "user"(
    id bigint unique,
    tg_username text unique,
    created_at timestamp,
    role role default 'user',
    primary key (id)
);

create table if not exists resident(
    id int generated always as identity,
    tg_username text null,
    firstname varchar(50) not null,
    lastname varchar(50) not null,
    patronymic varchar(50) null,
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

create table if not exists schedule(
    id int generated always as identity,
    photo_file_id varchar(150) null,
    created_at timestamp default current_timestamp not null,
    primary key (id)
);

create table if not exists service(
  id int generated always as identity,
  name varchar(200),
  primary key (id)
);

create table if not exists service_describe(
    id int generated always as identity,
    id_service int,
    name text not null,
    photo_file_id varchar(150) null,
    describe text not null,
    foreign key (id_service)
        references service (id) on delete cascade,
    primary key (id)
);


create table if not exists feedback(
    id int generated always as identity,
    message text,
    type varchar(50) not null,
    created_at timestamp default current_timestamp not null,
    tg_username text,
    primary key (id)
);

create table if not exists pptx(
  pptx_file_id varchar(150) unique
);
