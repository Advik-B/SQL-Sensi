from telegram import Update
from telegram.ext import ContextTypes
from backend.core import db
from backend.user_management import user_exists
from backend.misc import id_from_User

async def credentials(update: Update, context: ContextTypes) -> None:
    if not update.effective_user:
        return
    if not user_exists(db, update.effective_user):
        await update.message.reply_text("You need to register first to use this command, use /start to register")
        return

    with db as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                "SELECT sql_username, sql_password, sql_db_name FROM telegram.user_map WHERE sql_username=%s",
                (id_from_User(update.effective_user),),
            )
            result = cursor.fetchone()
            await update.message.reply_text(
                    f"Username: `{result[0]}`\nPassword: `{result[1]}`\nDatabase: `{result[2]}`\nHost: `{connection.server_host}",
                    parse_mode="MarkdownV2"
                )
            