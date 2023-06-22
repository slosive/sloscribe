#!/bin/sh

set -e

if ! command -v tar >/dev/null; then
	echo "Error: tar is required to install sloscribe" 1>&2
	exit 1
fi

if ! command -v curl >/dev/null; then
    echo "Error: curl is required to install sloscribe" 1>&2
    exit 1
fi


case $(uname -sm) in
"Darwin x86_64") target="sloscribe-darwin-amd64" ;;
"Darwin arm64")  target="sloscribe-darwin-arm64" ;;
  "Linux x86_64")  target="sloscribe-linux-amd64" ;;
"Linux aarch64") target="sloscribe-linux-arm64" ;;
  *)
      echo "Error: Unsupported operating system or architecture: $(uname -sm)" 1>&2
      exit 1 ;;
esac
  target_file="sloscribe"


sloscribe_uri="https://github.com/slosive/sloscribe/releases/latest/download/${target}.tar.gz"

sloscribe_install="${SLOSCRIBE_INSTALL:-$HOME/.sloscribe}"
bin_dir="$sloscribe_install/bin"
bin="$bin_dir/$target_file"

if current_install=$( command -v sloscribe ) && [ ! -x "$bin" ] && [ "$current_install" != "$bin" ]; then
    echo "failed to install sloscribe to \"$bin\"" >&2
    echo "sloscribe is already installed in another location: \"$current_install\"" >&2
    exit 1
fi

if [ ! -d "$bin_dir" ]; then
	mkdir -p "$bin_dir"
fi

curl --fail --location --progress-bar --output "$bin.tar.gz" "$sloscribe_uri"
tar xfO "$bin.tar.gz" "$target/$target_file" > "$bin"
chmod +x "$bin"
rm "$bin.tar.gz"

echo "SLOsive's sloscribe CLI was installed successfully to $bin"
if command -v sloscribe >/dev/null; then
	echo "Run 'sloscribe --help' to get started"
else
	if [ "$SHELL" = "/bin/zsh" ] || [ "$ZSH_NAME" = "zsh" ]; then
        shell_profile=".zshrc"
    else
        shell_profile=".bashrc"
	fi
    echo
	echo "Manually add the directory to your \$HOME/$shell_profile (or similar)"
	echo "  export SLOSCRIBE_INSTALL=\"$sloscribe_install\""
	echo "  export PATH=\"\$SLOSCRIBE_INSTALL/bin:\$PATH\""
    echo
	echo "And run \"source $HOME/$shell_profile\" to update your current shell"
    echo
fi
echo
echo "Checkout https://slotalk.fadey.io for more information"
