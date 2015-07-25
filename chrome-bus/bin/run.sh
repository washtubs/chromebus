#!/usr/bin/env zsh
if [ -z $CHROME_BUS_FILE ]; then
    mkdir $HOME/.chrome-bus
    CHROME_BUS_FILE=$HOME/.chrome-bus/events
fi
./main.js > $CHROME_BUS_FILE
