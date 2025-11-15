# xend - A local file-server over HTTP

A local file-server designed for the curious minds and anyone looking to quickly share files over HTTP. It's an excellent companion for services like [Cloudflare Tunnels](https://www.cloudflare.com/products/tunnel/), allowing you to expose a local directory to the internet securely.
My motivation was to experiment Cloudflare tunneling and HTTP.

## Features

-   **Secure by Default**: Automatically prevents hidden files and directories (like `.git` or `.env`) from being exposed.
-   **Gzip Compression**: Compresses assets for clients that support it, ensuring faster load times.

## Prerequisites (Development)

-   Go 1.24 or higher

## Installation

1.  **Clone the repository**

    ```bash
    git clone https://github.com/mrshabel/xend.git
    cd xend
    ```

2.  **Build the binary**

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

`xend` is perfect for exposing a local folder to the internet securely with Cloudflare Tunnels.

1.  **Start the `xend` server:**

    ```bash
    # Serve the 'example' folder on the default port 8000
    ./xend -dir /example
    ```

2.  **Expose it with `cloudflared`:**

    In a new terminal, run the `cloudflared` command to create a tunnel to your local server. This generates a random URL.
    You can configure your public domain to point to the instance on the cloudflare dashboard

    ```bash
    cloudflared tunnel --url http://localhost:8000
    ```

    This generates a random public URL that securely tunnels traffic to your local `xend` instance.
    You can configure your own custom public domain to point to the instance on the Cloudflare dashboard

## References

-   [Cloudflare Tunnels](https://developers.cloudflare.com/pages/how-to/preview-with-cloudflare-tunnel/)
