CREATE TABLE tbl_todo (
	id          SERIAL PRIMARY KEY,
    name        VARCHAR(256) NOT NULL,
    status      VARCHAR(12) NOT NULL,
    create_ts   TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC'),
    update_ts   TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC')
);