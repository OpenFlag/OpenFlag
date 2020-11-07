create table if not exists rules
(
    id            bigint      auto_increment primary key,
    description   text        not null,
    flag          text        not null,
    tag           text,
    entity        text        not null,
    rule          text        not null,
    apply_value   text        not null,
    default_value text        not null,
    created_at    timestamp   not null default now(),
    deleted_at    timestamp   null
);
