CREATE TABLE users
(
    id serial not null unique,
    email     text not null,
    password    text not null,
    full_name  text not null,
    subscription_type int not null,
    photo text not null,
    rating float not null
);

CREATE TABLE debts
(
    id serial not null unique,
    debtor_id  int not null,
    lender_id  int not null,
    status  int default 1,
    amount  int not null,
    description text not null,
    created_at timestamp without time zone default (now() at time zone 'utc'),
    updated_at timestamp without time zone default (now() at time zone 'utc')
);

CREATE TABLE friends
(
    my_id  int not null,
    friend_id  int not null
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
    reviewer_id  int not null,
    lender_id  int not null,
    comment text not null,
    rate  int not null,
    updated_at timestamp without time zone default (now() at time zone 'utc')
);

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = (now() at time zone 'utc');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER
    update_debts_status
    AFTER update
    on debts
    for EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();