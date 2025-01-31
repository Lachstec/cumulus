# when running local without compose:
    # docker build -t svelte .
    # docker run -p 3000:3000 svelte
    # docker run -p 3000:3000 --env-file ./.env svelte
    # http://localhost:3000/

FROM node:20 AS builder

WORKDIR /app

COPY . .

RUN npm install
RUN npm run build

FROM node:20

WORKDIR /app

COPY package*.json .
COPY --from=builder /app/build .

CMD [ "node", "index.js" ]