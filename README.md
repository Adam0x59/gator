# Gator

A simple cli tool for aggregating rss feeds.

## Features

- Handle multiple users with a user linked follow list.
- Register new rss-feeds.
- Aggregate RSS feeds and store in the database.
- Browse latest feeds for a user.

## Dependencies

- Go (1.24 or later)- Docs: https://go.dev/doc/install
- Postgres (16.9 or later) - Docs: https://www.postgresql.org/docs/current/tutorial-install.html

## Installation

1. Clone the repo, run this command where you want to clone the repo to: 

```
git clone https://github.com/adam0x59/gator.git
cd gator
```

2. Create a config file: Run the following command from the root of the cloned repo, this will create a config file in your home directory and open it for editing:

```
nano ~/.gatorconfig.json
```

3. Paste the following into the config file you just created and have open:

```
{
    "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=>,    
    "current_user_name": "kahya"
}
```

4. Change the db_url to the URL of your database replacing the username and password you want to use to access the database, also change the username to your preffered user name. Then press Ctrl+S to save and Ctrl+X to exit.

5. You should still be in the root of the repo on your system, so next run the following commands to run the database migrations, before running the, ensure you have created a DB named gator in postgres first:

```
cd sql/schema/
goose postgres <Your Database URL> up
cd ../..
```

6. Install gator by running the following command from the root of the repo on your machine:

```
go install .
```

Congratulations Gator Is Now Installed!

## Commands

To use gator the command structure is as follows:

```
gator <command> <arg1> <arg2>...
```

### Available Commands

- login       - Log in as another user. Takes a single argument: user name.
- register    - Register a user. Takes a single argument: user name.
- reset       - Resets the database. Takes no argument, deletes all users and feeds, database structure remains intact.
- users       - Lists all users. Takes no argument.
- agg         - Starts the rss feed aggregator loop. Takes one argument: Time between fetch requests "1s, 1m, 1h".
- addfeed     - Adds an rss feed (And adds the feed to current user's follow list). Takes two arguments: Argument1 - Feed name, Argument2 - Feed URL.
- feeds       - Lists rss feeds stored in the database. Takes no argument
- follow      - Adds a given feed to a users follow list. Takes one argument: Feed url (Must match a url in the feeds list.)
- following   - Lists feeds being followed by the logged in user. Takes no argument.
- unfollow    - Removes a feed from the logged in users follow list. Takes one argument: Feed url.
- browse      - Lists the most recent posts, number listed is the argument. Takes one argument: number of posts to list, (default is 2)


