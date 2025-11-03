# Twin Pick

Twin Pick is a tool designed to help users find a film based on their Letterboxd watchlists and specified criteria such as genres, platforms, and duration.

## Usage

Twin Pick provides :
- An API REST
- A Command Line Interface (CLI)
- A MCP Server

`Twin Pick` uses [Taskfile](https://taskfile.dev) and invites you to use too.

### API

First start the server :

```bash
task run:api
```

Then you can make requests :

```bash
curl -X GET "http://localhost:8080/api/v1/pick?usernames=abroudoux,potatoze&genres=action?limit=1"
```

> `usernames` param is mandatory.

### CLI

You can use Twin Pick directly from your terminal.

```bash
task run:cli -- --usernames abroudoux,potatoze --genres action --limit 1
```

## Todo

- [ ] Can't handle genres & platform at the same time
- [ ] Return film details
- [ ] Duration filter
- [ ] Spot mode (wip...)
- [ ] TUI mode

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
