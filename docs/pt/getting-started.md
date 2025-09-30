# Início Rápido

Bem-vindo ao Sloth-Runner! Este guia o ajudará a começar a usar a ferramenta rapidamente.

> **📝 Nota Importante:** A partir da versão atual, os arquivos de workflow do Sloth Runner usam a extensão `.sloth` em vez de `.lua`. A sintaxe Lua permanece a mesma - apenas a extensão do arquivo mudou para melhor identificação dos arquivos DSL do Sloth Runner.

## Instalação

Para instalar o `sloth-runner` em seu sistema, você pode usar o script `install.sh` fornecido. Este script detecta automaticamente seu sistema operacional e arquitetura, baixa a versão mais recente do GitHub e coloca o executável `sloth-runner` em `/usr/local/bin`.

```bash
bash <(curl -sL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/install.sh)
```

**Nota:** O script `install.sh` requer privilégios de `sudo` para mover o executável para `/usr/local/bin`.

## Uso Básico

### Gerenciamento de Stacks

```bash
# Criar um novo stack
sloth-runner stack new my-app --description "Stack de deployment da aplicação"

# Executar workflows em stacks
sloth-runner run my-app -f examples/basic_pipeline.sloth

# Listar todos os stacks
sloth-runner stack list

# Ver detalhes do stack
sloth-runner stack show my-app
```

### Execução Direta de Workflow

Para executar um arquivo de workflow diretamente:

```bash
sloth-runner run -f examples/basic_pipeline.sloth
```

Para listar as tarefas em um arquivo:

```bash
sloth-runner list -f examples/basic_pipeline.sloth
```

## Agendador de Tarefas (Novo!)

O Sloth-Runner agora inclui um poderoso agendador de tarefas que permite automatizar a execução de seus fluxos de trabalho em segundo plano usando sintaxe cron. Para mais detalhes sobre como configurar e usar o agendador, consulte a documentação completa em [Agendador de Tarefas](./scheduler.md).

## Próximos Passos

Agora que você tem o Sloth-Runner instalado e funcionando, explore os [Conceitos Essenciais](./core-concepts.md) para entender como definir suas tarefas, ou mergulhe diretamente nos novos [Módulos Built-in](../index.md#módulos-built-in) para automação avançada com Git, Pulumi e Salt.

---
[English](../en/getting-started.md) | [Português](./getting-started.md) | [中文](../zh/getting-started.md)