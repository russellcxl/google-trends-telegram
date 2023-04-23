# Google Trends API + Telegram bot + Redis

## Features
- Telegram bot that returns a list of daily trending topics from Google Trends when a user calls the command `/getdaily`
- Users can click on any number corresponding to a topic to get a list of related top articles (see screenshots below)
- Access control -- `t.me/RussellGTrendBot?start=123`. Tokens are stored in .env as MD5 hashes
- Trending topics are stored in Redis temporarily for reduced querying


## Usage
1. Build the images
```
docker compose build
```
2. Run the containers
```
docker compose up
```

## Screenshots
![get_daily_trending_topics](/images/get_daily.png)
![get_daily_trending_topic](/images/get_topic.png)
