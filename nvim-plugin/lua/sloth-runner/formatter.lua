-- ═══════════════════════════════════════════════════════════════════════
-- Sloth Runner - Code Formatter
-- ═══════════════════════════════════════════════════════════════════════

local M = {}

local utils = require("sloth-runner.utils")
local config_module = require("sloth-runner.config")

---Format current buffer
function M.format_buffer()
	local config = config_module.get()

	if not utils.is_sloth_file() then
		vim.notify("Current buffer is not a .sloth file", vim.log.levels.WARN)
		return
	end

	-- Check if external formatter exists
	if config.formatter.cmd ~= "" and utils.command_exists(config.formatter.cmd) then
		M.format_with_external(config)
	elseif config.formatter.use_builtin then
		M.format_with_builtin()
	else
		vim.notify("No formatter available", vim.log.levels.WARN)
	end
end

---Format with external formatter (e.g., stylua)
---@param config table Plugin configuration
function M.format_with_external(config)
	local bufnr = vim.api.nvim_get_current_buf()
	local lines = vim.api.nvim_buf_get_lines(bufnr, 0, -1, false)
	local content = table.concat(lines, "\n")

	-- Write to temp file
	local temp_file = vim.fn.tempname() .. ".sloth"
	local f = io.open(temp_file, "w")
	if not f then
		vim.notify("Failed to create temp file", vim.log.levels.ERROR)
		return
	end
	f:write(content)
	f:close()

	-- Build formatter command
	local cmd_parts = { config.formatter.cmd }
	for _, arg in ipairs(config.formatter.args) do
		table.insert(cmd_parts, arg)
	end
	table.insert(cmd_parts, vim.fn.shellescape(temp_file))

	local cmd = table.concat(cmd_parts, " ")

	-- Run formatter
	local success, output = utils.run_command(cmd, {})

	if not success then
		vim.notify("Formatting failed: " .. (output or "unknown error"), vim.log.levels.ERROR)
		vim.fn.delete(temp_file)
		return
	end

	-- Read formatted content
	f = io.open(temp_file, "r")
	if not f then
		vim.notify("Failed to read formatted file", vim.log.levels.ERROR)
		vim.fn.delete(temp_file)
		return
	end

	local formatted = f:read("*a")
	f:close()
	vim.fn.delete(temp_file)

	-- Update buffer
	local formatted_lines = vim.split(formatted, "\n")

	-- Remove trailing empty line if present
	if formatted_lines[#formatted_lines] == "" then
		table.remove(formatted_lines)
	end

	-- Save cursor position
	local cursor_pos = vim.api.nvim_win_get_cursor(0)

	-- Update buffer
	vim.api.nvim_buf_set_lines(bufnr, 0, -1, false, formatted_lines)

	-- Restore cursor position
	pcall(vim.api.nvim_win_set_cursor, 0, cursor_pos)

	vim.notify("✓ File formatted", vim.log.levels.INFO)
end

---Format with built-in indentation
function M.format_with_builtin()
	local bufnr = vim.api.nvim_get_current_buf()

	-- Save cursor position
	local cursor_pos = vim.api.nvim_win_get_cursor(0)
	local view = vim.fn.winsaveview()

	-- Format using Vim's built-in indentation
	vim.cmd("silent! normal! gg=G")

	-- Restore view and cursor
	vim.fn.winrestview(view)
	pcall(vim.api.nvim_win_set_cursor, 0, cursor_pos)

	vim.notify("✓ File formatted (built-in)", vim.log.levels.INFO)
end

---Setup auto-format on save
function M.setup_autoformat()
	local config = config_module.get()

	if not config.formatter.format_on_save then
		return
	end

	vim.api.nvim_create_autocmd("BufWritePre", {
		pattern = "*.sloth",
		callback = function()
			M.format_buffer()
		end,
		desc = "Auto-format Sloth files on save",
	})
end

---Manually enable/disable auto-format
---@param enabled boolean Enable or disable
function M.set_autoformat(enabled)
	local config = config_module.get()
	config.formatter.format_on_save = enabled

	if enabled then
		M.setup_autoformat()
		vim.notify("Auto-format enabled", vim.log.levels.INFO)
	else
		vim.notify("Auto-format disabled", vim.log.levels.INFO)
	end
end

return M
