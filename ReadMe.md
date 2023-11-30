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
`HOOK_TOKEN=hook_123_token`

## Using

The app support only one podcast per one instance. You can run multiple instances for multiple podcasts.
Also its means the app support only one rss feed, one youtube account, etc per one instance.

## open source todo

- Move a cast show to environment variable
- Add logout method
- Some UI improvements
