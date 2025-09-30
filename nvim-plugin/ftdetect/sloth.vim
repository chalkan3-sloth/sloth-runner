" Vim filetype detection for Sloth Runner DSL
" Add this to your ~/.config/nvim/ftdetect/sloth.vim

" Detect .sloth files as sloth filetype
autocmd BufRead,BufNewFile *.sloth setfiletype sloth
autocmd BufRead,BufNewFile *.sloth.sloth setfiletype sloth
autocmd BufRead,BufNewFile */sloth-runner/*.sloth setfiletype sloth
autocmd BufRead,BufNewFile */workflows/*.sloth setfiletype sloth
autocmd BufRead,BufNewFile */tasks/*.sloth setfiletype sloth

" Detect common sloth file patterns
autocmd BufRead,BufNewFile *task*.sloth setfiletype sloth
autocmd BufRead,BufNewFile *workflow*.sloth setfiletype sloth

" Legacy .lua files with sloth content (for backward compatibility)
autocmd BufRead,BufNewFile *.lua 
    \ if search('task\s*(.*).*:.*build()', 'nw') || 
    \    search('workflow\.define\s*("', 'nw') || 
    \    search(':command\s*\(', 'nw') || 
    \    search(':description\s*\(', 'nw') |
    \   setfiletype sloth |
    \ endif