FROM debian


RUN apt-get update
RUN apt-get install -y curl

ADD run.sh /run.sh


ENTRYPOINT ["/run.sh"]
