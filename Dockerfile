FROM dockette/debian:stretch-slim

RUN apt-get update && apt-get install -y nginx

WORKDIR /home

COPY ./text.txt /home/

EXPOSE 80 443

ENTRYPOINT ["nginx"]