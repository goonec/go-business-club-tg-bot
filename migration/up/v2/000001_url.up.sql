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

--select * from "user" where tg_username IN (select tg_username from resident);

delete from resident where tg_username = 'n3ksmrnv';
delete from "user" where id = 7;

INSERT INTO "user" (id, tg_username, create_at, role)
VALUES
    (1, 'user1', '2023-11-30 12:00:00', 'admin'),
    (2, 'user2', '2023-11-30 12:15:00', 'user'),
    (3, 'user3', '2023-11-30 12:30:00', 'user');

INSERT INTO resident (tg_username, firstname, lastname, patronymic, resident_data, photo_file_id)
VALUES
    ('resident1', 'Иван', 'Иванов', 'Иванович', 'Данные о резиденте 1', 'файл1'),
    ('resident2', 'Петр', 'Петров', 'Петрович', 'Данные о резиденте 2', 'файл2'),
    ('resident3', 'Мария', 'Маринина', 'Марковна', 'Данные о резиденте 3', 'файл3');
