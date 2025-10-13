# Goroutine Channels Documentation

## Overview

O módulo `goroutine` foi melhorado com suporte completo para **channels**, permitindo comunicação segura entre goroutines na DSL do Sloth Runner. Agora você pode usar channels para implementar padrões de concorrência avançados diretamente em Lua.

## Features Implementadas

### 1. **Channel Creation**
- Channels unbuffered (síncronos)
- Channels buffered (assíncronos com capacidade)
- Channels direcionais (send-only, receive-only, bidirectional)

### 2. **Channel Operations**
- `send()` - Envia valor para o channel (blocking)
- `receive()` - Recebe valor do channel (blocking)
- `try_send()` - Envia valor sem bloqueio (non-blocking)
- `try_receive()` - Recebe valor sem bloqueio (non-blocking)
- `close()` - Fecha o channel
- `len()` - Retorna o número de elementos no buffer
- `cap()` - Retorna a capacidade do channel
- `is_closed()` - Verifica se o channel está fechado

### 3. **Select Statement**
- Multiplexação de múltiplos channels
- Casos de receive e send
- Caso default para operações non-blocking
- Handlers para cada caso

## API Reference

### Creating Channels

```lua
-- Unbuffered channel (sincronização)
local ch = goroutine.channel()

-- Buffered channel with capacity 10
local ch = goroutine.channel(10)

-- Send-only channel
local ch_send = goroutine.channel(5, "send")

-- Receive-only channel
local ch_receive = goroutine.channel(5, "receive")
```

### Channel Methods

#### send(value)
Envia um valor para o channel (blocking).

```lua
local ch = goroutine.channel()
ch:send("Hello, World!")
```

**Returns**: `true` on success, `false, error_message` on failure

#### receive()
Recebe um valor do channel (blocking).

```lua
local ch = goroutine.channel()
local value, ok = ch:receive()
if ok then
  log.info("Received: " .. value)
else
  log.info("Channel closed")
end
```

**Returns**: `value, ok` where `ok` indicates if the value was successfully received

#### try_send(value)
Tenta enviar sem bloqueio.

```lua
local ch = goroutine.channel(2)
local ok = ch:try_send("message")
if ok then
  log.info("Sent successfully")
else
  log.info("Channel full, message not sent")
end
```

**Returns**: `true` if sent, `false` if channel is full

#### try_receive()
Tenta receber sem bloqueio.

```lua
local ch = goroutine.channel()
local value, ok = ch:try_receive()
if ok then
  log.info("Received: " .. value)
else
  log.info("No value available")
end
```

**Returns**: `value, ok` where `ok` indicates if a value was available

#### close()
Fecha o channel. Após fechado, não pode mais receber valores.

```lua
local ch = goroutine.channel()
ch:close()
```

**Returns**: `true` on success, `false, error_message` if already closed

#### len()
Retorna o número de elementos no buffer do channel.

```lua
local ch = goroutine.channel(10)
ch:send("a")
ch:send("b")
log.info("Length: " .. ch:len()) -- Output: Length: 2
```

**Returns**: `number` - current buffer length

#### cap()
Retorna a capacidade total do channel.

```lua
local ch = goroutine.channel(10)
log.info("Capacity: " .. ch:cap()) -- Output: Capacity: 10
```

**Returns**: `number` - channel capacity

#### is_closed()
Verifica se o channel está fechado.

```lua
local ch = goroutine.channel()
ch:close()
if ch:is_closed() then
  log.info("Channel is closed")
end
```

**Returns**: `boolean` - true if closed, false otherwise

### Select Statement

O `select` permite multiplexar operações em múltiplos channels:

```lua
goroutine.select({
  {
    channel = ch1,
    receive = true,
    handler = function(value)
      log.info("Received from ch1: " .. value)
    end
  },
  {
    channel = ch2,
    send = "message",
    handler = function()
      log.info("Sent to ch2")
    end
  },
  {
    default = true,
    handler = function()
      log.info("No channel ready")
    end
  }
})
```

**Select Cases**:
- **Receive case**: `{ channel = ch, receive = true, handler = function(value) ... end }`
- **Send case**: `{ channel = ch, send = value, handler = function() ... end }`
- **Default case**: `{ default = true, handler = function() ... end }`

## Usage Patterns

### 1. Producer-Consumer

```lua
local ch = goroutine.channel()

-- Producer
goroutine.spawn(function()
  for i = 1, 5 do
    ch:send(i)
    log.info("Produced: " .. i)
  end
  ch:close()
end)

-- Consumer
while true do
  local value, ok = ch:receive()
  if not ok then break end
  log.info("Consumed: " .. value)
end
```

### 2. Buffered Channel

```lua
local ch = goroutine.channel(3)

-- Send without blocking
ch:send("first")
ch:send("second")
ch:send("third")

log.info("Buffer status: " .. ch:len() .. " / " .. ch:cap())

-- Receive all
while ch:len() > 0 do
  local value, ok = ch:receive()
  if ok then
    log.info("Got: " .. value)
  end
end
```

### 3. Fan-Out Fan-In

```lua
local jobs = goroutine.channel(100)
local results = goroutine.channel(100)

-- Worker function
local worker = function(id)
  while true do
    local job, ok = jobs:receive()
    if not ok then break end

    -- Process job
    log.info("Worker " .. id .. " processing: " .. job)
    results:send("Result of " .. job)
  end
end

-- Spawn workers
for i = 1, 3 do
  local worker_id = i
  goroutine.spawn(function()
    worker(worker_id)
  end)
end

-- Send jobs
for i = 1, 9 do
  jobs:send("Job-" .. i)
end
jobs:close()

-- Collect results
for i = 1, 9 do
  local result, ok = results:receive()
  if ok then
    log.info("Result: " .. result)
  end
end
```

### 4. Pipeline Pattern

```lua
-- Stage 1: Generate numbers
local generate = function()
  local out = goroutine.channel(5)
  goroutine.spawn(function()
    for i = 1, 5 do
      out:send(i)
    end
    out:close()
  end)
  return out
end

-- Stage 2: Square numbers
local square = function(input)
  local out = goroutine.channel(5)
  goroutine.spawn(function()
    while true do
      local num, ok = input:receive()
      if not ok then break end
      out:send(num * num)
    end
    out:close()
  end)
  return out
end

-- Build and run pipeline
local numbers = generate()
local squared = square(numbers)

while true do
  local value, ok = squared:receive()
  if not ok then break end
  log.info("Result: " .. value)
end
```

### 5. Timeout Pattern

```lua
local data = goroutine.channel()
local timeout = goroutine.channel()

-- Slow operation
goroutine.spawn(function()
  goroutine.sleep(500)
  data:send("Data")
  data:close()
end)

-- Timeout timer
goroutine.spawn(function()
  goroutine.sleep(200)
  timeout:send("timeout")
  timeout:close()
end)

-- Wait with timeout
goroutine.select({
  {
    channel = data,
    receive = true,
    handler = function(value)
      log.info("Got data: " .. value)
    end
  },
  {
    channel = timeout,
    receive = true,
    handler = function()
      log.info("Operation timed out!")
    end
  }
})
```

### 6. Quit Channel (Graceful Shutdown)

```lua
local data = goroutine.channel()
local quit = goroutine.channel()

-- Worker
goroutine.spawn(function()
  local count = 0
  while true do
    -- Check quit signal
    local signal, ok = quit:try_receive()
    if ok then
      log.info("Shutting down...")
      break
    end

    -- Do work
    count = count + 1
    data:send("Work " .. count)
    goroutine.sleep(100)
  end
  data:close()
end)

-- Receive data
for i = 1, 5 do
  local value, ok = data:receive()
  if ok then
    log.info("Received: " .. value)
  end
end

-- Send quit signal
quit:send(true)
quit:close()
```

### 7. Rate Limiting

```lua
local requests = goroutine.channel(10)
local limiter = goroutine.channel(3)

-- Fill limiter with tokens
for i = 1, 3 do
  limiter:send(true)
end

-- Refill periodically
goroutine.spawn(function()
  while true do
    goroutine.sleep(100)
    limiter:try_send(true) -- Add token
  end
end)

-- Process requests with rate limiting
local handle_request = function(id)
  -- Wait for token
  local token, ok = limiter:receive()
  if ok then
    log.info("Processing request " .. id)
    goroutine.sleep(50)
    log.info("Request " .. id .. " complete")
  end
end

-- Send requests
for i = 1, 10 do
  local req_id = i
  goroutine.spawn(function()
    handle_request(req_id)
  end)
end
```

## Examples

Os seguintes exemplos estão disponíveis:

- **`examples/goroutine_channels_example.sloth`**: 8 exemplos de padrões com channels
- **`examples/goroutine_select_example.sloth`**: 7 exemplos de uso de select

Execute os exemplos com:

```bash
sloth-runner run producer_consumer --file examples/goroutine_channels_example.sloth --yes
sloth-runner run basic_select --file examples/goroutine_select_example.sloth --yes
```

## Best Practices

### 1. Always Close Channels

```lua
local ch = goroutine.channel()
-- Use channel
ch:close() -- Always close when done
```

### 2. Receiver Closes

Geralmente, quem **recebe** deve fechar o channel, não quem envia.

```lua
-- Producer
goroutine.spawn(function()
  for i = 1, 5 do
    ch:send(i)
  end
  ch:close() -- Producer closes after sending all
end)
```

### 3. Check Channel Closed

```lua
local value, ok = ch:receive()
if not ok then
  -- Channel is closed
  return
end
```

### 4. Use Buffered Channels for Performance

```lua
-- Unbuffered: blocking on every send/receive
local ch = goroutine.channel()

-- Buffered: allows N sends without blocking
local ch = goroutine.channel(10)
```

### 5. Select with Default for Non-Blocking

```lua
goroutine.select({
  {
    channel = ch,
    receive = true,
    handler = function(value)
      -- Process value
    end
  },
  {
    default = true,
    handler = function()
      -- No value available, continue
    end
  }
})
```

## Implementation Details

### Thread Safety

Todos os channels são thread-safe e podem ser usados por múltiplas goroutines simultaneamente.

### Memory Management

Channels são automaticamente garbage collected quando não há mais referências. O método `Cleanup()` fecha todos os channels abertos quando o módulo é finalizado.

### Goroutine Safety

Cada goroutine tem seu próprio Lua state, evitando race conditions. Os channels fornecem a comunicação segura entre goroutines.

## Limitations

1. **No channel of channels**: Atualmente, não é possível enviar channels dentro de channels.
2. **Value types**: Apenas valores Lua básicos (number, string, bool, table) podem ser enviados através de channels.
3. **Select limitations**: O select executa apenas o primeiro caso que estiver pronto. Não há prioridade entre casos.

## Migration from Old Code

Se você estava usando apenas goroutines sem channels:

**Antes**:
```lua
local wg = goroutine.wait_group()
wg:add(2)

goroutine.spawn(function()
  -- Do work
  wg:done()
end)

goroutine.spawn(function()
  -- Do work
  wg:done()
end)

wg:wait()
```

**Depois** (com channels para comunicação):
```lua
local results = goroutine.channel(2)

goroutine.spawn(function()
  local result = do_work()
  results:send(result)
end)

goroutine.spawn(function()
  local result = do_work()
  results:send(result)
end)

-- Collect results
for i = 1, 2 do
  local result, ok = results:receive()
  if ok then
    log.info("Got: " .. result)
  end
end
```

## Troubleshooting

### Channel Deadlock

```lua
-- WRONG: Unbuffered channel with no receiver
local ch = goroutine.channel()
ch:send("value") -- DEADLOCK! Nobody to receive

-- RIGHT: Use buffered channel or spawn receiver first
local ch = goroutine.channel(1)
ch:send("value") -- OK, buffer has space
```

### Channel Closed Panic

```lua
-- WRONG: Sending to closed channel
ch:close()
ch:send("value") -- ERROR: send on closed channel

-- RIGHT: Check if closed first
if not ch:is_closed() then
  ch:send("value")
end
```

### Resource Leaks

```lua
-- WRONG: Never closing channel
local ch = goroutine.channel()
-- ... use channel but never close

-- RIGHT: Always close when done
local ch = goroutine.channel()
-- ... use channel
ch:close()
```

## Conclusion

Os channels trazem o poder da programação concorrente do Go para o Sloth Runner, permitindo construir workflows complexos e altamente paralelos com segurança e elegância.

Para mais informações sobre concorrência em Go (que inspirou esta implementação):
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go by Example - Channels](https://gobyexample.com/channels)
