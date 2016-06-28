-- +migrate Up
create table people (
    id integer not null primary key,
    version int not null,
    created_at timestamp with time zone not null,
    changed_at timestamp with time zone not null,

    name text not null,
    title text not null,
    department text not null,
    email_address text not null,

    street text not null,
    postal_code text not null,
    state text not null,
    city text not null,
    country text not null,

    comment text not null
);

create table phone_numbers (
    id integer not null primary key,

    number text not null,
    type text not null,
    person_id int default null,

    foreign key (person_id) references people(id) on update cascade on delete cascade
);

create table users (
    id integer not null primary key,
    version int not null,
    created_at timestamp with time zone not null,
    changed_at timestamp with time zone not null,
    admin boolean not null,

    login text not null unique,
    password_hash text not null
);

create table sessions (
    token text not null primary key,
    "user" text not null references users(login),
    valid_until timestamp with time zone not null
);


-- +migrate Down
drop table phone_numbers;
drop table people;
drop table sessions;
drop table users;
