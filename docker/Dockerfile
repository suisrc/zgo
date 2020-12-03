FROM debian:buster-slim

RUN mkdir -p /www/
WORKDIR /www

EXPOSE 80
ENTRYPOINT ["./zgo", "web", "-c", "config.toml"]

#deploy files
ADD ["zgo", "config.toml", "rbac_*", "/www/"]