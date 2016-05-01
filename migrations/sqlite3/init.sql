-- +migrate Up
create table people (
    id integer not null primary key autoincrement,
    version int not null,
    created_at datetime not null,
    changed_at datetime not null,

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
    id integer not null primary key autoincrement,
    version int not null,
    created_at datetime not null,
    changed_at datetime not null,

    number text not null,
    type text not null,
    person_id int default null,

    foreign key (person_id) references people(id) on update cascade on delete cascade
);

-- +migrate Down
drop table people;
drop table phone_numbers;
