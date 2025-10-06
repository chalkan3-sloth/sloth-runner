-- ═══════════════════════════════════════════════════════════════════════
-- Sloth Runner - Keymap Configuration
-- ═══════════════════════════════════════════════════════════════════════

local M = {}

local config = nil

---Setup keymaps
---@param cfg table Plugin configuration
function M.setup(cfg)
	config = cfg

	-- Only setup keymaps for sloth filetype
	vim.api.nvim_create_autocmd("FileType", {
		pattern = "sloth",
		callback = function(event)
			M.setup_buffer_keymaps(event.buf)
		end,
		desc = "Setup Sloth Runner keymaps",
	})
end

---Setup buffer-specific keymaps
---@param bufnr number Buffer number
function M.setup_buffer_keymaps(bufnr)
	local opts = { buffer = bufnr, silent = true }
	local prefix = config.keymaps.prefix

	-- Runner commands
	vim.keymap.set("n", prefix .. config.keymaps.run, function()
		local runner = require("sloth-runner.runner")
		runner.run({
			file = vim.fn.expand("%"),
			dry_run = false,
		})
	end, vim.tbl_extend("force", opts, { desc = "Run Sloth workflow" }))

	vim.keymap.set("n", prefix .. config.keymaps.list, function()
		local runner = require("sloth-runner.runner")
		runner.list({
			file = vim.fn.expand("%"),
		})
	end, vim.tbl_extend("force", opts, { desc = "List Sloth tasks" }))

	vim.keymap.set("n", prefix .. config.keymaps.test, function()
		local runner = require("sloth-runner.runner")
		runner.run({
			file = vim.fn.expand("%"),
			dry_run = true,
		})
	end, vim.tbl_extend("force", opts, { desc = "Dry run Sloth workflow" }))

	vim.keymap.set("n", prefix .. config.keymaps.validate, function()
		local sloth = require("sloth-runner")
		sloth.validate()
	end, vim.tbl_extend("force", opts, { desc = "Validate Sloth file" }))

	vim.keymap.set("n", prefix .. config.keymaps.format, function()
		local sloth = require("sloth-runner")
		sloth.format()
	end, vim.tbl_extend("force", opts, { desc = "Format Sloth file" }))

	-- Text objects
	M.setup_text_objects(bufnr)
end

---Setup text object keymaps
---@param bufnr number Buffer number
function M.setup_text_objects(bufnr)
	local opts = { buffer = bufnr, silent = true }

	-- Task text object (it - inner task)
	vim.keymap.set({ "o", "x" }, config.keymaps.task_textobj, function()
		M.select_task_block()
	end, vim.tbl_extend("force", opts, { desc = "Select task block" }))

	-- Workflow text object (iw - inner workflow)
	vim.keymap.set({ "o", "x" }, config.keymaps.workflow_textobj, function()
		M.select_workflow_block()
	end, vim.tbl_extend("force", opts, { desc = "Select workflow block" }))
end

---Select task block text object
function M.select_task_block()
	local utils = require("sloth-runner.utils")
	local start_line = utils.find_task_start()

	if not start_line then
		vim.notify("No task block found", vim.log.levels.WARN)
		return
	end

	local end_line = utils.find_task_end(start_line)
	if not end_line then
		vim.notify("Could not find task end", vim.log.levels.WARN)
		return
	end

	-- Enter visual line mode and select the range
	vim.cmd("normal! V")
	vim.fn.cursor(start_line, 1)
	vim.cmd("normal! o")
	vim.fn.cursor(end_line, 1)
end

---Select workflow block text object
function M.select_workflow_block()
	local utils = require("sloth-runner.utils")
	local start_line = utils.find_workflow_start()

	if not start_line then
		vim.notify("No workflow block found", vim.log.levels.WARN)
		return
	end

	local end_line = utils.find_workflow_end(start_line)
	if not end_line then
		vim.notify("Could not find workflow end", vim.log.levels.WARN)
		return
	end

	-- Enter visual line mode and select the range
	vim.cmd("normal! V")
	vim.fn.cursor(start_line, 1)
	vim.cmd("normal! o")
	vim.fn.cursor(end_line, 1)
end

---Create which-key integration if available
---@param bufnr number Buffer number
function M.setup_which_key(bufnr)
	local ok, which_key = pcall(require, "which-key")
	if not ok then
		return
	end

	local prefix = config.keymaps.prefix

	which_key.add({
		{ prefix, group = "Sloth Runner", buffer = bufnr },
		{ prefix .. config.keymaps.run, desc = "Run workflow", buffer = bufnr },
		{ prefix .. config.keymaps.list, desc = "List tasks", buffer = bufnr },
		{ prefix .. config.keymaps.test, desc = "Dry run", buffer = bufnr },
		{ prefix .. config.keymaps.validate, desc = "Validate", buffer = bufnr },
		{ prefix .. config.keymaps.format, desc = "Format", buffer = bufnr },
	})
end

return M
