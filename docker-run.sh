docker run -it --rm \
        --name gochess-dev \
        -v $HOME/gochess:/app \
        gochess \
        go run .
