from .start import start
from telegram.ext import Application, CommandHandler

def register(app: Application):
    app.add_handler(CommandHandler("start", start))

