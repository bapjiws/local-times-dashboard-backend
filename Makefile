get_cities:
	mkdir cities && curl http://download.maxmind.com/download/worldcities/worldcitiespop.txt.gz | gunzip -c > cities/worldcities.txt



