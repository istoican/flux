FROM debian

ADD flux-server /usr/bin/flux-server

RUN ["flux-server"]