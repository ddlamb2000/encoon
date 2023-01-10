FROM golang:1.19

WORKDIR /usr/src/encoon

RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash
RUN . ~/.bashrc && nvm install node

COPY package.json package-lock.json babel.config.json ./
RUN . ~/.bashrc && npm install --save-dev @babel/core @babel/cli

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify


COPY . .

RUN go build -v
RUN . ~/.bashrc && npm run build

CMD ["./encoon"]