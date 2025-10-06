" ═══════════════════════════════════════════════════════════════════════
" Sloth Runner - Neovim Plugin Initialization
" ═══════════════════════════════════════════════════════════════════════

if exists('g:loaded_sloth_runner')
  finish
endif
let g:loaded_sloth_runner = 1

" Save user's cpoptions
let s:save_cpo = &cpo
set cpo&vim

" Auto-initialize plugin if not manually initialized
" This ensures the plugin works even without explicit setup() call
augroup SlothRunnerInit
  autocmd!
  autocmd FileType sloth lua require('sloth-runner').setup()
augroup END

" Restore cpoptions
let &cpo = s:save_cpo
unlet s:save_cpo
