CREATE TABLE users
(
    id serial not null unique,
    email     text not null unique,
    password    text not null,
    full_name  text not null,
    subscription_type int default 1,
    photo text not null,
    rating float default 5,
    marks_sum int default 0,
    marks_number int default 0
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

CREATE OR REPLACE FUNCTION update_debts_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now() at time zone 'utc';
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_debt_time
    BEFORE UPDATE ON debts
    FOR EACH ROW
EXECUTE FUNCTION update_debts_updated_at();

CREATE TRIGGER update_review_time
    BEFORE UPDATE ON reviews
    FOR EACH ROW
EXECUTE FUNCTION update_debts_updated_at();
