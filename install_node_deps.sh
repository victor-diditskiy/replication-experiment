#!/bin/sh

# Install ZSH
echo "PROMPT='%F{green}% $USER @ %{$fg[red]%}%M %{$fg_bold[red]%} ➜ %{$fg_bold[green]%}%p %{$fg[cyan]%}%c %{$fg_bold[blue]%}$(git_prompt_info)%{$fg_bold[blue]%} % %{$reset_color%}'" > ~/.oh-my-zsh/themes/robbyrussell.zsh-theme
echo "ZSH_THEME_GIT_PROMPT_PREFIX=\"%{$fg_bold[blue]%}git:(%{$fg[red]%}\"" >> ~/.oh-my-zsh/themes/robbyrussell.zsh-theme
echo "ZSH_THEME_GIT_PROMPT_SUFFIX=\"%{$reset_color%} \"" >> ~/.oh-my-zsh/themes/robbyrussell.zsh-theme
echo "ZSH_THEME_GIT_PROMPT_DIRTY=\"%{$fg[blue]%}) %{$fg[yellow]%}\"" >> ~/.oh-my-zsh/themes/robbyrussell.zsh-theme
echo "ZSH_THEME_GIT_PROMPT_CLEAN=\"%{$fg[blue]%})\"" >> ~/.oh-my-zsh/themes/robbyrussell.zsh-theme

# Install PG

sudo apt install -y postgresql

# Install Node exporter

wget https://github.com/prometheus/node_exporter/releases/download/v1.3.1/node_exporter-1.3.1.linux-amd64.tar.gz
tar xvf node_exporter-1.3.1.linux-amd64.tar.gz
cd node_exporter-1.3.1.linux-amd64/
sudo cp node_exporter /usr/local/bin
cd ..
rm -rf node_exporter-1.3.1.linux-amd64
rm -rf node_exporter-1.3.1.linux-amd64.tar.gz

sudo useradd --no-create-home --shell /bin/false node_exporter
sudo chown node_exporter:node_exporter /usr/local/bin/node_exporter
sudo vim /etc/systemd/system/node_exporter.service

SRV=$(<<-EndOfMessage
[Unit]
Description=Node Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=node_exporter
Group=node_exporter
Type=simple
ExecStart=/usr/local/bin/node_exporter

[Install]
WantedBy=multi-user.target
EndOfMessage
) && echo $SRV > sudo /etc/systemd/system/node_exporter.service

DATA_SOURCE_NAME="postgresql://user1:Q4QLpgywgXGtT6@84.252.142.120:5432/replication_experiment?sslmode=disable" ./postgres_exporter &