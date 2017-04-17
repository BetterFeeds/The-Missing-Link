# The Missing Link

This project aims to create feeds for online services that don't do so them selves. TML scrapes the web page(s) and returns a Atom feed that can be used in most feed/rss readers.

# Sources

| Source                 | Source file        | URL                                    |
|------------------------|--------------------|----------------------------------------|
| Open Rights Group blog | sources/org.go     | https://tml.betterfeeds.org/org/1.atom |
| IMDb Trailers          | sources/imdb.go    | https://tml.betterfeeds.org/imdb.atom  |
| iPlayer                | sources/iplayer.go |                                        |

# Set up

Install with your package manager or go.

```
go get github.com/BetterFeeds/The-Missing-Link
go install github.com/BetterFeeds/The-Missing-Link
```

You can start TML at boot with a systemd service:

```
cat /etc/systemd/system/tml.service
[Unit]
Description=Local TML Server

[Service]
User=tml
Type=simple
ExecStart=/home/tml/go/bin/The-Missing-Link

[Install]
WantedBy=multi-user.target
```

The Missing Link does not, and will not do any caching. So when set up the service should be behind [Varnish](https://www.varnish-cache.org/) so similar, by default TML listens on 127.0.0.1:8080.
