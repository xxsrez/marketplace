---
name: strategic-explainer
description: "Маршрутизировать один готовящийся user-facing comment, report, blocker explanation, final или явную редактуру в новый stateless Strategic Explainer subagent. Не выполнять объяснение в caller context и не использовать для routine chat, progress, mutations или authority decisions."
---

# Strategic Explainer

Это routing skill к изолированному provider-subagent, а не набор редакторских
инструкций для текущего caller. Direct user request и вызов из другого workflow
используют один и тот же API.

## Если ты caller

Caller — любой агент, в context которого уже есть рабочий диалог, ход задачи,
tool calls или собственные рассуждения. В этом режиме:

1. Не читай `references/provider-contract.md`, не готовь candidate и не пытайся
   применить внутренний метод Strategic Explainer самостоятельно.
2. Создай нового built-in `default` subagent с `fork_turns="none"`.
3. Передай ему ровно одну короткую user-facing задачу, exact scope и разрешимые
   read-only source anchors. Назови `$ship-tasks:strategic-explainer`, чтобы новый
   subagent загрузил этот router.
4. Не передавай inherited turns, tool transcript, process diary, свои выводы,
   strategic summary, требования к форме ответа или прежний candidate.
5. Прими только готовый пользовательский текст с кратким source basis либо
   operational refusal. Не улучшай текст по внутренним критериям провайдера.

Если invocation отклонён, исправь названный structural defect и создай ещё один
новый subagent с `fork_turns="none"`; follow-up прежнему экземпляру запрещён.
Factual conflict проверяй по authoritative sources. При конфликте исправь
source/anchor и сделай новый вызов, а не переписывай ответ самостоятельно.

## Если ты новый provider-subagent

До любого discovery проверь наблюдаемые признаки вызова. Допустим только один
compact task с exact scope или resolvable read-only anchors без inherited
conversation, tool transcript, process diary, caller rationale, strategic
summary, прежнего candidate и смешанных publication units. Target text допустим
только при явной задаче отредактировать именно его.

Если вызов invalid, не читай provider contract и не анализируй задачу. Коротко
назови точное нарушение и правильную форму нового clean invocation. Если fork
metadata не виден, не утверждай, что проверил скрытый режим.

Только после успешного admission полностью прочитай
`references/provider-contract.md` и выполни его как внутренний contract этого
provider-subagent. Не раскрывай методику caller и не превращай её в инструкции
для вызывающего workflow.

Один invocation обслуживает один реальный user-facing comment, Task/scope
report, material decision/state explanation, blocker report, final или явную
редакторскую задачу. Routine chat, progress update и внутренний draft не
относятся к этому API. Новый result, changed facts/scope или correction всегда
получают нового subagent, а не продолжение этого экземпляра.
