-- ═══════════════════════════════════════════════════════════════════════
-- Sloth Runner - Telescope Integration
-- ═══════════════════════════════════════════════════════════════════════

local M = {}

local utils = require("sloth-runner.utils")
local config_module = require("sloth-runner.config")

local has_telescope, telescope = pcall(require, "telescope")
if not has_telescope then
	return M
end

local pickers = require("telescope.pickers")
local finders = require("telescope.finders")
local conf = require("telescope.config").values
local actions = require("telescope.actions")
local action_state = require("telescope.actions.state")
local previewers = require("telescope.previewers")

---Setup telescope integration
---@param config table Plugin configuration
function M.setup(config)
	-- Register telescope extension
	telescope.register_extension({
		exports = {
			tasks = M.tasks,
			workflows = M.workflows,
			sloth = M.all,
		},
	})

	utils.log("Telescope integration loaded", vim.log.levels.DEBUG)
end

---Open task picker
---@param opts table|nil Telescope options
function M.tasks(opts)
	opts = opts or {}
	local config = config_module.get()

	-- Merge with default theme
	opts = vim.tbl_deep_extend("force", {
		prompt_title = "📋 Sloth Tasks",
		layout_strategy = config.telescope.theme,
		layout_config = config.telescope.layout_config,
	}, opts)

	-- Get tasks from buffer
	local tasks = utils.parse_tasks_from_buffer()

	if #tasks == 0 then
		vim.notify("No tasks found in current buffer", vim.log.levels.WARN)
		return
	end

	-- Create picker
	pickers
		.new(opts, {
			finder = finders.new_table({
				results = tasks,
				entry_maker = function(task)
					return {
						value = task,
						display = config.ui.icons.task
							.. " "
							.. task.name
							.. (task.description and (" - " .. task.description) or ""),
						ordinal = task.name .. " " .. (task.description or ""),
						lnum = task.line,
					}
				end,
			}),
			sorter = conf.generic_sorter(opts),
			previewer = M.create_task_previewer(),
			attach_mappings = function(prompt_bufnr, map)
				-- Run task on <CR>
				actions.select_default:replace(function()
					local selection = action_state.get_selected_entry()
					actions.close(prompt_bufnr)

					if selection then
						local runner = require("sloth-runner.runner")
						runner.run({
							file = utils.get_current_file(),
							task = selection.value.name,
						})
					end
				end)

				-- Go to task definition on <C-g>
				map("i", "<C-g>", function()
					local selection = action_state.get_selected_entry()
					actions.close(prompt_bufnr)

					if selection then
						vim.api.nvim_win_set_cursor(0, { selection.lnum, 0 })
					end
				end)

				return true
			end,
		})
		:find()
end

---Open workflow picker
---@param opts table|nil Telescope options
function M.workflows(opts)
	opts = opts or {}
	local config = config_module.get()

	opts = vim.tbl_deep_extend("force", {
		prompt_title = "🔄 Sloth Workflows",
		layout_strategy = config.telescope.theme,
		layout_config = config.telescope.layout_config,
	}, opts)

	-- Get workflows from buffer
	local workflows = utils.parse_workflows_from_buffer()

	if #workflows == 0 then
		vim.notify("No workflows found in current buffer", vim.log.levels.WARN)
		return
	end

	-- Create picker
	pickers
		.new(opts, {
			finder = finders.new_table({
				results = workflows,
				entry_maker = function(workflow)
					return {
						value = workflow,
						display = config.ui.icons.workflow
							.. " "
							.. workflow.name
							.. (workflow.description and (" - " .. workflow.description) or ""),
						ordinal = workflow.name .. " " .. (workflow.description or ""),
						lnum = workflow.line,
					}
				end,
			}),
			sorter = conf.generic_sorter(opts),
			previewer = M.create_workflow_previewer(),
			attach_mappings = function(prompt_bufnr, map)
				-- Run workflow on <CR>
				actions.select_default:replace(function()
					local selection = action_state.get_selected_entry()
					actions.close(prompt_bufnr)

					if selection then
						local runner = require("sloth-runner.runner")
						runner.run({
							file = utils.get_current_file(),
						})
					end
				end)

				-- Go to workflow definition on <C-g>
				map("i", "<C-g>", function()
					local selection = action_state.get_selected_entry()
					actions.close(prompt_bufnr)

					if selection then
						vim.api.nvim_win_set_cursor(0, { selection.lnum, 0 })
					end
				end)

				return true
			end,
		})
		:find()
end

---Open combined picker (tasks + workflows)
---@param opts table|nil Telescope options
function M.all(opts)
	opts = opts or {}
	local config = config_module.get()

	opts = vim.tbl_deep_extend("force", {
		prompt_title = "🚀 Sloth Tasks & Workflows",
		layout_strategy = config.telescope.theme,
		layout_config = config.telescope.layout_config,
	}, opts)

	-- Get both tasks and workflows
	local tasks = utils.parse_tasks_from_buffer()
	local workflows = utils.parse_workflows_from_buffer()

	local items = {}

	-- Add tasks
	for _, task in ipairs(tasks) do
		table.insert(items, {
			type = "task",
			name = task.name,
			description = task.description,
			line = task.line,
		})
	end

	-- Add workflows
	for _, workflow in ipairs(workflows) do
		table.insert(items, {
			type = "workflow",
			name = workflow.name,
			description = workflow.description,
			line = workflow.line,
		})
	end

	if #items == 0 then
		vim.notify("No tasks or workflows found in current buffer", vim.log.levels.WARN)
		return
	end

	-- Create picker
	pickers
		.new(opts, {
			finder = finders.new_table({
				results = items,
				entry_maker = function(item)
					local icon = item.type == "task" and config.ui.icons.task or config.ui.icons.workflow
					return {
						value = item,
						display = icon .. " " .. item.name .. (item.description and (" - " .. item.description) or ""),
						ordinal = item.name .. " " .. (item.description or ""),
						lnum = item.line,
					}
				end,
			}),
			sorter = conf.generic_sorter(opts),
			previewer = previewers.vim_buffer_cat.new(opts),
			attach_mappings = function(prompt_bufnr, map)
				actions.select_default:replace(function()
					local selection = action_state.get_selected_entry()
					actions.close(prompt_bufnr)

					if selection then
						vim.api.nvim_win_set_cursor(0, { selection.lnum, 0 })
					end
				end)

				return true
			end,
		})
		:find()
end

---Create task previewer
---@return table Previewer
function M.create_task_previewer()
	return previewers.new_buffer_previewer({
		title = "Task Preview",
		define_preview = function(self, entry)
			local bufnr = vim.api.nvim_get_current_buf()
			local lines = vim.api.nvim_buf_get_lines(bufnr, 0, -1, false)

			-- Extract task definition
			local start_line = entry.lnum
			local end_line = start_line

			-- Find end of task (look for :build())
			for i = start_line, #lines do
				if lines[i]:match(":build%s*%(") then
					end_line = i
					break
				end
			end

			local task_lines = {}
			for i = start_line, end_line do
				table.insert(task_lines, lines[i])
			end

			vim.api.nvim_buf_set_lines(self.state.bufnr, 0, -1, false, task_lines)
			vim.api.nvim_buf_set_option(self.state.bufnr, "filetype", "sloth")
		end,
	})
end

---Create workflow previewer
---@return table Previewer
function M.create_workflow_previewer()
	return previewers.new_buffer_previewer({
		title = "Workflow Preview",
		define_preview = function(self, entry)
			local bufnr = vim.api.nvim_get_current_buf()
			local lines = vim.api.nvim_buf_get_lines(bufnr, 0, -1, false)

			-- Extract workflow definition
			local start_line = entry.lnum
			local end_line = utils.find_workflow_end(start_line) or start_line + 10

			local workflow_lines = {}
			for i = start_line, math.min(end_line, #lines) do
				table.insert(workflow_lines, lines[i])
			end

			vim.api.nvim_buf_set_lines(self.state.bufnr, 0, -1, false, workflow_lines)
			vim.api.nvim_buf_set_option(self.state.bufnr, "filetype", "sloth")
		end,
	})
end

return M
