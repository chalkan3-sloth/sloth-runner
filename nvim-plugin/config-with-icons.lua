-- LunarVim configuration for Sloth Runner DSL with custom icon
-- Add this to your ~/.config/lvim/config.lua

vim.opt.mouse = ""
lvim.format_on_save.enabled = true

-- Configure nvim-web-devicons for .sloth files
require('nvim-web-devicons').setup {
  override = {
    sloth = {
      icon = "🦥",
      color = "#8B4513",
      cterm_color = "95", 
      name = "SlothDSL"
    }
  },
  -- Add file extension mapping
  override_by_extension = {
    ["sloth"] = {
      icon = "🦥",
      color = "#8B4513",
      cterm_color = "95",
      name = "SlothDSL"
    }
  }
}

-- Sloth Runner DSL Configuration (with icon support)
local ok, sloth = pcall(require, "sloth-runner")
if ok then
  local setup_ok, err = pcall(sloth.setup, {
    runner = {
      command = "sloth-runner",
      keymaps = {
        run_file = "<leader>sr",
        list_tasks = "<leader>sl", 
        dry_run = "<leader>st",
        debug = "<leader>sd",
      }
    },
    completion = {
      enable = true,
      snippets = true,
    },
    folding = {
      enable = true,
      auto_close = false,
    }
  })
  
  if not setup_ok then
    vim.notify("Sloth Runner setup failed: " .. tostring(err), vim.log.levels.WARN)
  end
else
  vim.notify("Sloth Runner plugin not found", vim.log.levels.WARN)
end

-- Enhanced file tree icons (if using nvim-tree)
if lvim.builtin.nvimtree then
  lvim.builtin.nvimtree.setup.renderer.icons.glyphs.extension = {
    sloth = "🦥"
  }
end

-- Simple key mappings with sloth context
lvim.keys.normal_mode["<leader>sr"] = function()
  local file = vim.api.nvim_buf_get_name(0)
  if file ~= "" then
    if file:match("%.sloth$") then
      vim.notify("🦥 Running Sloth workflow: " .. vim.fn.fnamemodify(file, ":t"), vim.log.levels.INFO)
    end
    vim.cmd("split | terminal sloth-runner run -f " .. vim.fn.shellescape(file))
  else
    vim.notify("No file to run", vim.log.levels.WARN)
  end
end

lvim.keys.normal_mode["<leader>sl"] = function()
  local file = vim.api.nvim_buf_get_name(0)
  if file ~= "" then
    if file:match("%.sloth$") then
      vim.notify("🦥 Listing tasks in: " .. vim.fn.fnamemodify(file, ":t"), vim.log.levels.INFO)
    end
    vim.cmd("split | terminal sloth-runner list -f " .. vim.fn.shellescape(file))
  else
    vim.notify("No file to analyze", vim.log.levels.WARN)
  end
end

-- Enhanced highlights with sloth branding
vim.api.nvim_create_autocmd("ColorScheme", {
  pattern = "*",
  callback = function()
    local highlights = {
      slothDSLKeyword = { fg = '#8B4513', bold = true },  -- Sloth brown
      slothModule = { fg = '#228B22', bold = true },      -- Forest green
      slothMethod = { fg = '#4682B4' },                   -- Steel blue
      slothFunction = { fg = '#DAA520' },                 -- Goldenrod
      slothEnvVar = { fg = '#DC143C', bold = true },      -- Crimson
      slothPath = { fg = '#20B2AA' }                      -- Light sea green
    }
    
    for group, opts in pairs(highlights) do
      vim.api.nvim_set_hl(0, group, opts)
    end
  end
})

-- Status line integration (show sloth icon for .sloth files)
vim.api.nvim_create_autocmd({"BufEnter", "BufWinEnter"}, {
  pattern = "*.sloth",
  callback = function()
    vim.b.sloth_file = true
    -- Update status line or show notification
    vim.notify("🦥 Sloth Runner DSL file loaded", vim.log.levels.INFO, {
      title = "Sloth Runner",
      timeout = 2000
    })
  end
})

-- Auto-commands for sloth file management
vim.api.nvim_create_autocmd("BufNewFile", {
  pattern = "*.sloth",
  callback = function()
    -- Insert template for new .sloth files
    local template = {
      "-- 🦥 Sloth Runner DSL Workflow",
      "-- Generated: " .. os.date("%Y-%m-%d %H:%M:%S"),
      "",
      "local example_task = task(\"example_task\")",
      "    :description(\"Example task description\")",
      "    :command(function(params, deps)",
      "        -- TODO: Implement task logic",
      "        print(\"🦥 Hello from Sloth Runner!\")",
      "        return true",
      "    end)",
      "    :build()",
      "",
      "workflow.define(\"example_workflow\", {",
      "    description = \"Example workflow\",",
      "    version = \"1.0.0\",",
      "    tasks = { example_task }",
      "})"
    }
    
    vim.api.nvim_buf_set_lines(0, 0, -1, false, template)
    vim.notify("🦥 New Sloth workflow template created!", vim.log.levels.SUCCESS)
  end
})

lvim.plugins = {
  {
    "Pocco81/auto-save.nvim",
    config = function()
      require("auto-save").setup({
        enabled = true,
        execution_message = {
          message = function()
            local ft = vim.bo.filetype
            local icon = ft == "sloth" and "🦥 " or ""
            return icon .. "AutoSave: " .. vim.fn.strftime("%H:%M:%S")
          end,
          dim = 0.18,
          cleaning_interval = 1250,
        },
        trigger_events = { "InsertLeave", "TextChanged" },
        conditions = {
          exists = true,
          filetype_is_not = {},
          modifiable = true,
        },
        write_all_buffers = false,
        debounce_delay = 135,
      })
    end,
  },
}