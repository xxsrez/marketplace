---
name: strategic-explainer
description: "Маршрутизировать один user-facing comment, report, blocker explanation, final или явную редактуру в fresh terminal Strategic Explainer provider. Не использовать для routine chat, mutations, authority decisions или orchestration."
---

# Strategic Explainer

Это opaque routing skill к изолированному provider-subagent, а не набор
редакторских инструкций для текущего caller.

## Сначала разреши роль

Ровно один режим определяется до любых других действий:

- Если compact task текущего агента содержит отдельную точную строку
  `STRATEGIC_EXPLAINER_PROVIDER_V1`, это терминальный provider. Не исполняй
  caller protocol, не создавай и не продолжай agents и не вызывай этот skill
  повторно. Сразу полностью прочитай
  `references/provider-entrypoint.md` и следуй только provider path.
- Во всех остальных случаях текущий агент является caller. Не читай
  `references/provider-entrypoint.md` или `references/provider-contract.md` и
  следуй только client protocol ниже.

Role lock не выводится из истории диалога, tool calls, parent/team metadata,
имени агента или предположения о чистоте context. Provider marker терминален:
получивший его агент никогда не переклассифицирует себя в caller.

## Opaque client protocol

Direct user request и вызов из другого workflow используют один API:

1. Не готовь explanation candidate и не применяй внутренний метод provider-а.
2. Создай нового built-in `default` subagent с `fork_turns="none"`.
3. Передай ему ровно одну compact task со следующими частями:
   - отдельная точная строка `STRATEGIC_EXPLAINER_PROVIDER_V1`;
   - указание применить `$strategic-explainer:strategic-explainer`;
   - одна реальная user-facing formulation или explicit editing task;
   - exact scope и resolvable read-only source anchors.
4. Не передавай inherited turns, tool transcript, process diary, caller
   candidate, strategic summary, свои рассуждения или format recipe.
5. Прими только publication-ready text с отдельно обозначенным коротким source
   basis либо `STRATEGIC_EXPLAINER_INVOCATION_ERROR`. Публикуй только текст: не
   сливай basis с сообщением и не переписывай provider result.

Если provider вернул invocation error, исправь названный structural defect и
создай один новый clean subagent с тем же role lock; follow-up прежнему
экземпляру запрещён. Factual conflict проверяй по authoritative sources и при
необходимости передавай исправленный source/anchor в новый clean invocation.

Один invocation обслуживает один готовящийся человеку comment, report, material
decision/state explanation, blocker explanation, final или явную редактуру.
Routine chat, progress update, planning, implementation, mutation, lifecycle,
authority decision и внутренняя оркестрация не относятся к этому API.
