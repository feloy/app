BEGIN;

CREATE TABLE IF NOT EXISTS comments (
   id serial primary key,
   article_id int not null,
   author_id int not null,
   body text not null,
   created_at timestamptz not null default now(),
   constraint fk_article foreign key(article_id) references articles(id) on delete cascade,
   constraint fk_author foreign key(author_id) references users(id) on delete cascade
);

COMMIT;