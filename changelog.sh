#!/usr/bin/bash
## CHANGELOG GENERATOR SCRIPT
# supply range of pull requests since last release as arguments for sequence
[[ $1 == "" ]] && cat $0 && exit
for _i in $(seq $1 $2 | tac) ; do
_merged="$(curl -s https://github.com/skycoin/skycoin/pull/${_i} | grep 'Status: Merged')"
if [[ $_merged != "" ]] ; then
_title="$(curl -s https://github.com/skycoin/skycoin/pull/${_i} | grep '<title>')"
_title="$(curl -s https://github.com/skycoin/skycoin/pull/${_i} | grep '<title>')"
_title=${_title//"<title>"/}
_title=${_title//"by"*/}
[[ ${_title} != "" ]] && echo "- ${_title} [#${_i}](https://github.com/skycoin/skycoin/pull/${_i})"
fi
done
