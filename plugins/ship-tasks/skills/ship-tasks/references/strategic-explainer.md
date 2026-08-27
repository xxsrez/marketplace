# Client protocol выбора Strategic Explainer

Этот reference описывает один communication mode ShipTask для двух отдельно
установленных generic providers и безопасного native writing. ShipTask не
содержит runtime providers и не выбирает их по файлам, памяти или качеству
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
2. Fast `$strategic-explainer-fast:strategic-explainer-fast`;
3. native ShipTask writing.

| Ordinary | Fast | Mode |
| --- | --- | --- |
| доступен | доступен | ordinary |
| доступен | отсутствует | ordinary |
| отсутствует | доступен | Fast |
| отсутствует | отсутствует | native |

Availability fallback применяется только при выборе mode. Для одной publication
unit вызывай не более одного provider. Ошибка уже выбранного ordinary не
запускает Fast; ошибка выбранного Fast не запускает ordinary. В обоих случаях
перейди в native mode для этой и следующих units текущего run.

Пользователь может отключить оба providers, только ordinary, только Fast или
Explainer для выбранных publications. Общий запрет создавать subagents исключает
ordinary: выбери Fast, если он разрешён и доступен, иначе native. Explicit full
opt-out выбирает native. Opt-out не разрешает применять внутренний метод
отключённого provider-а.

## Ordinary path

Для каждой publication unit создай нового built-in `default` read-only subagent
с `fork_turns="none"`. Compact task содержит:

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

`STRATEGIC_EXPLAINER_INVOCATION_ERROR` из-за structural defect исправь одним
новым clean subagent; follow-up старому запрещён. Повторный structural refusal
или другой provider failure переводит mode в native. Factual conflict получает
corrected source/anchor и новый clean invocation того же ordinary provider, пока
он остаётся исправимым factual conflict, а не provider failure.

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

Changed facts, scope или anchors получают новый focused Fast pass. Execution,
admission или unusable-result failure переводит mode в native без ordinary retry.

## Native path

ShipTask сам формулирует publication unit по собственным current
truth/lifecycle/reporting requirements и authoritative facts. Он не читает,
загружает, применяет или имитирует provider-only method ordinary или Fast и не
заявляет эквивалентное качество.

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
