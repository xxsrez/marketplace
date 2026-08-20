# Strategic Explainer handoff для ShipTask

Прочитать перед каждым новым ShipTask Task report comment и перед terminal
user-facing handoff. Общий behavioral contract находится в sibling-skill
[`$strategic-explainer`](../../strategic-explainer/SKILL.md).

## Когда вызывать

Вызвать свежий Strategic Explainer обязательно:

- перед каждым новым Task report comment: `COMPLETED`, `REWORK REQUIRED`,
  `BLOCKED` или `CANCELED`;
- перед scope-level `PARTIAL` или `BLOCKED` user-facing handoff;
- для продолжения действительно нужен user action, decision, новая authority,
  другой человек/аккаунт или external state change;
- для aggregate batch result, который не представлен одним task-scoped brief.

Простой success не освобождает Task comment от Explainer pipeline. Его exact
task-scoped brief можно переиспользовать в single-Task chat report при
совместимом material meaning. Planned comment/read-back и terminal status
reconciliation не делают его stale, если совпали с next-state contract и не
выявили drift. Обычная промежуточная red/green iteration не создаёт comment и
поэтому не является trigger.

## Сформировать Technical Brief

Перед spawn перечитать exact current state. Передать только:

```text
Target surface: TASK_COMMENT | RUN_REPORT
Artifact scope: <exact Task ref or aggregate run identity>
Authoritative report state: <state already chosen by ShipTask>
Audience and goal:
Reader purpose:
Confirmed outcome:
Unfinished or unknown:
User impact:
Evidence and confidence:
Current capability and attempts:
Constraints:
Candidate user dependency:
Next-state contract:
Useful identifiers:
Output language and channel:
Required authoritative envelope:

Scenario:
Expected user-visible behavior:
State: VERIFIED | FAILED | UNVERIFIED | NOT_APPLICABLE
Evidence basis:
User impact:
Needed input:
```

Это bounded handoff prompt, а не схема результата. Удалить неприменимые поля и
изложить факты связным текстом, если так яснее; не добавлять full conversation,
tool transcript или process diary.

`Candidate user dependency` включать только когда исходное evidence уже
показывает, что основной агент не может продолжить без человека или внешнего
state. Не просить Explainer решать, существует ли blocker, какой status выбрать
или какую authority считать достаточной.

`Reader purpose` формулирует, что пользователь должен понять или суметь сделать
после чтения. `Next-state contract` называет actor, минимальное действие,
причину, observable success signal и что ShipTask сможет продолжить после него.
Для material или multi-scenario ситуации заполнять отдельный scenario ledger;
`Needed input` не является разрешением просить пользователя и появляется
только для подтверждённой dependency.

`Target surface` и `Artifact scope` не дают Explainer authority над comment или
run. `Authoritative report state` уже выбран ShipTask. `Required authoritative
envelope` перечисляет refs/fields, которые ShipTask добавит после adaptation и
которые Explainer не должен превращать в process narrative.

Не передавать full conversation, process diary, raw logs, intended wording или
готовый вывод. Если два ограничения независимы, описать их отдельными
сценариями, а не одним техническим абзацем.

## Запустить свежий субагент

Использовать built-in `default` agent с уникальным task name
`strategic_explainer_<surface>_<scope>` и без унаследованной истории.
В orchestration API передавать точный `fork_turns="none"`; `fork_turns="all"`,
положительное число fork turns и продолжение старого thread запрещены. Не
задавать model override: наследовать текущий model и reasoning effort. В initial
task явно потребовать:

1. применить `$strategic-explainer`;
2. не выполнять writes, recovery, status/Goal decisions или external actions;
3. обработать только переданный Technical Brief и не вызывать tools;
4. до анализа проверить, что перед current handoff нет более ранних
   user/assistant turns или tool transcript;
5. при загрязнении вернуть только `CONTEXT_INTEGRITY_ERROR` с инструкцией
   перезапустить default subagent через `fork_turns="none"`;
6. при чистом контексте вернуть свободное стратегическое объяснение как
   смысловую основу для родителя, не structured result и не copy-ready comment.

Дождаться результата. Не переиспользовать старый Strategic Explainer thread для
другого состояния: чистый context является частью design. Если после brief был
выполнен recovery или изменилось evidence, сформировать новый brief и запустить
новый субагент с уникальным task-name suffix.

Если получен `CONTEXT_INTEGRITY_ERROR`, не переходить к local adaptation и не
делать comment/status/Goal writes. Исправить собственный вызов и один раз создать
новый default subagent с `fork_turns="none"`. Повторный context-integrity отказ
останавливает report workflow до любых Task Manager mutations и сообщается как
внутренняя orchestration failure, а не blocker Task.

## Использовать результат безопасно

Выполнить два независимых прохода: forward trace сверяет каждое существенное
утверждение стратегического объяснения с исходным evidence; reverse coverage
проверяет, что каждый decision-relevant факт handoff сохранён либо осознанно
исключён как не влияющий на outcome, impact/risk, action или confidence.
Strategic Explainer не создаёт facts, authority, lifecycle status или решение.
ShipTask самостоятельно выбирает recovery, user request, Task/Goal transition
и terminal outcome по действующему contract.

ShipTask обязан прочитать объяснение и самостоятельно написать Task comment или
user-facing handoff своими словами. Он может сокращать, перестраивать и
адаптировать форму под channel, но сохраняет outcome, impact, причинную границу,
confidence, подтверждённую dependency и next state. Нельзя копировать ответ
механически, противоречить ему, добавлять неподтверждённый смысл или заменять
его собственной technical summary/process diary.

Authoritative envelope, exact status и нужные evidence refs ShipTask добавляет
по исходным фактам. В пользовательском тексте не рассказывать о субагенте или
внутренней orchestration. Технические детали и identifiers оставлять только
когда они помогают человеку понять вывод, проверить его или выполнить действие;
не переносить полный execution transcript.

Переиспользовать brief между `TASK_COMMENT` и `RUN_REPORT` можно только когда
audience, artifact scope, facts, state, user dependency и next action совпадают
по смыслу. Planned comment/read-back и terminal status reconciliation допустимы,
если они совпали с next-state contract и не выявили drift. Recovery, divergent
write или другое material meaning change делают brief stale. Для aggregate
batch, нескольких blocked Tasks или другого audience запустить новый scope-level
Explainer.

Если subagent tools или `$strategic-explainer` действительно недоступны, не
создавать из этого новый Task/Goal blocker и не скрывать исходный outcome.
Применить тот же смысловой contract самостоятельно, выполнить те же
forward/reverse checks и зафиксировать `degraded-adaptation` только во
внутреннем evidence. `CONTEXT_INTEGRITY_ERROR` не является недоступностью и
всегда требует исправленного fresh invocation.
