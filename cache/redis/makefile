SHELL = /bin/bash

init: config
 for port in `seq 6380 6385`; do \
   mkdir -p ./${port}/config \
   && IP=10.10.4.42 PORT=${port} envsubst < ./redis-cluster.tmpl > ./${port}/config/redis.conf \
   && mkdir -p ./${port}/data && chmod 777 ./${port}/config/redis.conf;\
 done

