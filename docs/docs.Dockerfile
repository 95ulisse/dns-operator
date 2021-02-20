FROM alpine:3

COPY requirements.txt /mkdocs/
WORKDIR /mkdocs
VOLUME /mkdocs

RUN apk --no-cache --no-progress add py3-pip gcc musl-dev python3-dev \
  && pip3 install -r requirements.txt