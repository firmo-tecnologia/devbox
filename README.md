# devbox

CLI para executar o Claude Code dentro de um container Docker isolado.

https://github.com/user-attachments/assets/68875a8c-aa6a-42db-9768-d11742a5a7b9

## Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/)

## Instalação

Baixe o binário para sua plataforma na [página de releases](https://github.com/firmotecnologia/devbox/releases) e mova para um diretório no seu `PATH`:

```bash
# Exemplo para Linux amd64
curl -L https://github.com/firmotecnologia/devbox/releases/latest/download/devbox-linux-amd64 -o devbox
chmod +x devbox
sudo mv devbox /usr/local/bin/
```

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
- Cache do dotbins → `~/.devbox/dotbins` (se configurado)

A variável de ambiente `GITHUB_TOKEN` é automaticamente repassada do host para o container, quando disponível.

## Dotbins

O devbox usa o [dotbins](https://github.com/nikitabobko/dotbins) para disponibilizar ferramentas adicionais dentro do container (ex: `gh`, `jq`, `fzf`), sem precisar instalá-las na imagem base.

Se o arquivo `~/.dotbins/dotbins.yaml` existir no host (ou o caminho definido por `--dotbins-config`), ele é montado no container junto com o cache de binários em `~/.devbox/dotbins/`. Na primeira execução o dotbins baixa os binários; nas seguintes o cache é reaproveitado.

Exemplo de configuração (`~/.dotbins/dotbins.yaml`):

```yaml
tools:
  - repo: cli/cli
    binary: gh
  - repo: jqlang/jq
    binary: jq
```

Sem esse arquivo o container sobe normalmente, sem ferramentas extras.

## Licença

MIT — veja [LICENSE](./LICENSE)
