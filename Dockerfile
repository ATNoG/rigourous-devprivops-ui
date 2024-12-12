# Base image
FROM golang:1.23.3-alpine3.20 as build

# Install fuseki
RUN apk update && apk add \
	wget \
	unzip \
    git \
    nodejs \
    npm

WORKDIR /opt
RUN wget https://dlcdn.apache.org/jena/binaries/apache-jena-fuseki-5.2.0.tar.gz && \
    tar xzf apache-jena-fuseki-5.2.0.tar.gz && \
    rm apache-jena-fuseki-5.2.0.tar.gz && \
    mv apache-jena-fuseki-5.2.0 fuseki

COPY shiro.ini /opt/fuseki/shiro.ini

# Build PrivGuide
RUN git clone https://github.com/ATNoG/rigourous-devprivops.git /privguide_src
WORKDIR /privguide_src
RUN go mod tidy && \
    go build

# Build UI
COPY . /src
WORKDIR /src
RUN go install github.com/a-h/templ/cmd/templ@latest && \
    templ generate && \
    npm install -D tailwindcss && \
    npx tailwindcss -i static/css/source.css -o static/css/style.css --minify && \
    go mod tidy && \
    go build

# Execution environment
FROM alpine:latest

# Add runtime dependencies
RUN apk update && apk add \
    openjdk21-jre \
    git

# Get executables
COPY --from=build /opt/fuseki /opt/fuseki
COPY --from=build /privguide_src/devprivops /usr/local/bin/
COPY --from=build /src/devprivops-ui /usr/local/bin/
COPY --from=build /src/static/ /var/www/privguide/static

# Allow host directories to be used as git repositories, fixing "dubious ownership"
RUN git config --system --add safe.directory '*'

# Start fuseki server
CMD /opt/fuseki/fuseki-server --mem --port=3030 /tmp & devprivops-ui