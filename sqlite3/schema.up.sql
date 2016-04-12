CREATE TABLE IF NOT EXISTS users(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT NOT NULL UNIQUE,
	name TEXT NOT NULL,
	is_admin BOOLEAN DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS subjects(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	title TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS threads(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	subject_id INTEGER NOT NULL,
	creator_user_id INTEGER NOT NULL,
	time_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	is_pinned BOOLEAN DEFAULT 0 NOT NULL,
	is_visible BOOLEAN DEFAULT 1 NOT NULL,
	UNIQUE(id, subject_id),
	FOREIGN KEY(subject_id) REFERENCES subjects(id) ON DELETE CASCADE,
	FOREIGN KEY(creator_user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_threads_subject_id on threads(subject_id);

CREATE TABLE IF NOT EXISTS thread_votes(
	thread_id INTEGER NOT NULL,
	user_id TEXT  NOT NULL,
	PRIMARY KEY(thread_id, user_id),
	FOREIGN KEY(thread_id) REFERENCES threads(id) ON DELETE CASCADE,
	FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_thread_votes_thread_id thread_votes(thread_id);
CREATE INDEX idx_thread_votes_user_id thread_votes(user_id);

CREATE TABLE IF NOT EXISTS tags(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	subject_id INTEGER NOT NULL,
	UNIQUE(name, subject_id),
	UNIQUE(id, subject_id),
	FOREIGN KEY(subject_id) REFERENCES subjects(id) ON DELETE CASCADE
);

CREATE INDEX idx_tags_subject_id on tags(subject_id);

CREATE TABLE IF NOT EXISTS thread_tags(
	thread_id INTEGER NOT NULL,
	tag_id INTEGER NOT NULL,
	subject_id INTEGER NOT NULL,
	PRIMARY KEY (thread_id, tag_id),
	FOREIGN KEY(thread_id, subject_id) REFERENCES threads(id, subject_id) ON DELETE CASCADE,
	FOREIGN KEY(tag_id, subject_id) REFERENCES tags(id, subject_id) ON DELETE CASCADE
);

CREATE INDEX idx_thread_tags_thread_id tags(thread_id);
CREATE INDEX idx_thread_tags_tag_id tags(tag_id);
