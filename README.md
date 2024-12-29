# OpenVSCode-Server-DEB

A utility to build `.deb` packages for OpenVSCode-Server tailored for various architectures.

## Overview

This repository provides a simple tool to build Debian (`.deb`) packages of OpenVSCode-Server for specific versions and architectures. It leverages a Go-based script to automate the packaging process.

## Prerequisites

Ensure you have the following installed:

- [Go](https://go.dev/doc/install) (version 1.23.4 or higher recommended)
- `dpkg`, `curl`, `sed`, `tar`, `rm`, `mkdir`, `chmod` and related tools (typically available on Debian-based systems)

## Usage

### 1. Compile the `build.go` Script

To start, compile the `build.go` script into an executable named `build`:

```bash
go build -o build main.go
```

### 2. Build the `.deb` Package

Run the `build` executable with the desired version and architecture.

#### For `amd64`:

```bash
./build --version 1.96.0 --arch amd64
```

#### For `arm64`:

```bash
./build --version 1.96.0 --arch arm64
```

### 3. Clean Build Artifacts

To remove any generated files and clean up the working directory:

```bash
./build --clean
```

## Options

| Option         | Description                                          |
|----------------|------------------------------------------------------|
| `--version`    | Specifies the OpenVSCode-Server version to package. |
| `--arch`       | Sets the target architecture (`amd64`, `arm64`).     |
| `--clean`      | Cleans up build artifacts.                           |

## Example Workflow

1. Compile the build script:

    ```bash
    go build -o build main.go
    ```

2. Build the `.deb` package for `amd64`:

    ```bash
    ./build --version 1.96.0 --arch amd64
    ```

3. Clean up the workspace:

    ```bash
    ./build --clean
    ```

## License

This project is licensed under the [MIT License](LICENSE).

## Contributing

Contributions are welcome! Please open an issue or submit a pull request if you'd like to contribute to the project.

## Contact

For questions or feedback, feel free to open an issue in this repository. <br />
Join the discussion: [gitpod-io/openvscode-server#587](https://github.com/gitpod-io/openvscode-server/discussions/587)
