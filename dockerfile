FROM golang:1.22.0

WORKDIR /app

# Set environment variables
ENV NODE_ENV production
ENV PORT 80

# install bun
RUN apt-get update && apt-get install -y unzip
RUN curl -fsSL https://bun.sh/install | BUN_INSTALL=/usr bash

COPY package.json .
COPY bun.lockb .

RUN bun install --production

COPY tailwind.config.ts .
COPY tsconfig.json .
COPY . .

RUN bun run build
# RUN go generate
RUN go build -o ./tmp/main .

EXPOSE 80

CMD ["./tmp/main"]