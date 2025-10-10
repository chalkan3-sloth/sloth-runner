# Conceitos Essenciais - Modern DSL

Este documento explica os conceitos fundamentais do `sloth-runner` usando a **Modern DSL**, ajudando você a entender como definir e orquestrar fluxos de trabalho complexos com a nova API fluente.

---

## Visão Geral da Modern DSL

A Modern DSL substitui a abordagem legada `Modern DSLs` por uma API mais intuitiva e fluente para definir fluxos de trabalho. Em vez de grandes estruturas de tabela, você agora usa métodos encadeáveis para construir tarefas e definir fluxos de trabalho de forma declarativa.

```lua
-- meu_pipeline.sloth - Modern DSL
local minha_tarefa = task("nome_da_tarefa")
    :description("Descrição da tarefa")
    :command(function(this, params)
        -- Lógica da tarefa
        return true, "Tarefa concluída"
    end)
    :build()

workflow.define("nome_do_workflow")
    :description("Descrição do workflow - Modern DSL")
    :version("1.0.0")
    :tasks({ minha_tarefa })
```

---

## Definição de Tarefa com Modern DSL

As tarefas agora são definidas usando a função `task()` e métodos da API fluente:

### Estrutura Básica de Tarefa

```lua
local minha_tarefa = task("nome_da_tarefa")
    :description("O que esta tarefa faz")
    :command(function(params, deps)
        -- Lógica da tarefa aqui
        return true, "Mensagem de sucesso", { dados_de_saida = "valor" }
    end)
    :timeout("5m")
    :retries(3, "exponential")
    :build()
```

### Métodos do Task Builder

**Propriedades Principais:**
*   `:description(string)` - Descrição legível da tarefa
*   `:command(function|string)` - Lógica de execução da tarefa
*   `:timeout(string)` - Tempo máximo de execução (ex: "10s", "5m", "1h")
*   `:retries(number, strategy)` - Configuração de retry com estratégia ("exponential", "linear", "fixed")
*   `:depends_on(array)` - Array de nomes de tarefas das quais esta tarefa depende

**Recursos Avançados:**
*   `:async(boolean)` - Habilitar execução assíncrona
*   `:artifacts(array)` - Arquivos para salvar após execução bem-sucedida
*   `:consumes(array)` - Artefatos de outras tarefas para usar
*   `:run_if(function|string)` - Lógica de execução condicional
*   `:abort_if(function|string)` - Condição para abortar todo o workflow

**Hooks de Ciclo de Vida:**
*   `:on_success(function)` - Executar quando a tarefa for bem-sucedida
*   `:on_failure(function)` - Executar quando a tarefa falhar
*   `:on_timeout(function)` - Executar quando a tarefa atingir timeout
*   `:pre_hook(function)` - Executar antes do comando principal
*   `:post_hook(function)` - Executar após o comando principal

**Exemplo:**
```lua
-- Workflow que gerencia seu próprio diretório temporário
local setup_task = task("setup")
    :description("Setup inicial")
    :command(function(this, params)
        log.info("Configurando ambiente...")
        return true, "Setup completo"
    end)
    :build()

workflow.define("meu_grupo")
    :description("Um grupo que gerencia seu próprio diretório temporário")
    :version("1.0.0")
    :tasks({setup_task})
    :config({
        create_workdir_before_run = true,
        clean_workdir_after_run = "on_success"
    })
    :on_complete(function(success, results)
        if not success then
            log.warn("O workflow falhou. O diretório de trabalho será mantido para depuração.")
        end
    end)
```

---

## Tarefas Individuais

Uma tarefa é uma única unidade de trabalho. É definida como uma tabela com várias propriedades disponíveis para controlar seu comportamento.

### Propriedades Básicas

*   `name` (string): O nome único da tarefa dentro de seu grupo.
*   `description` (string): Uma breve descrição do que a tarefa faz.
*   `command` (string ou função): A ação principal da tarefa.
    *   **Como string:** É executada como um comando de shell.
    *   **Como função:** A função Lua é executada. Ela recebe dois argumentos: `params` (uma tabela com seus parâmetros) e `deps` (uma tabela contendo os outputs de suas dependências). A função deve retornar:
        1.  `booleano`: `true` para sucesso, `false` para falha.
        2.  `string`: Uma mensagem descrevendo o resultado.
        3.  `tabela` (opcional): Uma tabela de outputs da qual outras tarefas podem depender.

### Dependência e Fluxo de Execução

*   `depends_on` (string ou tabela): Uma lista de nomes de tarefas que devem ser concluídas com sucesso antes que esta tarefa possa ser executada.
*   `next_if_fail` (string ou tabela): Uma lista de nomes de tarefas a serem executadas *apenas se* esta tarefa falhar. Útil para tarefas de limpeza ou notificação.
*   `async` (booleano): Se `true`, a tarefa é executada em segundo plano, e o runner não espera que ela termine para iniciar a próxima tarefa na ordem de execução.

### Tratamento de Erros e Robustez

*   `retries` (número): O número de vezes que uma tarefa será tentada novamente se falhar. O padrão é `0`.
*   `timeout` (string): Uma duração (ex: `"10s"`, `"1m"`) após a qual a tarefa será encerrada se ainda estiver em execução.

### Execução Condicional

*   `run_if` (string ou função): A tarefa será pulada a menos que esta condição seja atendida.
    *   **Como string:** Um comando de shell. Um código de saída `0` significa que a condição foi atendida.
    *   **Como função:** Uma função Lua que retorna `true` se a tarefa deve ser executada.
*   `abort_if` (string ou função): Todo o fluxo de trabalho será abortado se esta condição for atendida.
    *   **Como string:** Um comando de shell. Um código de saída `0` significa abortar.
    *   **Como função:** Uma função Lua que retorna `true` para abortar.

### Hooks de Ciclo de Vida

*   `pre_exec` (função): Uma função Lua que é executada *antes* do `command` principal.
*   `post_exec` (função): Uma função Lua que é executada *após* o `command` principal ter sido concluído com sucesso.

### Reutilização

*   `uses` (tabela): Especifica uma tarefa pré-definida de outro arquivo (carregado via `import`) para usar como base. A definição da tarefa atual pode então sobrescrever propriedades como `params` ou `description`.
*   `params` (tabela): Um dicionário de pares chave-valor que podem ser passados para a função `command` da tarefa.
*   `artifacts` (string ou tabela): Um padrão de arquivo (glob) ou uma lista de padrões que especificam quais arquivos do `workdir` da tarefa devem ser salvos como artefatos após uma execução bem-sucedida.
*   `consumes` (string ou tabela): O nome de um artefato (ou uma lista de nomes) de uma tarefa anterior que deve ser copiado para o `workdir` desta tarefa antes que ela seja executada.

---

## Gerenciamento de Artefatos

O Sloth-Runner permite que as tarefas compartilhem arquivos entre si através de um mecanismo de artefatos. Uma tarefa pode "produzir" um ou mais arquivos como artefatos, e tarefas subsequentes podem "consumir" esses artefatos.

Isso é útil para pipelines de CI/CD, onde uma etapa de compilação pode gerar um binário (o artefato), que é então usado por uma etapa de teste ou de implantação.

### Como Funciona

1.  **Produzindo Artefatos:** Adicione a chave `artifacts` à sua definição de tarefa. O valor pode ser um único padrão de arquivo (ex: `"report.txt"`) ou uma lista (ex: `{"*.log", "app.bin"}`). Após a tarefa ser executada com sucesso, o runner procurará por arquivos no `workdir` da tarefa que correspondam a esses padrões e os copiará para um armazenamento de artefatos compartilhado para a pipeline.

2.  **Consumindo Artefatos:** Adicione a chave `consumes` à definição de outra tarefa (que normalmente `depends_on` da tarefa produtora). O valor deve ser o nome do arquivo do artefato que você deseja usar (ex: `"report.txt"`). Antes que esta tarefa seja executada, o runner copiará o artefato nomeado do armazenamento compartilhado para o `workdir` desta tarefa, tornando-o disponível para o `command`.

### Exemplo de Artefatos

```lua
local build_task = task("build")
    :description("Cria um binário e o declara como um artefato")
    :command(function(this, params)
        exec.run("echo 'conteudo_binario' > app.bin")
        return true, "Binário criado"
    end)
    :artifacts({"app.bin"})
    :build()

local test_task = task("test")
    :description("Consome o binário para executar testes")
    :depends_on({"build"})
    :consumes({"app.bin"})
    :command(function(this, params)
        -- Neste ponto, 'app.bin' existe no workdir desta tarefa
        local content = exec.run("cat app.bin")
        if content:find("conteudo_binario") then
            log.info("Artefato consumido com sucesso!")
            return true, "Artefato validado"
        else
            return false, "Conteúdo do artefato incorreto!"
        end
    end)
    :build()

workflow.define("ci_pipeline")
    :description("Demonstra o uso de artefatos")
    :version("1.0.0")
    :tasks({build_task, test_task})
    :config({
        timeout = "10m",
        create_workdir_before_run = true
    })
```

---

## Funções Globais

O `sloth-runner` fornece funções globais no ambiente Lua para ajudar a orquestrar os fluxos de trabalho.

### `import(path)`

Carrega outro arquivo Lua e retorna o valor que ele retorna. Este é o principal mecanismo para criar módulos de tarefas reutilizáveis. O caminho é relativo ao arquivo que chama `import`.

**Exemplo (`reusable_tasks.sloth`):**
```lua
-- Importa um módulo que retorna definições de tarefas
local docker_tasks = import("shared/docker.sloth")

-- Usa a tarefa importada com parâmetros personalizados
local build_app = docker_tasks.build_image("my-app")
    :description("Build da imagem Docker my-app")
    :timeout("10m")
    :build()

workflow.define("main")
    :description("Workflow principal usando tarefas reutilizáveis")
    :version("1.0.0")
    :tasks({build_app})
```

### `parallel(tasks)`

Executa uma lista de tarefas concorrentemente e espera que todas terminem.

*   `tasks` (tabela): Uma lista de tabelas de tarefas para executar em paralelo.

**Exemplo:**
```lua
command = function()
  log.info("Iniciando 3 tarefas em paralelo...")
  local results, err = parallel({
    { name = "short_task", command = "sleep 1" },
    { name = "medium_task", command = "sleep 2" },
    { name = "long_task", command = "sleep 3" }
  })
  if err then
    return false, "Execução paralela falhou"
  end
  return true, "Todas as tarefas paralelas terminaram."
end
```

### `export(table)`

Exporta dados de qualquer ponto de um script para a CLI. Quando a flag `--return` é usada, todas as tabelas exportadas são mescladas com o output da tarefa final em um único objeto JSON.

*   `table`: Uma tabela Lua a ser exportada.

**Exemplo:**
```lua
command = function()
  export({ valor_importante = "dado do meio da tarefa" })
  return true, "Tarefa concluída", { output_final = "algum resultado" }
end
```
Executar com `--return` produziria:
```json
{
  "valor_importante": "dado do meio da tarefa",
  "output_final": "algum resultado"
}
```