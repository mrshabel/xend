# Xend - A local file-server over HTTP

xend is a zero-config file server for developers, tinkerers, and anyone who needs to quickly share a local directory over HTTP.
It pairs perfectly with systems like [Cloudflare Tunnel](https://cloudflare.com/products/tunnel), letting you securely expose a local folder to the internet.

My motivation was to experiment with Cloudflare tunneling and HTTP while building something lightweight and useful.

## Features

-   Secure by Default
    Hidden files and directories (.git, .env, etc.) are automatically blocked from being served.
-   Gzip Compression
    Assets are compressed when the client supports it.
-   Zero Configuration
    Run it in any directory you want.

## Installation

There are two ways to install `xend`.

### 1. Download from Releases (Recommended)

You can download the latest pre-compiled binary for your operating system from the [Releases](https://github.com/mrshabel/xend/releases) page. This is the easiest way to get started.

### 2. Build from Source

If you prefer to build it yourself, you'll need Go installed.

**Prerequisites**

-   Go 1.24 or higher

**Steps**

1. **Clone the repository:**

    ```bash
    git clone https://github.com/mrshabel/xend.git
    cd xend
    ```

2. **Build the binary:**

    ```bash
    go build .
    ```

    This will create an executable named `xend` (or `xend.exe` on Windows) in the current directory.

## Usage

You can run the server by pointing it to the directory you wish to serve.

```bash
# serve the current directory on localhost:8000
./xend

# serve a specific directory on a different port
./xend -dir ./public -port 9090

# see all available options
./xend -h
```

### Options

| Flag    | Description                        | Default     |
| ------- | ---------------------------------- | ----------- |
| `-host` | HTTP network host to listen on     | `localhost` |
| `-port` | HTTP network port to listen on     | `8000`      |
| `-dir`  | Root directory to serve files from | `.`         |

## Use with Cloudflare Tunnel

`xend` is perfect for exposing a local folder to the internet securely with Cloudflare Tunnel.

1. **Start the `xend` server:**

    ```bash
    # Serve the 'example' folder on the default port 8000
    ./xend -dir /example
    ```

2. **Expose it with `cloudflared`:**

    In a new terminal, run the `cloudflared` command to create a tunnel to your local server.

    ```bash
    cloudflared tunnel --url http://localhost:8000
    ```

    This generates a random public URL that securely tunnels traffic to your local `xend` instance.
    You may also configure your own custom public domain to point to the instance on the Cloudflare dashboard

## References

-   [Cloudflare Tunnels](https://developers.cloudflare.com/pages/how-to/preview-with-cloudflare-tunnel/)
