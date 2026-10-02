FROM nginx:1.28-alpine@sha256:a8b39bd9cf0f83869a2162827a0caf6137ddf759d50a171451b335cecc87d236
COPY nginx.conf.template /etc/nginx/templates/default.conf.template
COPY index.html styles.css favicon.svg /usr/share/nginx/html/
COPY src /usr/share/nginx/html/src
EXPOSE 80
