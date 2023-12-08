## Selfpod

The system for podcast owners to manage their podcast by free and self-hosted.
You control your podcast, not the other way around.

### Run app in docker

```bash
docker build -t webhook .
```

```bash
docker run --name webhook -p 6000:5000 -v "$(pwd)/for_docker_mount/:/application/tmp_files/" --rm webhook
```

Environment variables:
`A_CAST_HOOK_TOKEN=hook_123_token`
`APP_PORT=5000`
`GOOGLE_REDIRECT_URL=http://localhost:6000/`

The google redirect url will be `http://localhost:6000/auth/google/callback`

## Using

The app support only one podcast per one instance. You can run multiple instances for multiple podcasts.
Also its means the app support only one rss feed, one youtube account, etc per one instance.

## open source todo

- Move a cast show to environment variable
- Add logout method
- Some UI improvements

## Build docker image

```sh
docker buildx create --use
docker buildx build --platform linux/arm64/v8 -t vovochkastelmashchuk/selfpod:0.0.1 --push .
```

Redis:
port: 6379
