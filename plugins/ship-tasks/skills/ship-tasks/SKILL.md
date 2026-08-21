---
name: ship-tasks
description: "Доставлять однозначно выбранный Task Manager scope до фактически проверенного результата. Использовать явно через $ship-tasks; implicit invocation — только при delivery intent с exact существующей Task вроде TM-123 либо уже выбранным Task Manager Project/Release/current scope. Одного delivery-глагола недостаточно. Create-and-deliver — только при явной просьбе создать ровно одну Task в Task Manager и сразу выполнить её. Bare $ship-tasks запускает batch по project memory; single работает без Goal, batch — с Goal. Не использовать для чтения, статуса, аудита, объяснения, планирования или backlog capture."
---

# Ship Tasks

Доведи exact Task Manager scope до результата, который подтверждён фактами и
правдиво отражён в Task. Сам выбирай инструменты, порядок работы, реализацию и
проверки. Не заменяй пользовательский результат внутренней процедурой.

## 1. Выбери mode и exact scope

- Exact Task или create-and-deliver одной Task → `single`, без Goal.
- Project, Release, несколько Tasks, current scope или bare `$ship-tasks` →
  `batch`, с Goal.
- Чтение, аудит, объяснение, planning/backlog capture → не delivery.
- Обычная просьба исправить код/продукт без Task Manager anchor → не ShipTask.

Для bare invocation прочитай
[project memory](references/project-memory.md). Prompt selector имеет приоритет,
но не обновляет memory. Через Task Manager adapter разреши canonical refs,
прочитай все страницы до `hasMore=false`, затем current detail, acceptance,
relations, status, `version` и materially relevant comments. Memory не заменяет
live state.

Создай Goal только для batch после exact scope resolution и до первой non-Goal
mutation. Goal хранит progress, но не определяет Task outcome, число попыток или
status. Несовместимый active Goal либо неразрешимый shared scope/authority →
`TASK CONTEXT ALARM` до mutation.

## 2. Соблюдай обязательные требования

### Правдивый статус и обязательный комментарий

Перед любым существенным status transition опубликуй понятный native Task
Manager comment и перечитай его. Затем измени status с current `version` и
перечитай Task. Исключение — очевидный старт `To Do → In Progress`.

Комментарий обязателен перед `In Progress → In Review`, `In Review → In
Progress`, `In Review → Done`, reopen из terminal status, новым
`Canceled`/`Duplicate` и необычной lifecycle-корректировкой. Material blocker
также требует comment, даже если status не меняется. Ответ в Codex не заменяет
Task comment; `description` и другие fields не являются fallback.

Комментарий простым языком объясняет, что установлено, почему статус меняется
или сохраняется, что это значит для пользователя, чем подтверждён вывод и что
произойдёт дальше. Пиши на языке пользователя; внутренние reason codes, смесь
жаргона и отчёт о процессе не являются объяснением.

### Необходимый инструмент сначала восстанови

Если нужный инструмент отсутствует или не работает:

1. установи требуемую операцию и фактическую поломку;
2. диагностируй и попытайся безопасно восстановить основной способ в текущем
   scope и authority;
3. после material repair повтори исходную операцию;
4. альтернативу используй только после repair attempt и только если она
   доказывает то же требование, не более слабое.

Не продолжай молча так, будто инструмент не нужен. Остановись, когда продолжение
требует новой authority, опасного действия или внешнего state change, и ясно
назови причину, impact и условие возобновления. Фиксированного числа попыток нет;
не повторяй действие без изменившегося состояния.

До material delivery mutations проверь comment create и comment list/read. Если
сломан сам comment channel и его нельзя восстановить, не выполняй существенный
status transition и не начинай новые delivery mutations, которые нельзя будет
объяснить в Task.

### Сохраняй authority boundary

Не расширяй scope. Production, destructive durable-data changes, secrets,
privacy/access-policy changes, external recipients и unbounded cost требуют
явной authority. Обычные необходимые local/dev/test/QA/UAT/staging/preview/
sandbox effects разрешены после надёжного определения non-production target.

Подробности autonomy и release: [reference](references/autonomy-and-release.md).

## 3. Выполни и проверь result

Следуй repository/project instructions, сохраняй unrelated user changes,
реализуй минимальный целостный in-scope result и выбирай проверки, которые
доказывают current Task contract. Установи exact result identity и проверь
required runtime/external effects. Недоступное называй `not-available`, а не
verified.

Если safe in-scope repair способен устранить найденную проблему, выполни его и
перечитай affected state. Out-of-scope finding не исправляй и не превращай в
новую Task без authority.

Когда candidate готов, через Strategic Explainer сформулируй comment о результате
и проведённой проверке, опубликуй и перечитай его, затем переведи
`In Progress → In Review`. Сразу проведи приёмку: `In Review` не является
ожиданием человека.

## 4. Разбери `In Review` по текущим фактам

История изменений acceptance и количество прошлых попыток сами по себе ничего
не доказывают.

- `task-contract-conflict`: current mandatory requirements противоречат друг
  другу или не определяют observable result. Однозначное исправление по accepted
  source выполни и перечитай; иначе оставь понятный comment и сохрани
  `In Review`.
- `verified-failure`: exact candidate в совместимой среде прямо нарушает current
  acceptance. Сначала comment с observed failure, impact и причиной возврата,
  затем `In Review → In Progress`; перечитай Task и продолжай rework в этом же
  run.
- `verification-blocked`: нельзя доказать ни success, ни failure. Сначала
  восстанови нужные инструменты/среду, если это безопасно возможно. Иначе comment
  объясняет границу знания и 2–4 способа приёмки с предпосылками, доказательной
  силой, trade-off, рекомендацией и success signal; Task остаётся `In Review`.
- `verified-success`: acceptance, checks, result identity и required effects
  доказаны. Сначала понятный completion comment и read-back, затем
  `In Review → Done` и Task read-back.

Падение общего batch gate без task-level attribution не доказывает defect каждой
Task. Сначала получи separating evidence; не возвращай весь batch в rework.

## 5. Используй Strategic Explainer для человеческих объяснений

Каждый обязательный lifecycle/blocker comment сформулируй с помощью
`$ship-tasks:strategic-explainer`. ShipTask сначала сам устанавливает facts,
outcome, status, authority и next action; Explainer этого не решает.

Передай semantic `Problem to solve`, current facts, известное/неизвестное,
impact и bounded strategic anchors. Используй свежего subagent с
`fork_turns="none"`; discovery только read-only и по необходимости. Проверь
ответ и сам напиши final comment без упоминания внутренней orchestration. Если
helper недоступен, сначала восстанови его по правилу необходимого инструмента.
Без Strategic Explainer не изображай это требование выполненным и не выполняй
зависящий от comment существенный transition.

Handoff contract: [reference](references/strategic-explainer.md). Требования к
comment: [delivery report](references/delivery-report.md).

## 6. Продолжай автономно и финализируй

Task-local blocker не останавливает независимую runnable работу. В batch перед
ожиданием пользователя перечитай полный inventory; пока остаётся безопасная
in-scope работа, продолжай её. Goal остаётся active при любой `To Do`,
`In Progress`, `In Review`, rework, незавершённом effect или unresolved in-scope
defect.

Перед финальным ответом перечитай affected Tasks, comments, statuses, Goal и
external effects. Выполни доступный safe repair. Затем дай outcome-first
`SHIPTASK RUN REPORT`: result и state, что доказано/не проверено, причина gap,
выполненный repair и точное условие продолжения. Не заменяй Task comments этим
ответом. Batch Goal заверши только после fresh full inventory без незавершённой
in-scope работы.
