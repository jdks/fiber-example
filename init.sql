\c users

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
	user_id uuid DEFAULT uuid_generate_v4 (),
	first_name VARCHAR NOT NULL,
	last_name VARCHAR NOT NULL,
	PRIMARY KEY (user_id)
);

CREATE TABLE user_events (
	event_id uuid NOT NULL,
	user_id uuid NOT NULL,
	created_at timestamp NOT NULL,
	payload jsonb NOT NULL,
  processed boolean DEFAULT false,

	associated_user_ids uuid[],
	PRIMARY KEY (event_id),
	CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

CREATE INDEX payload_idx ON user_events USING GIN (payload) WITH (fastupdate = off);
CREATE INDEX associated_users_idx ON user_events USING GIN (associated_user_ids) WITH (fastupdate = off);
CREATE INDEX unprocessed_events ON user_events USING BTREE (created_at) WHERE processed = false;
