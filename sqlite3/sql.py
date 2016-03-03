import sqlite3

if __name__ == '__main__':
	conn = sqlite3.connect('test.db')
	c = conn.cursor()


	with open('schema.up.sql') as f:
		for query in f.read().split(';'):
			c.execute(query)

	conn.commit()
	conn.close()
