from backend.core import DataBase
from backend.user_management import get_user_db
from telegram import BotCommand, Update
from telegram.ext import Application, CommandHandler, ContextTypes
from os import getenv as env

db = DataBase()
db.create_database("test")

with db as connection:
    with connection.cursor() as cursor:
        cursor.execute("SHOW DATABASES")
        print(cursor.fetchall())


token = env("TOKEN")
app = Application.builder().token(token).build()

async def start(update: Update, context: ContextTypes) -> None:
    # await update.message.reply_text("Hello World!")
    if update.effective_user:
        db = get_user_db(db, update.effective_user)
        with db as connection:
            with connection.cursor() as cursor:
                cursor.execute("SHOW DATABASES")
                await update.message.reply_text(str(cursor.fetchall()))


# app.bot.set_my_commands(
#     [
#         BotCommand("start", "Start the bot"),
#     ]
# )

app.add_handler(CommandHandler("start", start))

app.run_polling()