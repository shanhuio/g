// Copyright (C) 2022  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package bosinit

const bashrc = `
PS1='\[\033[01;32m\]\h:\[\033[01;34m\]\w\[\033[00m\]\$ '

alias ls='ls --color=auto'
alias cdp='cd -P'

# enable programmable completion features.
if ! shopt -oq posix; then
  if [[ -f /usr/share/bash-completion/bash_completion ]]; then
    source /usr/share/bash-completion/bash_completion
  fi
fi
`

// BashRC provides a the .bashrc file.
var BashRC = &WriteFile{
	Path:        "/home/rancher/.bashrc",
	Permissions: FilePerm(0644),
	Owner:       "rancher",
	Content:     bashrc,
}
