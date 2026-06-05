default:
    air

gen:
    go run ./cmd/internal gen

clean:
    rm -rf work
    rm -rf tmp
    rm -f result

test-build:
    nix build --no-link .#backend
    nix build --no-link .#web
