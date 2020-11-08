create table if not exists flags
(
    id              bigint        auto_increment primary key,
    tag             varchar(255),
    description     text          not null,
    flag            varchar(255)  not null,
    segments        json          not null,
    default_variant json,
    created_at      timestamp     not null default now(),
    deleted_at      timestamp     null
);

create index flags_tag_idx on flags(tag);
create index flags_flag_idx on flags(flag);
create index flags_created_at_idx on flags(created_at);
create index flags_deleted_at_idx on flags(deleted_at);
