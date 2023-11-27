create extension if not exists "uuid-ossp";

create type role as enum ('user','admin');

-- ФИО, возраст, регион, сфера деятельности, название компании, чем может быть полезен, увлечения, фото
create table resident(
    id bigint,
    firstname varchar(50) not null,
    lastname varchar(50) not null,
    patronymic varchar(50) not null,
    age int not null,
    region varchar(50) null,
    work_activity varchar(50) null,
    company_name varchar(50) null,
    advantage varchar(100) null,
    hobie varchar(150) null,
    photo_file_id varchar(150) null,
    primary key (id)
);

create table resident_hobie(
    id int,
    resident_id bigint,
    hobie varchar(50),
    primary key (id),
    foreign key (resident_id)
        references resident (id)
);

create table resident_role(
    id int,
    resident_id bigint,
    role role default 'user',
    primary key (id),
    foreign key (resident_id)
         references resident (id)
);