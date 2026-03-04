FROM nginx:mainline-alpine
COPY init_proxy.sh /init.sh
RUN apk update && apk add --no-cache openssl
COPY ./nginx.conf /etc/nginx/nginx.conf
ENTRYPOINT ["/bin/sh", "/init.sh"]
CMD ["nginx", "-g", "daemon off;"]