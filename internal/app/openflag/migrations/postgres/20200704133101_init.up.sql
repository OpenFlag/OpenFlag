create table if not exists rules
(
    id            bigserial,
    description   text         not null,
    flag          text         not null,
    tag           text,
    entity        text         not null,
    rule          jsonb        not null,
    apply_value   text         not null,
    default_value text         not null,
    created_at    timestamp    not null default now(),
    deleted_at    timestamp,
    primary key (id)
);

create index rules_flag_idx on rules(flag);
create index rules_tag_idx on rules(tag);
create index rules_entity_idx on rules(entity);
create index rules_created_at_idx on rules(created_at);
create index rules_deleted_at_idx on rules(deleted_at);
