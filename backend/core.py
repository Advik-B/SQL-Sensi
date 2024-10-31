from mysql.connector import connect, MySQLConnection
from mysql.connector.cursor import MySQLCursor
from mysql.connector.errors import ProgrammingError
from os import getenv as env

class DataBase:
    def __init__(self, host: str = None, user: str = None, password: str = None, db_name:str = None) -> None:
        self.db = connect(
            host=host or env("SQL_HOST"),
            user=user or env("SQL_USER"),
            password=password or env("SQL_PASSWORD"),
            database=db_name
        )
        self.last_database: str = None

    def __enter__(self) -> MySQLConnection:
        return self.db
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        self.db.commit()
    

    def create_database(self, name: str) -> None:
        with self.db.cursor() as cursor:
            cursor.execute(f"CREATE DATABASE IF NOT EXISTS {name}")
        self.last_database = name
    
    def drop_database(self, name: str) -> None:
        with self.db.cursor() as cursor:
            cursor.execute(f"DROP DATABASE {name}")
        self.last_database = None if self.last_database == name else self.last_database


db: DataBase = DataBase() # This is the database object that will be used throughout the application and will be shared across modules