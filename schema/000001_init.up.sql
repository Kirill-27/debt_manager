CREATE TABLE users
(
    id serial not null unique,
    email     text not null,
    password    text not null,
    full_name  text not null,
    subscription_type int not null,
    phone text not null,
    rating float not null
);

CREATE TABLE debts
(
    id serial not null unique,
    debtor_id  int not null,
    lender_id  int not null,
    full_name text not null,
    status  int not null,
    amount  int not null,
    description text not null,
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null
);

CREATE TABLE friends
(
    friend1_id  int not null,
    friend2_id  int not null
);

CREATE TABLE current_debts
(
    id serial not null unique,
    debtor_id  int not null,
    lender_id  int not null,
    amount  int not null
);

CREATE TABLE reviews
(
    id serial not null unique,
    debtor_id  int not null,
    lender_id  int not null,
    comment text not null,
    rate  int not null,
    created_at timestamp without time zone not null
);