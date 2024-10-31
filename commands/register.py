from .start import start
from .credentials import credentials
from telegram.ext import Application, CommandHandler

def register(app: Application):
    app.add_handler(CommandHandler("start", start))
    app.add_handler(CommandHandler("credentials", credentials))

