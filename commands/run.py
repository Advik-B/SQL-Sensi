from telegram import Update
from telegram.ext import ContextTypes
from backend.core import db
from backend.user_management import user_exists, get_user_db
import mysql.connector

MAX_MESSAGE_LENGTH = 4000

def wrap_code(text: str) -> str:
    return f"```\n{text}\n```"


async def run(update: Update, context: ContextTypes) -> None:
    """
    /run

    ```sql
    SELECT * FROM table_name
    ```

    -- or --

    /run

    ```
    SELECT * FROM table_name
    ```

    -- or --

    /run SELECT * FROM table_name

    -- or --

    /run
    SELECT * FROM table_name
    """
    if not update.effective_user:
        return
    
    if not user_exists(db, update.effective_user):
        await update.message.reply_text("You need to register first to use this command, use /start to register")
        return
    
    user_db = get_user_db(db, update.effective_user)

    query = " ".join(context.args)
    if not query:
        await update.message.reply_text("Please provide a query to run 👉👈")
        return

    # Parse the query
    if query.startswith("```sql") and query.endswith("```"):
        query = query[5:-3]
    elif query.startswith("```") and query.endswith("```"):
        query = query[3:-3]
    elif query.startswith("`") and query.endswith("`"):
        query = query[1:-1]
    elif query.startswith("```") or query.startswith("`"):
        query = query[3:] if query.startswith("```") else query[1:]
    elif query.endswith("```") or query.endswith("`"):
        query = query[:-3] if query.endswith("```") else query[:-1]
    elif "`" not in query:
        pass
    else:
        await update.message.reply_text("Invalid query format, please use the correct format 🥺")
        return

    try:
        cursor = user_db.db.cursor()
        cursor.execute(query)
        rows = cursor.fetchall()
        message_text = "\n".join([str(row) for row in rows])
        final_text = wrap_code(message_text)
        if len(final_text) > MAX_MESSAGE_LENGTH:
            await update.message.reply_text(f"Query executed successfully, but the output is too long to display in a single message")
            # Split the message into multiple messages
            for i in range(0, len(message_text), MAX_MESSAGE_LENGTH):
                await update.message.reply_text(wrap_code(message_text[i:i + MAX_MESSAGE_LENGTH]), parse_mode="MarkdownV2")
            return
        await update.message.reply_text(final_text, parse_mode="MarkdownV2")
    except mysql.connector.errors.ProgrammingError as e:
        await update.message.reply_text(f"Error 💀: ```\n{e}\n```", parse_mode="MarkdownV2")

