create table if not exists "user"(
    id bigint unique,
    tg_username text unique not null,
    create_at timestamp,
    role role default 'user',
    primary key (id)
);

create table if not exists resident(
    id int generated always as identity,
    tg_username text,
    firstname varchar(50) not null,
    lastname varchar(50) not null,
    patronymic varchar(50) not null,
    resident_data text null,
    photo_file_id varchar(150) null,
    primary key (id),
    foreign key (tg_username)
    references "user" (tg_username)
);


