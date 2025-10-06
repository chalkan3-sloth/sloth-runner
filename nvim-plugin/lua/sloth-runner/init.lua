-- ═══════════════════════════════════════════════════════════════════════
-- Sloth Runner - Modern Neovim Plugin
-- ═══════════════════════════════════════════════════════════════════════

local M = {}

-- Plugin state
M._initialized = false
M._config = nil

---Setup function to initialize the plugin
---@param opts table|nil Configuration options
function M.setup(opts)
	if M._initialized then
		vim.notify("Sloth Runner already initialized", vim.log.levels.WARN)
		return
	end

	-- Load config module and merge with user options
	local config = require("sloth-runner.config")
	M._config = config.setup(opts or {})

	-- Load core modules
	local commands = require("sloth-runner.commands")
	local keymaps = require("sloth-runner.keymaps")

	-- Register commands
	commands.setup(M._config)

	-- Setup keymaps if enabled
	if M._config.keymaps.enabled then
		keymaps.setup(M._config)
	end

	-- Setup completion if nvim-cmp is available
	if M._config.completion.enabled then
		local ok, completion = pcall(require, "sloth-runner.completion")
		if ok then
			completion.setup(M._config)
		else
			vim.notify(
				"Sloth Runner: nvim-cmp not available, completion disabled",
				vim.log.levels.DEBUG
			)
		end
	end

	-- Setup telescope integration if available
	if M._config.telescope.enabled then
		local ok = pcall(require, "telescope")
		if ok then
			local telescope = require("sloth-runner.telescope")
			telescope.setup(M._config)
		else
			vim.notify(
				"Sloth Runner: telescope.nvim not available, integration disabled",
				vim.log.levels.DEBUG
			)
		end
	end

	M._initialized = true

	vim.notify("Sloth Runner initialized", vim.log.levels.INFO)
end

---Get current configuration
---@return table Configuration
function M.get_config()
	return M._config or require("sloth-runner.config").defaults
end

---Check if plugin is initialized
---@return boolean
function M.is_initialized()
	return M._initialized
end

---Run a sloth workflow
---@param args table Arguments for the runner
function M.run(args)
	local runner = require("sloth-runner.runner")
	runner.run(args)
end

---List available tasks/workflows
---@param args table Arguments for listing
function M.list(args)
	local runner = require("sloth-runner.runner")
	runner.list(args)
end

---Validate current sloth file
---@return boolean, string|nil Success, error message
function M.validate()
	local runner = require("sloth-runner.runner")
	return runner.validate()
end

---Format current sloth file
function M.format()
	local formatter = require("sloth-runner.formatter")
	formatter.format_buffer()
end

return M
