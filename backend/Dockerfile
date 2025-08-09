FROM node:20-alpine

WORKDIR /app


COPY package*.json ./
RUN npm install -g @angular/cli && npm install


COPY . .


ENV NG_CLI_ANALYTICS=ci
ENV CHOKIDAR_USEPOLLING=true

EXPOSE 4200
CMD ["ng", "serve", "--host", "0.0.0.0", "--port", "4200", "--poll", "2000", "--proxy-config", "proxy.conf.json"]
