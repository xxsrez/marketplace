# Delivery report как Task comment

Использовать этот reference только после выбора exact Task и формирования
task-specific evidence. Report объясняет результат пользователю; он не заменяет
checks, automatic acceptance decision, source identity или external-effect
verification.

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

Начинать с outcome, а не с process diary. Report должен быть глубоким, но
компактным: сначала понять результат, причины и evidence, затем оставить только
то, что помогает человеку быстро понять Task. Масштабировать detail по task type
и risk. Plain text является безопасным baseline. Использовать Markdown только
если current comment renderer доказан; Mermaid — только при подтверждённом
rendering.

```text
SHIPTASK DELIVERY REPORT
State: COMPLETED | REWORK REQUIRED | BLOCKED | CANCELED
Task: <identifier> — <title>
Result: <commit/build/deploy/artifact identity or not-applicable>
Report key: shiptask/<task-ref>/<state>/<result-identity>

Outcome
<Что теперь получил пользователь или какое решение требуется.>

How it works / Why
<Короткое объяснение main flow либо failure narrative.>

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

Diagram или compact before/after добавлять только если они ускоряют понимание;
сложность Task сама по себе не делает визуализацию обязательной.

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
