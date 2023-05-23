# gomod2nix basic usage
# After you have entered your development shell you can generate a gomod2nix.toml using:
#   gomod2nix generate
# To speed up development and avoid downloading dependencies again in the Nix store you can import them directly from the Go cache using:
#   gomod2nix import
alias develop='nix develop'
