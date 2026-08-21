---
name: ship-tasks
description: "Доставлять однозначно выбранный Task Manager scope до фактически проверенного результата. Использовать явно через $ship-tasks; implicit invocation — только при delivery intent с exact существующей Task вроде TM-123 либо уже выбранным Task Manager Project/Release/current scope. Одного delivery-глагола недостаточно. Create-and-deliver — только при явной просьбе создать ровно одну Task в Task Manager и сразу выполнить её. Bare $ship-tasks разрешает mode по live inventory. Single и release работают без Goal; Goal нужен только для batch-implementation минимум двух Tasks. Не использовать для чтения, статуса, аудита, объяснения, планирования или backlog capture."
---

# Ship Tasks

Доведи exact Task Manager scope до результата, который подтверждён фактами и
правдиво отражён в Task. Сам выбирай инструменты, порядок работы, реализацию и
проверки. Не заменяй пользовательский результат внутренней процедурой.

## 1. Выбери mode и exact scope

- Exact Task или create-and-deliver одной Task → `single`, без Goal.
- Реальная implementation/rework минимум двух concrete Tasks →
  `batch-implementation`, с Goal.
- Commit/push/publish/deploy/smoke/rollback уже подготовленного candidate →
  `release`, без Goal, включая production release.
- Чтение, аудит, объяснение, planning/backlog capture → не delivery.
- Обычная просьба исправить код/продукт без Task Manager anchor → не ShipTask.

Project, Release, current scope, несколько Tasks и bare `$ship-tasks` — selectors,
а не modes. Определи фактическую работу после live inventory. Чтение, проверка,
приёмка или lifecycle reconciliation нескольких Tasks не создают Goal.

Для bare invocation прочитай
[project memory](references/project-memory.md). Prompt selector имеет приоритет,
но не обновляет memory. Через Task Manager adapter разреши полный exact scope и
перечитай live Task state, acceptance, relations и materially relevant comments.
Memory не заменяет live state.

Создавай Goal только после подтверждения `batch-implementation` минимум двух
Tasks и до первой implementation mutation. `single` и `release` работают без
Goal. Release-only не создаёт, не переиспользует, не ретаргетит и не завершает
Goal только ради release. Release может продолжить уже активный совместимый Goal,
только если был его исходным done criterion; сам release новый Goal не создаёт.
Goal не определяет Task outcome, число попыток или status. Неразрешимый конфликт
scope/shared state/authority → `TASK CONTEXT ALARM` до небезопасных writes.

## 2. Соблюдай обязательные требования

### Правдивый статус и обязательный комментарий

Перед любым существенным status transition опубликуй понятный native Task
Manager comment и перечитай его. Затем измени status и перечитай Task.
Исключение — очевидный старт `To Do → In Progress`.

Комментарий обязателен перед `In Progress → In Review`, `In Review → In
Progress`, `In Review → Done`, reopen из terminal status, новым
`Canceled`/`Duplicate` и необычной lifecycle-корректировкой. Material blocker
также требует comment, даже если status не меняется. Ответ в Codex не заменяет
Task comment; `description` и другие fields не являются fallback.

Комментарий простым языком объясняет, что установлено, почему статус меняется
или сохраняется, что это значит для пользователя, чем подтверждён вывод и что
произойдёт дальше. Пиши на языке пользователя; внутренние reason codes, смесь
жаргона и отчёт о процессе не являются объяснением.

### Доказательство важнее выбранного способа

Сам выбирай инструменты, способы диагностики и приёмки. Сбой одного выбранного
способа не делает его обязательным, не доказывает defect и не создаёт blocker
сам по себе. Можно исправить инструмент, заменить его, объединить несколько
источников evidence или перестроить проверку — в зависимости от ситуации.

Нельзя снижать current acceptance или называть результат проверенным без
достаточного evidence. Если в текущем scope и полномочиях доказать success или
failure не удаётся, это `verification-blocked`: назови границу знания, влияние
и условие, при котором приёмка станет возможной. Фиксированного числа попыток
или обязательного порядка действий нет.

Обязательный Task comment — требуемый результат, а не предписанный tool flow.
Существенный transition завершён только когда comment действительно существует
в Task и перечитан; как обеспечить это, решает агент.

### Сохраняй authority boundary

Не расширяй scope. Production, destructive durable-data changes, secrets,
privacy/access-policy changes, external recipients и unbounded cost требуют
явной authority. Обычные необходимые local/dev/test/QA/UAT/staging/preview/
sandbox effects разрешены после надёжного определения non-production target.

Подробности autonomy и release: [reference](references/autonomy-and-release.md).

## 3. Выполни и проверь result

Реализуй целостный in-scope result и получи evidence, достаточный для current
Task contract и обязательных runtime/external effects. Способ реализации,
диагностики и проверки выбирай сам. Недоступное называй `not-available`, а не
verified. Out-of-scope finding не исправляй и не превращай в новую Task без
authority.

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
- `verification-blocked`: выбранные агентом разумные способы не позволяют
  доказать ни success, ни failure в текущем scope и полномочиях. Comment
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

Передай решаемую проблему, current facts, известное/неизвестное и влияние.
Strategic discovery остаётся read-only и выполняется только по необходимости.
Проверь ответ и сам напиши final comment без упоминания внутренней orchestration.
Обязательное требование — применить Strategic Explainer и получить понятный
comment; способ выполнения этого требования выбирает агент.

Handoff contract: [reference](references/strategic-explainer.md). Требования к
comment: [delivery report](references/delivery-report.md).

## 6. Продолжай автономно и финализируй

Task-local blocker не останавливает независимую runnable работу. В
`batch-implementation` перед ожиданием пользователя проверь весь scope; пока
остаётся безопасная in-scope работа, продолжай её. Применимый Goal остаётся
active, пока implementation scope не завершён фактически.

Перед финальным ответом перечитай affected Tasks, comments, statuses, применимый
Goal и external effects. Затем дай outcome-first `SHIPTASK RUN REPORT`: result и
state, что доказано/не проверено, причина gap и точное условие продолжения. Не
заменяй Task comments этим ответом. Goal `batch-implementation` заверши только
после повторной проверки, что в scope нет незавершённой in-scope работы.
Release-only run Goal не создаёт и не финализирует.
