# devbox

CLI para executar o Claude Code dentro de um container Docker isolado.

## Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/)

## Instalação

```bash
task build
```

Isso compila o binário Go e move para `/usr/bin/devbox`.

## Uso

```bash
devbox [flags]
```

### Flags

| Flag | Atalho | Padrão | Descrição |
|------|--------|--------|-----------|
| `--image` | `-i` | `firmotecnologia/devbox:latest` | Imagem Docker a ser utilizada |
| `--no-pull` | | `false` | Pula o `docker pull` antes de executar |
| `--shell` | | `false` | Inicia um shell bash em vez do Claude Code |
| `--dotbins-config` | `-d` | `~/.dotbins/dotbins.yaml` | Caminho para a configuração do dotbins |

### Exemplos

```bash
# Executar o Claude Code no diretório atual
devbox

# Abrir um shell bash no container
devbox --shell

# Usar uma imagem customizada sem fazer pull
devbox --image minha-imagem:tag --no-pull
```

## Como funciona

O `devbox` monta o diretório atual como `/workspace` dentro do container, preservando as configurações do Claude entre sessões via volumes:

- Diretório atual → `/workspace`
- `~/.claude` → `/home/claude/.claude`
- `~/.claude.json` → `/home/claude/.claude.json`
- Cache do dotbins → `/home/claude/.devbox/dotbins` (se configurado)

A variável de ambiente `GITHUB_TOKEN` é automaticamente repassada do host para o container, quando disponível.

## Imagem Docker

Para construir e publicar a imagem:

```bash
# Apenas build
task docker:build

# Build + push
task docker:push

# Fluxo completo de desenvolvimento (push da imagem + build do binário)
task dev
```

## Licença

MIT — veja [LICENSE](./LICENSE)
