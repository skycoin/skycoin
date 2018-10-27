#!/bin/bash
set -e
set -o pipefail

# Install vscode extensions declared on EXTENSIONS
for ext in $VS_EXTENSIONS; do code --user-data-dir $HOME --install-extension $ext; done

#su user -p -c /usr/share/code/code

exec "$@"
