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

-- Insert test data into "user" table
INSERT INTO "user" (id, tg_username, create_at, role)
VALUES
    (1, 'john_doe', '2023-01-01', 'admin'),
    (2, 'jane_smith', '2023-01-02', 'user'),
    (3, 'bob_jones', '2023-01-03', 'user');

-- Insert test data into "resident" table
INSERT INTO resident (tg_username, firstname, lastname, patronymic, resident_data, photo_file_id)
VALUES
    ('john_doe', 'John', 'Doe', 'Pat', 'Some resident data for John', 'photo1.jpg'),
    ('jane_smith', 'Jane', 'Smith', 'Lee', 'Some resident data for Jane', 'photo2.jpg'),
    ('bob_jones', 'Bob', 'Jones', 'Kane', 'Some resident data for Bob', NULL);


