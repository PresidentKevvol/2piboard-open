-- for creating all the tables needed in 2piboard
CREATE TABLE board (
	url				varchar(24)		primary key,
	title			varchar(255)
);

CREATE TABLE tripcode (
	id				bigint			primary key,
	username		varchar(24)		not null,
	hash			varchar(70)		not null
);

CREATE TABLE post (
	id				bigint			primary key,
	title			varchar(80)		not null,
	body			varchar(400),
	time_posted		timestamp with time zone	not null,
	author			bigint,
	poster_ip		varchar(150),
	in_board		varchar(24)		not null,
	constraint post_by_poster_PK foreign key (author) references tripcode (id),
	constraint post_in_board foreign key (in_board) references board (url)
);

CREATE TABLE reply (
	id				bigint			primary key,
	body			varchar(400)	not null,
	time_posted		timestamp with time zone	not null,
	author			bigint,
	parent			bigint,
	poster_ip		varchar(150),
	constraint reply_by_poster_PK foreign key (author) references tripcode (id),
	constraint reply_to_post foreign key (parent) references post (id)
);

CREATE TABLE access_check (
	entry			varchar(255)	primary key
);

-- setup the basic attributes
INSERT INTO access_check (entry) VALUES ('true');
INSERT INTO tripcode (id, username, hash) VALUES (-1, 'Anonymous', '');
INSERT INTO board (url, title) VALUES
	('vg', 'Video Games'),
	('r', 'Random'),
	('a', 'Anime & Manga'),
	('m', 'Memes'),
	('f', 'Food'),
	('sc', 'Math, Science, & Coding'),
	('ar', 'Art & Media'),
	('an', 'Cute Animals'),
	('in', 'International');