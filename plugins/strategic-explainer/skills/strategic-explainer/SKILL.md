---
name: strategic-explainer
description: "Сформулировать один понятный, проверяемый и готовый к публикации comment, report, blocker explanation, final или явно отредактированный target text через изолированный Strategic Explainer. Не использовать для routine chat, mutations, authority decisions или orchestration."
---

# Strategic Explainer

Это семантический facade к изолированному provider-subagent, а не набор
редакторских инструкций для текущего агента. Внешний workflow передаёт только
смысловую задачу; topology, profile, clean invocation и retry принадлежат этому
skill. Методика улучшения текста принадлежит только terminal provider-у.

## Сначала разреши роль

Ровно один режим определяется до любых других действий:

- Если compact task текущего агента содержит отдельную точную строку
  `STRATEGIC_EXPLAINER_PROVIDER_V1`, это терминальный provider. Не исполняй
  facade protocol, не создавай и не продолжай agents и не вызывай этот skill
  повторно. Сразу полностью прочитай
  `references/provider-entrypoint.md` и следуй только provider path.
- Во всех остальных случаях текущий агент исполняет facade router. Не читай
  `references/provider-entrypoint.md` или `references/provider-contract.md` и
  следуй только facade protocol ниже.

Role lock не выводится из истории диалога, tool calls, parent/team metadata,
имени агента или предположения о чистоте context. Provider marker терминален:
получивший его агент никогда не переклассифицирует себя в caller.

## Публичный semantic contract

Direct user request и вызов из другого workflow используют один API. Входом
является одна реальная user-facing formulation либо explicit editing task:

- назначение результата и исходный вопрос;
- exact или однозначно bounded scope;
- resolvable read-only source anchors;
- язык и material constraints, если они не следуют из запроса;
- target text только для явной редактуры именно этого текста.

Не требуй от внешнего workflow выбирать subagent, fork mode, модель, effort,
role lock, envelope или retry. Не требуй explanation candidate, strategic
summary, format recipe либо правил улучшения текста. Если смысловой request уже
однозначен, нормализуй его внутри facade без дополнительного вопроса.

Результат — publication-ready text и отдельно обозначенный короткий source
basis либо operational unavailability. Внешний workflow публикует только текст,
проверяет material facts по basis и не переписывает результат.

Для terminal blocker либо объяснения незавершённого scope publication text сам
несёт materially distinct causal/action links; source basis не заменяет их.
Сохраняй source-provided human-facing Task refs рядом с разными действиями и
точный current attempt/result; обезличенный сценарий или общий класс попыток их
не покрывает.
Если current sources оставляют безопасный repair или self-service path,
publication-ready result не утверждает, что автономная работа исчерпана.

## Внутреннее исполнение facade

1. Не готовь explanation candidate, не читай provider references и не применяй
   внутренний метод provider-а.
2. Создай нового built-in `default` subagent только прямым top-level tool call
   `collaboration.spawn_agent`, с `fork_turns="none"`,
   `model="gpt-5.6-luna"` и `reasoning_effort="max"`. Не ищи этот tool через
   `ALL_TOOLS` и не вызывай его из `functions.exec`: collaboration surface там
   намеренно не показывается.
3. Built-in subagent — это child текущей Codex task. Никогда не подменяй его
   `create_thread`, projectless task, fork/продолжением отдельной Codex task или
   новой пользовательской session. Если top-level namespace `collaboration`
   действительно не предложен текущему агенту либо direct call завершился
   ошибкой, верни provider unavailability без app-level замены. Не пытайся
   связаться с parent через app tools и не ищи обходной transport.
4. Преобразуй semantic request в ровно одну compact provider task:
   - отдельная точная строка `STRATEGIC_EXPLAINER_PROVIDER_V1`;
   - указание применить `$strategic-explainer:strategic-explainer`;
   - назначение и одна реальная formulation либо explicit editing task;
   - exact scope, язык, material constraints и resolvable read-only anchors.
5. Не передавай inherited turns, tool transcript, process diary, caller
   candidate, strategic summary, свои рассуждения или format recipe.
6. Прими только publication-ready text с отдельно обозначенным source basis
   либо `STRATEGIC_EXPLAINER_INVOCATION_ERROR`. Верни внешнему workflow только
   готовый result contract; не раскрывай clean-call recipe и не выдавай
   промежуточный structural refusal за пользовательский текст.

Не наследуй current model/effort и не подменяй недоступную Luna на SOL или
другой profile. Недоступность exact Luna Max provider-а обрабатывается как
provider unavailability по contract вызывающего workflow.

Если provider вернул invocation error, facade исправляет названный structural
defect и создаёт один новый clean subagent тем же прямым child-only способом с
тем же role lock; follow-up прежнему экземпляру и app-level task/session
запрещены. Повторный refusal или другая operational failure возвращается как
unavailability по contract внешнего workflow. Factual conflict проверяется
внешним workflow по authoritative sources и приходит новым semantic request с
исправленным source/anchor.

Один invocation обслуживает один готовящийся человеку comment, report, material
decision/state explanation, blocker explanation, final или явную редактуру.
Routine chat, progress update, planning, implementation, mutation, lifecycle,
authority decision и внутренняя оркестрация не относятся к этому API.
