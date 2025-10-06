-- ═══════════════════════════════════════════════════════════════════════
-- Sloth Runner - Health Check
-- ═══════════════════════════════════════════════════════════════════════

local M = {}

local health = vim.health or require("health")
local utils = require("sloth-runner.utils")

---Run health check
function M.check()
	health.start("Sloth Runner")

	-- Check if plugin is initialized
	local sloth = require("sloth-runner")
	if sloth.is_initialized() then
		health.ok("Plugin initialized")
	else
		health.warn("Plugin not initialized", {
			"Call require('sloth-runner').setup() in your Neovim config",
		})
	end

	-- Check sloth-runner executable
	M.check_runner()

	-- Check optional dependencies
	M.check_optional_dependencies()

	-- Check configuration
	M.check_configuration()
end

---Check sloth-runner command
function M.check_runner()
	health.start("Sloth Runner Executable")

	local config = require("sloth-runner.config").get()
	local cmd = config.runner.cmd

	if utils.command_exists(cmd) then
		health.ok("Found " .. cmd .. " in PATH")

		-- Try to get version
		local success, output = utils.run_command(cmd .. " --version", {})
		if success and output then
			local version = output:match("(%d+%.%d+%.%d+)")
			if version then
				health.info("Version: " .. version)
			end
		end
	else
		health.error(cmd .. " not found in PATH", {
			"Install sloth-runner from: https://github.com/user/sloth-runner",
			"Make sure it's in your PATH",
		})
	end
end

---Check optional dependencies
function M.check_optional_dependencies()
	health.start("Optional Dependencies")

	-- Check nvim-cmp
	local has_cmp = pcall(require, "cmp")
	if has_cmp then
		health.ok("nvim-cmp installed (completion enabled)")
	else
		health.warn("nvim-cmp not installed", {
			"Install nvim-cmp for DSL completion: https://github.com/hrsh7th/nvim-cmp",
		})
	end

	-- Check telescope
	local has_telescope = pcall(require, "telescope")
	if has_telescope then
		health.ok("telescope.nvim installed (picker enabled)")
	else
		health.warn("telescope.nvim not installed", {
			"Install telescope for task/workflow picker: https://github.com/nvim-telescope/telescope.nvim",
		})
	end

	-- Check which-key
	local has_which_key = pcall(require, "which-key")
	if has_which_key then
		health.ok("which-key installed (keymap hints enabled)")
	else
		health.info("which-key not installed (optional)", {
			"Install which-key for keymap hints: https://github.com/folke/which-key.nvim",
		})
	end

	-- Check formatter (stylua)
	local config = require("sloth-runner.config").get()
	if config.formatter.cmd ~= "" then
		if utils.command_exists(config.formatter.cmd) then
			health.ok("Formatter (" .. config.formatter.cmd .. ") found")
		else
			health.warn("Formatter (" .. config.formatter.cmd .. ") not found", {
				"Install " .. config.formatter.cmd .. " or configure a different formatter",
				"Built-in formatter will be used as fallback",
			})
		end
	end
end

---Check configuration
function M.check_configuration()
	health.start("Configuration")

	local config = require("sloth-runner.config").get()

	-- Check runner config
	if type(config.runner.cmd) == "string" and config.runner.cmd ~= "" then
		health.ok("Runner command configured: " .. config.runner.cmd)
	else
		health.error("Invalid runner command configuration")
	end

	-- Check UI config
	if config.runner.use_float then
		health.info("Using floating window for output")
	else
		health.info("Using terminal split for output")
	end

	-- Check keymaps
	if config.keymaps.enabled then
		health.ok("Keymaps enabled (prefix: " .. config.keymaps.prefix .. ")")
	else
		health.info("Keymaps disabled")
	end

	-- Check completion
	if config.completion.enabled then
		if pcall(require, "cmp") then
			health.ok("Completion enabled")
		else
			health.warn("Completion enabled but nvim-cmp not installed")
		end
	else
		health.info("Completion disabled")
	end

	-- Check telescope integration
	if config.telescope.enabled then
		if pcall(require, "telescope") then
			health.ok("Telescope integration enabled")
		else
			health.warn("Telescope integration enabled but telescope.nvim not installed")
		end
	else
		health.info("Telescope integration disabled")
	end

	-- Check format on save
	if config.formatter.format_on_save then
		health.info("Auto-format on save enabled")
	else
		health.info("Auto-format on save disabled")
	end

	-- Check debug mode
	if config.debug.enabled then
		health.warn("Debug mode enabled", {
			"Logging to: " .. config.debug.log_file,
		})
	else
		health.info("Debug mode disabled")
	end
end

return M
