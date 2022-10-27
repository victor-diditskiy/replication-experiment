#!/bin/sh

sudo apt install -y zsh

mkdir .oh-my-zsh
tar -xf oh.tar -C .

sudo -k chsh -s "/usr/bin/zsh" "user"
