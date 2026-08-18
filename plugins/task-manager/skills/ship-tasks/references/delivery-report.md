# Delivery report как Task comment

Прочитать этот reference после разрешения exact delivery scope и до batch Goal
или первой delivery mutation. Capability gate применяется сразу; task-specific
форматы — только после формирования evidence. Report объясняет результат
пользователю и не заменяет checks, automatic acceptance decision, source
identity или external-effect verification.

## Содержание

- Capability gate
- Write и reconciliation
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
- connector предоставляет native list/read для idempotency и read-back;
- current authority позволяет вызвать этот write;
- operation создаёт Task Manager comment, а не imported/external annotation.

`not-available` означает, что operation отсутствует, unsupported, недоступна по
authority или сообщает, что feature ещё не работает. Imported comments из
`get_task_external_context` — только read-only provenance.

При `not-available`:

1. До Goal и delivery mutations выдать `TASK CONTEXT ALARM` с reason
   `terminal-report-channel-unavailable`.
2. Не вызывать Task/code/Git/deploy writes и не продолжать runnable queue.
3. Не менять `description`, acceptance, status text или другой Task field.
4. Показать known scope, missing operations/authority, last safe read-only
   checkpoint и один exact refresh/reconnect/resume step.

Это scope-wide blocker: mandatory report channel общий для всех изменяемых
Tasks. Production approval, passing checks или доступный UI не заменяют native
connector create/list/read. Узкое исключение допускает только exact `single`
Task, которая сама восстанавливает этот channel; после её минимального
non-terminal bootstrap checkpoint всё равно нужен fresh capability preflight.

Если `not-available`/`write-outcome-unknown` возник уже после delivery mutation,
остановить новый dispatch, reconciliate active lanes и сформировать scope-wide
`GOAL BLOCKER REPORT`; не продолжать «независимые» Tasks с общей сломанной
terminal dependency. Gap блокирует `Done` affected Tasks и Goal completion до
восстановления create/list/read и успешного report read-back.

## Write и reconciliation

- Публиковать `COMPLETED` report при terminal completion изменённой рабочей
  Task.
- Публиковать `REWORK REQUIRED`/`BLOCKED` report после material failure,
  changes-requested или blocker, который важно объяснить пользователю.
- Для каждого task-local defer при доступном общем channel обязательно
  публиковать `BLOCKED` report. Shared channel loss обрабатывается scope-wide
  barrier, а не серией недоставленных per-Task reports.
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
  `write-outcome-unknown`, оставить Task non-terminal и остановить новый
  dispatch до scope-wide reconciliation.
- Comment write и status update считать отдельными side effects, пока current
  connector явно не гарантирует atomicity.

Для дедупликации использовать стабильный видимый ключ:

```text
Report key: shiptask/<canonical-task-ref>/<STATE>/<exact-result-identity>
```

Не включать secrets, signed URLs или private raw logs в key либо body.

## Deferred/BLOCKED handoff

Deferred report не обязан быть incident postmortem. Его задача — показать один
точный material blocker и позволить следующему run безопасно продолжить.

```text
SHIPTASK DELIVERY REPORT
State: BLOCKED
Task: <identifier> — <title>
Result: <exact checkpoint identity or not-started>
Report key: shiptask/<task-ref>/BLOCKED/<checkpoint-identity>

Why this Task was deferred
<Reason code и почему safe default здесь недостаточен.>

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

Для `production-approval-required` указать exact production target и candidate
identity, но не формулировать Task/Goal/успешный UAT как уже выданное approval.
Не задавать вопрос в comment, если ещё остаётся runnable work: comment является
handoff, а consolidated decision request формируется после исчерпания queue.

## Общий формат

Начинать с outcome, а не с process diary. Масштабировать detail по task type и
risk. Plain text является безопасным baseline. Использовать Markdown только
если current comment renderer доказан; Mermaid — только при подтверждённом
rendering, иначе text diagram.

```text
SHIPTASK DELIVERY REPORT
State: COMPLETED | REWORK REQUIRED | BLOCKED | CANCELED
Task: <identifier> — <title>
Result: <commit/build/deploy/artifact identity or not-applicable>
Report key: shiptask/<task-ref>/<state>/<result-identity>

Outcome
<Что теперь получил пользователь или какое решение требуется.>

How it works / What happened
<Короткое объяснение main flow либо failure narrative.>

Diagram
<Одна полезная flow/boundary/causal diagram либо compact before/after.>

Evidence
- <acceptance criterion> -> <exact check/result>
- Batch: <identity and aggregate gate>
- External effect: <verified result or not-applicable>

Limits and next action
<Ограничения, remaining risk и optional verification/reopen guidance.>
```

## Success report

Объяснить:

- user outcome и наблюдаемое новое поведение;
- основной runtime/data flow простыми словами;
- ключевые implementation decisions и tradeoffs;
- exact targeted/batch/external evidence;
- limitations, remaining risk и optional user verification/reopen guidance.

Для trivial change вместо diagram использовать compact before/after. Для
non-trivial feature или cross-component change включить одну-две схемы, только
если они ускоряют понимание.

## Material failure report

Material failure — это user-visible impact, failed batch/external effect,
rollback/revert, invalidated candidate, repeated rework или настоящий blocker
после execution. Обычный красный тест, найденный и исправленный внутри
implementation loop, сам по себе не требует incident comment.

Добавить к общему формату:

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
