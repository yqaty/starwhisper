from centos:7
COPY main /web/
COPY config.json /web/
CMD ["/web/main"]
