from backend.create_db import create_db, create_table
from mysql.connector import connect


from os import getenv as env
conn = connect(
    host=env("SQL_HOST"),
    user=env("SQL_USER"),
    password=env("SQL_PASSWORD")
)

create_db(conn, "test_db")
conn.cursor().execute("USE test_db;")
conn.commit()
create_table(conn, "test_table", {"name": "VARCHAR(255)", "age": "INT"})


# Run the desc command.
cursor = conn.cursor()
cursor.execute("DEx SC test_table;")
print(cursor.fetchall())