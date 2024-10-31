from backend.user_management import get_user_db
from backend.core import db


from telegram import Update
from telegram.ext import ContextTypes


from os import getenv as env
from textwrap import dedent


async def start(update: Update, context: ContextTypes) -> None:
    # await update.message.reply_text("Hello World!")
    if not update.effective_user:
        return
    get_user_db(db, update.effective_user)
    sql_username = f"u{update.effective_user.id}"
    with db as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                "SELECT sql_username, sql_password, sql_db_name FROM telegram.user_map WHERE sql_username=%s",
                ( sql_username,),
            )
            result = cursor.fetchone()

            await update.message.reply_text(
                f"Hello <b>{update.effective_user.first_name}</b>, welcome to <b>SQL Sensi</b>\nYour database environment is <b>ready</b> ✅", parse_mode='HTML')

            await update.message.reply_text(dedent(
                f"""                 

Connection code:

```python
import mysql.connector
db = mysql.connector.connect(
    host="{env("SQL_HOST")}",
    user="{result[0]}",
    password="{result[1]}",
    database="{result[2]}"
)
cursor = db.cursor()

# Example query to create a table
cursor.execute("CREATE TABLE test_table(id INT PRIMARY KEY, name VARCHAR(20))")

# Example query to insert data
cursor.execute("INSERT INTO test_table(id, name) VALUES (1, 'John')")
db.commit()

```
You can now connect to your database using the connection code above\\.
"""), parse_mode="MarkdownV2")