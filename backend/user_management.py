from .core import DataBase
from bcrypt import hashpw, gensalt
from telegram import User
from os import getenv as env
from .misc import id_from_User, dbname_from_User

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
            sql_username = id_from_User(user)
            sql_db_name = dbname_from_User(user)
            password = hashpw(str(user.id).encode(), salt).decode()[:20] # Only the first 10 characters
            db.create_database(sql_db_name)
            # Grant all privileges to the user on the new database
            cursor.execute(f"GRANT ALL PRIVILEGES ON {sql_db_name}.* TO '{env("SQL_USER")}'@'%'")

            # Create a new MySQL user first
            cursor.execute(
                f"CREATE USER IF NOT EXISTS '{sql_username}'@'%' IDENTIFIED BY '{password}'"
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
                    password,
                    sql_db_name,
                    is_admin,
                    user.is_premium,
                    # None,  # Assuming gemini_api_key is None by default
                ),
            )
            # print(f"Created database {sql_db_name} for user {user.id} with password {password}")
            db.create_database(sql_db_name)
            cursor.execute(f"GRANT ALL PRIVILEGES ON {sql_db_name}.* TO '{sql_username}'@'%'")
            cursor.execute("FLUSH PRIVILEGES")
        

            return DataBase(user=sql_username, password=password, db_name=sql_db_name)

def get_user_db(db: DataBase, user: User) -> DataBase:
    with db as connection:
        with connection.cursor() as cursor:
            sql_username = id_from_User(user)
            cursor.execute(
                "SELECT sql_username, sql_password, sql_db_name FROM telegram.user_map WHERE sql_username=%s",
                ( sql_username,),
            )
            result = cursor.fetchone()
            # print(result)
            if result:
                return DataBase(user=sql_username, password=result[1], db_name=result[2])
            else:
                return create_db_for_user(db, user)

def user_exists(db: DataBase, user: User) -> bool:
    with db as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                "SELECT id FROM telegram.user_map WHERE sql_username=%s",
                (id_from_User(user),),
            )
            result = cursor.fetchone()
            print(result)
            return bool(result)
