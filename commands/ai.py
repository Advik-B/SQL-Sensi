from telegram import Update
from telegram.ext import ContextTypes
from backend.core import db
from backend.user_management import user_exists
from backend.misc import id_from_User
import google.generativeai as genai
from os import getenv as env
from google.generativeai.types import HarmCategory, HarmBlockThreshold

system_prompt = """
You are SQL Sensi, a SQL database learning assistant. You can be asked anything about SQL databases and you will provide the best possible answer.
You can also provide examples and code snippets to help the user understand better.
You write queries for the user and help them understand the results.
You can also provide tips and tricks to help the user write better queries.

Any off topic conversation will be ignored and you should not provide any personal information to the user.
Any off topic conversation will be ignored and you should ask the user to focus the conversation on SQL databases.
"""

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

            if result[0] is not None:
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
            system_instruction=system_prompt,
        )
    response = model.generate_content(message_text)
    await update.message.reply_text(response.text)

