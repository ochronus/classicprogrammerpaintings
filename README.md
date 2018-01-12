# A silly simple Slack command handler which returns images from http://classicprogrammerpaintings.com/ which are relevant to the searched term
                                                                                           
## how to use
                                                                                          
1. Register an account and an index named 'classicprogrammerpaintings' on https://www.algolia.com/. Save your app id and api key to the .env file
                                                                                          
2. [Register a Slack command](https://get.slack.help/hc/en-us/articles/201259356-Slash-Commands) and save the verification token (NOT the Oauth token!) to the .env file
                                                                                           
3. Crawl http://classicprogrammerpaintings.com/ with the crawler in crawl-n-index, e.g. go run crawl-n-index/main.go - this will fetch the descriptions and image urls and save them in the 'classicprogrammerpaintings' index in Algolia
                                                                                           
4. Run the webserver in the root folder. Configure the handler url in the Slack command in step 2. - you can either fire up a proper reverse proxy such as nginx in front of the go server or access it directly.
                                                                                           
## testing in dev mode
                                                                                           
Use [Localtunnel](https://localtunnel.github.io/www/) to pipe traffic from your dev machine

## protips

* if you call the slack command without parameters it returns a random image from the whole index
