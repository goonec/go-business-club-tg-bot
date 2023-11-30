create table if not exists "user"(
    id bigint unique,
    tg_username text,
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

create table if not exists user_resident(
    id int generated always as identity,
    user_id bigint unique,
    tg_username text,
    primary key (id),
    foreign key (user_id)
            references "user" (id),
    foreign key (tg_username)
        references resident (tg_username)
);

select id, user_id, tg_username from user_resident where user_id = 2 or tg_username = 'resident3';

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

INSERT INTO user_resident (user_id, tg_username)
VALUES
    (2, 'resident2'),
    (3, 'resident3');

INSERT INTO user_resident (tg_username)
VALUES
    ('resident1');