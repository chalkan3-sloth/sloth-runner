-- Sloth Runner DSL Formatter
-- Simple formatter for .sloth files

local M = {}

-- Format Sloth file using sloth-runner if available
function M.format_sloth()
  local bufnr = vim.api.nvim_get_current_buf()
  local filename = vim.api.nvim_buf_get_name(bufnr)
  
  -- Check if current file is a .sloth file
  if not filename:match("%.sloth$") then
    return
  end
  
  -- Check if sloth-runner is available
  local handle = io.popen("which sloth-runner 2>/dev/null")
  local sloth_runner_path = handle:read("*a"):gsub("%s+", "")
  handle:close()
  
  if sloth_runner_path == "" then
    -- Fallback to simple lua formatting
    vim.cmd("normal! gg=G``")
    return
  end
  
  -- Check if file exists and is readable
  local file = io.open(filename, "r")
  if not file then
    return
  end
  file:close()
  
  -- Try to format with sloth-runner (if it has a format command)
  local cmd = string.format("%s format %s 2>/dev/null", sloth_runner_path, vim.fn.shellescape(filename))
  local result = vim.fn.system(cmd)
  
  -- If command failed or returned empty, fallback to basic formatting
  if vim.v.shell_error ~= 0 or result:match("^%s*$") then
    -- Basic Lua-style formatting
    vim.cmd("normal! gg=G``")
  end
end

-- Setup autoformat on save
function M.setup_autoformat()
  vim.api.nvim_create_augroup("SlothFormatter", { clear = true })
  
  vim.api.nvim_create_autocmd("BufWritePre", {
    group = "SlothFormatter",
    pattern = "*.sloth",
    callback = function()
      -- Only format if explicitly enabled
      if vim.g.sloth_format_on_save then
        M.format_sloth()
      end
    end,
  })
end

-- Manual format command
function M.format_current_buffer()
  local bufnr = vim.api.nvim_get_current_buf()
  local filename = vim.api.nvim_buf_get_name(bufnr)
  
  if filename:match("%.sloth$") then
    M.format_sloth()
    print("Formatted " .. vim.fn.fnamemodify(filename, ":t"))
  else
    print("Not a .sloth file")
  end
end

return M