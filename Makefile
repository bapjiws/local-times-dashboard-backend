run:
	revel run timezones_mc/revel_app

get_cities:
	mkdir cities && curl http://download.maxmind.com/download/worldcities/worldcitiespop.txt.gz | gunzip -c > cities/worldcities.txt

get_elasticsearch:
	mkdir .elasticsearch && \
	curl https://download.elastic.co/elasticsearch/release/org/elasticsearch/distribution/tar/elasticsearch/2.3.5/elasticsearch-2.3.5.tar.gz |\
	tar -xvzf - -C .elasticsearch

start_elasticsearch:
	.elasticsearch/elasticsearch-2.3.5/bin/elasticsearch -d -p es_pid

stop_elasticsearch:
	kill `cat es_pid`

get_deps:
	go get -u gopkg.in/olivere/elastic.v3
	go get -u github.com/revel/cmd/revel
	go get -u github.com/satori/go.uuid
