create extension if not exists "uuid-ossp";

create type role as enum ('user','admin');

-- ФИО, возраст, регион, сфера деятельности, название компании, чем может быть полезен, увлечения, фото
create table if not exists resident(
    id int generated always as identity,
    tg_username text,
    phone_number varchar(11) null,
    firstname varchar(50) not null,
    lastname varchar(50) not null,
    patronymic varchar(50) not null,
    age int not null,
    region varchar(100) null,
    work_activity varchar(200) null,
    company_name varchar(200) null,
    advantage varchar(300) null,
    hobie varchar(500) null ,
    photo_file_id varchar(150) null,
    primary key (id),
    foreign key (tg_username)
         references "user" (tg_username)
);

create table if not exists "user"(
     id bigint unique,
     tg_username text unique not null,
     create_at timestamp,
     role role default 'user',
     primary key (id)
);

select r.id,r.tg_username,r.phone_number,r.firstname,r.patronymic,r.age,r.region,r.work_activity,
       r.company_name,r.advantage,r.hobie,r.photo_file_id
from resident r
order by r.id;

-- Insert test data into "user" table
INSERT INTO "user" (id, tg_username, create_at, role)
VALUES
    (1, 'john_doe', '2023-01-01', 'admin'),
    (2, 'jane_smith', '2023-01-02', 'user'),
    (3, 'bob_jones', '2023-01-03', 'user');

-- Insert test data into "resident" table
INSERT INTO resident (tg_username, phone_number, firstname, lastname, patronymic, age, region, work_activity, company_name, advantage, hobie, photo_file_id)
VALUES
    ('john_doe', '12345678901', 'John', 'Doe', 'Pat', 30, 'City1', 'Developer', 'ABC Inc.', 'Quick learner', 'Reading, Coding', 'photo1.jpg'),
    ('jane_smith', '98765432109', 'Jane', 'Smith', 'Lee', 25, 'City2', 'Designer', 'XYZ Corp.', 'Creative thinker', 'Painting, Traveling', 'photo2.jpg'),
    ('bob_jones', NULL, 'Bob', 'Jones', 'Kane', 35, 'City3', 'Manager', '123 Company', 'Detail-oriented', 'Sports, Cooking', NULL);


-- create table resident_hobie(
--     id int,
--     resident_id bigint,
--     hobie varchar(50),
--     primary key (id),
--     foreign key (resident_id)
--         references resident (id)
-- );
--
-- create table resident_role(
--     id int,
--     resident_id bigint,
--     primary key (id),
--     foreign key (resident_id)
--          references resident (id)
-- );

