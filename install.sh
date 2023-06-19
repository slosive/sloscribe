#!/bin/sh

set -e

if ! command -v tar >/dev/null; then
	echo "Error: tar is required to install slosive" 1>&2
	exit 1
fi

if ! command -v curl >/dev/null; then
    echo "Error: curl is required to install slosive" 1>&2
    exit 1
fi


case $(uname -sm) in
"Darwin x86_64") target="slosive-darwin-amd64" ;;
"Darwin arm64")  target="slosive-darwin-arm64" ;;
  "Linux x86_64")  target="slosive-linux-amd64" ;;
"Linux aarch64") target="slosive-linux-arm64" ;;
  *)
      echo "Error: Unsupported operating system or architecture: $(uname -sm)" 1>&2
      exit 1 ;;
esac
  target_file="xslosive"


slosive_uri="https://github.com/tfadeyi/slosive/releases/latest/download/${target}.tar.gz"

slosive_install="${SLOSIVE_INSTALL:-$HOME/.slosive}"
bin_dir="$slosive_install/bin"
bin="$bin_dir/$target_file"

if current_install=$( command -v xslosive ) && [ ! -x "$bin" ] && [ "$current_install" != "$bin" ]; then
    echo "failed to install slosive to \"$bin\"" >&2
    echo "slosive is already installed in another location: \"$current_install\"" >&2
    exit 1
fi

if [ ! -d "$bin_dir" ]; then
	mkdir -p "$bin_dir"
fi

curl --fail --location --progress-bar --output "$bin.tar.gz" "$slosive_uri"
tar xfO "$bin.tar.gz" "$target/$target_file" > "$bin"
chmod +x "$bin"
rm "$bin.tar.gz"

echo "SLOsive's xslosive CLI was installed successfully to $bin"
if command -v xslosive >/dev/null; then
	echo "Run 'xslosive --help' to get started"
else
	if [ "$SHELL" = "/bin/zsh" ] || [ "$ZSH_NAME" = "zsh" ]; then
        shell_profile=".zshrc"
    else
        shell_profile=".bashrc"
	fi
    echo
	echo "Manually add the directory to your \$HOME/$shell_profile (or similar)"
	echo "  export SLOSIVE_INSTALL=\"$slosive_install\""
	echo "  export PATH=\"\$SLOSIVE_INSTALL/bin:\$PATH\""
    echo
	echo "And run \"source $HOME/$shell_profile\" to update your current shell"
    echo
fi
echo
echo "Checkout https://slotalk.fadey.io for more information"
