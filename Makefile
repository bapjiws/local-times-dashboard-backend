run:
	revel run timezones_mc/revel_app

get_cities:
	mkdir cities && curl http://download.maxmind.com/download/worldcities/worldcitiespop.txt.gz | gunzip -c > cities/worldcities.txt

get_elasticsearch:
	mkdir .elasticsearch && \
	curl https://download.elastic.co/elasticsearch/release/org/elasticsearch/distribution/tar/elasticsearch/2.4.0/elasticsearch-2.4.0.tar.gz |\
	tar -xvz -C .elasticsearch --strip-components=1

start_elasticsearch:
	.elasticsearch/bin/elasticsearch -d -p es_pid

stop_elasticsearch:
	kill `cat es_pid`

clean_elasticsearch:
	curl -XDELETE "localhost:9200/*"

get_kibana:
	mkdir .kibana && \
	curl https://download.elastic.co/kibana/kibana/kibana-4.6.0-linux-x86_64.tar.gz |\
	tar -xvz -C .kibana --strip-components=1

get_deps:
	go get -u gopkg.in/olivere/elastic.v3
	go get -u github.com/revel/cmd/revel
	go get -u github.com/satori/go.uuid

reindex_cities:
	go run scripts/cities/main.go -file="cities/worldcities.txt"
