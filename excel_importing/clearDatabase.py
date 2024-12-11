import sqlite3

conn = sqlite3.connect('*/diploma_database.db')
cursor = conn.cursor()

cursor.execute("SELECT name FROM sqlite_master WHERE type='table';")
tables = cursor.fetchall()

for table in tables:
    cursor.execute(f"DELETE FROM {table[0]}")

conn.commit()
conn.close()
