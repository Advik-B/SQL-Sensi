from mysql.connector import MySQLConnection

def create_db(conn: MySQLConnection, db_name: str):
    cursor = conn.cursor()
    cursor.execute(f"CREATE DATABASE IF NOT EXISTS {db_name}")
    cursor.close()
    conn.commit()


def create_table(conn: MySQLConnection, table_name: str, columns: dict):
    cursor = conn.cursor()
    query = f"CREATE TABLE IF NOT EXISTS {table_name} ("
    for col_name, col_type in columns.items():
        query += f"{col_name} {col_type}, "
    query = query[:-2] + ")"
    cursor.execute(query)
    cursor.close()
    conn.commit()


def drop_table(conn: MySQLConnection, table_name: str):
    cursor = conn.cursor()
    cursor.execute(f"DROP TABLE {table_name}")
    cursor.close()
    conn.commit()


def drop_db(conn: MySQLConnection, db_name: str):
    cursor = conn.cursor()
    cursor.execute(f"DROP DATABASE {db_name}")
    cursor.close()
    conn.commit()

