FROM drone/ca-certs

ADD config/config.yaml /etc/skaioskit/config.yaml
ADD src/gateway /

CMD ["/gateway"]
