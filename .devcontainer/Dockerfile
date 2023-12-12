FROM mcr.microsoft.com/devcontainers/go:1-1.21-bullseye

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends htop curl

# [Optional] Uncomment the next lines to use go get to install anything else you need
USER vscode
# RUN go get -x

USER root
RUN curl -s https://ohmyposh.dev/install.sh | bash -s

USER vscode
RUN echo 'eval "$(oh-my-posh init bash)"' >> ~/.bashrc






# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1
