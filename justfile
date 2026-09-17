
just:
    just -l

run:
    CGO_ENABLED=0 GO111MODULE=on go run . --log-level=1

build_release_linux version="latest":
    GOOS=linux CGO_ENABLED=0 GO111MODULE=on go build -ldflags "-s -w -X 'main.buildVersion={{version}}' -X 'main.buildHash=$(git rev-parse HEAD)' -X 'main.buildDate=$(date -u)'" -trimpath -o dist/kanata-tray-linux

# e.g. "push_tag v0.1.0"
push_tag tag:
    git tag {{tag}}
    git push --tags