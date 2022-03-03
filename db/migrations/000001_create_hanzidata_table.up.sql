SELECT pg_try_advisory_lock(migration);

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

SELECT pg_advisory_unlock_all();