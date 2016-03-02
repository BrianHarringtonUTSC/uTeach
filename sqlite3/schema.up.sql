CREATE TABLE IF NOT EXISTS users(
	username TEXT PRIMARY KEY
)

CREATE TABLE IF NOT EXISTS subjects(
	name TEXT PRIMARY KEY,
	title TEXT NOT NULL
)

CREATE TABLE IF NOT EXISTS threads(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	subject_name TEXT NOT NULL,
	created_by_username TEXT NOT NULL,
	FOREIGN KEY(subject_name) REFERENCES subjects(name) ON DELETE CASCADE,
	FOREIGN KEY(created_by_username) REFERENCES users(username) ON DELETE CASCADE
)


CREATE TABLE IF NOT EXISTS upvotes(
	username TEXT  NOT NULL,
	thread_id INTEGER NOT NULL,
	PRIMARY KEY (username, thread_id),
	FOREIGN KEY(username) REFERENCES users(username) ON DELETE CASCADE,
	FOREIGN KEY(thread_id) REFERENCES threads(id) ON DELETE CASCADE
)
