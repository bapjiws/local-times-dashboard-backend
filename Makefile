RAW_DATA_FOLDER = .raw_data
TOOLS_FOLDER = .tools

CITIES_URL = http://download.maxmind.com/download/worldcities/worldcitiespop.txt.gz
ELASTICSEARCH_URL = https://download.elastic.co/elasticsearch/release/org/elasticsearch/distribution/tar/elasticsearch/2.4.0/elasticsearch-2.4.0.tar.gz
KIBANA_URL = https://download.elastic.co/kibana/kibana/kibana-4.6.1-linux-x86_64.tar.gz

start_app:
	go run app/app.go

get_cities:
	mkdir -p $(RAW_DATA_FOLDER) && mkdir -p $(RAW_DATA_FOLDER)/cities && curl $(CITIES_URL) | gunzip -c > $(RAW_DATA_FOLDER)/cities/worldcities.txt

get_elasticsearch:
	mkdir -p $(TOOLS_FOLDER) && mkdir -p $(TOOLS_FOLDER)/elasticsearch && \
	curl $(ELASTICSEARCH_URL) |\
	tar -xvz -C $(TOOLS_FOLDER)/elasticsearch --strip-components=1

start_elasticsearch:
	$(TOOLS_FOLDER)/elasticsearch/bin/elasticsearch -d -p $(TOOLS_FOLDER)/elasticsearch/elasticsearch.pid

stop_elasticsearch:
	kill `cat $(TOOLS_FOLDER)/elasticsearch/elasticsearch.pid`

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

run: start_elasticsearch start_kibana start_app

stop: stop_elasticsearch stop_kibana

deps:
	go get gopkg.in/olivere/elastic.v3
	go get github.com/gin-gonic/gin
	go get github.com/satori/go.uuid
	go get github.com/stretchr/testify
	go get github.com/vektra/mockery/.../
	go get github.com/tools/godep

test:
	go test -v $$(go list ./... | grep -v /vendor/)

vet:
	go vet -v $$(go list ./... | grep -v /vendor/)

build:
	go build $$(go list ./... | grep -v /vendor/)

check: vet test build

set_hooks:
	sudo chmod +x hooks/pre-commit && cd .git/hooks && ln -sf ../../hooks/pre-commit pre-commit


reindex_cities:
	go run scripts/cities/main.go -file="$(RAW_DATA_FOLDER)/cities/worldcities.txt"
