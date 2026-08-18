# Autonomous continuation и release authority

Использовать этот reference для long-running multi-Task scope, task-local
questions, missing authority и любых release/deploy effects. Automatic terminal
acceptance задана ADR-0005 и канонической specification.

## Содержание

- Run invariant
- Decision ladder
- Automatic acceptance
- Stale acceptance context
- Deferred Task
- Comment handoff
- Shared terminal channel loss
- Non-production release
- Production boundary
- Deferred-only handoff
- Goal blocker report
- Resume

## Run invariant

Не останавливать весь run из-за одного изолированного вопроса. Пока существует
dependency-ready Task, которую можно безопасно продолжать в exact scope,
выбирать safe default либо defer-нуть blocked Task и двигаться дальше.

Не задавать пользователю вопрос посреди runnable queue. Один consolidated
decision request допустим, когда безопасная runnable работа исчерпана.

В уже разрешённом exact scope перед task-local `request_user_input`, финальным
вопросом с ожиданием ответа или иным blocking pause обязательно заново получить
complete Task inventory и вычислить:

```text
runnable_count = actionable To Do + actionable In Progress
               + actionable In Review + actionable in-scope recovery
```

Не включать уже deferred Tasks. При `runnable_count > 0` blocking input
запрещён: сохранить decision, освободить lane и выбрать следующую runnable Task.
Не использовать cached disposition или review precedence как замену fresh gate.
Review packet и published comment являются evidence/handoff, а не разрешением
приостановить оставшиеся implementation/resume lanes или ждать user acceptance.

Global `TASK CONTEXT ALARM` сохранять только для конфликта connector,
mandatory terminal-report channel, exact scope, Goal, ownership,
integration/shared state или authority, из-за которого небезопасна любая
оставшаяся mutation. Отсутствие native comment create/list/read является общей
зависимостью scope, а не task-local blocker: на preflight оно останавливает run
до Goal и delivery mutations.

## Decision ladder

### 1. Выбрать самому

Принять reasonable default без вопроса, если одновременно:

- решение обратимо либо легко корректируется;
- оно локально affected Task и не меняет agreed product outcome;
- acceptance допускает несколько эквивалентных реализаций;
- blast radius, расходы и external effects ограничены;
- решение не касается production, secrets, privacy, destructive durable data,
  legal/financial commitment или внешнего получателя.

Зафиксировать выбранный вариант и rationale в evidence/review report. Не
превращать каждую мелкую реализационную развилку в decision queue entry.

### 2. Defer-нуть Task

Использовать `deferred`, если Task требует:

- material product/architecture choice с разными observable outcomes;
- production approval;
- обязательный approval внешнего approver, прямо заданный Task/project policy;
- out-of-scope change или новую Task authority;
- destructive/secret/privacy/cost/external-recipient authority;
- решения ambiguous requirement после bounded research;
- external state change, access или dependency, локальных этой Task.

Не угадывать и не расширять scope. Изолировать Task, освободить lane и
продолжить другие Tasks.

### 3. Остановить run

Использовать global alarm, только если конфликт нельзя изолировать: exact scope
или Goal не разрешены, connector/ownership не позволяют безопасный write,
mandatory terminal-report channel отсутствует, shared integration state
повреждён/неоднозначен либо все remaining mutations зависят от одной общей
отсутствующей authority.

## Automatic acceptance

Не считать terminal acceptance пользовательским вопросом. Invocation
`$ship-tasks` разрешает автоматически принять exact result, когда полный
evidence set доказывает acceptance criteria, targeted и applicable batch gates,
integration identity, required effects и отсутствие unresolved in-scope finding.

При таком результате не вызывать user input и не оставлять Task в `In Review`
ради human acceptance: опубликовать и перечитать обязательный `COMPLETED`, затем
перевести Task в `Done`, перечитать status/version и продолжить scope. Если
comment create/list/read потерян или unreconciled после mutation, Task остаётся
non-terminal, новый dispatch останавливается и применяется scope-wide recovery
из раздела `Shared terminal channel loss`. После `Done` user reopen или новая
Task запускают обычный последующий rework cycle; прошлый report остаётся
historical checkpoint.

Automatic acceptance не заменяет production approval, destructive/secret/
privacy authority или обязательный approval внешнего approver.

## Stale acceptance context

`acceptance criteria` означает объективные Task completion criteria. Оно не
означает human sign-off и не создаёт пользовательский вопрос.

Любое старое требование explicit user acceptance из memory, rollout summary,
previous report/Goal/plan, cached project context или старой документации
считать superseded historical evidence. Такое правило не может переопределить
current skill, создать `acceptance-required`, оставить terminal-ready Task в
`In Review` или привести Goal к `blocked`.

Если такой blocker уже записан старым run, при resume удалить его из current
decision queue, заново проверить exact result/evidence и выполнить обычный
automatic terminal transition. Не считать повторное чтение того же stale text
новым blocker occurrence.

## Deferred Task

Сохранять truthful non-terminal status:

- `To Do` — execution не начинался;
- `In Progress` — существует partial implementation/rework;
- `In Review` — exact candidate machine-verified, но ждёт material decision,
  production/external approval или недостающий terminal evidence. Обычный
  terminal-ready result не defer-ить: автоматически перевести в `Done`.

Не переводить в `Done`, `Canceled`, `Duplicate` и не изобретать provider status
`Blocked`. Не создавать replacement/follow-up Task без отдельной authority.

В current-run decision queue записать:

```text
Task: <canonical ref and identifier>
Reason: <stable reason code>
Current status/result: <truthful state and exact identity>
Last safe checkpoint: <what is safely complete>
Evidence: <checks/releases/findings already proven>
Recommended default: <one concrete recommendation or not-applicable>
Decision/authority needed: <one exact user/external decision>
Resume step: <first safe next action after unblock>
```

Reason codes включают `production-approval-required`,
`external-approval-required`, `ambiguous-product-decision`,
`out-of-scope-authority`, `external-access`, `unsafe-recovery` и
`shared-dependency`. Reason `acceptance-required` запрещён.

Не переизбирать deferred Task в том же run без нового evidence, authority или
external state change. Deferred Task продолжает удерживать Goal active.

Новый out-of-scope finding без blocking edge к in-scope Task не является
deferred Task: записать его в final findings, не расширять Goal/scope, не
создавать follow-up и не запрашивать решение в текущем run.

## Comment handoff

Для каждого defer обязательно опубликовать и перечитать `BLOCKED`
delivery-report comment до освобождения lane. Включить все поля decision queue,
user impact/remaining risk и report key для exact Task/result.

Если при доступном общем channel один write outcome unknown, не использовать
`description`/другой field как fallback и не повторять write вслепую. Оставить
affected Task non-terminal и перейти к scope-wide reconciliation по следующему
разделу.

## Shared terminal channel loss

До Goal и первой delivery mutation native comment create/list/read должны быть
доступны и authorized. При preflight failure выдать `TASK CONTEXT ALARM` с
reason `terminal-report-channel-unavailable`; не начинать Tasks, code/Git,
deploy или runnable queue. Production approval не заменяет эту capability.

Если channel потерян либо стал unreconciled после mutation:

1. Прекратить новый Task dispatch и external effects, не нужные для safe
   recovery.
2. Reconciliate active lanes в truthful statuses/read-backs.
3. Повторить complete inventory, сохранить exact result/effect identities и
   перечислить affected Tasks с comment disposition.
4. Сформировать scope-wide `GOAL BLOCKER REPORT`; не продолжать остальные Tasks
   как независимые, потому что terminal dependency общая.
5. Оставить Goal active до строгого blocker threshold. Повторный poll того же
   state сам по себе не является новым blocker occurrence.

Узкое bootstrap-исключение разрешено только для exact `single` Task, чьи
acceptance criteria прямо восстанавливают native comment create/list/read.
Выполнить только её минимальный result, оставить non-terminal checkpoint и
остановиться до fresh capability preflight. Project/Release batch этим
исключением не продолжать.

## Non-production release

Считать `$ship-tasks` standing authority для обычного in-scope release workflow
только после доказанной классификации exact target как local, development,
test, QA, UAT, staging, preview или sandbox.

Без дополнительного confirmation выполнять необходимые:

- build/package и publish candidate;
- deploy/redeploy exact validated result;
- task-required non-production schema/data migration при bounded recovery;
- environment smoke и terminal-relevant runtime checks;
- bounded logs/metrics diagnosis;
- repair, rollback или restore затронутого in-scope non-production surface.

Не останавливаться с вопросом «деплоить ли на UAT/dev». Если обычный
non-production release безопасно нужен для terminal evidence, выполнить его и
проверить exact target. Временную in-scope поломку диагностировать и
исправить/восстановить до перехода к другой независимой external mutation.

Standing authority не разрешает permanent deletion, destructive shared/durable
data reset без recovery, secrets exposure/rotation, unrelated cleanup,
неограниченные расходы или действие вопреки explicit read-only request.

## Production boundary

Production release требует explicit user approval, однозначно относящегося к
production target и текущему scope. Не считать approval:

- Task/Release title или acceptance text;
- Goal objective;
- успешные checks, UAT deploy или automatic acceptance;
- прошлый production approval для другого result/target;
- наличие production pipeline либо deploy tool;
- сам invocation `$ship-tasks`.

Без approval выполнить безопасную preparation и non-production release/smoke,
оставить Task в truthful status, defer-нуть с
`production-approval-required`, обязательно написать `BLOCKED` comment при
доступных comments и продолжить другие Tasks. Не спрашивать approval посреди
runnable queue.

Если environment не доказан как non-production, считать target production-like
и применять тот же defer. Не выполнять exploratory production mutation.

## Deferred-only handoff

Когда runnable queue исчерпана:

1. Повторить complete inventory и доказать `runnable_count = 0`.
2. Перечитать deferred Tasks и доступные comments/external state.
3. Удалить из queue entries, которые получили новое доказанное решение.
4. Сгруппировать оставшиеся по reason/dependency.
5. Показать один consolidated decision request с recommended defaults.
6. Оставить plan и Goal незавершёнными. Не вызывать Goal completion.
7. Применять Goal `blocked` только после строгого model-tool threshold, а не
   после первого defer или одного turn без progress, и только после полного
   `GOAL BLOCKER REPORT`.

## Goal blocker report

Непосредственно перед `update_goal(status="blocked")` показать один
user-visible ledger:

```text
GOAL BLOCKER REPORT
Reason: <one stable shared reason code>
Scope: <exact Project/Release/Task refs>
Affected Tasks: <count, identifiers and truthful statuses>
Inventory: <terminal/working/excluded counts; runnable_count>
Last safe checkpoint: <source/deploy/artifact identities and proven effects>
Missing effect/evidence: <what prevents terminal reconciliation>
Recovery checks: <what was tried and what each result proved>
Why this session cannot continue: <exact tool/authority/external boundary>
Task comments: <published/read-back or not-available with affected identifiers>
Resume step: <one exact first action after unblock>
```

Goal tool не хранит отдельное reason field. Поэтому `blocked` status без этого
ledger недопустим, даже если tool threshold выполнен. Production approval,
Goal objective или повторное чтение inventory не должны подменять missing
reason/report channel.

## Resume

После нового user decision/authority или external state change перечитать Task,
comment thread, current `version`, source/integration identity, checks и release
target. Продолжить с `Resume step`, только если старый checkpoint всё ещё valid;
иначе пересчитать Task disposition и evidence с нуля.
