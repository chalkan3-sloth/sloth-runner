" Vim filetype detection for Sloth Runner DSL
" Add this to your ~/.config/nvim/ftdetect/sloth.vim

" Detect .sloth files as sloth filetype
autocmd BufRead,BufNewFile *.sloth setfiletype sloth

" Detect common sloth file patterns
autocmd BufRead,BufNewFile *task*.sloth setfiletype sloth
autocmd BufRead,BufNewFile *workflow*.sloth setfiletype sloth

" Detect files in sloth-runner directory structure
autocmd BufRead,BufNewFile */sloth-runner/*.sloth setfiletype sloth
autocmd BufRead,BufNewFile */workflows/*.sloth setfiletype sloth
autocmd BufRead,BufNewFile */tasks/*.sloth setfiletype sloth