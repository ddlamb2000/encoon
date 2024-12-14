# εncooη : data structuration, presentation and navigation.
# Copyright David Lambert 2023

# Go build
FROM golang:1.20-buster as go-builder
WORKDIR /usr/src/encoon
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY apis apis
COPY configuration configuration
COPY database database
COPY model model
COPY utils utils
COPY encoon.go .
RUN go build -v

# Node build
FROM node:18.13.0 as node-builder
WORKDIR /usr/src/encoon
COPY frontend/package.json frontend/babel.config.json ./
RUN npm install --save-dev @babel/core @babel/cli
RUN npm install --save-dev babel-preset-minify
COPY frontend/javascript javascript
COPY frontend/react react
RUN npm run build

# Image
FROM almalinux:9.1
LABEL application="εncooη"
LABEL description="Data structuration, presentation and navigation."
LABEL version="2.0"
RUN useradd -ms /bin/bash encoon
USER encoon
WORKDIR /usr/encoon
COPY --from=go-builder /usr/src/encoon/encoon /usr/encoon
COPY frontend/lib frontend/lib
COPY frontend/templates frontend/templates
COPY frontend/encoon.css frontend/encoon.css
COPY frontend/favicon.ico frontend/favicon.ico
COPY --from=node-builder /usr/src/encoon/javascript frontend/javascript
COPY seedData.json .
CMD ["sh", "-c", "./encoon"]
