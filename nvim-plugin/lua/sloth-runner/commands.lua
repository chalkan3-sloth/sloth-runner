-- ═══════════════════════════════════════════════════════════════════════
-- Sloth Runner - Command Definitions
-- ═══════════════════════════════════════════════════════════════════════

local M = {}

local config = nil

---Setup commands
---@param cfg table Plugin configuration
function M.setup(cfg)
	config = cfg

	-- :SlothRun [task_name] - Run workflow or specific task
	vim.api.nvim_create_user_command("SlothRun", function(opts)
		local runner = require("sloth-runner.runner")
		runner.run({
			file = vim.fn.expand("%"),
			task = opts.args ~= "" and opts.args or nil,
			dry_run = false,
		})
	end, {
		nargs = "?",
		desc = "Run Sloth workflow or task",
		complete = function()
			return M.complete_tasks()
		end,
	})

	-- :SlothList - List all tasks/workflows in current file
	vim.api.nvim_create_user_command("SlothList", function()
		local runner = require("sloth-runner.runner")
		runner.list({
			file = vim.fn.expand("%"),
		})
	end, {
		desc = "List all Sloth tasks and workflows",
	})

	-- :SlothTest [task_name] - Dry run workflow or task
	vim.api.nvim_create_user_command("SlothTest", function(opts)
		local runner = require("sloth-runner.runner")
		runner.run({
			file = vim.fn.expand("%"),
			task = opts.args ~= "" and opts.args or nil,
			dry_run = true,
		})
	end, {
		nargs = "?",
		desc = "Dry run Sloth workflow or task",
		complete = function()
			return M.complete_tasks()
		end,
	})

	-- :SlothValidate - Validate current file
	vim.api.nvim_create_user_command("SlothValidate", function()
		local sloth = require("sloth-runner")
		local success, err = sloth.validate()
		if success then
			vim.notify("✓ Sloth file is valid", vim.log.levels.INFO)
		else
			vim.notify("✗ Sloth file validation failed: " .. (err or "unknown error"), vim.log.levels.ERROR)
		end
	end, {
		desc = "Validate current Sloth file",
	})

	-- :SlothFormat - Format current file
	vim.api.nvim_create_user_command("SlothFormat", function()
		local sloth = require("sloth-runner")
		sloth.format()
	end, {
		desc = "Format current Sloth file",
	})

	-- :SlothInfo - Show plugin information
	vim.api.nvim_create_user_command("SlothInfo", function()
		M.show_info()
	end, {
		desc = "Show Sloth Runner plugin information",
	})

	-- :SlothWelcome - Show welcome message again
	vim.api.nvim_create_user_command("SlothWelcome", function()
		local welcome = require("sloth-runner.welcome")
		welcome.show()
	end, {
		desc = "Show Sloth Runner welcome message",
	})

	-- :SlothAnimate - Show sloth animation (easter egg)
	vim.api.nvim_create_user_command("SlothAnimate", function()
		local welcome = require("sloth-runner.welcome")
		welcome.animate()
	end, {
		desc = "Show Sloth Runner animation",
	})

	-- :SlothTasks - Open Telescope task picker (if available)
	if config.telescope.enabled and pcall(require, "telescope") then
		vim.api.nvim_create_user_command("SlothTasks", function()
			local telescope = require("sloth-runner.telescope")
			telescope.tasks()
		end, {
			desc = "Open Telescope task picker",
		})
	end

	-- :SlothWorkflows - Open Telescope workflow picker (if available)
	if config.telescope.enabled and pcall(require, "telescope") then
		vim.api.nvim_create_user_command("SlothWorkflows", function()
			local telescope = require("sloth-runner.telescope")
			telescope.workflows()
		end, {
			desc = "Open Telescope workflow picker",
		})
	end
end

---Get completion items for task names
---@return string[] List of task names
function M.complete_tasks()
	local utils = require("sloth-runner.utils")
	local tasks = utils.parse_tasks_from_buffer()

	local completions = {}
	for _, task in ipairs(tasks) do
		table.insert(completions, task.name)
	end

	return completions
end

---Show plugin information in a float
function M.show_info()
	local utils = require("sloth-runner.utils")
	local cfg = config or require("sloth-runner.config").get()

	local lines = {
		"╔════════════════════════════════════════════════════════════════╗",
		"║                  🚀 SLOTH RUNNER PLUGIN                        ║",
		"╚════════════════════════════════════════════════════════════════╝",
		"",
		"📋 Configuration:",
		"  ├─ Runner:      " .. cfg.runner.cmd,
		"  ├─ Float UI:    " .. (cfg.runner.use_float and "enabled" or "disabled"),
		"  ├─ Completion:  " .. (cfg.completion.enabled and "enabled" or "disabled"),
		"  ├─ Telescope:   " .. (cfg.telescope.enabled and "enabled" or "disabled"),
		"  └─ Format:      " .. (cfg.formatter.format_on_save and "on save" or "manual"),
		"",
		"⌨️  Commands:",
		"  ├─ :SlothRun [task]       Run workflow or task",
		"  ├─ :SlothList             List all tasks/workflows",
		"  ├─ :SlothTest [task]      Dry run workflow or task",
		"  ├─ :SlothValidate         Validate current file",
		"  ├─ :SlothFormat           Format current file",
		"  └─ :SlothInfo             Show this information",
		"",
		"🔑 Keymaps (prefix: " .. cfg.keymaps.prefix .. "):",
		"  ├─ " .. cfg.keymaps.prefix .. cfg.keymaps.run .. "    Run workflow",
		"  ├─ " .. cfg.keymaps.prefix .. cfg.keymaps.list .. "    List tasks",
		"  ├─ " .. cfg.keymaps.prefix .. cfg.keymaps.test .. "    Dry run",
		"  ├─ " .. cfg.keymaps.prefix .. cfg.keymaps.validate .. "    Validate",
		"  └─ " .. cfg.keymaps.prefix .. cfg.keymaps.format .. "    Format",
		"",
		"📝 Text Objects:",
		"  ├─ " .. cfg.keymaps.task_textobj .. "    Select task block",
		"  └─ " .. cfg.keymaps.workflow_textobj .. "    Select workflow block",
		"",
		"Press q or <Esc> to close",
	}

	utils.create_float({
		lines = lines,
		title = " Sloth Runner Info ",
		border = "double",
		filetype = "sloth-info",
		mappings = {
			q = "close",
			["<Esc>"] = "close",
		},
	})
end

return M
