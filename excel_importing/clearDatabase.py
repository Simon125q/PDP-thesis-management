import sqlite3

PATH_TO_DATABASE = 'Path/To/diploma_database.db'

conn = sqlite3.connect(PATH_TO_DATABASE)
cursor = conn.cursor()

cursor.execute("SELECT name FROM sqlite_master WHERE type='table';")
tables = cursor.fetchall()

for table in tables:
    cursor.execute(f"DELETE FROM {table[0]}")

conn.commit()
conn.close()
