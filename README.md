
# Video Describer

This project uses OpenAI's Vision API to describe YouTube videos.


## Usage/Examples

```bash
curl --location 'http://localhost:8080/v1/video/describe' \
--header 'Content-Type: application/json' \
--data '{
    "videoID": "dQw4w9WgXcQ",
    "prompt": "These are frames of a video. Describe what happens in the video in one sentence, including a description of the environments where it takes place.",
    "end": "30s"
}'
```

Response:

`
The video features various scenes including a close-up of a person's shoe tapping, a man singing into a microphone with a graphic background, a smiling individual in a brick environment, a dancing woman with bright backlighting, a person posing in a denim outfit by a chain-link fence, a woman dancing in an urban outdoor setting, and the same man performing with different expressions in the brick environment and by the fence.
`


## Environment Variables

| Parameter | Default | Description                |
| :-------- | :-------- | :------------------------- |
| `OPENAI_API_KEY` | "" | **Required**. API key from your [OpenAI](https://platform.openai.com/api-keys) account. |
| `SRV_HOST` | 0.0.0.0 | Server host. |
| `SRV_PORT` | 8080 | Server port. |
| `SRV_DEBUG` | false | Enable debug logs. |


## Run Locally

Requirements:

* [Go](https://go.dev/doc/install)
* [ffmpeg](https://ffmpeg.org/download.html)

Clone the project

```bash
  git clone https://github.com/WillRoque/video-describer
```

Go to the api directory

```bash
  cd cmd/api
```

Start the server

```bash
  go run .
```


## API Reference

#### Describe a video

```
  POST /v1/video/describe
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `videoID` | `string` | **Required**. ID of a YouTube video. |
| `prompt` | `string` | Custom prompt to describe the video. The default one is "These are frames of a video. Describe what happens in the video in one sentence." |
| `start` | `string` | Describe the video starting from this position. Use the string representation of a [time.Duration](https://pkg.go.dev/maze.io/x/duration#ParseDuration). |
| `end` | `string` | Describe the video up until this point. Use the string representation of a [time.Duration](https://pkg.go.dev/maze.io/x/duration#ParseDuration). |
