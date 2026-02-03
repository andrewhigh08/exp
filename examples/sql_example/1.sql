create table video (
    id int primary key
    title text
    author_id int
)

create table author (
    id int primary key
    name text
)

SELECT COUNT(*)
FROM video
WHERE author_id = (
        SELECT id
        FROM author
        WHERE name = "искомое_имя"
    )