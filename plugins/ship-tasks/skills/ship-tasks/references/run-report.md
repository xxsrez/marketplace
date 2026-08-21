# Осмысленная финализация и ShipTask run report

Прочитать перед любым terminal outcome, blocking pause или финальным ответом.
Task comment объясняет одну Task; этот reference задаёт finalization pass и
объяснение всего run человеку. Он не меняет Task comment/status rules и не
заменяет checks, read-back или external evidence.

## Finalization pass

Не выбирать terminal outcome по последнему tool result. Сначала построить
целостную и проверенную картину:

1. Сопоставить requested outcome с фактически полученным результатом.
2. Перечитать current Task/Goal/source/effect state и проверить, что scope,
   counts, identities и существенное evidence согласованы.
3. Найти необъяснённые gaps и значимые проблемы. Для каждого material gap
   установить root cause настолько, насколько позволяет evidence; отделить факт
   от hypothesis и указать uncertainty только там, где она влияет на вывод.
4. Проверить current operations, authority и project recovery policy. Если
   безопасное in-scope действие может устранить проблему или завершить
   reconciliation, выполнить его самостоятельно.
5. После любого recovery перечитать affected state, повторить нужные checks и
   начать finalization pass заново. Не писать terminal report по устаревшему
   состоянию.
6. Выбрать outcome только когда состояние устойчиво: `COMPLETED`, `PARTIAL`,
   `BLOCKED` или `NO WORK`.

Этот анализ обязателен и при успехе. Не искать искусственный incident, но
проверить, что достигнут именно обещанный результат, причины ключевых решений
понятны, evidence относится к exact result, а ограничения названы честно.

## Blocker Task и незавершённый Goal

Blocker — фактическое условие, не позволяющее завершить результат. Он остаётся
blocker до устранения, даже если агент уже видит один допустимый способ repair.

Если repair безопасен, находится в scope и уже разрешён, выполнить его один раз,
перечитать результат и продолжить при material state change.

Task `BLOCKED`, run `PARTIAL`/`BLOCKED` и Goal status не равны друг другу.
Task-local blocker оставляет Goal активным и незавершённым. Не повторять тот же
repair, poll, acceptance scenario или Goal turn без нового result, evidence,
environment, access, authority либо Task contract.

Перед blocking handoff дать пользователю понятное объяснение причины и состояния.
Финальный ответ должен отражать фактический Goal status без попытки подогнать его
под Task report.

## Глубокий компактный отчёт

Глубина — это качество понимания, а не объём. Сначала установить факты,
причинность и disposition; затем сжать их до минимальной формы, из которой
человеку сразу понятны:

- что получилось или не получилось;
- какой сейчас статус и почему;
- чем существенные выводы подтверждены;
- осталось ли действие, решение или ограничение.

Начинать с результата. Писать простым языком на уровне задачи пользователя.
Technical terms оставлять только когда без них теряется смысл или проверяемость.
Не перечислять все tool calls, промежуточные ошибки, файлы, Tasks или checks;
использовать counts и исключения, а exact refs — только для полезной навигации.

Структуру адаптировать, а не заполнять механически. Обычно достаточно:

```text
SHIPTASK RUN REPORT — COMPLETED | PARTIAL | BLOCKED | NO WORK

Итог и статус
<Полученный результат или точная граница незавершённого.>

Почему так
<Ключевая причинная модель, важные решения и выполненный recovery.>

Подтверждение
<Минимальный набор evidence, на котором держится вывод.>

Осталось / следующий шаг
<Реальное ограничение и одно действие либо «ничего».>
```

Объединять или опускать секции, которые не добавляют смысла. При простом успехе
может хватить нескольких предложений и короткого evidence list. При сложном
успехе объяснить ключевой flow и tradeoffs. При partial/blocked обязательно
назвать impact, причину или честную границу знания, выполненный bounded repair,
почему агент больше не может продолжить сам и exact resume condition.

Report недопустим, если это только reason code, raw error, status inventory или
фраза о недостающем effect. Он должен передавать уже осмысленный вывод, а не
перекладывать диагностику на пользователя.

## Strategic Explainer

Каждый новый Task report comment проходит отдельный task-scoped
[Strategic Explainer handoff](strategic-explainer.md), включая простой
`COMPLETED` success. Новый субагент получает bounded `Strategic Handoff` с
обязательным `Problem to solve` через `fork_turns="none"`, самостоятельно
выполняет bounded read-only strategic discovery и возвращает свободное
объяснение с source basis. Он не выбирает outcome, report state, recovery или
authority. ShipTask читает результат и пишет окончательный report своими
словами, сохраняя material meaning.

Для terminal chat report:

- при exact single-Task result можно переиспользовать проверенный task-level
  explanation, если audience/facts/state/dependency/next action совпадают; planned
  comment/read-back и terminal status reconciliation не делают его stale, если
  совпали с переданным next-state contract и не выявили drift;
- для aggregate batch, нескольких blockers, material partial или другого
  audience нужен новый scope-level brief;
- до blocking handoff scope-level plain-language explanation уже должно быть
  подготовлено и показано пользователю;
- после recovery, divergent write или другого material meaning change старый
  explanation считать stale.

Если subagent tools или skill недоступны, применить тот же смысловой contract
самостоятельно и отметить `degraded-adaptation` только во внутреннем evidence.
Formatting layer не создаёт новый blocker, но raw reason code или technical
summary не являются допустимым fallback.

`CONTEXT_INTEGRITY_ERROR` не является недоступностью. ShipTask исправляет вызов
и один раз повторяет fresh subagent с `fork_turns="none"`; повторный отказ
переходит в локальную `degraded-adaptation`, не блокируя truthful rework status.

`PROBLEM_CONTEXT_ERROR` также не является недоступностью: ShipTask обязан
передать semantic beneficiary/desired outcome/exact scope, а не заменять задачу
identifier или техническим title. После одного исправленного fresh invocation
неустранённый problem gap становится `task-contract-conflict` и явно
возвращается пользователю.
