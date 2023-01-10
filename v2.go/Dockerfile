
# Go build container

FROM golang:1.19 as go-builder
WORKDIR /usr/src/encoon
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v

# Node build container

FROM node:18.13.0 as node-builder
WORKDIR /usr/src/encoon
COPY package.json babel.config.json ./
COPY frontend frontend
RUN npm install --save-dev @babel/core @babel/cli
RUN npm run build

# Production container

FROM debian:11-slim
WORKDIR /usr/encoon
COPY --from=go-builder /usr/src/encoon/encoon /usr/encoon
COPY --from=node-builder /usr/src/encoon/frontend frontend
COPY configuration.yml .
CMD ["./encoon", "-log", "encoon.log"]
