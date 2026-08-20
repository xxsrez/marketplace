# Strategic Explainer handoff для ShipTask

Прочитать перед user-facing handoff с material partial/blocked outcome, запросом
user action/authority или сложным technical terminal result. Общий behavioral
contract находится в sibling-skill
[`$strategic-explainer`](../../strategic-explainer/SKILL.md).

## Когда вызывать

Вызвать свежий Strategic Explainer, если выполняется хотя бы одно условие:

- ShipTask собирается сообщить `PARTIAL` или `BLOCKED` из-за material gap;
- для продолжения действительно нужен user action, decision, новая authority,
  другой человек/аккаунт или external state change;
- сложный terminal result содержит несколько независимых сценариев, акторов или
  ограничений;
- без адаптации в user-facing ответ попадут внутренние terms, account types,
  transport paths, tool mechanics или status inventory.

Для простого success, который ясно объясняется одной-двумя фразами, отдельный
model run не нужен. Обычная промежуточная red/green iteration также не является
trigger.

## Сформировать Technical Brief

Перед spawn перечитать exact current state. Передать только:

```text
Audience and goal:
Confirmed outcome:
Unfinished or unknown:
User impact:
Evidence and confidence:
Current capability and attempts:
Constraints:
Candidate user dependency:
Useful identifiers:
Output language and channel:
```

`Candidate user dependency` включать только когда исходное evidence уже
показывает, что основной агент не может продолжить без человека или внешнего
state. Не просить Explainer решать, существует ли blocker, какой status выбрать
или какую authority считать достаточной.

Не передавать full conversation, process diary, raw logs, intended wording или
готовый вывод. Если два ограничения независимы, описать их отдельными
сценариями, а не одним техническим абзацем.

## Запустить свежий субагент

Использовать built-in `default` agent с отдельным task name
`strategic_explainer` и без унаследованной истории (`fork_turns="none"`, когда
этот параметр доступен). Не задавать model override: наследовать текущий model
и reasoning effort. В initial task явно потребовать:

1. применить `$strategic-explainer`;
2. не выполнять writes, recovery, status/Goal decisions или external actions;
3. обработать только переданный Technical Brief;
4. вернуть `User Brief` родительскому агенту;
5. поместить недостающие факты только в `PARENT NOTES`.

Дождаться результата. Не переиспользовать старый Strategic Explainer thread для
другого состояния: чистый context является частью design. Если после brief был
выполнен recovery или изменилось evidence, сформировать новый brief и запустить
новый субагент с уникальным task-name suffix.

## Использовать результат безопасно

Сверить каждое существенное утверждение `User Brief` с исходным evidence.
Strategic Explainer не создаёт facts, authority, lifecycle status или решение.
ShipTask самостоятельно выбирает recovery, user request, Task/Goal transition
и terminal outcome по действующему contract.

User-facing handoff должен сохранить ясную формулировку Explainer, но можно
добавить минимальные evidence refs и фактический status после решений ShipTask.
Не включать `PARENT NOTES` и не возвращать technical jargon, который brief уже
перевёл.

Если subagent tools или `$strategic-explainer` недоступны либо вызов завершился
ошибкой, не создавать из этого новый Task/Goal blocker и не скрывать исходный
outcome. Применить тот же User Brief contract самостоятельно и зафиксировать
degraded adaptation только во внутреннем evidence.
