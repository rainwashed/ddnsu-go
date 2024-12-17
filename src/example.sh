#!/bin/bash
if [ ! -f ~/.config/ddnsu/ddnsu ]; then
    arch=$(uname -i)
    a=""
    if [[ $arch == x86_64* ]]; then
        echo "detected x64 architecture"
        a="x64"
    elif [[ $arch == i*86 ]]; then
       echo "detected 32/i386 architecture"
       a="x32"
    elif  [[ $arch == arm* ]]; then
        echo "detected arm64 architecture"
        a="arm64"
    else
        echo "idk what architecture you have bro..."
        exit 1
    fi

    echo "could not find ddnsu binary. it is being downloaded from the latest github release."
    curl -s https://api.github.com/repos/rainwashed/ddnsu-go/releases/latest \
    | grep "ddnsu.$(a)" \
    | cut -d : -f 2,3 \
    | tr -d \" \
    | wget -O ~/.config/ddnsu/ddnsu --progress=bar:force -i -

    chmod +x ~/.config/ddnsu/ddnsu
fi

exec ~/.config/ddnsu/ddnsu start