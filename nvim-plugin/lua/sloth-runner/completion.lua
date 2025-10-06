-- ═══════════════════════════════════════════════════════════════════════
-- Sloth Runner - nvim-cmp Completion Source
-- ═══════════════════════════════════════════════════════════════════════

local M = {}

local utils = require("sloth-runner.utils")
local config_module = require("sloth-runner.config")

---@class SlothCompletionItem
---@field label string Completion text
---@field kind number Completion kind
---@field detail string|nil Detail text
---@field documentation string|nil Documentation

---Setup nvim-cmp source
---@param config table Plugin configuration
function M.setup(config)
	local ok, cmp = pcall(require, "cmp")
	if not ok then
		return
	end

	-- Register source
	cmp.register_source("sloth", M.new())

	utils.log("Sloth completion source registered", vim.log.levels.DEBUG)
end

---Create new completion source
---@return table Completion source
function M.new()
	local source = {}

	---Get source name
	---@return string
	function source:get_keyword_pattern()
		return [[\k\+]]
	end

	---Check if source is available
	---@return boolean
	function source:is_available()
		return utils.is_sloth_file()
	end

	---Get completion items
	---@param params table Completion parameters
	---@param callback function Callback function
	function source:complete(params, callback)
		local config = config_module.get()
		local items = {}

		local line = params.context.cursor_before_line
		local cursor_col = params.context.cursor.col

		-- Get completion items based on context
		local completions = M.get_completions(line, cursor_col, config)

		-- Convert to nvim-cmp format
		for _, item in ipairs(completions) do
			table.insert(items, {
				label = item.label,
				kind = item.kind or require("cmp").lsp.CompletionItemKind.Keyword,
				detail = item.detail,
				documentation = config.completion.show_docs and item.documentation or nil,
				insertText = item.insertText or item.label,
			})
		end

		callback({ items = items, isIncomplete = false })
	end

	return source
end

---Get completion items based on context
---@param line string Current line
---@param col number Cursor column
---@param config table Plugin configuration
---@return SlothCompletionItem[] Completion items
function M.get_completions(line, col, config)
	local items = {}
	local cmp = require("cmp")

	-- Check context
	local before_cursor = line:sub(1, col - 1)

	-- Method chaining (after colon)
	if before_cursor:match(":(%w*)$") then
		return M.get_method_completions(config)
	end

	-- Module access (after dot)
	if before_cursor:match("%.(%w*)$") then
		return M.get_module_method_completions(before_cursor, config)
	end

	-- Default: keywords and modules
	return M.get_keyword_completions(config)
end

---Get DSL method completions
---@param config table Plugin configuration
---@return SlothCompletionItem[]
function M.get_method_completions(config)
	local methods = {
		{
			label = "command",
			kind = require("cmp").lsp.CompletionItemKind.Method,
			detail = "Define task command",
			documentation = "Define the command to execute for this task.\n\nExample:\n:command(function(params, deps)\n  -- implementation\n  return true\nend)",
		},
		{
			label = "description",
			kind = require("cmp").lsp.CompletionItemKind.Method,
			detail = "Task description",
			documentation = "Provide a human-readable description of the task.",
		},
		{
			label = "timeout",
			kind = require("cmp").lsp.CompletionItemKind.Method,
			detail = "Set timeout duration",
			documentation = "Set the maximum execution time for this task.",
		},
		{
			label = "retries",
			kind = require("cmp").lsp.CompletionItemKind.Method,
			detail = "Set retry count",
			documentation = "Number of times to retry on failure.",
		},
		{
			label = "depends_on",
			kind = require("cmp").lsp.CompletionItemKind.Method,
			detail = "Define dependencies",
			documentation = "Specify tasks that must complete before this one.",
		},
		{
			label = "on_success",
			kind = require("cmp").lsp.CompletionItemKind.Method,
			detail = "Success callback",
			documentation = "Function to call when task succeeds.",
		},
		{
			label = "on_failure",
			kind = require("cmp").lsp.CompletionItemKind.Method,
			detail = "Failure callback",
			documentation = "Function to call when task fails.",
		},
		{
			label = "build",
			kind = require("cmp").lsp.CompletionItemKind.Method,
			detail = "Build task definition",
			documentation = "Finalize and build the task definition.",
		},
		{
			label = "run_on",
			kind = require("cmp").lsp.CompletionItemKind.Method,
			detail = "Specify execution target",
			documentation = "Define where this task should run (local, remote, etc).",
		},
		{
			label = "agent",
			kind = require("cmp").lsp.CompletionItemKind.Method,
			detail = "Set agent/executor",
			documentation = "Specify which agent should execute this task.",
		},
		{
			label = "tags",
			kind = require("cmp").lsp.CompletionItemKind.Method,
			detail = "Add task tags",
			documentation = "Tag the task for categorization and filtering.",
		},
		{
			label = "artifacts",
			kind = require("cmp").lsp.CompletionItemKind.Method,
			detail = "Define artifacts",
			documentation = "Specify output artifacts produced by this task.",
		},
		{
			label = "condition",
			kind = require("cmp").lsp.CompletionItemKind.Method,
			detail = "Set execution condition",
			documentation = "Define a condition that must be met for execution.",
		},
		{
			label = "schedule",
			kind = require("cmp").lsp.CompletionItemKind.Method,
			detail = "Schedule task",
			documentation = "Set a cron-like schedule for this task.",
		},
	}

	return methods
end

---Get module method completions
---@param line string Current line
---@param config table Plugin configuration
---@return SlothCompletionItem[]
function M.get_module_method_completions(line, config)
	-- Detect which module
	local module = line:match("(%w+)%.%w*$")

	if module == "workflow" then
		return {
			{
				label = "define",
				kind = require("cmp").lsp.CompletionItemKind.Function,
				detail = "Define workflow",
				documentation = "Define a new workflow with tasks and configuration.",
			},
		}
	end

	return {}
end

---Get keyword and module completions
---@param config table Plugin configuration
---@return SlothCompletionItem[]
function M.get_keyword_completions(config)
	local items = {}
	local cmp = require("cmp")

	-- Core keywords
	local keywords = {
		{
			label = "task",
			kind = cmp.lsp.CompletionItemKind.Function,
			detail = "Create task",
			documentation = "Create a new task definition.\n\nExample:\nlocal my_task = task('task-name')\n  :description('My task')\n  :command(function(params, deps)\n    return true\n  end)\n  :build()",
		},
		{
			label = "workflow",
			kind = cmp.lsp.CompletionItemKind.Module,
			detail = "Workflow module",
			documentation = "Access workflow definition functions.",
		},
	}

	-- Add modules
	local modules = {
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
	}

	for _, module in ipairs(modules) do
		table.insert(keywords, {
			label = module,
			kind = cmp.lsp.CompletionItemKind.Module,
			detail = "Sloth module",
			documentation = "Access " .. module .. " module functions and utilities.",
		})
	end

	return keywords
end

return M
