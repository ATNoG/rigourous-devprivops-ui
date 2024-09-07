# UI

This repository holds the code for the web editor of the devprivops tool.

The web editor is a website allowing to edit the files in the local directory of the tool.  
It provides not only text editors (powered by [monaco](https://microsoft.github.io/monaco-editor/)) but also visual forms to aid configuration maintenance.

## Deployment

The UI can be deployed either through the docker container or building the binary from source.

### Docker Container

To use the docker container, we run the following command

```sh
docker run \ 
    --env-file .env \ 
    -p 8082:8082 \ 
    -v "/tmp/.devprivops:/tmp/.devprivops" \ 
    --name devprivops-ui \ 
    devprivops-ui 
```

- The environment variables can be passed at once through the `--env-file` argument.
- The UI requires an exposed port, dictated by the `PORT` variable.
- The tool's local directory can be on the host machine and be passed to the container through the `-v` argument

### Manual Compilation

To manually compile the UI and execute it:

1. Set an adequate `.env` file. All the documentation for the relevant variables is in the `.env.example` file
2. Build the editor with `go build`
2. Run `./devprivops-ui`
3. Access the endpoint specified in the `.env` file in any web browser

## Development

While developing, we provide a configuration for [air](https://github.com/air-verse/air) that not only reloads the go code, but also templates and styles.

To turn on the air server, run `air .` on the repository's root.
