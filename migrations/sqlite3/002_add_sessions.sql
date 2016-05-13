-- +migrate Up
create table sessions (
    token text not null primary key,
    user text not null,
    valid_until datetime not null
);

-- +migrate Down
drop table sessions;
