-- ═══════════════════════════════════════════════════════════════════════
-- Sloth Runner - Workflow Runner
-- ═══════════════════════════════════════════════════════════════════════

local M = {}

local utils = require("sloth-runner.utils")
local config_module = require("sloth-runner.config")

---Run a workflow or task
---@param args table Arguments { file, task, dry_run }
function M.run(args)
	local config = config_module.get()

	-- Validate file exists
	local file = args.file or utils.get_current_file()
	if not file or file == "" then
		vim.notify("No file specified", vim.log.levels.ERROR)
		return
	end

	if vim.fn.filereadable(file) ~= 1 then
		vim.notify("File not found: " .. file, vim.log.levels.ERROR)
		return
	end

	-- Build command
	local cmd_parts = { config.runner.cmd, "run", "-f", vim.fn.shellescape(file) }

	if args.task then
		table.insert(cmd_parts, vim.fn.shellescape(args.task))
	end

	if args.dry_run then
		table.insert(cmd_parts, "--dry-run")
	end

	-- Add default args
	for _, arg in ipairs(config.runner.default_args) do
		table.insert(cmd_parts, arg)
	end

	local cmd = table.concat(cmd_parts, " ")

	-- Log command
	utils.log("Running command: " .. cmd, vim.log.levels.DEBUG)

	-- Notify start
	if config.runner.notify then
		local icon = config.ui.icons.running
		local msg = args.dry_run and "Dry running workflow..." or "Running workflow..."
		vim.notify(icon .. " " .. msg, vim.log.levels.INFO)
	end

	-- Run command
	if config.runner.use_float then
		M.run_in_float(cmd, config)
	else
		M.run_in_terminal(cmd, config)
	end
end

---Run command in floating window
---@param cmd string Command to run
---@param config table Plugin configuration
function M.run_in_float(cmd, config)
	-- Create a buffer for output
	local bufnr = vim.api.nvim_create_buf(false, true)
	vim.api.nvim_buf_set_option(bufnr, "bufhidden", "wipe")
	vim.api.nvim_buf_set_option(bufnr, "filetype", "sloth-output")

	-- Calculate window size
	local width = math.floor(vim.o.columns * config.ui.float.width)
	local height = math.floor(vim.o.lines * config.ui.float.height)
	local row = math.floor((vim.o.lines - height) / 2)
	local col = math.floor((vim.o.columns - width) / 2)

	-- Window options
	local win_opts = {
		relative = "editor",
		width = width,
		height = height,
		row = row,
		col = col,
		border = config.ui.float.border,
		title = " Sloth Runner Output ",
		title_pos = config.ui.float.title_pos,
		style = "minimal",
	}

	-- Create window
	local winnr = vim.api.nvim_open_win(bufnr, true, win_opts)

	-- Setup close mappings
	local close_mappings = { "q", "<Esc>" }
	for _, key in ipairs(close_mappings) do
		vim.keymap.set("n", key, function()
			if vim.api.nvim_win_is_valid(winnr) then
				vim.api.nvim_win_close(winnr, true)
			end
		end, { buffer = bufnr, nowait = true, silent = true })
	end

	-- Run command asynchronously
	local lines = {}
	local job_id = vim.fn.jobstart(cmd, {
		on_stdout = function(_, data)
			if data then
				for _, line in ipairs(data) do
					if line ~= "" then
						table.insert(lines, line)
					end
				end
				-- Update buffer
				if vim.api.nvim_buf_is_valid(bufnr) then
					vim.api.nvim_buf_set_option(bufnr, "modifiable", true)
					vim.api.nvim_buf_set_lines(bufnr, 0, -1, false, lines)
					vim.api.nvim_buf_set_option(bufnr, "modifiable", false)
					-- Scroll to bottom
					if vim.api.nvim_win_is_valid(winnr) then
						vim.api.nvim_win_set_cursor(winnr, { #lines, 0 })
					end
				end
			end
		end,
		on_stderr = function(_, data)
			if data then
				for _, line in ipairs(data) do
					if line ~= "" then
						table.insert(lines, "ERROR: " .. line)
					end
				end
				-- Update buffer
				if vim.api.nvim_buf_is_valid(bufnr) then
					vim.api.nvim_buf_set_option(bufnr, "modifiable", true)
					vim.api.nvim_buf_set_lines(bufnr, 0, -1, false, lines)
					vim.api.nvim_buf_set_option(bufnr, "modifiable", false)
				end
			end
		end,
		on_exit = function(_, exit_code)
			local icon = exit_code == 0 and config.ui.icons.success or config.ui.icons.error
			local msg = exit_code == 0 and "Workflow completed successfully"
				or "Workflow failed (exit " .. exit_code .. ")"
			local level = exit_code == 0 and vim.log.levels.INFO or vim.log.levels.ERROR

			if config.runner.notify then
				vim.notify(icon .. " " .. msg, level)
			end

			-- Auto-close on success if configured
			if exit_code == 0 and config.runner.auto_close_on_success then
				vim.defer_fn(function()
					if vim.api.nvim_win_is_valid(winnr) then
						vim.api.nvim_win_close(winnr, true)
					end
				end, 1000)
			end
		end,
	})

	if job_id <= 0 then
		vim.notify("Failed to start job", vim.log.levels.ERROR)
		if vim.api.nvim_win_is_valid(winnr) then
			vim.api.nvim_win_close(winnr, true)
		end
	end
end

---Run command in terminal
---@param cmd string Command to run
---@param config table Plugin configuration
function M.run_in_terminal(cmd, config)
	-- Open terminal in split
	vim.cmd("botright split | terminal " .. cmd)

	-- Setup auto-close on success
	if config.runner.auto_close_on_success then
		vim.api.nvim_create_autocmd("TermClose", {
			buffer = vim.api.nvim_get_current_buf(),
			callback = function()
				local exit_code = vim.v.event.status
				if exit_code == 0 then
					vim.defer_fn(function()
						vim.cmd("bdelete!")
					end, 1000)
				end
			end,
			once = true,
		})
	end
end

---List all tasks and workflows
---@param args table Arguments { file }
function M.list(args)
	local config = config_module.get()

	local file = args.file or utils.get_current_file()
	if not file or file == "" then
		vim.notify("No file specified", vim.log.levels.ERROR)
		return
	end

	-- Build command
	local cmd_parts = { config.runner.cmd, "list", "-f", vim.fn.shellescape(file) }

	for _, arg in ipairs(config.runner.default_args) do
		table.insert(cmd_parts, arg)
	end

	local cmd = table.concat(cmd_parts, " ")

	-- Run and get output
	local success, output = utils.run_command(cmd, {})

	if not success then
		vim.notify("Failed to list tasks: " .. (output or "unknown error"), vim.log.levels.ERROR)
		return
	end

	-- Show in float
	local lines = vim.split(output or "", "\n")
	utils.create_float({
		lines = lines,
		title = " Sloth Tasks & Workflows ",
		border = "rounded",
		filetype = "sloth-list",
		mappings = {
			q = "close",
			["<Esc>"] = "close",
		},
	})
end

---Validate current sloth file
---@return boolean Success
---@return string|nil Error message
function M.validate()
	local config = config_module.get()
	local file = utils.get_current_file()

	if not file or file == "" then
		return false, "No file specified"
	end

	if not utils.is_sloth_file() then
		return false, "Current file is not a .sloth file"
	end

	-- Check if runner command exists
	if not utils.command_exists(config.runner.cmd) then
		return false, "sloth-runner command not found in PATH"
	end

	-- Run validation
	local cmd = config.runner.cmd .. " validate -f " .. vim.fn.shellescape(file)
	local success, output = utils.run_command(cmd, {})

	if not success then
		return false, output or "Validation failed"
	end

	return true, nil
end

return M
