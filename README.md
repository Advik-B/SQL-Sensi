# SQL-Sensi
[![Static Badge](https://img.shields.io/badge/Docker_Image-%20?logo=Docker&logoSize=auto&color=white)](https://github.com/Advik-B/SQL-Sensi/pkgs/container/sql-sensi)
[![Static Badge](https://img.shields.io/badge/written%20in-go-%20?logoColor=white&logoSize=auto&color=%2300ADD8)](https://go.dev/)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/Advik-B/SQL-Sensi/total?label=total%20downloads&color=violet&link=https%3A%2F%2Fgithub.com%2FAdvik-B%2FSQL-Sensi%2Freleases)
![GitHub commits since latest release](https://img.shields.io/github/commits-since/Advik-B/SQL-Sensi/latest?include_prereleases&)
[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/Advik-B/SQL-Sensi/publish.yml?branch=main&)](https://github.com/Advik-B/SQL-Sensi/actions/workflows/publish.yml)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/Advik-B/SQL-Sensi)


SQL-Sensi is a Telegram bot designed to help users learn and practice MySQL. It provides an interactive platform where users can connect to a database, execute SQL commands, and share tables with other users. The bot also offers sample tables and data to help users get started quickly.

## Features
Here's what SQL-Sensi currently offers and what's coming soon:

### Currently Available ✅
- **Direct Database Connection**: Connect directly to your MySQL database with secure credentials
- **SQL Command Execution**: Run SQL queries directly through Telegram with support for all major SQL operations
- **Sample Data**: Pre-built tables and sample data for practice and learning
- **AI Assistance**: Get help from AI to solve MySQL problems and generate queries
- **Command Reference**: Comprehensive help system listing all available commands and their usage
- **Interactive Learning**: Practice SQL with immediate feedback and results

### Coming Soon 🚧
- **Enhanced Table Sharing**: Share your tables and data structures with other users for collaborative learning
- **Query Templates**: Pre-built query templates for common database operations
- **Query History**: Track and revisit your previously executed queries
- **Performance Analytics**: Get insights into your query performance and optimization suggestions
- **Custom Datasets**: Import your own datasets for practice
- **Interactive Tutorials**: Step-by-step SQL learning modules with hands-on exercises

## Deploy

SQL-Sensi is also available as a [![Docker Image](https://img.shields.io/badge/Docker_Image-%20?logo=Docker&logoSize=auto&color=white)](https://github.com/Advik-B/SQL-Sensi/pkgs/container/sql-sensi)

## Deploy with Docker Compose

To deploy SQL-Sensi using Docker Compose, follow these steps:

1. Download [docker-compose.yml](https://raw.githubusercontent.com/Advik-B/SQL-Sensi/refs/heads/main/docker-compose.yml) to a folder and `cd` into it

2. Create a `.env` file in the same folder and add the required environment variables:

    ```env
    DB_HOST=your_database_host
    DB_USER=your_database_user
    DB_PASS=your_database_password
    TELEGRAM_TOKEN=your_telegram_api_key
    GEMINI_API_KEY=your_gemini_api_key
    ```

3. Run the following command to start the services:

    ```sh
    docker-compose up -d
    ```

## Build and run

### Direct download

> [!NOTE]
> You can directly download the linux binary from [releases](https://github.com/Advik-B/SQL-Sensi/releases/latest)
>
> 🐧 Linux Binary: [sql.sensi](https://github.com/Advik-B/SQL-Sensi/releases/latest/download/sql.sensi)
>
> 🗜️ Linux Binary (7-Zipped): [sql.sensi.7z](https://github.com/Advik-B/SQL-Sensi/releases/latest/download/sql.sensi.7z)

### Prerequisites
- Go compiler installed

### Download the source code
```
git clone https://github.com/Advik-B/SQL-Sensi.git
cd sql-sensi
```

### Build the bot
```
go build -gcflags="all=-l -B" -ldflags="-s -w -extldflags '-static'" -trimpath -o sql.sensi
```

### Setup Enviroment variables in a `.env` file or directly thru the environment

| Variable Name | Description | Importance |
|---------------|-------------|----------|
| `DB_HOST` | The mysql database host | REQUIRED |
| `DB_USER` | The mysql database user to connect with | REQUIRED |
| `DB_PASS` | The mysql database password | REQUIRED |
| `TELEGRAM_API_KEY` | The telegram bot token | REQUIRED |
| `GEMINI_API_KEY` | The Gemini API key for ai features | OPTIONAL |

### And run the bot

```
./sql.sensi
```
