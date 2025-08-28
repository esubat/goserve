# goserve
A simple, lightweight server with reverse proxy and static file serving.


## Configuration

Create a `config.yaml` file in the project root to customize the server:

```yaml
port: 8080
target: "http://localhost:3000"
static_path: "./public"
```

If you don’t provide a config file, the server uses these defaults:

- `port`: 8080
- `target`: `http://localhost:3000`
- `static_path`: `./public`

## Running the Server

**Proxy mode** (forward requests to a target server):

```bash
go run . -serve=proxy -config=config.yaml
```

**Static server mode** (serve files from a directory):

```bash
go run . -serve=static -config=config.yaml
```

### Notes

- Omitting `-serve` defaults to static server mode.
- Omitting `-config` defaults to `config.yaml` or uses the hardcoded values above.
