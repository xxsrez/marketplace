# Opaque client protocol Strategic Explainer

Этот reference описывает только вызов sibling
`$ship-tasks:strategic-explainer`. Внутренняя expertise принадлежит
provider-subagent и не является знанием ShipTask.

## Когда вызывать

Пока effective topology rule не отключает роль, каждый Task Manager comment,
отдельный Task/scope report, blocker explanation и final является новым
publication unit. Routine chat, progress update и внутренний draft Explainer не
запускают.

## Как вызывать

Для каждой единицы создай нового built-in `default` read-only subagent с
`fork_turns="none"`. Передай:

- одну короткую user-facing задачу;
- exact scope;
- разрешимые read-only anchors к current Task/relations, evidence/candidate,
  bounded session state и project/repository sources;
- идентификатор `$ship-tasks:strategic-explainer`, чтобы subagent загрузил
  provider router.

Не передавай inherited conversation, tool transcript, process diary, caller
analysis/rationale, strategic summary, инструкции о структуре или стиле ответа и
прежний candidate. Не читай provider-internal reference и не составляй текст за
Explainer.

## Как обработать result

Допустимы три результата:

1. Готовый пользовательский текст и короткий source basis. ShipTask проверяет
   material factual claims по authoritative sources и публикует текст без
   самостоятельного editorial improvement.
2. Operational refusal из-за invalid invocation. Исправь названный structural
   defect и создай новый clean subagent; follow-up старому запрещён.
3. Capability/source failure. Mandatory comment и зависящий lifecycle effect
   остаются незавершёнными; независимая safe работа продолжается.

Если authoritative fact или anchor изменился, новый publication result получает
новый clean invocation. Caller не использует internal quality checklist для
оценки или ремонта текста. Явный user opt-out также не переносит provider method
в caller: ShipTask сообщает только обязательные lifecycle facts по собственному
truth contract и не заявляет эквивалент Strategic Explainer.

## Reflection до blocker

Готовый candidate blocker и source basis ShipTask читает как independent
reflection input. Он заново проверяет current primary sources, исходную цель,
primary/cascade cause и всю safe in-scope
diagnostic/repair/verification/reconciliation frontier. Result Explainer не
является evidence, authority, scope или status decision.

Если reflection показывает новый путь, ShipTask проверяет его по current
acceptance. Достаточный путь отменяет stale blocker и работа продолжается; новый
user-facing result позже получает отдельный invocation. Если blocker остаётся,
его current publication также является новым unit. Неизменившийся blocker
получает один reflection pass; повтор нужен только после material change или
invalid-call correction.
