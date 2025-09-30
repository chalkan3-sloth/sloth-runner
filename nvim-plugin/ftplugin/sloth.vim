" Vim filetype plugin for Sloth Runner DSL
" ~/.config/nvim/ftplugin/sloth.vim

if exists("b:did_ftplugin")
  finish
endif
let b:did_ftplugin = 1

" Use Lua settings as base
runtime! ftplugin/lua.vim

" Sloth-specific settings
setlocal commentstring=--\ %s
setlocal comments=:--
setlocal suffixesadd=.sloth,.lua

" Enhanced indentation for DSL chaining
setlocal indentexpr=GetSlothIndent()
setlocal indentkeys+=:,0),0},0],0=end,0=then,0=else,0=elseif,0=until

" Keywords for word boundaries and completion
setlocal iskeyword+=:

" Set up folding
setlocal foldmethod=expr
setlocal foldexpr=SlothFoldExpr(v:lnum)
setlocal foldtext=SlothFoldText()

" Define abbreviations for common DSL patterns
if has("autocmd") && exists("+omnifunc")
  if &omnifunc == ""
    setlocal omnifunc=SlothComplete
  endif
endif

" Abbreviations for quick DSL writing
iabbrev <buffer> _task local task_name = task("")<CR>:description("")<CR>:command(function(params, deps)<CR>-- TODO: implement<CR>return true<CR>end)<CR>:build()<Esc>7k$2h
iabbrev <buffer> _workflow workflow.define("", {<CR>description = "",<CR>version = "1.0.0",<CR>tasks = {<CR>-- tasks here<CR>}<CR>})<Esc>5k$h
iabbrev <buffer> _cmd :command(function(params, deps)<CR>-- TODO: implement<CR>return true<CR>end)<Esc>2k$

" Key mappings for DSL development
nnoremap <buffer> <leader>sr :!sloth-runner run -f %<CR>
nnoremap <buffer> <leader>sl :!sloth-runner list -f %<CR>
nnoremap <buffer> <leader>st :!sloth-runner run -f % --dry-run<CR>

" Text objects for DSL blocks
" Select task block: vit (visual in task)
onoremap <buffer> it :<C-u>call <SID>SelectTaskBlock()<CR>
vnoremap <buffer> it :<C-u>call <SID>SelectTaskBlock()<CR>

" Select workflow block: viw (visual in workflow)
onoremap <buffer> iw :<C-u>call <SID>SelectWorkflowBlock()<CR>
vnoremap <buffer> iw :<C-u>call <SID>SelectWorkflowBlock()<CR>

" Function to provide intelligent indentation
function! GetSlothIndent()
  let lnum = prevnonblank(v:lnum - 1)
  if lnum == 0
    return 0
  endif

  let line = getline(lnum)
  let ind = indent(lnum)

  " Increase indent after certain patterns
  if line =~ '\v(task\s*\(|workflow\.define\s*\(|function\s*\(|:\w+\s*\(|\{\s*$)'
    let ind += &shiftwidth
  endif

  " Increase indent for DSL chaining
  if line =~ '\v:\w+\s*\('
    let ind += &shiftwidth
  endif

  " Decrease indent for closing patterns
  let cline = getline(v:lnum)
  if cline =~ '\v^\s*(end|\}|\)|:build\(\s*\))'
    let ind -= &shiftwidth
  endif

  return ind
endfunction

" Folding expression
function! SlothFoldExpr(lnum)
  let line = getline(a:lnum)
  
  " Start fold for task definitions
  if line =~ '\v^\s*local\s+\w+\s*\=\s*task\s*\('
    return '>1'
  endif
  
  " Start fold for workflow definitions
  if line =~ '\v^\s*workflow\.define\s*\('
    return '>1'
  endif
  
  " Start fold for function definitions
  if line =~ '\v^\s*function\s+'
    return '>1'
  endif
  
  " End fold
  if line =~ '\v^\s*(end|\})\s*$'
    return '<1'
  endif
  
  return '='
endfunction

" Custom fold text
function! SlothFoldText()
  let line = getline(v:foldstart)
  let nucolwidth = &fdc + &number * &numberwidth
  let windowwidth = winwidth(0) - nucolwidth - 3
  let foldedlinecount = v:foldend - v:foldstart
  
  " Extract meaningful text from the fold
  let onetab = strpart('          ', 0, &tabstop)
  let line = substitute(line, '\t', onetab, 'g')
  
  if line =~ 'task\s*('
    let line = substitute(line, '\v.*task\s*\(\s*"([^"]*)".*', '📋 Task: \1', '')
  elseif line =~ 'workflow\.define'
    let line = substitute(line, '\v.*workflow\.define\s*\(\s*"([^"]*)".*', '🔄 Workflow: \1', '')
  elseif line =~ 'function'
    let line = substitute(line, '\v.*function\s+(\w+).*', '⚡ Function: \1', '')
  endif
  
  let line = line . ' (' . foldedlinecount . ' lines) '
  let fillcharcount = windowwidth - len(line)
  return line . repeat('⋯', fillcharcount)
endfunction

" Text object functions
function! s:SelectTaskBlock()
  let start = search('\v^\s*local\s+\w+\s*\=\s*task\s*\(', 'bcW')
  if start == 0
    return
  endif
  
  normal! V
  let end = search('\v:build\s*\(\s*\)', 'W')
  if end == 0
    let end = search('\v^\s*$', 'W')
  endif
endfunction

function! s:SelectWorkflowBlock()
  let start = search('\v^\s*workflow\.define\s*\(', 'bcW')
  if start == 0
    return
  endif
  
  normal! V
  " Find matching closing parenthesis/brace
  let end = searchpair('\v\{', '', '\v\}', 'W')
  if end == 0
    let end = search('\v\)\s*$', 'W')
  endif
endfunction

" Simple completion function
function! SlothComplete(findstart, base)
  if a:findstart
    " Find the start of the current word
    let line = getline('.')
    let start = col('.') - 1
    while start > 0 && line[start - 1] =~ '\v[a-zA-Z0-9_:]'
      let start -= 1
    endwhile
    return start
  else
    " Return completion matches
    let completions = []
    
    " DSL methods
    let methods = ['command', 'description', 'timeout', 'retries', 'depends_on', 
                  \'on_success', 'on_failure', 'build', 'run_on', 'agent',
                  \'tags', 'artifacts', 'condition', 'schedule', 'retry_count',
                  \'backoff_strategy', 'circuit_breaker']
    
    " Modules
    let modules = ['exec', 'fs', 'net', 'data', 'log', 'state', 'metrics',
                  \'aws', 'gcp', 'azure', 'docker', 'kubernetes', 'terraform', 'pulumi']
    
    " Functions
    let functions = ['task', 'workflow.define', 'require', 'print', 'type']
    
    for item in methods + modules + functions
      if item =~ '^' . a:base
        call add(completions, item)
      endif
    endfor
    
    return completions
  endif
endfunction

" Enable spell checking in comments and strings
syntax spell toplevel

" Set up auto-commands for this buffer
augroup SlothRunner
  autocmd! * <buffer>
  
  " Auto-format on save (disabled by default, enable with g:sloth_format_on_save = 1)
  if exists('g:sloth_format_on_save') && g:sloth_format_on_save
    autocmd BufWritePre <buffer> call s:FormatSlothFile()
  endif
  
  " Highlight matching DSL constructs
  autocmd CursorMoved <buffer> call s:HighlightMatchingConstruct()
augroup END

function! s:FormatSlothFile()
  " Only format .sloth files and only if explicitly enabled
  if expand('%:e') !=# 'sloth' || !exists('g:sloth_format_on_save') || !g:sloth_format_on_save
    return
  endif
  
  " Check if buffer is modified and file exists
  if !&modified || !filereadable(expand('%'))
    return
  endif
  
  " Simple indentation fix - safer than external commands
  let save_pos = getpos('.')
  try
    normal! gg=G
  catch
    " Ignore any errors during formatting
  finally
    call setpos('.', save_pos)
  endtry
endfunction

function! s:HighlightMatchingConstruct()
  " Clear previous highlights
  if exists('b:sloth_match_id')
    call matchdelete(b:sloth_match_id)
    unlet b:sloth_match_id
  endif
  
  " Highlight matching task/workflow constructs
  let line = getline('.')
  if line =~ '\v:build\s*\(\s*\)'
    let b:sloth_match_id = matchadd('MatchParen', '\v^\s*local\s+\w+\s*\=\s*task\s*\(.*$')
  elseif line =~ '\v^\s*local\s+\w+\s*\=\s*task\s*\('
    let b:sloth_match_id = matchadd('MatchParen', '\v:build\s*\(\s*\)')
  endif
endfunction

let b:undo_ftplugin = "setlocal commentstring< comments< suffixesadd< indentexpr< indentkeys< iskeyword< foldmethod< foldexpr< foldtext< omnifunc<"