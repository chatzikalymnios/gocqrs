-- create user accounts;
-- create database accounts;
-- grant all privileges on database accounts to accounts;

create table events
(
    event_id     varchar(20) primary key,
    event_type   text        not null,
    event_data   jsonb       not null,
    entity_type  text        not null,
    entity_id    varchar(20) not null,
    published    boolean     not null default false
);

create index events_idx on events (entity_type, entity_id, event_id);
create index events_published_idx on events (published, event_id);

create table entities
(
    entity_type    text,
    entity_id      varchar(20),
    entity_version integer not null,
    primary key (entity_type, entity_id)
);
