---
name: strategic-explainer
description: "Сформулировать один понятный, проверяемый и готовый к публикации comment, report, blocker explanation, final или явно отредактированный target text через изолированный Strategic Explainer. Не использовать для routine chat, mutations, authority decisions или orchestration. При активной модели Astra (gpt-6-astra) не вызывать автоматически: Astra сама формулирует текст; явный прямой запрос пользователя остаётся допустимым."
---

# Strategic Explainer

Это семантический facade к изолированному provider-subagent, а не набор
редакторских инструкций для текущего агента. Внешний workflow передаёт только
смысловую задачу; topology, profile, clean invocation и retry принадлежат этому
skill. Методика улучшения текста принадлежит только terminal provider-у.

## Astra guard для автоматических вызовов

Если активная модель — Astra (`gpt-6-astra`), этот skill не выбирается для
автоматической или delegated publication unit: вызывающий coordinator должен
сформулировать текст native сам и не запускать provider path. Явный прямой
пользовательский запрос к `$strategic-explainer:strategic-explainer` остаётся
отдельным разрешённым вызовом.

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

Source basis — вторая opaque-часть provider result, а не черновик публикации. В
нём могут намеренно оставаться точные audit-only значения или разрешимые anchors,
которые удалены из пользовательского текста. Facade и внешний workflow передают
его без смыслового сокращения: не заменяют значения общей отметкой об
исключённых деталях и не строят собственный basis. Непубликация basis не означает
его потерю.

Для незавершённого результата или невозможности продолжить publication text сам
связывает фактическое состояние и наблюдаемый result с причиной или boundary,
оставшимся условием, доступным действием и continuation signal; source basis не
заменяет эту причинность. Точное human-facing название остаётся рядом только
когда нужно читателю для различения или навигации. Если current sources оставляют
безопасный repair, retry или self-service path, result не утверждает, что
возможность продолжить исчерпана.

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
   готовый result contract: обе части provider result без смысловой переработки;
   не сокращай и не реконструируй basis, даже если он содержит audit-only
   значения, которые не войдут в публикацию. Не раскрывай clean-call recipe и
   не выдавай промежуточный structural refusal за пользовательский текст.

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
decision/state/obstacle explanation, final или явную редактуру. Routine chat,
progress update, planning, implementation, mutation, lifecycle, authority
decision и внутренняя оркестрация не относятся к этому API.
