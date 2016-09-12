RAW_DATA_FOLDER = .raw_data
TOOLS_FOLDER = .tools

CITIES_URL = http://download.maxmind.com/download/worldcities/worldcitiespop.txt.gz
ELASTICSEARCH_URL = https://download.elastic.co/elasticsearch/release/org/elasticsearch/distribution/tar/elasticsearch/2.4.0/elasticsearch-2.4.0.tar.gz
KIBANA_URL = https://download.elastic.co/kibana/kibana/kibana-4.6.1-linux-x86_64.tar.gz

run:
	revel run timezones_mc/revel_app

get_cities:
	mkdir -p $(RAW_DATA_FOLDER) && mkdir -p $(RAW_DATA_FOLDER)/cities && curl $(CITIES_URL) | gunzip -c > $(RAW_DATA_FOLDER)/cities/worldcities.txt

get_elasticsearch:
	mkdir -p $(TOOLS_FOLDER) && mkdir -p $(TOOLS_FOLDER)/elasticsearch && \
	curl $(ELASTICSEARCH_URL) |\
	tar -xvz -C $(TOOLS_FOLDER)/elasticsearch --strip-components=1

start_elasticsearch:
	$(TOOLS_FOLDER)/elasticsearch/bin/elasticsearch -d -p elasticsearch.pid

stop_elasticsearch:
	kill `cat elasticsearch.pid`

clean_elasticsearch:
	curl -XDELETE "localhost:9200/*"

get_kibana:
	mkdir -p $(TOOLS_FOLDER) && mkdir -p $(TOOLS_FOLDER)/kibana && \
	curl $(KIBANA_URL) |\
	tar -xvz -C $(TOOLS_FOLDER)/kibana --strip-components=1

start_kibana:
	$(TOOLS_FOLDER)/kibana/bin/kibana > /dev/null 2>&1 &

stop_kibana:
	ps aux | grep "kibana" | awk '{print $$2}' | xargs kill

get_sense:
	$(TOOLS_FOLDER)/kibana/bin/kibana plugin --install elastic/sense

get_tools: get_elasticsearch get_kibana get_sense

start_tools: start_elasticsearch start_kibana

stop_tools: stop_elasticsearch stop_kibana

deps:
	go get gopkg.in/olivere/elastic.v3
	go get github.com/revel/cmd/revel
	go get github.com/satori/go.uuid
	go get github.com/stretchr/testify
	go get github.com/vektra/mockery/.../
	go get github.com/tools/godep

test:
	go test -v datastore/*.go

reindex_cities:
	go run scripts/cities/main.go -file="$(RAW_DATA_FOLDER)/cities/worldcities.txt"
