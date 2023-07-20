DROP TRIGGER IF EXISTS update_debt_time ON debts;
DROP TRIGGER IF EXISTS update_review_time ON reviews;
DROP TRIGGER IF EXISTS update_stripe_payments_time ON reviews;
DROP FUNCTION IF EXISTS update_debts_updated_at();
DROP TABLE users;
DROP TABLE debts;
DROP TABLE friends;
DROP TABLE current_debts;
DROP TABLE reviews;
DROP TABLE stripe_payments