[![CI](https://github.com/axelrindle/limacity-dns-update/actions/workflows/ci.yml/badge.svg)](https://github.com/axelrindle/limacity-dns-update/actions/workflows/ci.yml)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/a9f4c2dcd7e047bd9b7f2f3e877dd210)](https://app.codacy.com/gh/axelrindle/limacity-dns-update/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/a9f4c2dcd7e047bd9b7f2f3e877dd210)](https://app.codacy.com/gh/axelrindle/limacity-dns-update/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

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
