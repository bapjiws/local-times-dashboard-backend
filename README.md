[![Build Status](https://travis-ci.org/bapjiws/timezones_mc.svg?branch=master)](https://travis-ci.org/bapjiws/timezones_mc)
[![Go Report Card](https://goreportcard.com/badge/github.com/bapjiws/timezones_mc)](https://goreportcard.com/report/github.com/bapjiws/timezones_mc)

This is the MC part (the V part abides [here](https://github.com/bapjiws/timezones_v)) of a simple webapp that allows to create a list of cities with their respective local times (updated in real time) via auto-completion. It consists of a Go server that's hooked up to an ElasticSearch service. Both are dockerized, and the entire backed can be fired up with the ```docker-compose up``` command run from the root directory.
