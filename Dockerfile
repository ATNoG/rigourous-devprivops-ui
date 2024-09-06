# Use image with golang and docker
FROM 192.168.56.1:5000/devprivops:latest

COPY . /src/
WORKDIR /src/ 

# Build the application
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN apt update -y
RUN apt install nodejs npm -y
RUN npm install -D tailwindcss
RUN npx tailwindcss -i static/css/source.css -o static/css/style.css --minify
RUN go mod tidy
RUN go build

# Cleanup
RUN mv devprivops-ui /bin/devprivops-ui
RUN rm -rf /src/

