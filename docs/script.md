# Script

## Meta

Script can contains meta for redefine some default data.

Meta should be placed only in top of script. Before any code line.

Meta format: `-- @metaName <params>`

- `--`  lua comment prefix
- `@metaName` name of meta, started by `@`
- `<params>` meta params, if required

### @interval

Redefine default interval for run the script

Format

```lua
-- @interval <interval value>
```

Example. Run interval will be set to 30 seconds

```lua
-- @interval 30s
```

### @name

Redefine default script name.
Default script name is filename (if use script source 'filesystem folder')

Format

```lua
-- @name <new name>
```

Example. Set script name 'New Name'

```lua
-- @name New Name
```

### @ignore

Not run this script

Format and example
```
-- @ignore
```
