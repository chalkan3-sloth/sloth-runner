-- Sloth Runner DSL plugin for Neovim
-- ~/.config/nvim/lua/sloth-runner.lua

local M = {}

-- Plugin configuration
M.config = {
  -- Highlight configuration
  highlights = {
    enable = true,
    use_treesitter = true,
  },
  
  -- Auto-completion
  completion = {
    enable = true,
    snippets = true,
  },
  
  -- Runner integration
  runner = {
    command = "sloth-runner",
    auto_run_on_save = false,
    keymaps = {
      run_file = "<leader>sr",
      list_tasks = "<leader>sl", 
      dry_run = "<leader>st",
      debug = "<leader>sd",
    }
  },
  
  -- Folding
  folding = {
    enable = true,
    auto_close = false,
  }
}

-- Setup function
function M.setup(opts)
  M.config = vim.tbl_deep_extend("force", M.config, opts or {})
  
  -- Set up file type detection
  vim.filetype.add({
    extension = {
      ['sloth.lua'] = 'sloth',
    },
    pattern = {
      ['.*task.*%.lua'] = function(path, bufnr)
        local content = vim.api.nvim_buf_get_lines(bufnr, 0, 50, false)
        for _, line in ipairs(content) do
          if line:match('task%s*%(') or line:match('workflow%.define') then
            return 'sloth'
          end
        end
      end,
      ['.*workflow.*%.lua'] = 'sloth',
    }
  })
  
  -- Set up keymaps
  if M.config.runner.keymaps then
    M.setup_keymaps()
  end
  
  -- Set up auto commands
  M.setup_autocmds()
  
  -- Set up completion
  if M.config.completion.enable then
    M.setup_completion()
  end
end

-- Set up key mappings
function M.setup_keymaps()
  local keymaps = M.config.runner.keymaps
  
  vim.api.nvim_create_autocmd("FileType", {
    pattern = "sloth",
    callback = function()
      local opts = { buffer = true, silent = true }
      
      vim.keymap.set('n', keymaps.run_file, function()
        M.run_current_file()
      end, opts)
      
      vim.keymap.set('n', keymaps.list_tasks, function()
        M.list_tasks()
      end, opts)
      
      vim.keymap.set('n', keymaps.dry_run, function()
        M.dry_run()
      end, opts)
      
      vim.keymap.set('n', keymaps.debug, function()
        M.debug_workflow()
      end, opts)
    end
  })
end

-- Set up auto commands
function M.setup_autocmds()
  vim.api.nvim_create_augroup("SlothRunner", { clear = true })
  
  -- Auto-run on save if enabled
  if M.config.runner.auto_run_on_save then
    vim.api.nvim_create_autocmd("BufWritePost", {
      group = "SlothRunner",
      pattern = "*.sloth.lua",
      callback = function()
        M.run_current_file()
      end
    })
  end
  
  -- Enhanced syntax highlighting
  vim.api.nvim_create_autocmd("FileType", {
    group = "SlothRunner",
    pattern = "sloth",
    callback = function()
      -- Set up enhanced highlighting
      vim.cmd([[
        syntax match SlothChain '\v:\w+\ze\s*\('
        highlight link SlothChain Function
        
        syntax match SlothTaskName '\v(task|workflow\.define)\s*\(\s*"\zs[^"]*\ze"'
        highlight link SlothTaskName String
        
        syntax keyword SlothBuiltin params deps context error result output
        highlight link SlothBuiltin Identifier
      ]])
    end
  })
end

-- Set up completion
function M.setup_completion()
  vim.api.nvim_create_autocmd("FileType", {
    pattern = "sloth",
    callback = function()
      -- Add DSL-specific completion items
      local completion_items = {
        -- DSL methods
        { label = "command", kind = vim.lsp.protocol.CompletionItemKind.Method },
        { label = "description", kind = vim.lsp.protocol.CompletionItemKind.Method },
        { label = "timeout", kind = vim.lsp.protocol.CompletionItemKind.Method },
        { label = "retries", kind = vim.lsp.protocol.CompletionItemKind.Method },
        { label = "depends_on", kind = vim.lsp.protocol.CompletionItemKind.Method },
        { label = "on_success", kind = vim.lsp.protocol.CompletionItemKind.Method },
        { label = "on_failure", kind = vim.lsp.protocol.CompletionItemKind.Method },
        { label = "build", kind = vim.lsp.protocol.CompletionItemKind.Method },
        
        -- Modules
        { label = "exec", kind = vim.lsp.protocol.CompletionItemKind.Module },
        { label = "fs", kind = vim.lsp.protocol.CompletionItemKind.Module },
        { label = "net", kind = vim.lsp.protocol.CompletionItemKind.Module },
        { label = "state", kind = vim.lsp.protocol.CompletionItemKind.Module },
        { label = "metrics", kind = vim.lsp.protocol.CompletionItemKind.Module },
        
        -- Functions
        { label = "task", kind = vim.lsp.protocol.CompletionItemKind.Function },
        { label = "workflow.define", kind = vim.lsp.protocol.CompletionItemKind.Function },
      }
      
      -- Register completion source (if using nvim-cmp)
      local has_cmp, cmp = pcall(require, 'cmp')
      if has_cmp then
        cmp.register_source('sloth', {
          complete = function(self, request, callback)
            callback({ items = completion_items })
          end
        })
        
        cmp.setup.filetype('sloth', {
          sources = cmp.config.sources({
            { name = 'sloth' },
            { name = 'lua_ls' },
            { name = 'buffer' },
          })
        })
      end
    end
  })
end

-- Runner functions
function M.run_current_file()
  local file = vim.api.nvim_buf_get_name(0)
  if file == "" then
    vim.notify("No file to run", vim.log.levels.WARN)
    return
  end
  
  local cmd = string.format("%s run -f %s", M.config.runner.command, vim.fn.shellescape(file))
  
  -- Run in terminal
  vim.cmd("split")
  vim.cmd("terminal " .. cmd)
  vim.cmd("startinsert")
end

function M.list_tasks()
  local file = vim.api.nvim_buf_get_name(0)
  if file == "" then
    vim.notify("No file to analyze", vim.log.levels.WARN)
    return
  end
  
  local cmd = string.format("%s list -f %s", M.config.runner.command, vim.fn.shellescape(file))
  
  -- Show in split
  vim.cmd("split")
  vim.cmd("terminal " .. cmd)
end

function M.dry_run()
  local file = vim.api.nvim_buf_get_name(0)
  if file == "" then
    vim.notify("No file to test", vim.log.levels.WARN)
    return
  end
  
  local cmd = string.format("%s run -f %s --dry-run", M.config.runner.command, vim.fn.shellescape(file))
  
  vim.cmd("split")
  vim.cmd("terminal " .. cmd)
end

function M.debug_workflow()
  local file = vim.api.nvim_buf_get_name(0)
  if file == "" then
    vim.notify("No file to debug", vim.log.levels.WARN)
    return
  end
  
  -- Get task/workflow name under cursor
  local line = vim.api.nvim_get_current_line()
  local task_name = line:match('task%s*%(%s*"([^"]*)"') or
                   line:match('workflow%.define%s*%(%s*"([^"]*)"')
  
  local cmd = string.format("%s run -f %s --verbose", M.config.runner.command, vim.fn.shellescape(file))
  if task_name then
    cmd = cmd .. " " .. task_name
  end
  
  vim.cmd("split")
  vim.cmd("terminal " .. cmd)
end

-- Snippet functions
function M.insert_task_snippet()
  local snippet = [[local ${1:task_name} = task("${2:task_name}")
    :description("${3:Task description}")
    :command(function(params, deps)
        ${4:-- TODO: implement task logic}
        return true
    end)
    :build()]]
    
  -- Insert snippet (requires snippet engine)
  vim.snippet.expand(snippet)
end

function M.insert_workflow_snippet()
  local snippet = [[workflow.define("${1:workflow_name}", {
    description = "${2:Workflow description}",
    version = "${3:1.0.0}",
    
    tasks = {
        ${4:-- Add tasks here}
    },
    
    on_success = function(results)
        ${5:-- Success handler}
    end,
    
    on_failure = function(error, context)
        ${6:-- Error handler}
    end
})]]
    
  vim.snippet.expand(snippet)
end

-- Commands
vim.api.nvim_create_user_command('SlothRun', M.run_current_file, {})
vim.api.nvim_create_user_command('SlothList', M.list_tasks, {})
vim.api.nvim_create_user_command('SlothDryRun', M.dry_run, {})
vim.api.nvim_create_user_command('SlothDebug', M.debug_workflow, {})
vim.api.nvim_create_user_command('SlothTaskSnippet', M.insert_task_snippet, {})
vim.api.nvim_create_user_command('SlothWorkflowSnippet', M.insert_workflow_snippet, {})

return M