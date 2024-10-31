from telegram import Update
from telegram.ext import ContextTypes
from backend.core import db
from backend.misc import id_from_User
import google.generativeai as genai
from os import getenv as env
from google.generativeai.types import HarmCategory, HarmBlockThreshold


async def ai(update: Update, context: ContextTypes) -> None:
    if not update.effective_user:
        return
    
    if not user_exists(db, update.effective_user):
        await update.message.reply_text("You need to register first to use this command, use /start to register")
        return
    
    # see if the user has a custom api key for gemini
    with db as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                "SELECT gemini_api_key FROM telegram.user_map WHERE sql_username=%s",
                (id_from_User(update.effective_user),),
            )
            result = cursor.fetchone()
            if result[0]:
                genai.configure(api_key=result[0])
            else:
                genai.configure(api_key=env("GENAI_API_KEY"))
            
    # get the message text
    message_text = update.message.text
    # remove the command from the message text
    message_text = message_text.replace("/ai", "")

    model = genai.GenerativeModel(
            model_name="gemini-1.5-pro",
            generation_config={
            "temperature": 1,
            "top_p": 0.95,
            "top_k": 64,
            "max_output_tokens": 8192,
            "response_mime_type": "text/plain",
        },
            safety_settings={
                HarmCategory.HARM_CATEGORY_HATE_SPEECH: HarmBlockThreshold.BLOCK_NONE,
                HarmCategory.HARM_CATEGORY_HARASSMENT: HarmBlockThreshold.BLOCK_NONE,
                HarmCategory.HARM_CATEGORY_SEXUALLY_EXPLICIT: HarmBlockThreshold.BLOCK_NONE,
                HarmCategory.HARM_CATEGORY_DANGEROUS_CONTENT: HarmBlockThreshold.BLOCK_NONE,
            },
            # system_instruction=,
        )
    response = model.generate_content(message_text)
    await update.message.reply_text(response.text)

