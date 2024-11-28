<!--
    - icon
    - buttons
    - project description
    - features
    - use cases
-->

# Web UI

This repository holds the code for the web editor of the devprivops tool.

The web editor is a website allowing to edit the files in the local directory of the tool.  
It provides not only text editors (powered by [monaco](https://microsoft.github.io/monaco-editor/)) but also visual forms to aid configuration development and maintenance.
It also provides a plugin system to create custom visualizers for each kind of system description.

# Deployment

The UI can be deployed either through the docker container or as a native binary.

## Docker Container

To use the docker container, we run the following command

```sh
docker run \ 
    --env-file .env \ 
    -p 8082:8082 \ 
    -v ".devprivops:/tmp/.devprivops" \ 
    --name devprivops-ui \ 
    devprivops-ui 
```

- The environment variables can be passed at once through the `--env-file` argument, configured according to the documentation in the `.env.example` file.
- The UI requires an exposed port, dictated by the `PORT` variable.
- The tool's local directory can be on the host machine and be passed to the container as a bind mount.

## Manual Compilation

To manually compile the UI and execute it:

1. Install [PrivGuide](https://github.com/ATNoG/rigourous-devprivops) **natively** and copy the executable to the root of this repository.
1. Set an adequate `.env` file. All the documentation for the relevant variables is in the `.env.example` file
2. Build the editor with `go build`
    - On the first execution, it is needed to generate the styles and templates with `tailwindcss -i static/css/source.css -o static/css/style.css --minify` and `templ generate`
    - Running `air` will also execute these commands and is advised during development
2. Run `./devprivops-ui`
3. Access the endpoint specified in the `.env` file in any web browser

# Development

While developing, we provide a configuration for [air](https://github.com/air-verse/air) that not only reloads the go code, but also the templates and styles.

To turn on the air server, run `air` on the repository's root.
Then, whenever files change the tool will automatically be rerun and ready to test.

To manually generate these files and compile the code, run the following commands:

```sh
tailwindcss -i static/css/source.css -o static/css/style.css --minify
templ generate
go build
```

For better dependency management, we provide a `shell.nix` file with all needed dependencies.
To use it, install `nix` and execute `nix-shell --extra-experimental-features "flakes"`.
All dependencies can also be installed manually through other means.

# Features

This interface allows to:

- Edit the raw files for trees, descriptions, reasoner rules, regulations, extra data queries, requirements and schemas;
- Edit those files visually with user-friendly editors;
- Colaborate seamlessly in the development and maintenance of the files;
- Safely resolve conflicts arising from the colaboration;
- Use custom visualizers for groups of descriptions;
- Run analysis and tests;

# Usage

...
