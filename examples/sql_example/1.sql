create table video (
    id int primary key,
    title text,
    author_id int
);

create table author (
    id int primary key,
    name text
);

-- IN, а не "=": name не уникален, и если авторов с таким именем несколько,
-- скалярный подзапрос с "=" упадет (PostgreSQL, MySQL) или молча вернет
-- неверное число (SQLite берет только первую строку подзапроса).
SELECT COUNT(*)
FROM video
WHERE author_id IN (
        SELECT id
        FROM author
        WHERE name = 'искомое_имя'
    );
