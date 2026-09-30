# Task Manager из облачного окружения

## Выбор подключения

Запрос «проверь доступ к задачам» требует Task Manager `get_workspace` и
`list_tasks`, а не Codex `list_threads`. Не останавливайся на отсутствии
нативных tools, пока не проверены доступный shell runtime и SDK-путь ниже.
Если самого executor нет, назови это ограничение отдельно от отсутствия
плагина или токена. Установка плагина, наличие файлов скилла и регистрация
native tools в конкретном чате — разные состояния.

1. Если в текущем агенте уже доступны Task Manager MCP tools, вызови
   `get_workspace` и используй их. Не запрашивай дополнительный токен.
2. Иначе проверь текущий runtime: при наличии инструмента вызови
   `cloud_environment.environment_status`, прочитай его runtime/network skill
   и `/etc/codex/network-policy.json`. Проверяй наличие переменной без вывода
   значения. Placeholder от credential proxy может быть нормальным значением;
   не проверяй действительность токена по его префиксу.
3. Для прямого MCP используй runtime secret **`TASK_MANAGER_TOKEN`** и endpoint
   **`https://task-manager.xxsrez-work.chatgpt.site/api/mcp`**. Источник адреса —
   `mcpServers.task-manager.url` в `.mcp.json` установленного плагина; при
   изменении пакета сверяй адрес с ним. Отправляй `Authorization: Bearer …`
   только этому endpoint, без перенаправлений на другие адреса.

Секрет для remote MCP отличается от `TASK_MANAGER_LOCAL_TOKEN`: последний
нужен Linux companion для загрузки локальных файлов. Настройка remote MCP
не доказывает доступность companion или native file ingress.

## Если токен ещё не добавлен

Сначала проверь, отсутствует ли секрет в runtime или лишь не применена его
настройка. Draft, сохранённое значение и переменная работающего процесса —
разные состояния. При доступной диагностике требуй current observations и
`ready`; после применения снова проверяй runtime и реальный запрос.

Объясни пользователю действия его языком, например:

> Для доступа из этого облачного окружения нужен личный токен Task Manager.
> Если он уже есть, новый создавать не нужно. В настройках **этого окружения**
> добавьте секрет с именем **TASK_MANAGER_TOKEN**, вставьте токен в защищённое
> поле и примените конфигурацию через **Publish**, если интерфейс предлагает
> этот шаг. Разрешённый домен: **task-manager.xxsrez-work.chatgpt.site**.
> Это публикация настроек окружения, не развёртывание сайта. Токен в чат
> присылать не нужно. После применения я проверю доступ чтением данных.

Если доступны инструменты настройки окружения, можно подготовить только
декларацию секрета: `name=TASK_MANAGER_TOKEN`, `source=user_provided`,
`target.type=environment_variable`, `target.name=TASK_MANAGER_TOKEN`,
`target.allowed_domains=["task-manager.xxsrez-work.chatgpt.site"]`.
Сначала прочитай draft и актуальную схему инструмента; сохрани остальные
настройки и secrets. Пользователь сам вводит значение и подтверждает
предусмотренное платформой применение. Секрет должен быть доступен во время
работы агента, а не только в setup script. Не выдавай наличие draft за успех.

### Как получить личный токен

Пользователь входит в [Task Manager](https://task-manager.xxsrez-work.chatgpt.site)
под своим аккаунтом. В проверенной версии интерфейса от 2026-09-30 нет отдельной
страницы создания токенов; не направляй пользователя в выдуманное меню
«Settings → API tokens». Доступен авторизованный endpoint
`POST /api/settings/api-credentials`.

Если готового токена нет, предложи пользователю выполнить следующую команду
**самостоятельно в консоли разработчика своей авторизованной вкладки Task
Manager**. Она создаёт токен на 90 дней для чтения и изменения задач — подходит
для Grinder. Для задачи только на чтение замени scopes на `["api:read"]`.
Не выполняй выдачу credential автоматически и не получай ответ с секретом
через agent/browser tool.

```javascript
await (async () => {
  if (location.origin !== "https://task-manager.xxsrez-work.chatgpt.site") {
    throw new Error("Откройте авторизованную вкладку Task Manager");
  }
  const response = await fetch("/api/settings/api-credentials", {
    method: "POST",
    credentials: "same-origin",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      name: "Codex Cloud",
      scopes: ["api:read", "api:write"],
      expiresInDays: 90
    })
  });
  if (!response.ok) throw new Error(`Создание токена: HTTP ${response.status}`);
  const result = await response.json();
  return result.token;
})();
```

Консоль покажет токен пользователю один раз. Он переносит его непосредственно
в защищённое поле секрета окружения; не присылает вывод или скриншот агенту.
Не сохраняй токен в репозитории, `.env`, конфиге, аргументах команды, логах или
snapshot и не копируй OAuth caches. При 401 на выдаче сначала нужен вход в
Task Manager. Если пользователь не может воспользоваться консолью, прямо
назови отсутствие UI выдачи; не предлагай несуществующую кнопку.

## Как агент вызывает MCP

### Штатный MCP-клиент

Когда агент действительно запускается через настраиваемый Codex CLI,
подключение использует ссылку на переменную, без значения секрета:

```toml
[mcp_servers.task-manager]
url = "https://task-manager.xxsrez-work.chatgpt.site/api/mcp"
bearer_token_env_var = "TASK_MANAGER_TOKEN"
```

Меняй только конфиг фактического запускающего процесса. Файл в произвольном
каталоге плагина не регистрирует tools в уже работающем cloud-чате. Не меняй
host-managed `CODEX_HOME` и не копируй credentials ради вложенного CLI.
Если native registration недоступна, используй SDK в текущем shell runtime.

### Официальный SDK из shell

Используй `@modelcontextprotocol/sdk`, `Client` из `client/index.js` и
`StreamableHTTPClientTransport` из `client/streamableHttp.js`. В проверенном
облачном подключении использовались SDK 1.31.0 и `undici.EnvHttpProxyAgent`.
Размещай зависимости/адаптер в доступном для записи служебном каталоге,
не в коде пользовательского продукта; существующий проверенный адаптер
предпочтительнее повторного создания. Фиксируй зависимости lockfile.

В ранее настроенном окружении адаптер мог находиться в
`/workspace/.task-manager-check/task-manager.mjs` с инструкцией рядом в
`README.md`. Это путь конкретного окружения, не файл поставки плагина:
сначала проверь наличие и прочитай интерфейс. Другой cloud-чат может иметь
другой executor и не наследовать этот файл. Если адаптера нет, создай его
по примеру ниже в своём writable каталоге; не сообщай пользователю, что
доступ невозможен только из-за отсутствия старого файла.

Ключевые параметры транспорта:

```javascript
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StreamableHTTPClientTransport } from "@modelcontextprotocol/sdk/client/streamableHttp.js";
import { EnvHttpProxyAgent, fetch as proxyFetch } from "undici";

const endpoint = "https://task-manager.xxsrez-work.chatgpt.site/api/mcp";
const token = process.env.TASK_MANAGER_TOKEN;
if (!token) throw new Error("Добавьте runtime secret TASK_MANAGER_TOKEN");
const dispatcher = new EnvHttpProxyAgent();
const transport = new StreamableHTTPClientTransport(new URL(endpoint), {
  requestInit: { headers: { Authorization: `Bearer ${token}` } },
  fetch: (url, init) => {
    if (new URL(String(url)).href !== endpoint) throw new Error("Unexpected MCP URL");
    return proxyFetch(url, { ...init, dispatcher, redirect: "error" });
  }
});
const client = new Client({ name: "task-manager-cloud", version: "1.0.0" });
try {
  await client.connect(transport);
  const catalog = await client.listTools();
  const result = await client.callTool({ name: "get_workspace", arguments: {} });
  if (result.isError) throw new Error("get_workspace returned a tool error");
  console.log(JSON.stringify({ connected: true, firstPageTools: catalog.tools.length }));
} finally {
  await client.close();
  await dispatcher.close();
}
```

Это пример проверки подключения, не полный адаптер. При создании адаптера
получай все страницы `tools/list`, проверяй аргументы по текущей `inputSchema`
и вызывай канонический `tools/call`. Для работы возвращай агенту `content` и
`structuredContent`, включая refs, versions, cursors и ошибки; для smoke
выводи только результат проверки без пользовательских данных. Удаляй значение
секрета из диагностических сообщений перед выводом. Не повторяй writes
автоматически: при неизвестном результате сначала сверяй состояние чтением.
Правила полномочий, pagination, comments и transitions основного скилла и
вызывающего workflow сохраняются. SDK не заменяет native OpenAI file bridge;
не передавай локальные пути инструментам с `openai/fileParams`.

## Подтверждение и ошибки

- Сохраняй унаследованные HTTP/HTTPS proxy, CA trust и TLS verification.
  Прямое соединение в обход sidecar и имитация браузера не нужны.
- Успех: MCP initialize, discovery и реальный `get_workspace` без `isError`.
  Наличие переменной, `mcp get`, каталог tools или HTTP 200 сами по себе
  не доказывают принятие credential. Read-only smoke не доказывает write scope.
  Если пользователь проверяет именно видимость задач, дополнительно получи
  текущую схему `list_tasks` и вызови его с небольшим `limit` (например, 5),
  без project/release filters, если они не заданы. Покажи несколько identifiers
  и факт наличия следующей страницы; не выдавай размер страницы за полный
  объём workspace и не выводи описания задач для проверки подключения.
- Cloudflare 403/1010 от одного ad hoc клиента не означает общий запрет MCP.
  Проверь штатный клиент/SDK через разрешённый proxy. Такой SDK-путь успешно
  проверен 2026-09-30 после отказа Python urllib. Если отказ сохраняется,
  сообщи точный слой и ошибку; не меняй Site audience, Cloudflare, production
  deployment или сетевую политику как автоматическое «исправление».
- При API 401 проверь готовность binding и применённую конфигурацию, затем
  expiry/revoke; для нового токена пользователь повторяет защищённую настройку.
  При `insufficient_scope` объясни требуемый scope; повторное создание токена
  не является первым действием при транспортном отказе.
- Нативные MCP tools и SDK-доступ — разные способы исполнения. Укажи,
  какой проверен. Для Grinder достаточно проверенного вызова канонических
  tools из текущего агента; другие требования Grinder этим не проверяются.
