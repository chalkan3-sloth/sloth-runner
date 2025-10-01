" Vim syntax file for Sloth Runner DSL
" Language: Sloth Runner DSL (Lua-based)
" Maintainer: Sloth Runner Team
" Latest Revision: 2025-09-30

if exists("b:current_syntax")
  finish
endif

" Inherit from Lua syntax
runtime! syntax/lua.vim

" Clear any existing syntax
syntax clear

" Comments
syntax match slothComment "--.*$"
syntax region slothComment start="--\[\[" end="\]\]" contains=@Spell

" Strings with template interpolation
syntax region slothString start='"' end='"' skip='\\"' contains=slothInterpolation
syntax region slothString start="'" end="'" skip="\\'"
syntax region slothString start='\[=*\[' end='\]=*\]' contains=@Spell

" Template interpolation
syntax region slothInterpolation contained start='\${' end='}' contains=ALL

" Numbers
syntax match slothNumber '\v<\d+>'
syntax match slothNumber '\v<\d+\.\d+>'
syntax match slothNumber '\v<0x\x+>'
syntax match slothNumber '\v<\d+[eE][\+\-]?\d+>'

" Booleans and nil
syntax keyword slothBoolean true false nil

" Core Lua keywords
syntax keyword slothKeyword local function end if then else elseif while for do repeat until in break return and or not

" Sloth Runner DSL specific keywords
syntax keyword slothDSLKeyword task workflow define
syntax keyword slothDSLKeyword command description timeout retries depends_on
syntax keyword slothDSLKeyword on_success on_failure build run_on agent
syntax keyword slothDSLKeyword tags artifacts condition schedule
syntax keyword slothDSLKeyword retry_count backoff_strategy circuit_breaker

" Built-in modules
syntax keyword slothModule exec fs net data log state metrics
syntax keyword slothModule aws gcp azure digitalocean
syntax keyword slothModule docker kubernetes terraform pulumi
syntax keyword slothModule git notification crypto utils

" DSL chaining methods (highlighted differently)
syntax match slothMethod '\v:\w+'

" Function calls
syntax match slothFunction '\v<\w+\ze\s*\('

" Constants (uppercase identifiers)
syntax match slothConstant '\v<[A-Z][A-Z0-9_]*>'

" Operators
syntax match slothOperator '\v[\+\-\*/%\^#]'
syntax match slothOperator '\v[\=\<\>!\~]'
syntax match slothOperator '\v\.\.'

" Delimiters
syntax match slothDelimiter '\v[\(\)\[\]\{\}]'
syntax match slothDelimiter '\v[,;:]'

" Special workflow structure
syntax region slothWorkflowBlock start='workflow\.define\s*(' end=')' contains=ALL fold
syntax region slothTaskBlock start='task\s*(' end='\.build()' contains=ALL fold

" Environment variables
syntax match slothEnvVar '\v\$\{[A-Z_][A-Z0-9_]*\}'
syntax match slothEnvVar '\v\$[A-Z_][A-Z0-9_]*'

" File paths
syntax match slothPath '\v"[/~][\w/\-\.]*"'
syntax match slothPath "\v'[/~][\\w/\\-\\.]*'"

" Define highlight groups
highlight default link slothComment     Comment
highlight default link slothString      String
highlight default link slothInterpolation Special
highlight default link slothNumber      Number
highlight default link slothBoolean     Boolean
highlight default link slothKeyword     Keyword
highlight default link slothDSLKeyword  Statement
highlight default link slothModule      Include
highlight default link slothMethod      Function
highlight default link slothFunction    Function
highlight default link slothConstant    Constant
highlight default link slothOperator    Operator
highlight default link slothDelimiter   Delimiter
highlight default link slothEnvVar      PreProc
highlight default link slothPath        String

" Enhanced highlights for modern terminals
if has('gui_running') || &t_Co >= 256
  highlight slothDSLKeyword  guifg=#569cd6 ctermfg=75  gui=bold
  highlight slothModule      guifg=#c586c0 ctermfg=176 gui=bold
  highlight slothMethod      guifg=#f9e79f ctermfg=222 gui=bold
  highlight slothFunction    guifg=#dcdcaa ctermfg=187
  highlight slothEnvVar      guifg=#ff6b6b ctermfg=203 gui=bold
  highlight slothPath        guifg=#98d8c8 ctermfg=116
endif

" Folding
syntax region slothFold start='{' end='}' transparent fold
syntax region slothFold start='function' end='end' transparent fold
syntax region slothFold start='workflow\.define' end='})' transparent fold

" Set folding
setlocal foldmethod=syntax
setlocal foldlevel=1

let b:current_syntax = "sloth"