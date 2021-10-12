FROM scratch
COPY revive /usr/bin/revive
ENTRYPOINT ["/usr/bin/revive"]
