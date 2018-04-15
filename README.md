# Editconf
A "script" for editing config files. It does this:

- Look in `~/.config/editconf.yaml` for configs/options/dmenu-options
- Spawn dmenu with the defined entries and options
- Upon selection, spawn editor in `editconf.yaml` (with `$EDITOR` as fallback) with that file
- If the file ends with `.h` (like the suckless utilities), run `make` in the folder that the file is in

Depends on dmenu. See example config in `editconf.def.yaml`.