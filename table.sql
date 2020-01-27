CREATE SEQUENCE tl_user_sequence
  start 100000000
  increment 1;

CREATE TABLE tl_user (
  id SERIAL PRIMARY KEY
);