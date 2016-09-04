CREATE TYPE status AS ENUM ('Pending', 'Sent', 'Rejected');

CREATE TABLE notification (
	id serial PRIMARY KEY,
	user_email varchar(256) NOT NULL,
	location_name varchar(100) NOT NULL,
	status status default 'Pending',
	wave_height_m real NOT NULL,
	forecast_time timestamp NOT NULL,
	created timestamp NOT NULL default current_timestamp,
	scheduled timestamp NOT NULL default current_timestamp,
	sent timestamp
);