#!/usr/bin/env bash
set -euo pipefail

# Load .env safely
if [ -f .env ]; then
  echo "[ping-many] Loading environment variables from .env"
  set -a; . ./.env; set +a
fi

# DB defaults
: "${DB_HOST:=localhost}"
: "${DB_PORT:=5432}"
: "${DB_NAME:=postgres}"
: "${DB_USER:=postgres}"
: "${DB_PASS:=2020}"

# Build DB_URL with application_name (easier to see in pg_stat_activity)
if [ "${DB_URL-}" = "" ]; then
  DB_URL="postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&application_name=hostmonitor"
else
  case "$DB_URL" in
    *\?*) DB_URL="${DB_URL}&application_name=hostmonitor" ;;
    *)    DB_URL="${DB_URL}?application_name=hostmonitor" ;;
  esac
fi
export DB_URL

# Full default HOSTS (your original)
HOSTS="${1:-actiontarget.com,google.com,facebook.com,openai.com,cloudflare.com,github.com,\
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
auburn.edu,duke.edu,unc.edu,psu.edu,msu.edu}"
PORT="${2:-80}"
INTERVAL="${3:-5s}"

echo "[ping-many] Using DB_URL=$DB_URL"
echo "[ping-many] Monitoring hosts: $HOSTS"
echo "[ping-many] Port: $PORT | Interval: $INTERVAL"

# Quick DB preflight (fails fast if wrong)
psql "$DB_URL" -v ON_ERROR_STOP=1 -P pager=off -c 'SELECT 1;' >/dev/null
echo "[ping-many] DB connectivity OK"

# Ensure checks table exists (runs your migration if needed) - non-destructive
if [ "$(psql "$DB_URL" -tA -P pager=off -c "SELECT to_regclass('public.checks') IS NOT NULL;")" != "t" ]; then
  echo "[ping-many] Creating public.checks via migrations/001_init.sql"
  psql "$DB_URL" -v ON_ERROR_STOP=1 -P pager=off -f migrations/001_init.sql
fi

# Run Go backend (pass DB_URL; add --db-url if your app supports it)
# backend/scripts/ping-many.sh — change the last line to:
go run ./cmd/monitor --hosts="$HOSTS" --port="$PORT" --interval="$INTERVAL"

