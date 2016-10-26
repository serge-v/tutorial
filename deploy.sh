host=$(cat host~)
rsync -rvaz content template ${host}:
ssh ${host} sudo cp -r content /usr/local/www/wet/
ssh ${host} sudo cp -r template /usr/local/www/wet/
