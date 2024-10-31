from .core import DataBase
from bcrypt import hashpw, gensalt
from telegram import User
from os import getenv as env

def init_db(db: DataBase) -> None:
    db.create_database("telegram")
    with db as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                CREATE TABLE IF NOT EXISTS telegram.user_map (
                    id INT PRIMARY KEY,
                    username VARCHAR(32),
                    first_name VARCHAR(32),
                    last_name VARCHAR(32),
                    language_code VARCHAR(8),
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    is_admin BOOLEAN DEFAULT FALSE,
                    is_premium BOOLEAN DEFAULT FALSE,
                    sql_username VARCHAR(32) NOT NULL,
                    sql_db_name VARCHAR(32) NOT NULL,
                    sql_password VARCHAR(32) NOT NULL,
                    gemini_api_key VARCHAR(32)
                )
                """
            )


def create_db_for_user(db: DataBase, user: User, is_admin: bool = False) -> DataBase:
    with db as connection:
        with connection.cursor() as cursor:
            salt = gensalt()
            sql_username = f"u{user.id}"
            sql_db_name = f"u_{user.id}"
            password = hashpw(str(user.id).encode(), salt)
            db.create_database(sql_db_name)
            # Grant all privileges to the user on the new database
            cursor.execute(f"GRANT ALL PRIVILEGES ON {sql_db_name}.* TO '{env("SQL_USER")}'@'%'")

            # Create a new MySQL user first
            cursor.execute(
                f"CREATE USER IF NOT EXISTS '{sql_username}'@'%' IDENTIFIED BY '{password.decode()}'"
            )
            cursor.execute(
                "INSERT IGNORE INTO telegram.user_map(id, username, first_name, last_name, language_code, sql_username, sql_password, sql_db_name, is_admin, is_premium) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)",
                (
                    user.id,
                    user.username,
                    user.first_name,
                    user.last_name,
                    user.language_code,
                    sql_username,
                    password.decode(),
                    sql_db_name,
                    is_admin,
                    user.is_premium,
                    # None,  # Assuming gemini_api_key is None by default
                ),
            )
            db.create_database(sql_db_name)
            cursor.execute(f"GRANT ALL PRIVILEGES ON {sql_db_name}.* TO '{sql_username}'@'%'")
            cursor.execute("FLUSH PRIVILEGES")
        

            return DataBase(user=sql_username, password=password.decode(), db_name=sql_db_name)


def get_user_db(db: DataBase, user: User) -> DataBase:
    with db as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                "SELECT sql_username, sql_password, sql_db_name FROM telegram.user_map WHERE id=%s",
                (user.id,),
            )
            result = cursor.fetchone()
            if result:
                return DataBase(user=user, password=result[1], db_name=result[2])
            else:
                return create_db_for_user(db, user)


def get_user_by_id(db: DataBase, user_id: int) -> User:
    with db as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                "SELECT id, username, first_name, last_name, language_code FROM telegram.user_map WHERE id=%s",
                (user_id,),
            )
            result = cursor.fetchone()
            if result:
                return User(
                    id=result[0],
                    username=result[1],
                    first_name=result[2],
                    last_name=result[3],
                    language_code=result[4],
                )
            else:
                return None
