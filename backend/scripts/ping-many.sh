#!/bin/bash


HOSTS="actiontarget.com,google.com,apple.com,facebook.com,
amazon.com,microsoft.com,cloudflare.com,openai.com,
github.com,linkedin.com,youtube.com,twitter.com,
stackoverflow.com,angular.io,netflix.com,zoom.us,
adobe.com,oracle.com,intel.com,ibm.com,
tesla.com,airbnb.com,uber.com,lyft.com,
pinterest.com,snapchat.com,quora.com,medium.com,
wordpress.org,blogger.com,tumblr.com,wikipedia.org,
reddit.com,instagram.com,whatsapp.com,tiktok.com,
spotify.com,discord.com,slack.com,trello.com,
atlassian.com,gitlab.com,bitbucket.org,docker.com,
kubernetes.io,ubuntu.com,debian.org,centos.org,
redhat.com,digitalocean.com,linode.com,vultr.com,
heroku.com,aws.amazon.com,azure.microsoft.com,cloud.google.com,
alibabacloud.com,godaddy.com,namecheap.com,wordpress.com,
shopify.com,walmart.com,target.com,bestbuy.com,
ebay.com,etsy.com,paypal.com,stripe.com,
visa.com,mastercard.com,americanexpress.com,wellsfargo.com,
bankofamerica.com,chase.com,citi.com,capitalone.com,
bloomberg.com,cnn.com,bbc.com,reuters.com,
theguardian.com,nytimes.com,washingtonpost.com,wsj.com,
mozilla.org,gnu.org,apache.org,gnuplot.info,
nginx.com,postgresql.org,mariadb.org,mysql.com,
sqlite.org,mongodb.com,redis.io,elastic.co,
cloudflarestatus.com,weather.com,timeanddate.com,accuweather.com,
duckduckgo.com,brave.com,opera.com,yandex.com,
baidu.com,ask.com,archive.org,sourceforge.net,
codepen.io,jsfiddle.net,stackexchange.com,producthunt.com,
npmjs.com,pypi.org,ruby-lang.org,python.org,
golang.org,dev.to,freecodecamp.org,geeksforgeeks.org,sentinelonecareers.com"

# Port to use (80 = HTTP)
PORT=80

# Ping interval (every 5 seconds)
INTERVAL=5s

# Run the Go monitor
go run ./cmd/monitor --hosts="$HOSTS" --port=$PORT --interval=$INTERVAL
