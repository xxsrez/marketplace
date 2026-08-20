# Delivery report как Task comment

Использовать этот reference только после выбора exact Task и формирования
task-specific evidence. Report объясняет результат пользователю; он не заменяет
checks, automatic acceptance decision, source identity или external-effect
verification.

## Содержание

- Capability gate
- Write и reconciliation
- Strategic Explainer composition
- Deferred/BLOCKED handoff
- Общий формат
- Success report
- Material failure report
- Task-specific boundary

## Capability gate

На каждом запуске проверить current Task Manager tool contract.

`available` означает одновременно:

- connector явно предоставляет native comment-create operation для canonical
  Task ref;
- current authority позволяет вызвать этот write;
- operation создаёт Task Manager comment, а не imported/external annotation.

`not-available` означает, что operation отсутствует, unsupported, недоступна по
authority или сообщает, что feature ещё не работает. Imported comments из
`get_task_external_context` — только read-only provenance.

При `not-available`:

1. Не вызывать `update_task` ради report.
2. Не менять `description`, acceptance, status text или другой Task field.
3. Сохранить report и blocker в review/interaction output.
4. Оставить Task truthful non-terminal, классифицировать её как
   `completion-remains`/`deferred` и продолжить независимый workflow.

Этот gap не является global `TASK CONTEXT ALARM`, но блокирует `Done` affected
Task и Goal completion, пока Task остаётся в scope.

## Write и reconciliation

- Публиковать `COMPLETED` report при terminal completion изменённой рабочей
  Task.
- Публиковать `REWORK REQUIRED`/`BLOCKED` report после material failure,
  changes-requested или blocker, который важно объяснить пользователю.
- Перед каждым новым report comment независимо от state выполнить обязательную
  [Strategic Explainer composition](strategic-explainer.md). Сначала ShipTask
  выбирает truthful state и фиксирует evidence, затем Explainer формирует
  человеческий narrative.
- Для каждого task-local defer обязательно публиковать `BLOCKED` report; если
  comments недоступны, сам comment delivery становится частью blocker.
- Не публиковать `ACCEPTANCE READY`: terminal-ready evidence автоматически
  приводит к `COMPLETED`; не писать comment на каждую внутреннюю red/green
  iteration.
- Не backfill-ить старые terminal Tasks и не писать отдельный report в
  `Duplicate` без explicit authority.
- При доступном comment list/read до write искать тот же report key, после
  write перечитывать созданный comment. Не считать один success response
  read-back, если connector предоставляет отдельное чтение.
- При unknown write outcome не повторять create вслепую. Сначала искать report
  через comment list/read; если это невозможно, показать
  `write-outcome-unknown`, оставить Task non-terminal и продолжить только
  независимую работу без duplicate risk.
- Comment write и status update считать отдельными side effects, пока current
  connector явно не гарантирует atomicity.

Для дедупликации использовать стабильный видимый ключ:

```text
Report key: shiptask/<canonical-task-ref>/<STATE>/<exact-result-identity>
```

Не включать secrets, signed URLs или private raw logs в key либо body.

## Strategic Explainer composition

Любой новый ShipTask report comment со state `COMPLETED`, `REWORK REQUIRED`,
`BLOCKED` или `CANCELED` должен содержать narrative, прошедший Strategic
Explainer contract. Это относится и к простому success.

Порядок обязателен:

1. ShipTask выполняет task-level finalization: перечитывает current
   Task/result/effects, проверяет exact evidence и доступный safe self-recovery.
   Затем самостоятельно фиксирует report state, verified/unverified scenarios,
   user impact, evidence/confidence, constraints и уже допустимый next action.
   Любое изменение состояния перезапускает этот шаг и инвалидирует старый brief.
2. Сформировать task-scoped `Technical Brief` с `Target surface: TASK_COMMENT`
   по [handoff contract](strategic-explainer.md). Не просить Explainer выбрать
   state, blocker, authority или terminal transition.
3. Получить `User Brief` из свежего субагента. Перед составлением comment
   выполнить forward trace и reverse coverage.
4. Собрать final comment из двух слоёв:
   - authoritative envelope: report type, `State`, Task, короткий result
     identity и report key; exact evidence остаётся внутренним, кроме
     минимального идентификатора, который действительно нужен человеку для
     навигации или действия;
   - Explainer narrative: outcome, impact, понятная причина/граница, требуемое
     действие и next state.
5. Удалить только явное дублирование. Не возвращать process diary или
   необъяснённый jargon и не менять смысл brief. Не сочинять narrative заново
   из raw evidence после успешного ответа Explainer.
6. Выполнить обычный duplicate search, write и read-back.

Для `TASK_COMMENT` действует presentation gate: narrative обычно не длиннее
1 600 символов, весь comment — не более трёх bullets/строк сверх envelope.
По умолчанию запрещены URL, query parameters, endpoint paths, signed URLs,
`file_id`/`download_url`, полные UUID, хэши, provider/environment IDs, raw tool
errors и полные inventories. Если assembled body нарушает gate, его нельзя
публиковать: нужен новый узкий brief или локальный `degraded-adaptation` с теми
же forward/reverse checks. Raw technical summary не является fallback.

В user-visible body не упоминать Strategic Explainer, субагента, Technical
Brief, delegation или orchestration. Если subagent/skill недоступен, ShipTask
локально применяет тот же contract и фиксирует `degraded-adaptation` только во
внутреннем evidence; это не разрешает опубликовать raw technical summary.

Если Technical Brief противоречив или не содержит decision-relevant факт,
вернуться к finalization. Нельзя компенсировать неполное состояние красивой
формулировкой.

## Deferred/BLOCKED handoff

Deferred report не обязан быть incident postmortem. Его задача — показать один
точный material blocker и позволить следующему run безопасно продолжить.
Task-level Strategic Explainer explanation должен существовать до write. Если
тот же blocker удерживает весь scope, до user-facing blocking handoff и
`update_goal(status="blocked")` требуется также scope-level explanation из
совместимого либо нового brief.

```text
SHIPTASK DELIVERY REPORT
State: BLOCKED
Task: <identifier> — <title>
Result: <exact checkpoint identity or not-started>
Report key: shiptask/<task-ref>/BLOCKED/<checkpoint-identity>

Why this Task was deferred
<Одно понятное предложение: что осталось непроверенным/недоступным и почему
без этого результат нельзя считать завершённым. Внутренний reason code — только
при необходимости и после человеческого объяснения.>

Last safe checkpoint
<Что завершено, проверено и не будет повторяться без причины.>

Impact
<Что остаётся недоступным/незавершённым; production impact только если доказан.>

Recommended default
<Один конкретный вариант и tradeoff либо not-applicable.>

Decision or authority needed
<Одно точное material decision, access или production/external approval.>

Resume step
<Первое безопасное действие после unblock.>
```

В `BLOCKED` comment не переносить список URL, deployment/archive IDs, OAuth
метаданные, Attachment refs или полный журнал smoke. Если человеку нужен
конкретный адрес или ref для следующего действия, оставить только этот один
идентификатор и объяснить его роль.

Для `production-approval-required` указать exact production target и candidate
identity, но не формулировать Task/Goal/успешный UAT как уже выданное approval.
Не задавать вопрос в comment, если ещё остаётся runnable work: comment является
handoff, а consolidated decision request формируется после исчерпания queue.

## Общий формат

Начинать с outcome, а не с process diary. Report должен быть глубоким, но
компактным: сначала понять результат, причины и evidence, затем оставить только
то, что помогает человеку быстро понять Task. Масштабировать detail по task type
и risk, но размер evidence сам по себе не является основанием для user-visible
dump. Plain text является безопасным baseline. Использовать Markdown только
если current comment renderer доказан; Mermaid — только при подтверждённом
rendering.

```text
SHIPTASK DELIVERY REPORT
State: COMPLETED | REWORK REQUIRED | BLOCKED | CANCELED
Task: <identifier> — <title>
Result: <commit/build/deploy/artifact identity or not-applicable>
Report key: shiptask/<task-ref>/<state>/<result-identity>

Outcome and impact
<Task-scoped User Brief: что получил пользователь или какая граница осталась.>

How it works / Why
<Только material часть User Brief; убрать секцию, если она дублирует outcome.>

Evidence
- <не более трёх компактных user-relevant фактов; exact refs — только если
  человеку нужен этот один идентификатор для навигации или действия>
- Batch: <краткий aggregate result либо not-applicable>
- External effect: <verified result или not-applicable>

Limits and next action
<Ограничения, remaining risk и optional verification/reopen guidance.>
```

## Success report

Task-scoped User Brief должен объяснить:

- user outcome и наблюдаемое новое поведение;
- основной runtime/data flow простыми словами;
- ключевые implementation decisions и tradeoffs;
- exact targeted/batch/external evidence;
- limitations, remaining risk и optional user verification/reopen guidance.

Diagram или compact before/after добавлять только если они ускоряют понимание;
сложность Task сама по себе не делает визуализацию обязательной.

## Material failure report

Material failure — это user-visible impact, failed batch/external effect,
rollback/revert, invalidated candidate, repeated rework или настоящий blocker
после execution. Обычный красный тест, найденный и исправленный внутри
implementation loop, сам по себе не требует incident comment.

Task-scoped User Brief должен передать impact, причинную модель и recovery на
уровне пользователя. Authoritative confidence и exact evidence ShipTask
сохраняет в envelope. При необходимости добавить к общему формату:

```text
Impact
<Что доказанно затронуто; если production impact не было, так и написать.>

Detection and causal chain
<Как обнаружили и какая подтверждённая цепочка привела к failure.>

Cause
Confidence: CONFIRMED | PROBABLE | UNKNOWN
<Root cause либо честная граница знания.>

Recovery and prevention
<Что исправлено/reopened, какие checks повторены, что предотвращает повтор.>
```

Писать blameless. Не превращать temporal sequence в доказанную causal chain,
не скрывать uncertainty и не создавать follow-up Task без отдельной authority.

## Task-specific boundary

В каждый comment включать только evidence конкретной Task и краткую shared batch
identity. Не копировать полный batch log всем members. При rework публиковать
delta относительно прошлого user-visible state: исправленные findings, новая
result identity, повторённые checks и remaining gaps.
