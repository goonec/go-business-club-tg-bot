create table if not exists "user"(
    id bigint unique,
    tg_username text,
    create_at timestamp,
    role role default 'user',
    primary key (id),
    foreign key (tg_username)
        references resident (tg_username)
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

create table if not exists

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


UPDATE resident set photo_file_id = 'AgACAgIAAxkBAAEQb89lZh-mLafu3vScLH0AAc4yXIV9sZYAAorXMRsgWjFL2S4QnfRotRoBAAMCAANzAAMzBA'
where firstname = 'John';

INSERT INTO resident (tg_username, firstname, lastname, patronymic, resident_data, photo_file_id)
VALUES
    ('vovchik', 'Владимир', 'Кочнев', 'К', 'Максим Тришин
Дата рождения: 25.03.1989 г.
Город: Москва
Дата вступления: 31 июля 2023 г.

🖥Основатель компании BAUM ZINDECH (бытовая техника под собственным брендом), компания на рынке РФ третий год. Обогнали по многим показателям крупные мировые бренды по отдельным видам товаров. Сайт: www.BaumZindech.ru
💰Текущий оборот: 10млн/мес, интересная рентабельность
Количество сотрудников: 6 чел +
🙌Опыт: в бизнесе 10 лет плюс опыт в правительственных структурах. Опыт в построение и управлении: Kirby, Desheli, нерудные материалы, АвтоДор, тендеры, производство в Китае

🏍🌊🌍Увлечения: мотоциклы (будет круто сделать пробег или даже мотоклуб в клубе), путешествия и море, психология, эстетика и промышленный дизайн, маркетинг. Воспитание сыновей с правильными ценностями. Христианство.
🎯Ценности в команде = мои ценности: 1)интеллект 2)эмоциональный интеллект 3)личная эффективность

👥В клуб вступил для сильного окружения, масштабирования. Для обмена опытом и энергией. Ищу единомышленников, друзей, возможно инвесторов. Особенно интересны форум и мастермайнды.

💪Сильные стороны и польза:
интернет коммерция, Китай, маркетплейсы, промышленный дизайн, СТМ, продажи, диджитал маркетинг, контент и продакшн

Tel: +79254837826
Telegram: @mt_bzgroup
Instagram: https://instagram.com/trishin.maxim?igshid=NTc4MTIwNjQ2YQ==', 'AgACAgIAAxkBAAEQb89lZh-mLafu3vScLH0AAc4yXIV9sZYAAorXMRsgWjFL2S4QnfRotRoBAAMCAANzAAMzBA');