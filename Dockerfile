FROM ubuntu:24.04

RUN apt-get update \ 
    && apt-get install -y ffmpeg \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY cmd/api/api /api

EXPOSE 8080

CMD [ "/api" ]
