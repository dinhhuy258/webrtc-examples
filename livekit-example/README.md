# Livekit example

## Frontend

```
npm start
```

## Backend

Start server

```
docker run --rm \
    -p 7880:7880 \
    -p 7881:7881 \
    -p 7882:7882/udp \
    -v $PWD/livekit.yaml:/livekit.yaml \
    livekit/livekit-server \
    --config /livekit.yaml \
    --node-ip=127.0.0.1
```

Generate token

```
docker run --rm -e LIVEKIT_KEYS="APIuPexPx9YaYu6: OO6mCCdZGdsfAFxty91YU0qeafuirffOloLwY4J9CNyH" \
    livekit/livekit-server create-join-token \
    --room "Demo livekit" \
    --identity "Huy Duong - 1"
```
