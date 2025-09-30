" Vim filetype detection for Sloth Runner DSL
" Add this to your ~/.config/nvim/ftdetect/sloth.vim

" Detect .lua files in specific directories as sloth files
autocmd BufRead,BufNewFile *.sloth.lua setfiletype sloth
autocmd BufRead,BufNewFile */sloth-runner/*.lua setfiletype sloth
autocmd BufRead,BufNewFile */workflows/*.lua setfiletype sloth
autocmd BufRead,BufNewFile */tasks/*.lua setfiletype sloth

" Detect common sloth file patterns
autocmd BufRead,BufNewFile *task*.lua if search('task\s*("', 'nw') | setfiletype sloth | endif
autocmd BufRead,BufNewFile *workflow*.lua if search('workflow\.define\s*("', 'nw') | setfiletype sloth | endif

" Detect by content - if file contains DSL keywords
autocmd BufRead,BufNewFile *.lua 
    \ if search('task\s*(.*).*:.*build()', 'nw') || 
    \    search('workflow\.define\s*("', 'nw') || 
    \    search(':command\s*\(', 'nw') || 
    \    search(':description\s*\(', 'nw') |
    \   setfiletype sloth |
    \ endif