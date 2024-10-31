from telegram.ext import Application
from backend.user_management import init_db
from backend.core import db
from commands import register as register_commands
from os import getenv as env
from dotenv import load_dotenv

def main():
    load_dotenv()
    token = env("TOKEN")
    init_db(db)
    app = Application.builder().token(token).build()
    register_commands(app)
    app.run_polling()


if __name__ == "__main__":
    main()

