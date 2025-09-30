# 🦥 Sloth Runner DSL - Icon Configuration

## File Icons for .sloth Extension

### VS Code Icon Configuration

Add to your VS Code `settings.json`:

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
    "*.sloth": "sloth"
  }
}
```

### Neovim/LunarVim Icon Configuration

For nvim-web-devicons:

```lua
require('nvim-web-devicons').setup {
  override = {
    sloth = {
      icon = "🦥",
      color = "#8B4513",
      cterm_color = "95",
      name = "Sloth"
    }
  }
}
```

### Terminal Icon Configuration

For file managers and terminals that support custom icons:

```bash
# Add to your shell config (~/.zshrc or ~/.bashrc)
alias ls='ls --color=auto'
export LS_COLORS="$LS_COLORS:*.sloth=38;5;95"

# For exa
alias exa='exa --icons'
```

### Icon Suggestions

#### Unicode Options:
- 🦥 (U+1F9A5) - Sloth emoji (primary choice)
- 🐨 (U+1F428) - Koala (alternative)
- 🌳 (U+1F333) - Tree (environment themed)
- ⚡ (U+26A1) - Lightning bolt (runner themed)

#### ASCII Art Options:
```
   🦥    Simple sloth emoji
  /o o\   ASCII sloth face
 (  -  )  
  \___/   

    ⚡     Lightning (speed)
   🌿🦥    Sloth in nature
```

### File Manager Icons

#### Finder (macOS)
Create a custom icon for .sloth files:
1. Find a .sloth file
2. Get Info (Cmd+I)
3. Drag sloth icon to the file icon in Get Info

#### Windows Explorer
Associate .sloth files with a custom icon through file type associations.

### Editor Integration

#### VS Code Extensions
Popular icon themes that can be extended:
- Material Icon Theme
- VSCode Icons
- Seti File Icons

#### JetBrains IDEs
Add custom file type with sloth icon in Settings > Editor > File Types

#### Sublime Text
Add to Packages/User/FileIcons.sublime-settings:
```json
{
  "file_types": {
    "sloth": {
      "icon": "🦥",
      "syntax": "Packages/Lua/Lua.sublime-syntax"
    }
  }
}
```

### Icon Design Specifications

For custom icon creation:

#### SVG Template:
```xml
<svg width="16" height="16" viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
  <!-- Sloth face -->
  <circle cx="8" cy="8" r="7" fill="#8B4513" stroke="#654321"/>
  <!-- Eyes -->
  <circle cx="6" cy="6" r="1" fill="#000"/>
  <circle cx="10" cy="6" r="1" fill="#000"/>
  <!-- Nose -->
  <circle cx="8" cy="8" r="0.5" fill="#000"/>
  <!-- Mouth -->
  <path d="M 6 10 Q 8 11 10 10" stroke="#000" stroke-width="0.5" fill="none"/>
  <!-- Tree branch -->
  <rect x="12" y="3" width="3" height="1" fill="#654321"/>
  <rect x="13" y="1" width="1" height="5" fill="#654321"/>
</svg>
```

#### Color Palette:
- Primary: #8B4513 (Saddle Brown)
- Secondary: #654321 (Dark Brown)
- Accent: #228B22 (Forest Green)
- Text: #000000 (Black)

### Implementation Examples

#### Neovim with nvim-tree:
```lua
require("nvim-tree").setup({
  renderer = {
    icons = {
      glyphs = {
        default = "",
        symlink = "",
        git = {
          unstaged = "✗",
          staged = "✓",
          unmerged = "",
          renamed = "➜",
          untracked = "★",
          deleted = "",
          ignored = "◌"
        },
        folder = {
          arrow_open = "",
          arrow_closed = "",
          default = "",
          open = "",
          empty = "",
          empty_open = "",
          symlink = "",
          symlink_open = "",
        },
        extension = {
          sloth = "🦥"
        }
      }
    }
  }
})
```

#### File Manager Integration:
```bash
# For ranger file manager
echo 'ext sloth = echo "🦥 Sloth Runner DSL"; cat "$1"' >> ~/.config/ranger/scope.sh

# For lf file manager
echo 'ext sloth
    mime text/x-lua
    !echo "🦥 Sloth Runner DSL" && cat "$f"' >> ~/.config/lf/previewer
```

## Usage Examples

### In File Explorers:
```
📁 workflows/
  🦥 ci-pipeline.sloth
  🦥 deployment.sloth
  🦥 testing.sloth
  📄 README.md
  📁 scripts/
```

### In Editor Tabs:
```
[🦥 main.sloth] [📄 config.yaml] [🐍 script.py]
```

### In Terminal Listings:
```bash
$ ls -la
🦥 build-pipeline.sloth
🦥 deploy-workflow.sloth
📁 scripts/
📄 README.md
```

This gives the `.sloth` extension a distinctive and memorable visual identity! 🦥⚡