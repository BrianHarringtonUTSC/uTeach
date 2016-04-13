CREATE TABLE IF NOT EXISTS users(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT NOT NULL UNIQUE,
	name TEXT NOT NULL,
	is_admin BOOLEAN DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS topics(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	title TEXT NOT NULL,
	description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS posts(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	topic_id INTEGER NOT NULL,
	creator_user_id INTEGER NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	is_pinned BOOLEAN DEFAULT 0 NOT NULL,
	is_visible BOOLEAN DEFAULT 1 NOT NULL,
	UNIQUE(id, topic_id),
	FOREIGN KEY(topic_id) REFERENCES topics(id) ON DELETE CASCADE,
	FOREIGN KEY(creator_user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_posts_topic_id ON posts(topic_id);

CREATE TABLE IF NOT EXISTS post_votes(
	post_id INTEGER NOT NULL,
	user_id TEXT  NOT NULL,
	PRIMARY KEY(post_id, user_id),
	FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
	FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_post_votes_post_id ON post_votes(post_id);
CREATE INDEX IF NOT EXISTS idx_post_votes_user_id ON post_votes(user_id);

CREATE TABLE IF NOT EXISTS tags(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	topic_id INTEGER NOT NULL,
	UNIQUE(name, topic_id),
	UNIQUE(id, topic_id),
	FOREIGN KEY(topic_id) REFERENCES topics(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_tags_topic_id ON tags(topic_id);

CREATE TABLE IF NOT EXISTS post_tags(
	post_id INTEGER NOT NULL,
	tag_id INTEGER NOT NULL,
	topic_id INTEGER NOT NULL,
	PRIMARY KEY (post_id, tag_id),
	FOREIGN KEY(post_id, topic_id) REFERENCES posts(id, topic_id) ON DELETE CASCADE,
	FOREIGN KEY(tag_id, topic_id) REFERENCES tags(id, topic_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_post_tags_post_id ON post_tags(post_id);
CREATE INDEX IF NOT EXISTS idx_post_tags_tag_id ON post_tags(tag_id);
