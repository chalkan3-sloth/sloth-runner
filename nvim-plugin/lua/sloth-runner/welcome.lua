-- ═══════════════════════════════════════════════════════════════════════
-- Sloth Runner - Welcome Banner
-- ═══════════════════════════════════════════════════════════════════════

local M = {}

---Show welcome banner when opening .sloth file
function M.show()
	local config = require("sloth-runner.config").get()

	-- Don't show if disabled in config
	if config.ui.show_welcome == false then
		return
	end

	-- ASCII art sloth banner
	local banner = {
		"",
		"        🦥 Sloth Runner DSL",
		"",
		"    ⚡ Ready to automate workflows!",
		"",
	}

	-- Alternative: larger ASCII art version
	local banner_large = {
		"",
		"     ╔══════════════════════════════════╗",
		"     ║     🦥  SLOTH RUNNER DSL  🦥     ║",
		"     ╚══════════════════════════════════╝",
		"",
		"        ⚡ Workflow automation made easy",
		"",
		"     💡 Tip: Use <leader>sr to run workflow",
		"",
	}

	-- Simple notification version (default)
	local banner_simple = {
		"🦥 Sloth Runner DSL - Ready to automate!",
	}

	-- Choose banner style based on config
	local style = config.ui.welcome_style or "notification"
	local lines = style == "large" and banner_large or style == "banner" and banner or banner_simple

	if style == "notification" then
		-- Simple notification
		vim.notify("🦥 Sloth Runner DSL", vim.log.levels.INFO, {
			title = "Welcome",
			timeout = 2000,
		})
	elseif style == "float" then
		-- Float window with banner
		M.show_float(lines)
	else
		-- Echo to command line
		for _, line in ipairs(lines) do
			vim.cmd('echohl String | echo "' .. line .. '" | echohl None')
		end
	end
end

---Show banner in floating window
---@param lines string[] Banner lines
function M.show_float(lines)
	local utils = require("sloth-runner.utils")

	-- Create float with banner
	local bufnr, winnr = utils.create_float({
		lines = lines,
		title = "",
		border = "rounded",
		width = 50,
		height = #lines,
		mappings = {
			q = "close",
			["<Esc>"] = "close",
			["<CR>"] = "close",
		},
	})

	-- Auto-close after 3 seconds
	vim.defer_fn(function()
		if vim.api.nvim_win_is_valid(winnr) then
			vim.api.nvim_win_close(winnr, true)
		end
	end, 3000)
end

---Show sloth animation (fun easter egg)
function M.animate()
	local frames = {
		"🦥",
		"🦥💤",
		"🦥💤💤",
		"🦥⚡",
		"🦥✨",
	}

	local frame = 1
	local timer = vim.loop.new_timer()

	timer:start(
		0,
		200,
		vim.schedule_wrap(function()
			if frame > #frames then
				timer:stop()
				timer:close()
				vim.notify("🦥 Ready to work!", vim.log.levels.INFO)
				return
			end

			vim.notify(frames[frame], vim.log.levels.INFO, {
				title = "Sloth Runner",
				timeout = 200,
				replace = frame > 1,
			})

			frame = frame + 1
		end)
	)
end

return M
