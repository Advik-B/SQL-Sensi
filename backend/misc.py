from telegram import User

def id_from_User(user: User) -> str:
    return f"u{user.id}"

def dbname_from_User(user: User) -> str:
    return f"u_{user.id}"