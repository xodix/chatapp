FROM node:26.3-alpine AS build
WORKDIR /app
COPY . .
RUN npm install
RUN npm run build

FROM nginx:1.31

WORKDIR /usr/local/apache2
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/

EXPOSE 80
ENTRYPOINT ["nginx", "-g", "daemon off;"]