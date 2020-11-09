create table if not exists flags
(
    id              bigserial,
    tags            jsonb,
    description     text         not null,
    flag            varchar(255) not null,
    segments        jsonb        not null,
    created_at      timestamp    not null default now(),
    deleted_at      timestamp,
    primary key (id)
);

create index flags_tags_idx on flags(tags);
create index flags_flag_idx on flags(flag);
create index flags_created_at_idx on flags(created_at);
create index flags_deleted_at_idx on flags(deleted_at);
