---
name: ship-tasks
description: "Доставлять однозначно выбранный Task Manager scope до terminal outcome. Использовать явно через $ship-tasks; implicit invocation — только когда delivery intent сопровождается exact существующей Task (например TM-123) либо уже выбранным Task Manager Project/Release/current scope. Одного delivery-глагола недостаточно: обычные просьбы исправить продукт, код, repository или plugin без Task Manager anchor не активируют ShipTask. Create-and-deliver использовать только при явной просьбе создать ровно одну Task именно в Task Manager и сразу начать/выполнить её; также использовать для явной настройки ShipTask project memory. Bare $ship-tasks запускает batch по memory current_scope; exact Task — single без Goal, multi-task/Project/Release — batch с Goal. Не использовать для чтения, статуса, аудита, объяснения, планирования или backlog capture."
---

# Ship Tasks

ShipTask владеет business delivery policy. Task Manager skill и connector
использовать как technical adapter. Project memory даёт selectors и
project profile, но не live Task state и не authority.

## 1. Проверить запуск и выбрать mode

Invocation gate пройден только в одном из случаев:

1. Пользователь явно вызвал `$ship-tasks`.
2. Запрос одновременно содержит delivery intent и один однозначный
   Task Manager anchor: exact Task, уже выбранный Project/Release/current
   scope либо явный create-and-deliver ровно одной Task.
3. Пользователь явно просит настроить ShipTask project memory.

Один delivery verb, текущий repository или software project не являются
Task Manager anchor. Не читать memory и не вызывать Task Manager/Goal tools,
чтобы придумать anchor задним числом. Если gate не пройден,
выполнить обычный workflow без ShipTask lifecycle, Goal и reports.

| Запрос | Mode | Goal |
|---|---|---|
| Bare `$ship-tasks` | `batch` по memory `current_scope` | обязателен |
| Exact одна Task | `single` | нет |
| Ровно одна новая Task в Task Manager и сразу выполнить | `single create-and-deliver` | нет |
| Несколько Tasks, Project, Release или current scope | `batch` | обязателен |
| Явно оформить memory | `memory-maintenance` | нет |
| Read/status/audit/explain/plan/backlog capture | `non-delivery` | нет |

«Просто добавь», «положи в backlog» и «запланируй» не запускают
delivery. Explicit read-only/local-only/no-deploy ограничивает более общий
delivery verb. Exact selector текущего запроса выше memory default только
для текущего run.

В `memory-maintenance` полностью прочитать
[project-memory reference](references/project-memory.md) и писать memory только
по явной просьбе. В `non-delivery` выполнить только точный read или planning
request через adapter, если он нужен.

## 2. Разрешить exact scope и live state

Для `batch`, bare invocation и remembered defaults полностью прочитать
[project-memory reference](references/project-memory.md). Приоритет:

1. exact selector текущего запроса;
2. memory `current_scope` для bare/current-scope intent;
3. current canonical lookup;
4. complete live inventory и full Task detail.

Memory не доказывает current status, `version`, relations, comments, access,
connector capability, deploy result или approval. Если обязательный context
отсутствует, неоднозначен или противоречит live evidence так, что write
небезопасен, выдать global `TASK CONTEXT ALARM` и не мутировать ничего.

Через Task Manager adapter:

1. Проверить workspace identity, read/write capabilities и status catalog.
2. Разрешить Project/Release/Task только в current canonical refs; проверить
   Release membership и ACL.
3. Для multi-task boundary пройти все страницы до `hasMore=false`.
4. До dependency, acceptance, relation или write reasoning получить full detail.
5. Для `duplicate_of` прочитать canonical Task и все incoming duplicates как
   review scenarios, не как отдельную execution work.
6. Проверить native comment create/list/read; imported comments не являются
   write capability.

Перед каждым Task update перечитать detail и передать current `version`.
При `version_conflict` перечитать и повторить только всё ещё применимый
intent. После write обязателен read-back. Не повторять create/comment вслепую
после unknown outcome; сначала искать duplicate.

Необязательно можно один раз назвать свежий Codex thread `ShipTask · ...`,
но только если surface надёжно подтверждает первый ход, title пуст или
очевидно сгенерирован, а live scope уже разрешён. Отсутствие title tool
никогда не блокирует delivery.

## 3. Создать Goal только для batch

В `single`, `memory-maintenance` и `non-delivery` не вызывать Goal tools.
В `batch` после exact refs и complete read-only inventory, но до первой
non-Goal mutation:

1. Вызвать `get_goal`. Если нужные Goal tools недоступны, выдать
   `TASK CONTEXT ALARM`.
2. Создать Goal без token budget, если user не задал его явно, либо
   продолжить только совместимый незавершённый Goal того же exact scope.
3. Зафиксировать objective, observable done criteria, verification, allowed writes и
   blocker rules.
4. Оставлять Goal активным при любой in-scope `To Do`, `In Progress`,
   `In Review`, deferred Task, rework/completion remainder или unresolved defect.

Task-level `BLOCKED`, run `PARTIAL`/`BLOCKED` и Goal status — разные состояния.
Task-local проблема оставляет Goal активным. Не повторять действие
и не создавать дополнительные Goal turns ради Goal `blocked`.

## 4. Согласовать lifecycle с фактами

```text
Backlog                         excluded; no writes
  └─ exact just-created Task → To Do   create-and-deliver recovery only
To Do → In Progress → In Review → Done
                         └──────→ Canceled
Duplicate                       terminal review context only
```

- Pre-existing `Backlog` не выполнять и не менять. Exact Task, только что
  созданную current create-and-deliver, можно вывести из default `Backlog`
  в current `To Do` после read-back.
- `To Do → In Progress` только после preflight; `In Progress → In Review`
  только после passing targeted gate и integration.
- `In Review` означает незавершённую классификацию result. Сам status не
  доказывает ни success, ни failure, ни пройденную проверку.
- `Done` разрешён только после полного terminal evidence и published/read-back
  `COMPLETED`. Новый `Canceled` transition также требует `CANCELED` report.

## 5. Разобрать каждую `In Review` за один проход

Перечитать current Task, acceptance, dependencies, material duplicates, accepted
project decisions и явные уточнения пользователя. История изменений acceptance
сама по себе не является conflict. После одного bounded diagnostic pass выбрать
ровно один исход:

| Исход | Достаточное основание | Действие |
|---|---|---|
| `task-contract-conflict` | Current mandatory requirements противоречат друг другу или не задают наблюдаемый result | Если exact correction однозначно следует из current accepted source и write authority уже есть, исправить contract и перечитать. Иначе оставить `In Review`, подготовить `BLOCKED` с точным противоречием и нужным решением. |
| `verified-success` | Полный evidence set доказывает current acceptance на exact result | Опубликовать/read-back `COMPLETED`, перевести в `Done`, перечитать. |
| `verified-failure` | Надёжное наблюдение exact candidate в совместимой среде прямо нарушает acceptance | Подготовить `REWORK REQUIRED`, перевести `In Review → In Progress`, перечитать и начать rework. |
| `verification-blocked` | Доступная проверка не доказывает ни success, ни failure | Оставить `In Review`, подготовить `BLOCKED` и завершить попытку до нового result/evidence/environment/access/authority/Task contract. |

Ошибка test harness, недоступный actor/access, несовместимая среда или
отсутствие наблюдаемости не доказывают product failure. Падение aggregate gate
без task-level attribution не делает весь batch defective: неразличимые members
остаются `In Review` до separating diagnostic.

Для `verification-blocked` попросить Strategic Explainer дать 2–4 способа
провести приёмку: prerequisites, что каждый способ докажет, tradeoff,
observable success signal и recommended next attempt.

Не повторять тот же acceptance scenario, poll, comment или Goal turn без
конкретного изменения входных условий.

## 6. Выполнить и проверить

До execution определить acceptance/dependencies, repository/integration identity,
allowed writes, targeted gate, batch gate и triggers, external effects, release target,
environment class, smoke/recovery и external approvals.

Disposition precedence для batch:

```text
global-conflict
→ actionable completion-remains
→ resume
→ work-remains
→ deferred-only
→ no-work
```

Global connector/scope/Goal/ownership/shared-state conflict даёт
`TASK CONTEXT ALARM`. Task-local ambiguity откладывает только affected Task.

`single` всегда serial. Для batch выбирать isolated lanes по dependencies,
conflicts, verification, integration и review capacity. Один writer владеет одной
Task; integration owner выполняет fan-in. Без реально запущенных isolated
workers active write target равен `1`.

Перед каждым новым `To Do → In Progress` выполнить status-reconciliation:
законченный implementation/rework должен сначала перейти в `In Review`,
честно отложенный partial result — освободить lane, terminal result — пройти
terminal read-back. Unknown status write не освобождает lane.

Для каждой Task:

1. Зафиксировать canonical ref, detail/version/status, acceptance, dependencies,
   result base и release target.
2. `To Do` начать после preflight; `In Progress` продолжить из coherent
   checkpoint; `In Review` сначала классифицировать, а не реализовывать
   заново без finding.
3. Реализовать минимальный целостный in-scope result; покрыть material
   duplicate scenarios; выполнить per-Task targeted gate.
4. Интегрировать exact candidate, перевести его в `In Review` и перечитать.
5. По trigger выполнить independent review и review-batch gate на exact integrated
   identity. Не запускать дорогой full gate для каждой Task без risk reason.
6. Выполнить required external effects и независимо проверить exact target.
7. Failed gate локализовать. Только proven failure возвращает Task в
   `In Progress`; unclear attribution оставляет неразличимые members в `In Review`
   как `verification-blocked`.

## 7. Продолжать автономно и соблюдать release boundary

Перед execution, первым defer или release decision полностью прочитать
[autonomy and release reference](references/autonomy-and-release.md).

- Не задавать task-local вопрос пока есть runnable work. Обратимый локальный
  default выбирать самостоятельно; material decision/authority blocker откладывать
  с checkpoint и продолжать независимые Tasks.
- Перед blocking input повторить complete inventory. При `runnable_count > 0`
  blocking input запрещён.
- Для defer подготовить и попытаться опубликовать/read-back `BLOCKED`.
  Без comment capability показать тот же handoff в run report, оставить
  communication remainder и не использовать Task fields как fallback.
- ShipTask delivery intent разрешает ordinary exact verified non-production workflow:
  local/dev/test/QA/UAT/staging/preview/sandbox build, deploy, bounded migration, smoke,
  repair и rollback.
- Production release требует explicit user authority для exact target/current scope.
  Unknown environment считать production-like. Без approval выполнить безопасную
  preparation/non-production verification, затем defer с
  `production-approval-required`.

Automatic acceptance не заменяет production, destructive-data, secrets, privacy,
access-policy, external-recipient или external-approver authority.

## 8. Никогда не ждать ручную приёмку

Acceptance criteria — проверяемые completion criteria, а не human sign-off.
После passing acceptance, targeted gate, applicable exact review-batch gate, integration,
required effects и resolved findings автоматически принять exact result,
опубликовать/read-back `COMPLETED`, перевести Task в `Done` и перечитать.

Не просить «принять» result, не создавать `acceptance-required`, не публиковать
`ACCEPTANCE READY` и не оставлять terminal-ready Task в `In Review`. Старые
memory/rollout/report/Goal/plan правила о ручной приёмке считать superseded
historical evidence.

## 9. Опубликовать понятный Task report

Перед первым report decision полностью прочитать
[delivery-report reference](references/delivery-report.md) и
[Strategic Explainer handoff](references/strategic-explainer.md).

| Report state | Status rule при недоступном comment channel |
|---|---|
| `COMPLETED` | Не переводить в `Done`; приёмку не повторять. |
| `CANCELED` для нового transition | Не выполнять terminal status write. |
| `REWORK REQUIRED` | Выполнить truthful `In Review → In Progress`; оставить communication remainder. |
| `BLOCKED` | Сохранить status, который диктует Task outcome; показать handoff в run report. |

Перед каждым новым `COMPLETED`, `REWORK REQUIRED`, `BLOCKED` или `CANCELED`:

1. ShipTask перечитывает Task/result/effects, проверяет evidence и доступный
   bounded repair, затем сам выбирает truthful state, authority и next action.
2. Сформулировать semantic `Problem to solve`: beneficiary, desired observable outcome и
   exact scope. Task ref/title не заменяют проблему.
3. Создать bounded `Strategic Handoff` и запустить fresh built-in `default`
   subagent с точным `fork_turns="none"`, поручив применить
   `$ship-tasks:strategic-explainer`. Не передавать inherited history, raw logs,
   process diary или готовый strategic conclusion.
4. Прочитать свободное problem-first explanation и source basis. Explainer не
   выбирает status, recovery, authority или Goal transition.
5. ShipTask пишет final comment своими словами: сохраняет problem, outcome,
   impact, причинную границу, confidence и next state; добавляет authoritative
   envelope и evidence из исходных фактов. Выполнить forward trace и reverse coverage.
6. Найти equivalent report по stable key, записать без blind retry и перечитать.

При `CONTEXT_INTEGRITY_ERROR` или `PROBLEM_CONTEXT_ERROR` исправить handoff и один
раз запустить новый fresh subagent. Повторный orchestration failure переходит
в локальную `degraded-adaptation`; это не блокирует truthful rework status.
Если после перечитывания Task нельзя сформулировать саму проблему, это
`task-contract-conflict`, а не повод придумывать текст.

Для `verification-blocked` передать `Decision support request` на 2–4 способа
приёмки. Никогда не писать report в description, acceptance, status text или другое
Task field. Не публиковать `ACCEPTANCE READY`.

## 10. Финализировать run

Перед terminal outcome, blocking pause или финальным ответом полностью прочитать
[run-report reference](references/run-report.md) и выполнить finalization pass:

1. Сопоставить обещанный и фактический result.
2. Перечитать current Task/Goal/source/effect state и material evidence.
3. Объяснить gaps и отделить доказанную причину от предположения.
4. Выполнить один safe bounded in-scope repair, если authority уже есть;
   перечитать state, повторить нужные checks и начать finalization заново.
5. Без material state change не повторять repair, приёмку, poll или Goal turn.

Для `single` доказать exact result, acceptance, singleton gates, effects, report и
truthful status. Goal — `not-applicable`.

Для `batch` перед terminal claim повторить complete inventory. Goal завершается
последним lifecycle write только когда нет `To Do`, `In Progress`, `In Review`,
deferred Tasks, unresolved in-scope defects или incomplete required effects. Если
остались только deferred Tasks, показать одну consolidated queue, оставить
Goal активным и остановить текущий run без искусственных повторов.

Каждый exit (`COMPLETED`, `PARTIAL`, `BLOCKED`, `NO WORK`) заканчивается глубоким,
но компактным `SHIPTASK RUN REPORT`: сначала итог, текущий статус и причина
простыми словами; затем только нужные evidence, ограничения и следующий шаг.
Не выгружать process diary, raw tool output или полный inventory и не выдавать
unknown/skipped effect за verified completion.
