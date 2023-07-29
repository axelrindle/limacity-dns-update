[![CI](https://github.com/axelrindle/limacity-dns-update/actions/workflows/ci.yml/badge.svg)](https://github.com/axelrindle/limacity-dns-update/actions/workflows/ci.yml)

# limacity-dns-update

🤖 Updates DNS entries on [Lima-City](https://www.lima-city.de/).

## Usage

### Docker

The Docker container is the preferred way to use this:

```
docker run -d \
    --name limacity-dns-update \
    --env-file .env \
    --network host \
    ghcr.io/axelrindle/limacity-dns-update:latest
```

The `--network host` flag is required for IPv6 support.

## Development

## License

[MIT](LICENSE)
