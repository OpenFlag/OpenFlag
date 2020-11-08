create table if not exists flags
(
    id              bigserial,
    tag             varchar(255),
    description     text         not null,
    flag            varchar(255) not null,
    segments        jsonb        not null,
    default_variant jsonb,
    created_at      timestamp    not null default now(),
    deleted_at      timestamp,
    primary key (id)
);

create index flags_tag_idx on flags(tag);
create index flags_flag_idx on flags(flag);
create index flags_created_at_idx on flags(created_at);
create index flags_deleted_at_idx on flags(deleted_at);
