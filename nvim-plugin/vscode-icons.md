# VS Code Icon Configuration for .sloth files

## File Icons Extension

Add this to your VS Code `settings.json`:

```json
{
  "vsicons.associations.files": [
    {
      "icon": "sloth",
      "extensions": ["sloth"],
      "format": "svg"
    }
  ],
  "material-icon-theme.files.associations": {
    "*.sloth": "../../nvim-plugin/icons/sloth-icon"
  },
  "workbench.iconTheme": "material-icon-theme"
}
```

## Material Icon Theme

For Material Icon Theme extension, create a custom association:

```json
{
  "material-icon-theme.files.associations": {
    "*.sloth": "sloth"
  },
  "material-icon-theme.folders.associations": {
    "workflows": "sloth",
    "sloth-runner": "sloth"
  }
}
```

## Custom Icon Package

Create your own icon theme package:

### package.json
```json
{
  "name": "sloth-runner-icons",
  "displayName": "Sloth Runner Icons",
  "description": "Icons for Sloth Runner DSL files",
  "version": "1.0.0",
  "engines": {
    "vscode": "^1.74.0"
  },
  "categories": ["Themes"],
  "contributes": {
    "iconThemes": [
      {
        "id": "sloth-runner",
        "label": "Sloth Runner Icons",
        "path": "./icons/sloth-icon-theme.json"
      }
    ]
  }
}
```

### sloth-icon-theme.json
```json
{
  "iconDefinitions": {
    "sloth": {
      "iconPath": "./sloth-icon.svg"
    }
  },
  "fileExtensions": {
    "sloth": "sloth"
  },
  "fileNames": {},
  "folderNames": {
    "workflows": "sloth",
    "sloth-runner": "sloth"
  }
}
```

This will make all .sloth files display with a cute sloth icon! 🦥