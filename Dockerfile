FROM drone/ca-certs

ADD src/gateway /

CMD ["/gateway"]
