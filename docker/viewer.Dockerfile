FROM node:25.8.0 AS build
WORKDIR /app
COPY www/package*.json ./
RUN npm install
COPY www/ ./
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist/frontend/browser/ /usr/share/nginx/html/
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
