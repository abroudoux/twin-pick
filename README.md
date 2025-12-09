# Twin Pick

Twin Pick est un outil pour trouver le film parfait basé sur les watchlists Letterboxd. Fini les heures passees a se demander quoi regarder entre amis !

## Concept

Twin Pick propose deux modes :

### Pick Mode

Compare les watchlists de plusieurs utilisateurs Letterboxd et trouve les films en commun. Ideal pour une soiree cine entre amis.

**Logique de selection :**
- Films presents dans **toutes** les watchlists
- Films presents dans **au moins la moitie** des watchlists

### Spot Mode

Suggere des films populaires sur Letterboxd, filtres par genre. Parfait pour decouvrir de nouveaux films.

## Installation

### Prerequis

- Go 1.24+
- [Taskfile](https://taskfile.dev) (optionnel mais recommande)

### Build

```bash
# Avec Taskfile
task build:all

# Sans Taskfile
go build -o bin/twinpick-api cmd/twinpick-api/main.go
go build -o bin/twinpick-cli cmd/twinpick-cli/main.go
go build -o bin/twinpick-mcp cmd/twinpick-mcp/main.go
```

## Utilisation

Twin Pick propose 3 interfaces : API REST, CLI et serveur MCP.

### API REST

```bash
task run:api
# ou
./bin/twinpick-api
```

Le serveur demarre sur `http://localhost:8080`.

#### Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /api/v1/pick` | Trouve les films communs entre watchlists |
| `GET /api/v1/spot` | Suggere des films populaires |

#### Parametres

| Parametre | Type | Description | Requis |
|-----------|------|-------------|--------|
| `usernames` | string | Noms d'utilisateurs Letterboxd separes par virgule | **Oui** (pick) |
| `genres` | string | Genres separes par virgule (ex: `action,thriller`) | Non |
| `platform` | string | Plateforme de streaming (ex: `netflix-fr`) | Non |
| `duration` | string | Filtre de duree : `short`, `medium`, `long` | Non |
| `limit` | int | Nombre maximum de films retournes | Non |

**Durees :**
- `short` : <= 100 minutes
- `medium` : <= 120 minutes
- `long` : pas de limite (defaut)

#### Exemples

```bash
# Trouver les films communs entre deux utilisateurs
curl "http://localhost:8080/api/v1/pick?usernames=alice,bob"

# Filtrer par genre et duree
curl "http://localhost:8080/api/v1/pick?usernames=alice,bob&genres=horror&duration=short"

# Limiter le nombre de resultats
curl "http://localhost:8080/api/v1/pick?usernames=alice,bob,charlie&limit=5"

# Films populaires (spot mode)
curl "http://localhost:8080/api/v1/spot"

# Films populaires d'un genre specifique
curl "http://localhost:8080/api/v1/spot?genres=animation"
```

#### Reponse

```json
{
  "films": [
    {
      "title": "Parasite",
      "duration": "132 min",
      "directors": "Bong Joon-ho",
      "year": "2019"
    }
  ]
}
```

### CLI

```bash
task run:cli -- [options]
# ou
./bin/twinpick-cli [command] [options]
```

#### Commandes

```bash
# Pick mode
./bin/twinpick-cli pick --usernames alice,bob --genres action --limit 5

# Spot mode
./bin/twinpick-cli spot --genres horror --duration medium
```

#### Options

| Option | Description |
|--------|-------------|
| `--usernames` | Noms d'utilisateurs (requis pour pick) |
| `--genres` | Genres a filtrer |
| `--platform` | Plateforme de streaming |
| `--duration` | Filtre de duree (short/medium/long) |
| `--limit` | Limite de resultats |

### Serveur MCP

Twin Pick peut fonctionner comme serveur [Model Context Protocol](https://modelcontextprotocol.io/) pour integration avec des LLMs.

```bash
task run:mcp
# ou
./bin/twinpick-mcp
```

**Tools disponibles :** `pick`, `spot`

## Architecture

```
twin-pick/
├── cmd/
│   ├── twinpick-api/      # Point d'entree API
│   ├── twinpick-cli/      # Point d'entree CLI
│   └── twinpick-mcp/      # Point d'entree MCP
├── internal/
│   ├── adapters/
│   │   ├── http/          # Handlers HTTP (Gin)
│   │   ├── cli/           # Interface CLI (Cobra)
│   │   └── mcp/           # Serveur MCP
│   ├── application/       # Services metier
│   ├── domain/            # Modeles et logique metier
│   └── infrastructure/
│       ├── scrapper/      # Scraping Letterboxd (Colly)
│       ├── client/        # Client HTTP pour details films
│       └── cache/         # Cache en memoire avec TTL
```

## Performance

Twin Pick integre plusieurs optimisations :

### Cache

- **Details des films** : cache 24h (evite de refetch les memes films)
- **Films populaires** : cache 24h
- **Watchlists** : cache 10 minutes

### Concurrence

- Scraping des pages de watchlist en parallele (max 15 concurrent)
- Fetch des details films en parallele (max 25 concurrent)
- Connection pooling HTTP

### Optimisation des requetes

Le filtre de duree est applique **avant** de limiter les resultats, ce qui permet de recuperer les details uniquement sur les films filtres.

## Subtilites

### Genres multiples

Les genres sont combines avec un `+` :
```bash
curl "http://localhost:8080/api/v1/pick?usernames=alice,bob&genres=action,thriller"
# => Scrape: /alice/watchlist/genre/action+thriller
```

### Platforms

Les identifiants de plateforme suivent le format Letterboxd :
- `netflix-fr`, `netflix-us`
- `amazon-prime-video-fr`
- `disney-plus-fr`
- etc.

### Films en commun

Un film est considere "en commun" s'il apparait dans :
- **Toutes** les watchlists, OU
- **Au moins 50%** des watchlists

Cela permet d'avoir des resultats meme si un utilisateur n'a pas exactement les memes films.

### Details des films

Les details (duree, annee, realisateurs) sont recuperes via l'API JSON de Letterboxd (`/film/{slug}/json/`). Si un detail n'est pas trouve, la valeur `"not found"` est retournee.

## Developpement

```bash
# Lancer les tests
task test

# Tests avec couverture
task test:coverage

# Tests d'integration API
task test:api
```

## Stack technique

- **Go 1.24**
- **Gin** - Framework HTTP
- **Colly** - Web scraping
- **Cobra** - CLI framework
- **Taskfile** - Task runner

## License

MIT License - voir [LICENSE](LICENSE)
