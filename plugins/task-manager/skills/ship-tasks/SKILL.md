---
name: ship-tasks
description: "Доставлять Task Manager Tasks до terminal outcome при явном $ship-tasks и по естественным просьбам выполнить, исправить или довести одну Task, несколько Tasks, текущий Project/Release/scope. Использовать также для просьбы создать одну Task и сразу начать или выполнить её: это single create-and-deliver, а не backlog capture; и для явной настройки project memory ShipTask. Bare $ship-tasks запускает batch по memory current_scope; exact Task — single без Goal, multi-task/Project/Release — batch с Goal. Не использовать для чтения, статуса, аудита, объяснения, планирования, backlog capture или «просто добавь/положи в backlog/запланируй» без delivery intent. Task Manager — live task source; ShipTask владеет delivery policy, memory — project selectors/profile. После полного evidence автоматически принять result, опубликовать native COMPLETED до Done и выпускать verified non-production без confirmation. Production, destructive data, secrets, privacy, access policy и external recipients требуют explicit authority."
---

# Ship Tasks

Применять одну business delivery policy к explicit и implicit delivery intent.
Task Manager skill/connector использовать как technical adapter; project memory
использовать как selector/profile hint, но не как live task state.

## Классифицировать intent до mutations

Выбрать ровно один mode:

| Intent | Mode | Goal |
|---|---|---|
| Bare `$ship-tasks` | `batch` по memory `current_scope` | обязательный |
| Exact одна Task, явно или natural language | `single` | нет |
| Создать одну Task и сразу начать/выполнить её | `single` (`create-and-deliver`) | нет |
| Несколько Tasks, Project, Release или current scope | `batch` | обязательный |
| Явно оформить/обновить ShipTask memory | `memory-maintenance` | нет |
| Read/status/audit/explain/plan/backlog capture | `non-delivery` | нет |

- Delivery verbs: выполнить, исправить, реализовать, довести, доставить,
  протестировать и выпустить выбранную Task/scope.
- Комбинация create intent и delivery intent, например «создай задачу и начинай
  делать», является `single create-and-deliver`. Не считать её backlog capture
  и не завершать flow после одного `create_task`.
- «Просто добавь», «положи в backlog», «запланируй» не запускают delivery и не
  меняют lifecycle существующей Task.
- Explicit read-only/local-only/no-deploy ограничивает более общий delivery
  verb.
- `single` требует ровно одну canonical Task до implementation. Обычный single
  разрешает существующую Task. `Create-and-deliver` сначала разрешает exact
  parent/content/acceptance/write authority, создаёт ровно одну Task, выполняет
  canonical read-back и использует только её. При нуле или нескольких
  совпадениях обычного single остановиться до mutation; не создавать замену и
  не выбирать похожую Task. Terminal Task не reopen-ить без explicit request.
- Exact selector в текущем prompt выше memory default и действует только в
  этом run. Не объединять scopes и не переписывать memory молча.

Для `non-delivery` не исполнять этот delivery workflow: выполнить только exact
read/planning Task Manager request через adapter. Для `memory-maintenance`
прочитать [project-memory reference](references/project-memory.md) полностью,
собрать source-grounded profile и записать memory только по explicit request.

## Разрешить project context и scope

Перед `batch`, bare invocation или использованием remembered project defaults
прочитать [project-memory reference](references/project-memory.md) полностью.

Scope precedence:

1. exact selector текущего prompt;
2. project memory `current_scope` для bare/current-scope intent;
3. current Task Manager canonical lookup;
4. complete live inventory и full Task detail.

Memory хранит selectors, repository, branches, commands, environments,
verification/release/recovery profile и authority hints. Никогда не считать её
authority для current Task detail/status/`version`/relations/comments/access,
connector capabilities, deploy result или user approval.

Если обязательная memory недоступна, не загружена на этой Codex surface,
неоднозначна, устарела или противоречит current evidence так, что mutation
небезопасна, выдать global `TASK CONTEXT ALARM`: known scope, conflicting facts
и sources, last safe checkpoint, одно exact решение. Не создавать Goal, не
менять Task, code, Git, deploy state или memory до разрешения alarm.

## Получить authoritative Task Manager state

Использовать current Task Manager adapter instructions и tool descriptions для
exact payloads. Сохранить минимальные invariants даже если adapter skill не был
отдельно загружен:

1. Проверить workspace identity, read/write capabilities и status catalog.
2. Разрешить Project/Release/Task selectors только в current canonical refs;
   проверить Release membership и resource ACL.
3. Для multi-task boundary получить все страницы inventory до `hasMore=false`.
   Counts/первая страница не являются complete inventory.
4. До dependency, acceptance, relation или write reasoning получить full detail
   каждой рабочей Task. External context читать только при material provenance.
5. Для `duplicate_of` прочитать canonical Task и все её duplicates как review
   scenarios, но не выделять duplicates отдельной execution work.
6. Проверить native Task comment create/list/read и authority. Imported
   comments не являются write capability.

Перед каждым Task update перечитать detail и передать current `version`. При
`version_conflict` перечитать и повторить только всё ещё применимый intent. Не
перетирать unrelated edits. После write выполнить read-back. Не повторять
create/comment write вслепую после unknown outcome; сначала искать duplicate.

## Установить lifecycle и рабочий scope

```text
Backlog                         excluded; no ShipTask writes
  \- exact just-created Task -> To Do   create-and-deliver recovery only
To Do -> In Progress -> In Review -> Done
                         \-----> Canceled
Duplicate                       terminal review context only
```

- Рабочие statuses: `To Do`, `In Progress`, `In Review` по current canonical
  catalog. Terminal аналоги распознавать по current category/status.
- Pre-existing `Backlog` исключён из delivery даже при Project/Release scope.
  Не выполнять и не менять его. Узкое исключение — exact Task, только что
  созданная текущим `create-and-deliver`: если adapter вернул её в `Backlog`,
  перечитать current `version`, перевести только её в current `To Do` и сделать
  read-back. Это не разрешает выбирать другую Backlog Task.
- Для `create-and-deliver` создать Task сразу в current `To Do`, когда adapter
  это поддерживает. После canonical/full-detail read-back и успешного preflight
  перевести её в `In Progress`, перечитать и только затем начать implementation.
  Не выдавать один create или оставшийся `Backlog`/`To Do` за начатое execution.
- `To Do` переводить в `In Progress` после preflight; `In Review` — только
  после targeted gate; `Done` — только после полного terminal evidence и
  published/read-back `COMPLETED` comment.
- Любая `In Review` Task означает `completion-remains`, а не human acceptance.

## Создать Goal только для batch

В `single`, `memory-maintenance` и `non-delivery` не вызывать Goal tools.

В `batch` после exact refs и complete read-only inventory, но до code, Git,
Task Manager, ownership или external writes:

1. Вызвать `get_goal`; при отсутствии обязательных Goal tools выдать
   `TASK CONTEXT ALARM`.
2. Создать Goal без token budget, если пользователь отдельно не задал numeric
   budget, либо продолжить только совместимый незавершённый Goal того же exact
   scope/outcome. Несовместимый Goal не заменять и не завершать.
3. Зафиксировать objective, observable done criteria, verification, exact
   scope/allowed writes и blocker rules.
4. Удерживать Goal активным при любой in-scope `To Do`, `In Progress`,
   `In Review`, rework/completion remnant, deferred Task или unresolved in-scope
   defect. Завершить Goal последним lifecycle write после final reconciliation.

## Провести preflight

До execution определить exact scope, acceptance/dependencies, repository and
integration identity, allowed writes, targeted gate, batch gate/triggers,
external effects, exact release target/environment class, smoke/recovery и
отдельные external approval gates.

Disposition precedence для batch:

```text
global-conflict
-> actionable completion-remains
-> resume
-> work-remains
-> deferred-only
-> no-work
```

Global connector/scope/applicable-Goal/ownership/shared-state conflict даёт
`TASK CONTEXT ALARM`. Task-local ambiguity даёт `deferred`; продолжать другие
runnable Tasks. В `single` при material blocker/failure опубликовать и
перечитать truthful report, сохранить non-terminal status и завершить flow, не
выбирая другую Task.

## Никогда не требовать ручную приёмку

- Acceptance criteria — проверяемые completion criteria, не human sign-off.
- Не просить «принять» result, не создавать `acceptance-required` и не оставлять
  terminal-ready Task в `In Review`.
- После passing acceptance, targeted gate, applicable exact review-batch gate,
  integration, required effects и resolved findings автоматически принять exact
  result, опубликовать/read-back `COMPLETED`, перевести Task в `Done` и
  перечитать её.
- Старые memory/rollout/report/Goal/plan правила о ручной приёмке считать
  superseded historical evidence. Никогда не переводить Goal в `blocked` по этой причине.
- Automatic acceptance не заменяет production, destructive-data, secrets,
  privacy, access-policy или external-recipient authority.

## Продолжать автономно и соблюдать release boundary

Перед execution или первым defer/release decision прочитать
[autonomy and release reference](references/autonomy-and-release.md) полностью.

- Не задавать пользователю вопрос посреди runnable queue. Выбирать безопасный
  обратимый in-scope default; material decision/authority blocker defer-ить с
  checkpoint и продолжать независимую work.
- Перед blocking input повторить complete inventory. При
  `runnable_count > 0` blocking input запрещён.
- Для каждого defer обязательно опубликовать/read-back `BLOCKED` handoff. Без
  comment capability оставить Task non-terminal и включить comment delivery в
  blocker; никогда не писать report в `description`.
- ShipTask delivery intent является standing authority для ordinary exact
  verified non-production target: local/dev/test/QA/UAT/staging/preview/sandbox.
  Выполнить required build/deploy/bounded migration/smoke/repair/rollback без
  дополнительного confirmation.
- Никогда не выполнять production release без explicit user authority для
  exact production target/current scope. Unknown environment считать
  production-like. Подготовить non-production evidence, defer-ить с
  `production-approval-required` и продолжить runnable work.

## Выполнить, интегрировать и проверить

`single` всегда serial и использует risk-appropriate singleton review batch.
Для batch выбирать adaptive isolated lanes по dependency, conflict,
verification, integration и review capacity; один writer владеет одной Task,
integration owner выполняет fan-in в dependency order. Не максимизировать
agent count ради числа.

Различать active write target и review batch target. Без фактически запущенных
isolated concurrent workers active write target равен `1`, даже для большого
batch. Вести run-scoped active lane set: подтверждённый `To Do → In Progress`
занимает lane; общий commit/UAT/full gate/comment не делает `In Progress`
waiting room для уже targeted-verified integrated candidate.

Перед каждым новым `To Do → In Progress` выполнить status-reconciliation
barrier. Task текущего run, на которой implementation/rework закончились,
сначала перевести в `In Review` и перечитать; незавершённый partial blocker —
truthful defer; terminal result — terminal read-back. Unknown/unreconciled
status write оставляет lane занятой. Не начинать новую Task, если занятые lanes
достигли active write target; targeted-verified candidate нельзя defer-ить как
partial только ради освобождения capacity.

Для каждой Task:

1. Зафиксировать canonical ref, current detail/version/status, acceptance,
   dependencies, result base и release target.
2. `To Do` начать после preflight; `In Progress` продолжить из coherent
   checkpoint; `In Review` сначала проверять как exact candidate, не
   реализовывать заново без finding.
3. Реализовать минимальный целостный in-scope result и выполнить per-Task targeted gate.
4. Покрыть material duplicate scenarios. Отдельную проблему вынести как
   finding/scope decision без silent work expansion.
5. Интегрировать exact candidate, перевести его в `In Review` и сформировать
   evidence identity. Завершить status write/read-back до освобождения lane и
   старта replacement Task; ожидание общего batch gate/effect этого не отменяет.
6. По batch trigger выполнить independent review и review-batch gate на exact
   integrated members/result. Не запускать дорогой full gate для каждого
   member без risk reason.
7. Выполнить required external effects и независимо проверить exact target.
8. Failed gate локализовать: затронутые Tasks вернуть в `In Progress`, выполнить
   rework и повторные gates. При неясной attribution invalidated весь связанный
   batch.

## Опубликовать Task report

Перед первым report decision прочитать
[delivery-report reference](references/delivery-report.md) полностью.

- Native comment является обязательным terminal effect для completed,
  materially failed/rework или blocked Task.
- Success: `COMPLETED` после checks/effects и до `Done`. Failure/rework/blocker:
  `REWORK REQUIRED`/`BLOCKED` с impact, checkpoint, evidence, cause confidence,
  remaining risk и exact resume decision.
- Использовать stable report identity, искать equivalent comment, выполнять
  read-back. `not-available` или `write-outcome-unknown` блокирует terminal
  transition affected Task.
- Никогда не писать report в description, acceptance, status text или другой
  Task field. Не публиковать `ACCEPTANCE READY`.

## Ограничивать defects и recovery

- Автоматически исправлять только acceptance/project-policy in-scope defect.
- Blocking out-of-scope finding defer-ит affected Task; non-blocking finding
  остаётся final-only. Не создавать follow-up Task без explicit authority.
- При resume перечитать Task, workspace/Git, checks и external state. Старый
  report — checkpoint, не current proof.
- Не выполнять reset/clean/stash/force-push/takeover или удаление чужой work
  при dirty/shared drift. Изолировать affected lane или alarm shared conflict.

## Завершить по mode

Перед blocking input, утверждением «продолжить невозможно» либо terminal Goal
decision прочитать [run-report reference](references/run-report.md) полностью.
Провести causal analysis, проверить доступные safe in-scope recovery actions и
выполнить их самостоятельно. Если имеющаяся operation/authority устраняет
причину, `blocked` запрещён: перечитать state и продолжить workflow.

Для `single` перечитать Task и доказать: exact result, acceptance, targeted и
singleton batch gates, required effects, automatic acceptance, published
report и truthful terminal status. Goal identity — `not-applicable`.

Для `batch` перед terminal claim повторить complete inventory и проверить:

- zero unfinished/deferred in-scope working Tasks и unresolved in-scope defects;
- `Backlog` исключён и показан отдельно;
- final review batch относится к exact final integrated result;
- duplicate scenarios, checks и external effects имеют exact evidence;
- для каждой изменённой/materially failed Task report опубликован и перечитан;
- Task Manager projection и result identities reconciled;
- decision queue пуста.

Если остаётся runnable work — продолжить. Если остаются только deferred Tasks —
показать одну consolidated queue и не завершать Goal. `blocked` применять только
после строгого tool threshold, проведённого causal analysis и user-visible
`SHIPTASK RUN REPORT`. Один reason code, raw error или status write без
понятного объяснения, recovery checks и exact resume step запрещён. При полном
evidence вызвать `update_goal(status="complete")` последним.

Каждый terminal exit (`complete`, `blocked`, partial/deferred, `no-work`)
заканчивается `SHIPTASK RUN REPORT`: сначала простыми словами итог, causal
explanation и выполненный self-recovery, затем scope/statuses, result/effects,
comments, Goal, remaining work и technical evidence. Не выдавать
unknown/skipped effect за verified completion.
