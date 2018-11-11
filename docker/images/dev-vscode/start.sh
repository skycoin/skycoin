#!/bin/bash
set -e
set -o pipefail

# Install vscode extensions declared on VS_EXTENSIONS
if [[ -n "$VS_EXTENSIONS" ]]; then
    for ext in $VS_EXTENSIONS; do code --user-data-dir $HOME --install-extension $ext; done
fi

# Check if skycoindev-vscode:dind image has been started
if [[ -n "$DIND_COMMIT" ]]; then
    # no arguments passed
    # or first arg is `-f` or `--some-option`
    if [[ "$#" -eq 0 -o "${1#-}" != "$1" ]]; then
        # add our default arguments
        set -- dockerd \
            --host=unix:///var/run/docker.sock \
            --host=tcp://0.0.0.0:2375 \
            "$@"
    fi

    if [[ "$1" = 'dockerd' ]]; then
        # if we're running Docker, let's pipe through dind
        # (and we'll run dind explicitly with "sh" since its shebang is /bin/bash)
        set -- sh "$(which dind)" "$@"
    fi
fi



# If user no pass a command when run docker image VS Code will run
# else, run user command
if [[ -n "$@" ]]; then
    su user -p -c /usr/share/code/code
else
    exec "$@"
fi
