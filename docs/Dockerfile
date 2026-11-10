FROM node:24-alpine AS build
WORKDIR /site
COPY docs/package*.json ./
RUN npm ci
COPY docs/ ./
RUN npm run build
FROM nginx:stable-alpine
COPY docs/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /site/dist /usr/share/nginx/html
EXPOSE 80
