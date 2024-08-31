CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    email varchar(255) unique not null,
    name varchar(255) not null,
    password_hash varchar(255) not null,
    coins BIGINT
);
CREATE TABLE requests
(
    id SERIAL PRIMARY KEY,
    sender int references users(id) on delete cascade not null,
    recipient int references users(id) on delete cascade not null,
    text varchar(255) not null
);
CREATE TABLE rooms
(
    id SERIAL PRIMARY KEY,
    users_quantity BIGINT not null,
    type varchar(255) not null
);
CREATE TABLE players
(
    id SERIAL PRIMARY KEY,
    user_id int references users(id),
    user_name varchar(255) not null,
    room_id int references rooms(id) on delete cascade not null
);
CREATE TABLE superpowers
(
    id SERIAL PRIMARY KEY,
    user_id int references users(id) on delete cascade not null,
    name varchar(255) not null
);