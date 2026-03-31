Application that helps to create a folder in specific path. The folder could be created:
- Empty, with only the name chosen
- From a github repo (cloning it in the folder)
The folder could be accessed writing `tryl <folder-name>` if already exists, otherwise the system should recommend to create it.

The interface is a `TUI` (Terminal User Interface), that can allow the user to:
- Find a folder in the specific path
- Create a new folder in the specific path
- Remove folder interactively (with a confirmation)
- Update the folder name interactively
- Open the config file with VIM as default (or others if configured in config file) and update it.
## Language & Frameworks
For this `TUI` application I choose `go` as language and [Bubble Tea](https://github.com/charmbracelet/bubbletea) with [Lip Gloss](https://github.com/charmbracelet/lipgloss) and [Bubbles](https://github.com/charmbracelet/bubbles) 

## Functionality
### Open the TUI
The user can exec in terminal
```bash
tryl

# (optional) open with folder name
tryl folder_name
```
and the application should open with:
- A search bar on top (with filter filled if opened with folder_name)
- A list of all folders already in specified path
- A `bottom bar` with commands for actions like:
	- Select for delete folder under the cursor (Ctrl + D)
	- Rename folder under the cursor ( Ctrl + R )
	- Move the folder under the cursor (Ctrl + M)
	- Close the application (Ctrl + C)
### Select a folder
The user can move through the list of folder with arrow keys, up and down. The cursor should be placed on the element where the arrow key lands.
### Delete a folder
The user can press `Ctrl + D` to mark a folder in state to delete. The user can mark more folders in `delete` state. When there are folders in this state the `bottom bar` actions changes:
- Confirm action (Enter)
- Cancel action (Esc)
### Rename folder
The user can press `Ctrl + R` to rename the folder under the cursor. The `TUI` should open an interactivity dialog (or new page), that permits to update the folder name. This dialog/page should have bottom bar with action Confirm and Cancel
### Move folder
The user can press `Ctrl + M` to move the folder in another location in the scope of current user `$HOME`. The user should write a destination path and if destination not exists should be created.
### Close the application
The user can press `Ctrl + C` and the application should shut down.
### Searching for folder
The user can write any letter (not key bindings that trigger other actions), and the TUI should focus on search bar input. The folder list should update list of folders based on value in the search input.

## Key Actions

| Command    | Action                                                                                 |
| ---------- | -------------------------------------------------------------------------------------- |
| `Ctrl + D` | Mark the folder under the cursor to be deleted                                         |
| `Ctrl + R` | Allow the user to rename the folder under the cursor                                   |
| `Ctrl + M` | Allow the user to move the folder under the cursor in a location in the scope of $HOME |
| `Ctrl + C` | Close the application                                                                  |
| `Enter`    | Confirm the actions below                                                              |
| `Esc`      | Cancel the actions below                                                               |
## Model
This application only works with one model `folder`, that should be described in go struct like this:
```go
type model struct {
	choices  []string         // items on the folder list
	cursor   int              // which folder item our cursor is pointing at
	selected map[int]struct{} // which folder item are selected
}
```

