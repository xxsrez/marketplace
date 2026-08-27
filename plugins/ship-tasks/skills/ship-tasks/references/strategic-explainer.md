# Client protocol выбора Strategic Explainer

Этот reference описывает availability-based routing между двумя отдельно
установленными generic providers. ShipTask не содержит их runtime и не выбирает
provider по файлам, памяти или предпочтению автора текста.

## Когда вызывать

Пока effective user rule не отключает Explainer для exact publication unit,
каждый Task Manager comment, отдельный Task/scope report, blocker explanation и
final получает один provider pass. Routine chat, progress update и внутренний
draft Explainer не запускают.

## Как выбрать provider

Перед каждой publication unit прочитай live skill catalog:

1. Если доступен `$strategic-explainer-fast:strategic-explainer-fast`, выбери
   Fast.
2. Иначе, если доступен `$strategic-explainer:strategic-explainer`, выбери
   ordinary.
3. Иначе зафиксируй capability failure.

Выбор действует для всей publication unit. Не вызывай оба provider и не
переходи с уже выбранного Fast на ordinary после factual, comprehension,
source или execution error: ordinary является availability fallback, а не
скрытым quality retry.

Пользователь может отключить оба provider, только Fast, только ordinary
fallback или Explainer для узкой группы publications. Общий запрет создавать
subagents не отключает Fast; при отсутствии Fast он делает ordinary provider
недоступным. Явный полный opt-out не разрешает ShipTask применять методы
отключённых providers.

## Fast path

Текущий ShipTask agent загружает
`$strategic-explainer-fast:strategic-explainer-fast` для одной user-facing task,
exact scope и resolvable read-only anchors. Skill выполняется in-context и сам
читает собственный provider reference. Не создавай для него subagent.

Inherited conversation и tool history физически остаются доступны, но не
считаются evidence. Fast самостоятельно отделяет authoritative sources от
process diary, прежнего candidate, caller summary и гипотез. Результат — готовый
publication text и отдельно обозначенный source basis. ShipTask проверяет
material factual conflict и публикует только text без второго rewrite. Не
заявляй независимость, statelessness или clean-context guarantee.

Changed facts, scope или anchors получают новый focused Fast pass. Если Fast не
может выполнить обязательную publication unit, она остаётся незавершённой;
ordinary provider автоматически не вызывается.

## Ordinary path

Для каждой publication unit создай нового built-in `default` read-only subagent
с `fork_turns="none"`. Передай:

- одну короткую user-facing задачу;
- exact scope;
- resolvable read-only anchors к current Task/relations, evidence/candidate,
  bounded session state и project/repository sources;
- идентификатор `$strategic-explainer:strategic-explainer`.

Не передавай inherited conversation, tool transcript, process diary, caller
analysis/rationale, strategic summary, format rules и прежний candidate. Не
читай ordinary provider-internal reference и не составляй текст за него.

Operational refusal из-за invalid invocation исправь новым clean subagent;
follow-up старому запрещён. Повторный structural refusal после исправления —
orchestration failure. Factual conflict получает corrected anchors и новый
clean invocation.

## Общая обработка результата

ShipTask принимает готовый text и отдельно обозначенный source basis либо
явный refusal/failure. Material claims сверяются по authoritative sources;
publication не переписывается caller-ом. Missing provider или непригодный
result оставляет mandatory comment и зависящий lifecycle effect
незавершёнными; независимая безопасная работа продолжается.

При explicit opt-out ShipTask сообщает только обязательные lifecycle facts по
собственному truth contract и не заявляет эквивалентное provider quality.
Финальный ответ при capability failure всё равно честно показывает фактическое
состояние и сам gap.

## Reflection до blocker

Готовый candidate blocker и source basis ShipTask читает как reflection input.
Он заново проверяет current primary sources, исходную цель, primary/cascade
cause и всю safe in-scope diagnostic/repair/verification/reconciliation
frontier. Provider result не является evidence, authority, scope или status
decision.

Достаточный найденный путь отменяет stale blocker и работа продолжается. Если
blocker остаётся, current publication формируется новым pass выбранного provider.
Неизменившийся blocker получает один reflection pass; повтор нужен только после
material change или исправления invalid ordinary invocation.
