# Автономность и release authority

Этот reference уточняет границы, но не задаёт пошаговый универсальный flow.

## Продолжение работы

- Самостоятельно выбирай обратимые локальные решения, которые не меняют
  согласованный observable outcome.
- Сам выбирай достаточный способ продвинуть Task и проверить результат.
- Task-local blocker не останавливает независимые Tasks. В
  `batch-implementation` перед ожиданием пользователя перечитай полный current
  inventory live selector, включая Tasks, появившиеся после старта; при наличии
  runnable work продолжай без нового approval.
- Не повторяй неизменившуюся операцию, проверку или poll. Повтор нужен после
  material change в result, tool, environment, access, authority или Task
  contract.
- Нет фиксированного числа попыток. Остановка допустима, когда следующий шаг
  требует новой authority, внешнего state change или небезопасного действия.

Global `TASK CONTEXT ALARM` нужен только при конфликте exact scope, connector,
Goal, shared integration state или общей authority, из-за которого любые
оставшиеся writes небезопасны.

Exact Task/явный список refs — closed selector. Project, Release и resolved
current scope — live selectors: стартовый inventory не замораживает membership.
Matching non-Backlog Task входит автоматически; `Backlog → To Do` является
началом уже разрешённой in-scope работы, а не scope expansion. Goal хранит
identity/predicate live selector, а не стартовый count или диапазон refs.

## Execution topology

Сначала разреши применимые natural-language rules пользователя. Они могут
задавать exact или relative число субагентов, роли, общий либо role-scoped
opt-out, а также condition по длительности, сложности или другому названному
признаку. Root/coordinator не входит в явно названное число субагентов.
Совместимые constraints объединяются; позднее более конкретное rule заменяет
прежнее того же scope. Exact count является обязательным, «побольше» сдвигает
решение к большей полезной delegation относительно automatic baseline, а
condition проверяется по указанной пользователем метрике.

Только без применимого rule после live inventory сам раздели scope на bounded
packets и реши, где delegation реально ускоряет или усиливает результат.
Учитывай dependency-ready работу, writable ownership, доступную изоляцию,
runtime capacity и стоимость fan-in, но не превращай это в обязательную
публичную scheduler формулу.

Несколько безопасных полезных lanes могут идти одновременно. Safety, authority,
useful ownership, worktree isolation и проверяемый fan-in сильнее topology
preference. Не создавай фиктивные subtasks ради quota. Если обязательное rule
конфликтует с этими границами или capacity, не подменяй его молча: сообщи exact
конфликт, фактическую topology и влияние на result.

Основной агент является единственным integration owner и владельцем Goal, Task
Manager comments/status/version writes, целостного candidate и Task attribution.
Каждый одновременно пишущий implementation subagent до первой mutation получает
собственные feature branch и Git worktree для exact Task. Writable worktree
принадлежит одному writer; writer не меняет integration target, чужую branch,
чужой worktree или пересекающийся shared surface. Read-only роли отдельного
worktree не требуют. Только integration owner делает fan-in и проверяет exact
объединённый candidate.

Пользовательский profile override для субагентов имеет приоритет. Без него
`gpt-5.6-luna`/`max` получает только genuinely simple packet: bounded,
self-contained, с ясными acceptance/evidence и без material creative,
architectural, authority/risk или environment uncertainty. Остальные packets и
Strategic Explainer наследуют current model/effort.

Luna не выполняет corrective recovery mutations. При ambiguity,
context/contract conflict, unexpected environment/tool state, scope expansion
или proof gap она прекращает packet, но bounded read-only inspection/read-back
обязательно устанавливает exact completed work, effects, evidence и unknown.
Integration owner reconciles state и сам продолжает тот же packet current
profile; явное ограничение profile субагентов этому не мешает. Повторный cheap
Luna loop запрещён. Если current primary сама Luna и broader context не снимает
uncertainty, сообщи
`luna-escalation=not-available` без скрытой подмены Sol. Ожидаемый результат
проверки, включая намеренно воспроизводимый failure, сам по себе не является
surprising environment.

Для model override выбирай compatible self-contained bounded context; выбранная
coordinator форма context не создаёт unavailability. Genuine unavailable auto
Luna Max допускает current-profile fallback. Явный unavailable user profile не
подменяется и уменьшает соответствующую role/capacity.

Общее «не используй субагентов»/«без субагентов» отключает всех субагентов на
весь run, включая research, implementation, review и comment Explainer.
Role-scoped rule меняет только названную роль: например, запрет implementation
workers не отключает Explainer. Если Explainer отключён user rule, обязательный
comment получает тот же quality contract напрямую от основного агента без claim
независимой проверки. Если effective rule сохраняет Explainer, его
недоступность блокирует только comment-dependent effects. Material capacity gap
или deviation от user rule сообщается явно; внутренняя target/width accounting
не требуется.

## Свобода способа и качество evidence

Выбирай, меняй и сочетай инструменты по ситуации. Отказ одного средства не
создаёт обязанности чинить именно его и не является готовым выводом о Task.
Решение оценивается по тому, доказывает ли итоговый evidence current acceptance.

До `verification-blocked` пройди автономную self-test frontier: создай сам
обычные synthetic fixtures (PDF, ZIP, PNG, изображения, Markdown), seed data и
доступные temporary/mock state, затем проверь их через поддерживаемые ingress.
Наличие user-provided файла не является blocker, если эквивалентный input можно
сгенерировать и безопасно передать самому. Отдельно проверь, не требует ли
criterion независимого principal, второй authenticated session, внешнего
account, provider-side evidence или access-policy effect: это может быть
настоящей authority boundary, если synthetic substitute в scope отсутствует.

Такой blocker не завершается голым «нужен файл/principal». Перед defer подготовь
через Strategic Explainer grounded decision report с self-service attempts,
primary/cascade cause, recommended test path, material alternatives,
prerequisites/authority, success signal, safe continuation и exact resume
condition. Explainer не выполняет mutation и не принимает status/authority
решение.

Browser/controller/session switch допустим как диагностика, но не является
repair продукта. Если exact candidate или server path уже доказал product
failure, сбой другого browser login либо MFA не отменяет incident и не создаёт
stop condition, пока остаётся безопасная in-scope работа с продуктом.

Не объявляй разные способы равноценными без основания и не снижай требование
ради удобного результата. Если достаточный способ не найден в текущем scope и
полномочиях, объясни границу знания, влияние и условие возобновления. Это итог
приёмки, а не политика использования отдельного инструмента.

Native Task comment нельзя заменить полем description или ответом в Codex.
Существенный lifecycle transition завершён только при фактическом comment и
read-back, но выбор технического пути к этому результату остаётся за агентом.

## Task-local defer

Сохраняй правдивый status:

- `To Do`, если execution не начинался;
- `In Progress`, если есть partial implementation или rework;
- `In Review`, если candidate предъявлен, но приёмка объективно заблокирована
  или current contract требует material решения.

Перед defer опубликуй и перечитай понятный comment. Он сообщает, что уже
установлено, что мешает продолжить, влияние и exact resume condition. Не
создавай provider status `Blocked` и не создавай follow-up Task без authority.
Task-local blocker оставляет применимый Goal массовой имплементации активным.

## Non-production

Обычный необходимый release в local/dev/test/QA/UAT/staging/preview/sandbox
входит в delivery authority после надёжной проверки target. Разрешены build,
publish, deploy/redeploy, smoke и bounded repair/rollback затронутого
non-production surface.

Это не разрешает permanent deletion, destructive durable-data reset без
recovery, secrets/privacy/access-policy changes, external-recipient action,
unbounded cost или unrelated cleanup.

## Production

Production workflow требует явного approval для exact target и candidate.
Без него оставь понятный Task comment, сохрани правдивый non-terminal status и
продолжай независимую работу. Approval не отменяет checks, comment/read-back и
terminal evidence.

Production approval меняет authority, но не Goal policy. Release-only не создаёт
Goal из-за production target, Project/Release selector или количества прочитанных
Tasks. Уже активный Goal допустим только если он возник из массовой имплементации
минимум двух Tasks и production release был его исходным done criterion.

## Resume

После ошибки, interrupted run, смены агента/сессии, нового evidence, authority
или внешнего state сначала восстанови current checkpoint. Перечитай Task,
comments, status/version, acceptance, result identity и affected environment;
старый handoff является указателем, а не current evidence.

До новой implementation surface проверь релевантные Git worktrees, branches,
commits, staged/unstaged/untracked changes и partial effects. Если exact
task-owned worktree содержит unfinished candidate, прежний writer доказанно
остановлен и concurrent owner отсутствует, прими exclusive ownership этого же
worktree и продолжай в нём. Не создавай replacement checkout и не повторяй
готовую работу только из-за новой Codex-сессии. Если осталась только branch, по
возможности восстанови worktree из неё.

Не делай takeover при живом writer, неизвестном ownership, недоказанной связи с
Task или unsafe state. Не очищай и не откатывай такой artifact; reconciliate или
продолжай другую safe работу. В первом содержательном chat update назови
unresolved acceptance incidents и material resume boundary. После adoption
перепроверь existing candidate по current acceptance и доведи обычный lifecycle
до terminal result.
