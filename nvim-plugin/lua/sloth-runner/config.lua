-- ═══════════════════════════════════════════════════════════════════════
-- Sloth Runner - Configuration Management
-- ═══════════════════════════════════════════════════════════════════════

local M = {}

---@class SlothConfig
---@field runner table Runner configuration
---@field formatter table Formatter configuration
---@field completion table Completion configuration
---@field keymaps table Keymap configuration
---@field telescope table Telescope configuration
---@field ui table UI configuration

---Default configuration
M.defaults = {
	-- Runner configuration
	runner = {
		-- Path to sloth-runner executable
		cmd = "sloth-runner",
		-- Default arguments for runner
		default_args = {},
		-- Use floating window for output
		use_float = true,
		-- Auto-close output window on success
		auto_close_on_success = false,
		-- Show notifications
		notify = true,
	},

	-- Formatter configuration
	formatter = {
		-- Enable auto-format on save
		format_on_save = false,
		-- Formatter command (can be "stylua" or custom)
		cmd = "stylua",
		-- Formatter arguments
		args = { "--indent-type", "Spaces", "--indent-width", "2" },
		-- Use built-in formatter if external not available
		use_builtin = true,
	},

	-- Completion configuration
	completion = {
		-- Enable nvim-cmp integration
		enabled = true,
		-- Completion priority
		priority = 100,
		-- Show documentation in completion
		show_docs = true,
		-- DSL keywords to complete
		keywords = {
			-- Core DSL
			"task",
			"workflow",
			"define",
			-- Task methods
			"command",
			"description",
			"timeout",
			"retries",
			"depends_on",
			"on_success",
			"on_failure",
			"build",
			"run_on",
			"agent",
			"tags",
			"artifacts",
			"condition",
			"schedule",
			"retry_count",
			"backoff_strategy",
			"circuit_breaker",
			-- Modules
			"exec",
			"fs",
			"net",
			"data",
			"log",
			"state",
			"metrics",
			"aws",
			"gcp",
			"azure",
			"digitalocean",
			"docker",
			"kubernetes",
			"terraform",
			"pulumi",
			"git",
			"notification",
			"crypto",
			"utils",
		},
	},

	-- Keymap configuration
	keymaps = {
		-- Enable default keymaps
		enabled = true,
		-- Prefix for all keymaps
		prefix = "<leader>s",
		-- Individual keymaps
		run = "r", -- Run workflow
		list = "l", -- List tasks
		test = "t", -- Dry run
		validate = "v", -- Validate file
		format = "f", -- Format file
		-- Text objects
		task_textobj = "it", -- Inner task
		workflow_textobj = "iw", -- Inner workflow
	},

	-- Telescope integration
	telescope = {
		-- Enable telescope integration
		enabled = true,
		-- Default theme
		theme = "dropdown",
		-- Layout configuration
		layout_config = {
			width = 0.8,
			height = 0.6,
		},
	},

	-- UI configuration
	ui = {
		-- Icons for different elements
		icons = {
			task = "📋",
			workflow = "🔄",
			running = "⚡",
			success = "✓",
			error = "✗",
			warning = "⚠",
		},
		-- Float window configuration
		float = {
			border = "rounded",
			title_pos = "center",
			width = 0.8,
			height = 0.8,
		},
		-- Highlight groups
		highlights = {
			task = "Function",
			workflow = "Type",
			method = "Method",
			module = "Include",
			keyword = "Keyword",
		},
		-- Welcome banner
		show_welcome = true, -- Show welcome message on first .sloth file
		welcome_style = "notification", -- Style: "notification", "banner", "large", "float"
	},

	-- Debug configuration
	debug = {
		-- Enable debug logging
		enabled = false,
		-- Log file path
		log_file = vim.fn.stdpath("cache") .. "/sloth-runner.log",
	},
}

---Current configuration
M.current = vim.deepcopy(M.defaults)

---Merge user configuration with defaults
---@param user_config table User configuration
---@return table Merged configuration
local function merge_config(user_config)
	return vim.tbl_deep_extend("force", M.defaults, user_config or {})
end

---Validate configuration
---@param config table Configuration to validate
---@return boolean, string|nil Valid, error message
local function validate_config(config)
	-- Check runner.cmd is string
	if type(config.runner.cmd) ~= "string" then
		return false, "runner.cmd must be a string"
	end

	-- Check formatter.format_on_save is boolean
	if type(config.formatter.format_on_save) ~= "boolean" then
		return false, "formatter.format_on_save must be a boolean"
	end

	-- Check keymaps.enabled is boolean
	if type(config.keymaps.enabled) ~= "boolean" then
		return false, "keymaps.enabled must be a boolean"
	end

	return true, nil
end

---Setup configuration with user options
---@param opts table User configuration
---@return table Final configuration
function M.setup(opts)
	local config = merge_config(opts)

	-- Validate configuration
	local valid, err = validate_config(config)
	if not valid then
		vim.notify(
			"Sloth Runner: Invalid configuration: " .. (err or "unknown error"),
			vim.log.levels.ERROR
		)
		return M.defaults
	end

	M.current = config
	return M.current
end

---Get current configuration
---@return table Current configuration
function M.get()
	return M.current
end

---Update configuration at runtime
---@param updates table Configuration updates
function M.update(updates)
	M.current = merge_config(updates)
end

---Reset configuration to defaults
function M.reset()
	M.current = vim.deepcopy(M.defaults)
end

return M
