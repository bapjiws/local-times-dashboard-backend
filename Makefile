run:
	revel run timezones_mc/revel_app

get_cities:
	mkdir cities && curl http://download.maxmind.com/download/worldcities/worldcitiespop.txt.gz | gunzip -c > cities/worldcities.txt

get_elasticsearch:
	mkdir -p .tools && mkdir -p .tools/elasticsearch && \
	curl https://download.elastic.co/elasticsearch/release/org/elasticsearch/distribution/tar/elasticsearch/2.4.0/elasticsearch-2.4.0.tar.gz |\
	tar -xvz -C .tools/elasticsearch --strip-components=1

start_elasticsearch:
	.tools/elasticsearch/bin/elasticsearch -d -p elasticsearch.pid

stop_elasticsearch:
	kill `cat elasticsearch.pid`

clean_elasticsearch:
	curl -XDELETE "localhost:9200/*"

get_kibana:
	mkdir -p .tools && mkdir -p .tools/kibana && \
	curl https://download.elastic.co/kibana/kibana/kibana-4.6.0-linux-x86_64.tar.gz |\
	tar -xvz -C .tools/kibana --strip-components=1

start_kibana:
	.tools/kibana/bin/kibana > /dev/null 2>&1 &

stop_kibana:
	ps aux | grep "kibana" | awk '{print $$2}' | xargs kill

get_tools: get_elasticsearch get_kibana

get_deps:
	go get -u gopkg.in/olivere/elastic.v3
	go get -u github.com/revel/cmd/revel
	go get -u github.com/satori/go.uuid

reindex_cities:
	go run scripts/cities/main.go -file=".raw_data/cities/worldcities.txt"
