# Client protocol выбора Strategic Explainer

Этот reference описывает communication mode ShipTask для одного отдельно
установленного generic provider и безопасного native writing. ShipTask не
содержит runtime provider и не выбирает его по файлам, памяти или качеству
предыдущего текста.

## Когда выбирать и вызывать

В начале run прочитай live skill catalog и effective user rule. Зафиксируй mode
для следующих publication units; перечитай выбор только после явного изменения
правила или failure выбранного provider-а.

Каждый Task Manager comment, отдельный Task/scope report, blocker explanation и
final является самостоятельной publication unit. Routine chat, progress update
и внутренний draft provider не запускают.

## Матрица выбора

Выбирай первый доступный и разрешённый mode:

1. ordinary `$strategic-explainer:strategic-explainer`;
2. native ShipTask writing.

| Ordinary | Mode |
| --- | --- |
| доступен и разрешён | ordinary |
| отсутствует или отключён | native |

Availability fallback применяется только при выборе mode. Для одной publication
unit вызывай не более одного provider. Ошибка уже выбранного ordinary переводит
run в native mode для этой и следующих units.

Пользователь может отключить ordinary полностью или для выбранных publications.
Общий запрет создавать subagents также исключает ordinary и выбирает native.
Opt-out не разрешает применять внутренний метод отключённого provider-а.

## Ordinary path

Для каждой publication unit создай нового built-in `default` read-only subagent
с `fork_turns="none"`, `model="gpt-5.6-luna"` и
`reasoning_effort="max"`. Не наследуй current model/effort. Compact task
содержит:

- отдельную точную строку `STRATEGIC_EXPLAINER_PROVIDER_V1`;
- identifier `$strategic-explainer:strategic-explainer`;
- одну короткую user-facing formulation task;
- exact scope;
- resolvable read-only anchors к current Task/relations, evidence/candidate,
  bounded session state и project/repository sources.

Не передавай inherited conversation, tool transcript, process diary, caller
analysis/rationale, strategic summary, format rules и прежний candidate. Не
читай ordinary provider-only entrypoint или internal contract и не составляй
текст за provider-а.

Если exact Luna Max profile недоступен, не заменяй его скрыто на Sol или другой
profile: считай ordinary provider недоступным и переводи run в native mode.

`STRATEGIC_EXPLAINER_INVOCATION_ERROR` из-за structural defect исправь одним
новым clean subagent; follow-up старому запрещён. Повторный structural refusal
или другой provider failure переводит mode в native. Factual conflict получает
corrected source/anchor и новый clean invocation того же ordinary provider, пока
он остаётся исправимым factual conflict, а не provider failure.

## Native path

ShipTask сам формулирует publication unit по собственным current
truth/lifecycle/reporting requirements и authoritative facts. Он не читает,
загружает, применяет или имитирует provider-only method ordinary и не заявляет
эквивалентное качество.

Native является нормальным mode, а не degraded capability incident. Не добавляй
в comment или final предупреждение только из-за отсутствия, opt-out или failure
Explainer. Обязательный comment всё равно публикуется и перечитывается; после
этого разрешённый lifecycle transition продолжается.

## Общая обработка результата

В provider mode ShipTask принимает готовый text и отдельно обозначенный source
basis либо явный refusal/failure. Material claims сверяются по authoritative
sources; publication не переписывается caller-ом. В native mode ShipTask
формулирует text сразу по собственному contract.

## Reflection до blocker

Готовый provider candidate и source basis ShipTask читает как reflection input.
Он заново проверяет current primary sources, исходную цель, primary/cascade cause
и всю safe in-scope diagnostic/repair/verification/reconciliation frontier.
Provider result не является evidence, authority, scope или status decision.

В native mode ShipTask выполняет тот же frontier review непосредственно по
primary facts перед публикацией blocker-а. Достаточный найденный путь отменяет
stale blocker и работа продолжается. Changed facts аннулируют прежнюю
publication unit; unchanged state не создаёт retry loop.
