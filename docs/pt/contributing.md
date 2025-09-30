# 🤝 Contribuindo para o Sloth Runner

**Obrigado pelo seu interesse em contribuir para o Sloth Runner!**

Acolhemos contribuições de desenvolvedores de todos os níveis de habilidade. Seja corrigindo bugs, adicionando recursos, melhorando a documentação ou criando plugins, sua ajuda torna o Sloth Runner melhor para todos.

## 🚀 Início Rápido

### Pré-requisitos

- **Go 1.21+** para desenvolvimento principal
- **Node.js 18+** para desenvolvimento de UI  
- **Lua 5.4+** para desenvolvimento DSL
- **Git** para controle de versão

### Configuração de Desenvolvimento

```bash
# Clonar o repositório
git clone https://github.com/chalkan3-sloth/sloth-runner.git
cd sloth-runner

# Instalar dependências
go mod download
npm install  # para componentes UI

# Executar testes
make test

# Construir o projeto
make build
```

## 📋 Formas de Contribuir

### 🐛 Relatórios de Bug

Encontrou um bug? Por favor, nos ajude a corrigi-lo:

1. **Pesquise issues existentes** para evitar duplicatas
2. **Use nosso template de relatório de bug** com:
   - Versão do Sloth Runner
   - Sistema operacional
   - Passos para reproduzir
   - Comportamento esperado vs real
   - Logs de erro (se houver)

### 💡 Solicitações de Recurso

Tem uma ideia para melhoria?

1. **Verifique o roadmap** para recursos planejados
2. **Abra uma solicitação de recurso** com:
   - Descrição clara do recurso
   - Casos de uso e benefícios
   - Possível abordagem de implementação

### 🔧 Contribuições de Código

Pronto para programar? Aqui está como:

1. **Faça fork do repositório**
2. **Crie uma branch de recurso** (`git checkout -b feature/recurso-incrivel`)
3. **Faça suas alterações** seguindo nossos padrões de código
4. **Adicione testes** para nova funcionalidade
5. **Atualize documentação** se necessário
6. **Commit com mensagens claras**
7. **Push e crie um Pull Request**

### 📚 Documentação

Ajude a melhorar nossa documentação:

- Corrija erros de digitação e explicações confusas
- Adicione exemplos e tutoriais
- Traduza conteúdo para outros idiomas
- Atualize documentação da API

### 🔌 Desenvolvimento de Plugin

Crie plugins para a comunidade:

- Siga nosso [Guia de Desenvolvimento de Plugin](plugin-development.md)
- Submeta ao registro de plugins
- Mantenha compatibilidade com versões principais

## 📐 Diretrizes de Desenvolvimento

### Estilo de Código

#### Código Go

Siga convenções padrão do Go:

```go
// Bom: Nomes de função claros e comentários
func ProcessWorkflowTasks(ctx context.Context, workflow *Workflow) error {
    if workflow == nil {
        return fmt.Errorf("workflow não pode ser nil")
    }
    
    for _, task := range workflow.Tasks {
        if err := processTask(ctx, task); err != nil {
            return fmt.Errorf("falhou ao processar tarefa %s: %w", task.ID, err)
        }
    }
    
    return nil
}
```

#### DSL Lua

Mantenha código DSL limpo e legível:

```lua
-- Bom: Definição clara de tarefa com encadeamento adequado
local tarefa_deploy = task("deploy_aplicacao")
    :description("Fazer deploy da aplicação para produção")
    :command(function(params, deps)
        local resultado = exec.run("kubectl apply -f deployment.yaml")
        if not resultado.success then
            log.error("Deploy falhou: " .. resultado.stderr)
            return false
        end
        return true
    end)
    :timeout(300)
    :retries(3)
    :build()
```

### Padrões de Teste

#### Testes Unitários

Escreva testes para toda nova funcionalidade:

```go
func TestProcessWorkflowTasks(t *testing.T) {
    tests := []struct {
        name     string
        workflow *Workflow
        wantErr  bool
    }{
        {
            name:     "workflow nil deve retornar erro",
            workflow: nil,
            wantErr:  true,
        },
        {
            name: "workflow válido deve processar com sucesso",
            workflow: &Workflow{
                Tasks: []*Task{{ID: "test-task"}},
            },
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ProcessWorkflowTasks(context.Background(), tt.workflow)
            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessWorkflowTasks() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Padrões de Documentação

- **Mantenha simples** - Use linguagem clara e concisa
- **Inclua exemplos** - Mostre, não apenas diga
- **Atualize com mudanças** - Mantenha docs sincronizados com código
- **Teste exemplos** - Garanta que todos os exemplos de código funcionem

## 🔄 Processo de Pull Request

### Antes de Submeter

- [ ] **Execute testes** - `make test`
- [ ] **Execute linting** - `make lint`
- [ ] **Atualize docs** - Se adicionando/alterando recursos
- [ ] **Adicione entrada no changelog** - Em `CHANGELOG.md`
- [ ] **Verifique compatibilidade** - Com recursos existentes

### Template de PR

Use nosso template de pull request:

```markdown
## Descrição
Breve descrição das alterações

## Tipo de Mudança
- [ ] Correção de bug
- [ ] Novo recurso
- [ ] Mudança disruptiva
- [ ] Atualização de documentação

## Testes
- [ ] Testes unitários adicionados/atualizados
- [ ] Testes de integração passam
- [ ] Teste manual completado

## Checklist
- [ ] Código segue diretrizes de estilo
- [ ] Documentação atualizada
- [ ] Changelog atualizado
```

## 🏗️ Estrutura do Projeto

Compreendendo a base de código:

```
sloth-runner/
├── cmd/                    # Comandos CLI
├── internal/              # Pacotes internos
│   ├── core/             # Lógica de negócio principal
│   ├── dsl/              # Implementação DSL
│   ├── execution/        # Motor de execução de tarefas
│   └── plugins/          # Sistema de plugins
├── pkg/                   # Pacotes públicos
├── plugins/              # Plugins integrados
├── docs/                 # Documentação
├── web/                  # Componentes de UI web
└── examples/             # Workflows de exemplo
```

## 🎯 Áreas de Contribuição

### Alta Prioridade

- **🐛 Correções de bug** - Sempre bem-vindas
- **📈 Melhorias de performance** - Oportunidades de otimização
- **🧪 Cobertura de teste** - Aumentar cobertura de teste
- **📚 Documentação** - Manter docs abrangentes

### Média Prioridade

- **✨ Novos recursos** - Seguindo prioridades do roadmap
- **🔌 Ecossistema de plugin** - Mais plugins e integrações
- **🎨 Melhorias de UI** - Melhor experiência do usuário

## 🏆 Reconhecimento

Contribuidores são reconhecidos em:

- **CONTRIBUTORS.md** - Todos os contribuidores listados
- **Notas de release** - Contribuições importantes destacadas
- **Showcase da comunidade** - Contribuições em destaque
- **Badges de contribuidor** - Reconhecimento no perfil GitHub

## 📞 Obtendo Ajuda

### Questões de Desenvolvimento

- **💬 Discord** - canal `#development`
- **📧 Lista de Email** - dev@sloth-runner.io
- **📖 Wiki** - Guias de desenvolvimento e FAQs

### Mentoria

Novo em open source? Oferecemos mentoria:

- **👥 Pareamento de mentor** - Pareado com contribuidores experientes
- **📚 Recursos de aprendizado** - Materiais de aprendizado curados
- **🎯 Contribuições guiadas** - Issues amigáveis para iniciantes

## 📜 Código de Conduta

Estamos comprometidos em fornecer um ambiente acolhedor e inclusivo. Por favor, leia nosso [Código de Conduta](https://github.com/chalkan3-sloth/sloth-runner/blob/main/CODE_OF_CONDUCT.md).

### Nossos Padrões

- **🤝 Seja respeitoso** - Trate todos com respeito
- **💡 Seja construtivo** - Forneça feedback útil
- **🌍 Seja inclusivo** - Acolha perspectivas diversas
- **📚 Seja paciente** - Ajude outros a aprender e crescer

---

**Pronto para contribuir?**

Comece explorando nossas [Good First Issues](https://github.com/chalkan3-sloth/sloth-runner/labels/good%20first%20issue) ou junte-se à nossa [comunidade Discord](https://discord.gg/sloth-runner) para se apresentar!

Obrigado por ajudar a tornar o Sloth Runner melhor! 🦥✨