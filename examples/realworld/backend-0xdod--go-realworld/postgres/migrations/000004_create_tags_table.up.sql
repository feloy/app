BEGIN;

CREATE TABLE IF NOT EXISTS tags (
    id serial primary KEY,
    name citext unique NOT NULL
);

CREATE TABLE IF NOT EXISTS article_tags (
    article_id INT NOT NULL,
    tag_id int not null,
    primary key(article_id, tag_id),
    constraint fk_article foreign key(article_id) references articles(id) on delete cascade ,
    constraint fk_tag foreign key(tag_id) references tags(id) on delete cascade
);

COMMIT;