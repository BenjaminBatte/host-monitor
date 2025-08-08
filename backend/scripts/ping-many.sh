#!/bin/bash

HOSTS="actiontarget.com,google.com,facebook.com,openai.com,cloudflare.com,github.com,\
linkedin.com,youtube.com,stackoverflow.com,netflix.com,zoom.us,\
oracle.com,intel.com,tesla.com,airbnb.com,reddit.com,\
spotify.com,discord.com,slack.com,docker.com,kubernetes.io,\
ubuntu.com,debian.org,heroku.com,aws.amazon.com,azure.microsoft.com,\
cloud.google.com,wikipedia.org,nytimes.com,cnn.com,reuters.com,\
bbc.com,weather.com,timeanddate.com,accuweather.com,\
whitehouse.gov,usa.gov,nasa.gov,nih.gov,nps.gov,\
irs.gov,fbi.gov,dhs.gov,cdc.gov,\
treasury.gov,senate.gov,house.gov,ssa.gov,uscourts.gov,\
cumberlands.edu,maharishi.edu,harvard.edu,mit.edu,stanford.edu,\
berkeley.edu,princeton.edu,yale.edu,columbia.edu,umich.edu,\
ucla.edu,utexas.edu,cmu.edu,purdue.edu,gatech.edu,bloomberg.com,forbes.com,adp.com,indeed.com,glassdoor.com,\
expedia.com,booking.com,tripadvisor.com,aircanada.com,delta.com,\
intel.com,amd.com,nvidia.com,cisco.com,dell.com,\
hp.com,lenovo.com,samsung.com,lg.com,panasonic.com,\
sony.com,hitachi.com,siemens.com,philips.com,ericsson.com,\
ups.com,usps.com,fedex.com,dhl.com,ikea.com,\
lego.com,pepsi.com,cocacola.com,nestle.com,unilever.com,\
pwc.com,deloitte.com,kpmg.com,ey.com,honeywell.com,\
shell.com,chevron.com,exxonmobil.com,boeing.com,lockheedmartin.com,\
raytheon.com,northropgrumman.com,3m.com,ge.com,basf.com,who.int,unesco.org,un.org,worldbank.org,imf.org,\
nasa.gov,nsf.gov,energy.gov,cdc.gov,nih.gov,\
noaa.gov,usa.gov,data.gov,whitehouse.gov,house.gov,\
senate.gov,uscourts.gov,nps.gov,irs.gov,ssa.gov,\
fao.org,oecd.org,wipo.int,wto.org,itu.int,\
mit.edu,harvard.edu,stanford.edu,yale.edu,princeton.edu,\
caltech.edu,cmu.edu,utexas.edu,berkeley.edu,ucla.edu,\
uncc.edu,purdue.edu,gatech.edu,rice.edu,upenn.edu,\
nyu.edu,columbia.edu,uic.edu,uic.edu,northwestern.edu,\
auburn.edu,duke.edu,unc.edu,psu.edu,msu.edu"

PORT=80
INTERVAL=5s

go run ./cmd/monitor --hosts="$HOSTS" --port=$PORT --interval=$INTERVAL
