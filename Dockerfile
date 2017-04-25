FROM golang:alpine
WORKDIR ./src/github.com/bapjiws/local_times_dashboard_backend
COPY app ./app
COPY models ./models
COPY vendor ./vendor
COPY datastore ./datastore
COPY utils ./utils

CMD go run app/app.go # See: https://docs.docker.com/compose/startup-order/
#CMD /bin/sh

### ALTERNATIVE: ###
#FROM ubuntu

#  -m, --create-home             create the user's home directory
# -s, --shell SHELL             login shell of the new account
#RUN useradd -ms /bin/bash admin # non-priviledged user
#USER admin

#ADD https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz .
#RUN tar -C /usr/local -xzf go1.8.linux-amd64.tar.gz
#RUN echo "export PATH=$PATH:/usr/local/go/bin" >> $HOME/.bashrc
#RUN rm go1.8.linux-amd64.tar.gz

### USEFUL COMMANDS: ###
#docker build -t bapjiws/timezones_api:0.0.1 .
#docker run -it --rm --name api bapjiws/timezones_api:0.0.1
#docker rmi -f $(docker images | grep none)