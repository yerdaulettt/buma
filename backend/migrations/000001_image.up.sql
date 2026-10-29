create table if not exists users (
    id serial primary key,
    name varchar(20) not null,
    email varchar(40) unique not null,
    password text not null
);

create table if not exists images (
    id serial primary key,
    image_type varchar(20) not null,
    s3_key varchar(50) unique not null,
    added_at timestamptz default now(),
    user_id integer references users(id) on delete cascade
);