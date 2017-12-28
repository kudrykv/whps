# Webhook Proxy Server

This application helps to route webhooks to localhost.

```
  ╔══════════════════╗            ┌──────────────────┐
  ║                  ║            │                  │
  ║       whps       ║<──webhook──│     Service      │
  ║                  ║            │                  │
  ╚══════════════════╝            └──────────────────┘
            ^
            ┃
       websocket
            ┃
┌───────────v──────────────────────────────Local machine
│ ┌──────────────────┐            ┌──────────────────┐ │
│ │                  │            │                  │ │
│ │       whpc       │──webhook──>│   Application    │ │
│ │                  │            │                  │ │
│ └──────────────────┘            └──────────────────┘ │
└──────────────────────────────────────────────────────┘
```

## # whps (Webhook Proxy Server)
This is a server. Client (whpc) opens websocket connection
to the server. Then, `Service` can send webhook to the server and
it will route request to the client via websocket passing request
headers and body. Now the client will be able to reproduce the webhook
to the local application.

## # whpc (Webhook Proxy Client)
[Client implementation](https://github.com/kudrykv/whpc) to establish
websocket connection to the server and relay message to the local
application.

# Usage

Service can push webhooks to the
`https://whps.herokuapp.com/webhook/<channame>`.
The `<channame>` can be arbitrary alphanumeric name of your choice.
E.g. `https://whps.herokuapp.com/webhook/betazoid`.

To receive webhook locally one should open websocket connection beforehand
to the `https://whps.herokuapp.com/websocket/betazoid`. Here, `betazoid`
is the `<channame>` as in the previous url.

Now, as soon as `whps` receives the webhook, it will relay it to the client.
