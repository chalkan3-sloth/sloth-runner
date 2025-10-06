-- ═══════════════════════════════════════════════════════════════════════
-- Sloth Runner - Utility Functions
-- ═══════════════════════════════════════════════════════════════════════

local M = {}

---Parse tasks from current buffer
---@return table[] List of tasks with name, line, description
function M.parse_tasks_from_buffer()
	local bufnr = vim.api.nvim_get_current_buf()
	local lines = vim.api.nvim_buf_get_lines(bufnr, 0, -1, false)

	local tasks = {}

	for i, line in ipairs(lines) do
		-- Match: local task_name = task("task-name")
		local task_name = line:match('local%s+(%w+)%s*=%s*task%s*%(%s*["\']([^"\']+)["\']%s*%)')
		if task_name then
			-- Look ahead for description
			local desc = nil
			if i < #lines then
				local next_line = lines[i + 1]
				desc = next_line:match(':description%s*%(%s*["\']([^"\']+)["\']%s*%)')
			end

			table.insert(tasks, {
				name = task_name,
				line = i,
				description = desc,
			})
		end
	end

	return tasks
end

---Parse workflows from current buffer
---@return table[] List of workflows with name, line, description
function M.parse_workflows_from_buffer()
	local bufnr = vim.api.nvim_get_current_buf()
	local lines = vim.api.nvim_buf_get_lines(bufnr, 0, -1, false)

	local workflows = {}

	for i, line in ipairs(lines) do
		-- Match: workflow.define("workflow-name", {
		local workflow_name = line:match('workflow%.define%s*%(%s*["\']([^"\']+)["\']%s*,%s*{')
		if workflow_name then
			-- Look ahead for description in the next few lines
			local desc = nil
			for j = i, math.min(i + 5, #lines) do
				local check_line = lines[j]
				desc = check_line:match('description%s*=%s*["\']([^"\']+)["\']')
				if desc then
					break
				end
			end

			table.insert(workflows, {
				name = workflow_name,
				line = i,
				description = desc,
			})
		end
	end

	return workflows
end

---Find start of task block from cursor position
---@return number|nil Line number or nil
function M.find_task_start()
	local cursor = vim.api.nvim_win_get_cursor(0)
	local current_line = cursor[1]

	-- Search backwards for task definition
	for i = current_line, 1, -1 do
		local line = vim.api.nvim_buf_get_lines(0, i - 1, i, false)[1]
		if line:match('local%s+%w+%s*=%s*task%s*%(') then
			return i
		end
	end

	return nil
end

---Find end of task block from start line
---@param start_line number Starting line number
---@return number|nil End line number or nil
function M.find_task_end(start_line)
	local lines = vim.api.nvim_buf_get_lines(0, start_line - 1, -1, false)

	-- Search for :build()
	for i, line in ipairs(lines) do
		if line:match(':build%s*%(') then
			return start_line + i - 1
		end
	end

	return nil
end

---Find start of workflow block from cursor position
---@return number|nil Line number or nil
function M.find_workflow_start()
	local cursor = vim.api.nvim_win_get_cursor(0)
	local current_line = cursor[1]

	-- Search backwards for workflow definition
	for i = current_line, 1, -1 do
		local line = vim.api.nvim_buf_get_lines(0, i - 1, i, false)[1]
		if line:match('workflow%.define%s*%(') then
			return i
		end
	end

	return nil
end

---Find end of workflow block from start line
---@param start_line number Starting line number
---@return number|nil End line number or nil
function M.find_workflow_end(start_line)
	local lines = vim.api.nvim_buf_get_lines(0, start_line - 1, -1, false)

	-- Track brace depth
	local depth = 0
	local found_opening = false

	for i, line in ipairs(lines) do
		-- Count opening and closing braces
		for char in line:gmatch(".") do
			if char == "{" then
				depth = depth + 1
				found_opening = true
			elseif char == "}" then
				depth = depth - 1
				if found_opening and depth == 0 then
					-- Check if this line also has closing paren
					if line:match("}%s*%)") then
						return start_line + i - 1
					end
				end
			end
		end
	end

	return nil
end

---Create a floating window with content
---@param opts table Options (lines, title, border, filetype, mappings)
---@return number Buffer number
---@return number Window number
function M.create_float(opts)
	opts = opts or {}
	local config = require("sloth-runner.config").get()

	-- Create buffer
	local bufnr = vim.api.nvim_create_buf(false, true)
	vim.api.nvim_buf_set_lines(bufnr, 0, -1, false, opts.lines or {})

	-- Set buffer options
	vim.api.nvim_buf_set_option(bufnr, "modifiable", false)
	vim.api.nvim_buf_set_option(bufnr, "bufhidden", "wipe")

	if opts.filetype then
		vim.api.nvim_buf_set_option(bufnr, "filetype", opts.filetype)
	end

	-- Calculate window size
	local width = opts.width or config.ui.float.width
	local height = opts.height or config.ui.float.height

	if width < 1 then
		width = math.floor(vim.o.columns * width)
	end
	if height < 1 then
		height = math.floor(vim.o.lines * height)
	end

	-- Calculate position
	local row = math.floor((vim.o.lines - height) / 2)
	local col = math.floor((vim.o.columns - width) / 2)

	-- Window options
	local win_opts = {
		relative = "editor",
		width = width,
		height = height,
		row = row,
		col = col,
		border = opts.border or config.ui.float.border,
		title = opts.title or "",
		title_pos = config.ui.float.title_pos,
		style = "minimal",
	}

	-- Create window
	local winnr = vim.api.nvim_open_win(bufnr, true, win_opts)

	-- Setup mappings
	if opts.mappings then
		for key, action in pairs(opts.mappings) do
			if action == "close" then
				vim.keymap.set("n", key, function()
					vim.api.nvim_win_close(winnr, true)
				end, { buffer = bufnr, nowait = true, silent = true })
			end
		end
	end

	return bufnr, winnr
end

---Run shell command and return output
---@param cmd string Command to run
---@param args table Arguments
---@return boolean Success
---@return string|nil Output or error message
function M.run_command(cmd, args)
	local full_cmd = cmd
	if args and #args > 0 then
		full_cmd = cmd .. " " .. table.concat(args, " ")
	end

	local handle = io.popen(full_cmd .. " 2>&1")
	if not handle then
		return false, "Failed to execute command"
	end

	local output = handle:read("*a")
	local success = handle:close()

	return success, output
end

---Check if command exists in PATH
---@param cmd string Command name
---@return boolean
function M.command_exists(cmd)
	local handle = io.popen("which " .. cmd .. " 2>/dev/null")
	if not handle then
		return false
	end

	local result = handle:read("*a")
	handle:close()

	return result ~= ""
end

---Get current file path
---@return string File path
function M.get_current_file()
	return vim.fn.expand("%:p")
end

---Check if current buffer is a sloth file
---@return boolean
function M.is_sloth_file()
	local filetype = vim.bo.filetype
	return filetype == "sloth"
end

---Debug log function
---@param message string Log message
---@param level number|nil Log level (vim.log.levels)
function M.log(message, level)
	local config = require("sloth-runner.config").get()

	if not config.debug.enabled then
		return
	end

	level = level or vim.log.levels.DEBUG

	-- Log to file
	local log_file = io.open(config.debug.log_file, "a")
	if log_file then
		local timestamp = os.date("%Y-%m-%d %H:%M:%S")
		log_file:write(string.format("[%s] %s\n", timestamp, message))
		log_file:close()
	end

	-- Also show in notify if INFO or higher
	if level >= vim.log.levels.INFO then
		vim.notify(message, level)
	end
end

---Highlight text in buffer
---@param bufnr number Buffer number
---@param ns_id number Namespace ID
---@param line number Line number (0-indexed)
---@param col_start number Start column
---@param col_end number End column
---@param hl_group string Highlight group
function M.highlight_range(bufnr, ns_id, line, col_start, col_end, hl_group)
	vim.api.nvim_buf_add_highlight(bufnr, ns_id, hl_group, line, col_start, col_end)
end

---Clear highlights in namespace
---@param bufnr number Buffer number
---@param ns_id number Namespace ID
function M.clear_highlights(bufnr, ns_id)
	vim.api.nvim_buf_clear_namespace(bufnr, ns_id, 0, -1)
end

return M
