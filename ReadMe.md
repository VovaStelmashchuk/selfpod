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