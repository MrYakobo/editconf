# Editconf
A "script" for editing config files. It does all of this:

- Look in `~/.config/editconf.yaml` for configs/options/dmenu-options
- Spawn dmenu with the defined entries and options
- Upon selection, spawn editor in `editconf.yaml` (with `$EDITOR` as fallback) with that file
- If the directory that the edited file resides in has a Makefile (like the suckless utilities), run `make` in that directory.

Depends on dmenu. See example config in `editconf.def.yaml`.