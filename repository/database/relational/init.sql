CREATE SCHEMA IF NOT EXISTS stream_anime;

CREATE TABLE IF NOT EXISTS stream_anime.anime (
	id int NOT NULL,
	title varchar NOT NULL,
	"type" varchar NOT NULL,
	summary varchar NOT NULL,
	genre varchar NOT NULL,
	airing_year varchar NOT NULL,
	status varchar NOT NULL,
	image_url varchar NOT NULL,
	latest_episode varchar NOT NULL,
    updated_at timestamp with time zone NOT NULL,
	CONSTRAINT anime_pk PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stream_anime.episode (
	id serial NOT NULL,
	anime_id int NOT NULL,
	text varchar NOT NULL,
	endpoint varchar NOT NULL,
	CONSTRAINT episode_pk PRIMARY KEY (id),
	CONSTRAINT episode_anime_anime_id_fk FOREIGN KEY (anime_id) REFERENCES stream_anime.anime(id)
);

CREATE TABLE IF NOT EXISTS stream_anime.user (
	id serial NOT NULL,
	user_token varchar NOT NULL,
	CONSTRAINT user_pk PRIMARY KEY (id),
	CONSTRAINT user_un UNIQUE (user_token)
);

CREATE TABLE IF NOT EXISTS stream_anime.user_anime_xref (
	user_id int NOT NULL,
	anime_id int NOT NULL,
	bookmarked_latest_episode varchar NOT NULL,
	CONSTRAINT user_anime_xref_pk PRIMARY KEY (user_id,anime_id),
	CONSTRAINT user_anime_xref_user_fk FOREIGN KEY (user_id) REFERENCES stream_anime."user"(id),
	CONSTRAINT user_anime_xref_anime_fk FOREIGN KEY (anime_id) REFERENCES stream_anime.anime(id)
);