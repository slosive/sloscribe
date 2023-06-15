#!/bin/sh

set -e

if ! command -v tar >/dev/null; then
	echo "Error: tar is required to install slotalk" 1>&2
	exit 1
fi

if ! command -v curl >/dev/null; then
    echo "Error: curl is required to install slotalk" 1>&2
    exit 1
fi


case $(uname -sm) in
"Darwin x86_64") target="slotalk-darwin-amd64" ;;
"Darwin arm64")  target="slotalk-darwin-arm64" ;;
  "Linux x86_64")  target="slotalk-linux-amd64" ;;
"Linux aarch64") target="slotalk-linux-arm64" ;;
  *)
      echo "Error: Unsupported operating system or architecture: $(uname -sm)" 1>&2
      exit 1 ;;
esac
  target_file="slotalk"


slotalk_uri="https://github.com/tfadeyi/slotalk/releases/latest/download/${target}.tar.gz"

slotalk_install="${SLOTALK_INSTALL:-$HOME/.slotalk}"
bin_dir="$slotalk_install/bin"
bin="$bin_dir/$target_file"

if current_install=$( command -v slotalk ) && [ ! -x "$bin" ] && [ "$current_install" != "$bin" ]; then
    echo "failed to install slotalk to \"$bin\"" >&2
    echo "slotalk is already installed in another location: \"$current_install\"" >&2
    exit 1
fi

if [ ! -d "$bin_dir" ]; then
	mkdir -p "$bin_dir"
fi

curl --fail --location --progress-bar --output "$bin.tar.gz" "$slotalk_uri"
tar xfO "$bin.tar.gz" "$target/$target_file" > "$bin"
chmod +x "$bin"
rm "$bin.tar.gz"

echo "slotalk was installed successfully to $bin"
if command -v slotalk >/dev/null; then
	echo "Run 'slotalk --help' to get started"
else
	if [ "$SHELL" = "/bin/zsh" ] || [ "$ZSH_NAME" = "zsh" ]; then
        shell_profile=".zshrc"
    else
        shell_profile=".bashrc"
	fi
    echo
	echo "Manually add the directory to your \$HOME/$shell_profile (or similar)"
	echo "  export SLOTALK_INSTALL=\"$slotalk_install\""
	echo "  export PATH=\"\$SLOTALK_INSTALL/bin:\$PATH\""
    echo
	echo "And run \"source $HOME/$shell_profile\" to update your current shell"
    echo
fi
echo
echo "Checkout https://slotalk.fadey.io for more information"
