package cli

import (
	"bufio"
	"bytes"
	"os"

	"github.com/spf13/cobra"
)

func completionCmd() *cobra.Command {
	completionCmd := &cobra.Command{
		Use:   "completion",
		Short: "Output shell completion code for the given shell (bash or zsh)",
		Long: `
Output shell completion code for bash or zsh
This command prints shell code which must be evaluated to provide interactive
completion of skycoin-cli commands.
Bash
	$ source <(skycoin-cli completion bash)
will load the skycoin-cli completion code for bash. Note that this depends on the
bash-completion framework. It must be sourced before sourcing the skycoin-cli
completion, e.g. on macOS:
	$ brew install bash-completion
	$ source $(brew --prefix)/etc/bash_completion
	$ source <(skycoin-cli completion bash)
	(or, if you want to preserve completion within new terminal sessions)
	$ echo 'source <(skycoin-cli completion bash)' >> ~/.bashrc
Zsh
	$ source <(skycoin-cli completion zsh)
	(or, if you want to preserve completion within new terminal sessions)
	$ echo 'source <(skycoin-cli completion zsh)' >> ~/.zshrc`,
	}

	completionCmd.AddCommand(bashAutocompleteCmd)
	completionCmd.AddCommand(zshAutocompleteCmd)

	return completionCmd
}

var bashAutocompleteCmd = &cobra.Command{
	Use:   "bash",
	Short: "Output shell completion code for bash",
	Long: `
Output shell completion code for bash.
This command prints shell code which must be evaluated to provide interactive
completion of skycoin-cli commands.
	$ source <(skycoin-cli completion bash)
will load the skycoin-cli completion code for bash. Note that this depends on the
bash-completion framework. It must be sourced before sourcing the skycoin-cli
completion, e.g. on macOS:
	$ brew install bash-completion
	$ source $(brew --prefix)/etc/bash_completion
	$ source <(skycoin-cli completion bash)
	(or, if you want to preserve completion within new terminal sessions)
	$ echo 'source <(skycoin-cli completion bash)' >> ~/.bashrc`,
	RunE: runCompletionBash,
}

var zshAutocompleteCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Output shell completion code for zsh",
	Long: `
Output shell completion code for zsh.
This command prints shell code which must be evaluated to provide interactive
completion of skycoin-cli commands.
	$ source <(skycoin-cli completion zsh)
	(or, if you want to preserve completion within new terminal sessions)
	$ echo 'source <(skycoin-cli completion zsh)' >> ~/.zshrc
zsh completions are only supported in versions of zsh >= 5.2`,
	RunE: runCompletionZsh,
}

func runCompletionBash(_ *cobra.Command, _ []string) error {
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	return skyCLI.GenBashCompletion(out)
}

// Copyright 2016 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
func runCompletionZsh(_ *cobra.Command, _ []string) error {
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	zshInitialization := `
__skycoin-cli_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand
	source "$@"
}
__skycoin-cli_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift
		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__skycoin-cli_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}
__skycoin-cli_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?
	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}
__skycoin-cli_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}
__skycoin-cli_declare() {
	if [ "$1" == "-F" ]; then
		whence -w "$@"
	else
		builtin declare "$@"
	fi
}
__skycoin-cli_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}
__skycoin-cli_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}
__skycoin-cli_filedir() {
	local RET OLD_IFS w qw
	__debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi
	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"
	IFS="," __debug "RET=${RET[@]} len=${#RET[@]}"
	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__skycoin-cli_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}
__skycoin-cli_quote() {
    if [[ $1 == \'* || $1 == \"* ]]; then
        # Leave out first character
        printf %q "${1:1}"
    else
    	printf %q "$1"
    fi
}
autoload -U +X bashcompinit && bashcompinit
# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi
__skycoin-cli_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__skycoin-cli_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__skycoin-cli_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__skycoin-cli_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__skycoin-cli_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__skycoin-cli_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/__skycoin-cli_declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__skycoin-cli_type/g" \
	<<'BASH_COMPLETION_EOF'
`
	out.Write([]byte(zshInitialization))

	buf := new(bytes.Buffer)
	skyCLI.GenBashCompletion(buf)
	out.Write(buf.Bytes())

	zshTail := `
BASH_COMPLETION_EOF
}
__skycoin-cli_bash_source <(__skycoin-cli_convert_bash_to_zsh)
`
	out.Write([]byte(zshTail))
	return nil
}
