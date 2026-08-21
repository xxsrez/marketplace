# Project memory contract

Этот reference задаёт логическую структуру project-specific context для
ShipTask. Читать его полностью в `batch-implementation`, при bare `$ship-tasks`,
при `memory-maintenance` и всякий раз, когда delivery/release зависит от
сохранённого project profile.

## Содержание

- Граница ответственности
- Поиск и приоритет источников
- Логическая схема
- Bootstrap и update
- Проверка перед delivery
- TASK CONTEXT ALARM
- Ограничения переносимости
- Пример

## Граница ответственности

Project memory хранит устойчивые selectors и правила конкретного проекта:
какой Task Manager Project/Release обычно считается текущим, где лежит repo,
какие ветки, команды, environments, verification gates и release boundaries
применяются.

Project memory не хранит как authority:

- текущие Task status, `version`, detail, relations, comments или access;
- факт успешного build/deploy/smoke;
- current connector capabilities;
- approval, которого пользователь не давал в текущем допустимом scope;
- secrets, tokens, signed URLs или private raw logs.

Эти данные всегда перечитываются из Task Manager, workspace/config и внешнего
provider state.

## Поиск и приоритет источников

Использовать context, доступный текущей Codex surface для открытого проекта.
Не считать наличие memory доказанным по прошлому rollout или другой surface.

Приоритет scope:

1. exact selector в текущем пользовательском запросе;
2. `current_scope` из найденного project memory для bare `$ship-tasks` или
   фразы «текущий scope/release/project»;
3. current Task Manager lookup, который превращает selector в canonical refs;
4. complete live inventory и full Task detail, определяющие рабочий scope.

Prompt override действует только в текущем run и не переписывает memory.
Никогда не объединять prompt scope и memory default молча.

## Логическая схема

Физический формат memory может меняться; ShipTask опирается на следующие
логические поля:

```yaml
shiptask_project_context:
  schema_version: 1
  project_identity:
    name: <human-readable project>
    workspace_hint: <optional Task Manager workspace hint>
  current_scope:
    kind: project | release | tasks
    project_selector: <name or remembered ref hint>
    release_selector: <optional name or remembered ref hint>
    task_selectors: [<optional identifiers>]
  repository:
    path: <absolute or workspace-resolvable path>
    branch_policy: <default/integration rules>
    worktree_policy: <isolation rules>
  verification:
    targeted_commands: [<commands or project-doc pointers>]
    batch_commands: [<commands or project-doc pointers>]
    batch_triggers: <risk/WIP/final-flush policy>
  delivery:
    provider: <provider or project-doc pointer>
    environments:
      non_production: [<known targets>]
      production: [<known targets>]
    smoke: [<commands, URLs without secrets, or doc pointers>]
    recovery: <bounded rollback/repair policy>
  authority:
    allowed_writes: <project-scoped normal writes>
    explicit_only: [production, destructive_data, secrets, privacy,
      access_policy, external_recipients]
  sources:
    - <canonical project document or config path>
  refreshed_at: <timestamp>
```

Поля могут быть опущены, когда они не нужны текущему scope. Не придумывать
значения ради полноты. Для обязательного неизвестного поля остановиться до
связанной mutation; для необязательного указать `unknown`/`not-applicable`.

## Bootstrap и update

Memory maintenance начинается только по явной просьбе вроде «настрой ShipTask
для проекта», «оформи project memory» или «обнови текущий release в памяти».

1. Прочитать доступные project instructions, repository config и Task Manager
   Project/Release metadata без mutations.
2. Показать resolved facts, источники, неизвестные поля и предлагаемый
   `current_scope`.
3. Запросить только material значения, которые нельзя безопасно вывести и
   которые меняют scope, environment class или authority.
4. Записать memory через доступный memory mechanism только после explicit
   request; не писать secrets и live task state.
5. Перечитать сохранённую memory, проверить обязательные поля и сообщить, на
   какой Codex surface она будет доступна.

Обычный `single`, `batch-implementation` или `release` delivery может сообщить
полезное предложение по memory, но не должен менять её как побочный эффект.

## Проверка перед delivery

Memory является selector/profile hint. Перед mutation:

- разрешить remembered Project/Release/Task selector через current Task
  Manager adapter и заменить его live canonical refs;
- перечитать Project/Release membership и полный Task detail;
- подтвердить repository path и project instructions в текущем workspace;
- проверить environment class по current config/provider state;
- проверить, что команды и policy sources ещё существуют;
- отделить remembered defaults от explicit current user authority.

`refreshed_at` помогает оценить риск, но свежая дата не заменяет live lookup.
Если remembered ref больше не разрешается, попробовать однозначный human name;
не выбирать похожий Project/Release автоматически.

## TASK CONTEXT ALARM

До любой delivery mutation выдать global `TASK CONTEXT ALARM`, если:

- bare `$ship-tasks` не находит ровно один применимый `current_scope`;
- memory и current prompt требуют несовместимые scopes;
- selector разрешается в несколько Project/Release/Tasks;
- remembered repository или Project/Release membership противоречит current
  evidence;
- environment/authority неоднозначны так, что следующая mutation небезопасна.

Alarm должен содержать известный scope, конфликтующие факты и источники,
последний безопасный checkpoint и одно точное решение/исправление context.
Не создавать Goal, не менять Task, code, Git, deploy state или memory до
разрешения global alarm.

Изолированный blocker уже разрешённой Task обрабатывается обычным ShipTask
`deferred` flow и не превращается в global alarm.

## Ограничения переносимости

Codex Memories могут быть отключены, ещё не загружены либо доступны только в
конкретной surface. Local Codex memory и память другой ChatGPT/Work surface не
считаются автоматически общим хранилищем.

Поэтому skill:

- не обещает cross-surface visibility без отдельной синхронизации;
- не считает прошлую беседу доказательством наличия current memory;
- допускает exact selector в prompt как безопасный временный override;
- при недоступной обязательной memory не угадывает default scope.

## Пример

```yaml
shiptask_project_context:
  schema_version: 1
  project_identity:
    name: Example App
  current_scope:
    kind: release
    project_selector: Example App
    release_selector: August UAT
  repository:
    path: /workspace/example-app
    branch_policy: integrate into the current release branch
    worktree_policy: one writable Task per isolated worktree
  verification:
    targeted_commands: ["npm test -- <affected suite>"]
    batch_commands: ["npm test", "npm run build"]
    batch_triggers: final flush or shared release candidate
  delivery:
    provider: project runbook
    environments:
      non_production: [UAT]
      production: [production]
    smoke: ["project runbook UAT smoke"]
    recovery: bounded redeploy or rollback per runbook
  authority:
    allowed_writes: ordinary in-scope local and UAT workflow
    explicit_only: [production, destructive_data, secrets, privacy,
      access_policy, external_recipients]
  sources: ["AGENTS.md", "docs/operations-runbook.md"]
  refreshed_at: 2026-08-16
```
