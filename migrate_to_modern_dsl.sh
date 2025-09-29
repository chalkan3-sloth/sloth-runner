#!/bin/bash

# Script para migrar todos os arquivos para Modern DSL apenas
# Remove completamente TaskDefinition e task_definition

set -e

echo "🚀 Starting complete migration to Modern DSL..."

# Função para converter arquivo individual
convert_file() {
    local file="$1"
    echo "Converting: $file"
    
    # Backup original
    cp "$file" "$file.backup"
    
    # Conversão de TaskDefinition para task() com builder pattern
    sed -i '' 's/TaskDefinition({/local task_name = task("task_name")\
    :description("Task description")\
    :command(function(params, deps)/g' "$file"
    
    # Remove fechamento de TaskDefinition e adiciona :build()
    sed -i '' 's/})/    :build()/g' "$file"
    
    # Converte campos de TaskDefinition para métodos fluent
    sed -i '' 's/name = "\([^"]*\)"/-- Task name: \1/g' "$file"
    sed -i '' 's/description = "\([^"]*\)"/:description("\1")/g' "$file"
    sed -i '' 's/command = /-- Using command function instead/g' "$file"
    sed -i '' 's/depends_on = {\([^}]*\)}/:depends_on({\1})/g' "$file"
    sed -i '' 's/timeout = "\([^"]*\)"/:timeout("\1")/g' "$file"
    sed -i '' 's/retries = \([0-9]*\)/:retries(\1)/g' "$file"
    sed -i '' 's/artifacts = {\([^}]*\)}/:artifacts({\1})/g' "$file"
    
    # Converte hooks
    sed -i '' 's/pre_hook = /-- Using pre_hook function/g' "$file"
    sed -i '' 's/post_hook = /-- Using post_hook function/g' "$file"
    sed -i '' 's/on_success = /:on_success(/g' "$file"
    sed -i '' 's/on_failure = /:on_failure(/g' "$file"
    
    echo "✅ Converted: $file"
}

# Converter todos os arquivos .lua com TaskDefinition
echo "🔍 Finding files with TaskDefinition..."
files_to_convert=$(find examples -name "*.lua" -exec grep -l "TaskDefinition\|task_definition" {} \;)

if [ -z "$files_to_convert" ]; then
    echo "✅ No files found with TaskDefinition - migration may already be complete!"
else
    echo "📝 Found files to convert:"
    echo "$files_to_convert"
    echo ""
    
    # Converter cada arquivo
    while IFS= read -r file; do
        convert_file "$file"
    done <<< "$files_to_convert"
fi

echo ""
echo "✅ Migration completed!"
echo "📁 Backup files created with .backup extension"
echo "🧪 Please test the converted files to ensure they work correctly"