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

Только основной агент выполняет Task Manager comment/status/version write.
Worker, reviewer, scout и любой другой subagent получают явный запрет на такой
write и возвращают только facts и evidence, verdict и рекомендацию. Прямая
запись субагента не считается выполнением обязательной publication unit:
основной агент перечитывает историю и перед следующим связанным effect
публикует через выбранный mode корректирующий comment, не скрывая уже
произошедшую запись.

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

Для каждой publication unit сделай один semantic call
`$strategic-explainer:strategic-explainer`. Передай только:

- назначение и исходный пользовательский вопрос;
- одну реальную formulation либо explicit editing task;
- exact scope и язык;
- material constraints;
- resolvable read-only anchors к current Task/relations, evidence, bounded
  session state и project/repository sources.

Не передавай explanation candidate, caller analysis/rationale, strategic
summary, format rules или методику улучшения текста. Не выбирай и не упоминай
никакие другие invocation parameters или provider instructions: внутренним
исполнением полностью владеет facade Strategic Explainer.

Facade возвращает готовый text и отдельно обозначенный source basis либо
operational unavailability. Final failure переводит run в native. Factual
conflict получает corrected source/anchor и новый semantic call того же ordinary
provider, пока он остаётся исправимым factual conflict, а не provider failure.

Следующий lifecycle comment, scope-level blocker report и final всегда получают
новую publication unit. Не продолжай и не переиспользуй Task comment provider;
совпадающие facts доступны только как новые read-only anchors.

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

Для blocked Goal до `update_goal(status=blocked)` после fresh full inventory
создай новую scope-level publication unit. Она не переиспользует Task comment
provider и называет первичную причину невозможности продолжать, исчерпанную
безопасную frontier, влияние, отдельно user-controlled и environment-controlled
prerequisites, первый безопасный шаг и наблюдаемый сигнал возобновления. Только
после factual check report допустим Goal write; затем Goal перечитывается, а
report возвращается человеку без замены перечнем симптомов.
