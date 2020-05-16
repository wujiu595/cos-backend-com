FROM hub.jimu.io/jimu/alpine

COPY hack/build/account /usr/local/bin/cores

ENV RUN_MODE=prod HTTP_ADDR=0.0.0.0 HTTP_PORT=80

EXPOSE 80

CMD ["cores"]
