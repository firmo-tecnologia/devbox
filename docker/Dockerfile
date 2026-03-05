FROM ubuntu:24.04
ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y curl git ca-certificates python3-pip xclip wl-clipboard
RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash -
RUN apt-get install -y nodejs && rm -rf /var/lib/apt/lists/*

RUN npm install -g @anthropic-ai/claude-code

RUN pip3 install dotbins --break-system-packages

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

RUN usermod -l claude -md /home/claude ubuntu && groupmod -n claude ubuntu
USER claude
WORKDIR /workspace

ENTRYPOINT ["entrypoint.sh"]
CMD ["claude", "--dangerously-skip-permissions"]
