# when running local without compose:
    # docker build -t svelte .
    # docker run -p 3000:3000 svelte
    # http://localhost:3000/

FROM node:20 AS builder

ENV PUBLIC_AUTH_DOMAIN="dev-boc0av4c0bacnon6.us.auth0.com"
ENV PUBLIC_AUTH_CLIENT_ID="si9esdik2n4z3VjzzMZR78H0mwUMsWFW"
ENV PUBLIC_AUTH_AUDIENCE="http://localhost:10000"
ENV PUBLIC_AUTH_CACHE_LOCATION="localstorage"
ENV PUBLIC_BACKEND_URL="http://192.168.178.84:10000"
ENV PUBLIC_REQUESTER_NAME="dev_janos_desktop"

WORKDIR /app

COPY . .

RUN npm install
RUN npm run build

FROM node:20

WORKDIR /app

COPY package*.json .
COPY --from=builder /app/build .

CMD [ "node", "index.js" ]