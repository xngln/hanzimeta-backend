BEGIN;

SELECT pg_try_advisory_lock(1);

CREATE TABLE IF NOT EXISTS characters (
	char_id serial PRIMARY KEY,
	simplified VARCHAR(50) NOT NULL,
    pinyin VARCHAR(100) NOT NULL,
	traditional VARCHAR(50) NOT NULL,
	japanese VARCHAR(50) NOT NULL,
	junda_freq SMALLINT,
	gs_num SMALLINT,
	hsk_lvl SMALLINT
);

CREATE INDEX junda_freq_idx ON characters (junda_freq);
CREATE INDEX gs_num_idx ON characters (gs_num);
CREATE INDEX hsk_lvl_idx ON characters (hsk_lvl);

SELECT pg_advisory_unlock_all();

COMMIT;