BEGIN;

CREATE TABLE IF NOT EXISTS favorites (
    article_id int not null,
    user_id int not null,
    created_at timestamptz not null default now(),
    primary key(article_id, user_id),
    constraint fk_article foreign key(article_id) references articles(id) on delete cascade,
    constraint fk_user foreign key(user_id) references users(id) on delete cascade
);

COMMIT;