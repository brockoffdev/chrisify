# quay.io/wework/faceify

FROM ubuntu:16.04

RUN apt-get update \
 && apt-get install python-dev libopencv-dev -y \
 && mkdir -p /app/www \
 && mkdir -p /app/faces

COPY faceify/faceify /app
COPY faceify/haarcascade_frontalface_alt.xml /app
COPY slackbot/slackbot /app

# Must provide slack token, location of faceify, haar file, location of face files, and url to serve photos from
CMD ["/app/slackbot, "$SLACK_API_TOKEN", "/app/faceify", "/app/haarcascade_frontalface_alt.xml", "/app/faces", "https://url.to.photos/"]
