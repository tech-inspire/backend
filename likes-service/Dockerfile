FROM node:20-alpine

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

COPY healthcheck.sh /usr/local/bin/healthcheck.sh
RUN chmod +x /usr/local/bin/healthcheck.sh

HEALTHCHECK \
  --interval=30s \
  --start-interval=2s \
  --timeout=3s \
  --start-period=5s \
  --retries=5 \
  CMD /usr/local/bin/healthcheck.sh

CMD ["npm", "start"]