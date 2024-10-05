## Moved to Archive
[AndroidStory](https://androidstory.dev/) moved to new solution which provide admin pannel for podcast hosts and some automation around podcast hosting. 
New solution also open source [YourPod](https://github.com/VovaStelmashchuk/yourpod)

## Selfpod

The system for podcast owners to duplicate their podcast to youtube automatically by power of open source; free and
self-hosted. You control your podcast, not the other way around.

The application is a simple go app which receive the webhook from a cast and download the audio and cover files.
Then the app create a video from the audio and cover files by ffmpeg, generate the youtube based on description from the
rss feed. Then the app upload the video to youtube.

## The project goal

The one reason to create automation pipeline for Android Story podcast workflow. The vector of the project will be
changed in case of Android story podcast workflow changes. **Keep in mind, the project is not a product. I recommend to
research the project and get the idea which can be applied into your project.**

### Run app by docker compose

```bash
version: '3'

services:
  app:
    container_name: selfpod
    image: vovochkastelmashchuk/selfpod:0.2.0
    ports:
      - "6000:5000"
    environment:
      - A_CAST_HOOK_TOKEN=<your token, using verify is a cast execute the rest method>
      - APP_PORT=<app port>
      - GOOGLE_REDIRECT_URL=<your redirect url>
      - A_CAST_SHOW_ID=<a cast show id>
      - YOUTUBE_CHANNEL_ID=<youtube channel id>
      - DISCORD_WEBHOOK_URL=<discord webhook url>
      - DISCORD_BOT_NAME=<discord bot name>
    volumes:
      - ./for_docker_mount/:/application/tmp_files/
      - ./files/:/application/files/
```

The google redirect url will be `<your redirest url>/auth/google/callback`

Create file in the root `client_secret.json` with the following content:

```json
{
  "web": {
    "client_id": "<your client id>",
    "project_id": "<your project id>",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
    "client_secret": "<your client secret>",
    "redirect_uris": [
      "<your redirect url>/auth/google/callback"
    ]
  }
}
``` 

## Using

The app support only one podcast per one instance. You can run multiple instances for multiple podcasts.
Also its means the app support only one rss feed, one youtube account, etc per one instance.

Currently, the app using by [Android Story podcast](https://www.youtube.com/channel/UC6-NFk4uOGsKvyisL1QC3rw). The app
run on Raspberry Pi 4 with 2GB RAM at my home. And we successfully downgrade our Acast plan to free.

## TODO

- Some UI improvements
- Write some tool for generating the public and private podcasts.
- Auto podcast editing tool.

## Build docker image

Check the building process in ci files.

### API

Look to api in file api.http
