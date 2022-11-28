FROM ubuntu

COPY ./gok8s ./gok8s

ENTRYPOINT [ "./gok8s" ]