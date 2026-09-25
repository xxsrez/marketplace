# Run, scope и Goal

Применяет `IG-GOAL-01..09`, `IG-SCOPE-01..03`, `IG-UI-01` и continuity-часть
`IG-AUTO-01`. Читай этот reference только после того, как чистый справочный
запрос исключён и Issue Grinder действительно начинает либо продолжает
delivery.

## Invocation и continuity

Отделяй явный `$issue-grinder` в prompt, запустившем текущий run, от
автоматической загрузки skill. Явный run сохраняет максимальную автономность до
своего terminal результата через turns, compaction и interruption. Старый
завершённый вызов не разрешает новый run; одного совпадения selector-а для
continuity недостаточно.

Текущий prompt имеет приоритет. Явный вызов без selector-а может использовать
однозначно доказанный current Release. Если он не доказан, остановись до
mutations и попроси scope. Неявной загрузке нужен конкретный issue, Release,
Project или иной ограниченный selector; current Release по памяти не подставляй.

Режим разрешается один раз для нового run до стратегической декомпозиции по
[Execution modes](execution-modes.md). Восстановленный mode record не
пересчитывается после compaction, interruption или смены модели.

До Goal, title и иных mutations обязателен current-root admission `IG-MODE-20`
из SKILL.md и execution-modes, в том числе при восстановленном режиме.

## Canonical live scope

Разреши Project, Release, statuses и relations через live Task Manager и
дочитай весь paginated inventory. Работай только с `To Do`, `In Progress` и
`In Review`; `Backlog` не реализуй и не меняй, terminal и foreign statuses не
обслуживай.

Scope хранится как live predicate, а не как стартовый список. При замеченном
изменении membership, после packet result, перед новой dispatch wave и перед
terminal decision перечитай inventory и перестрой frontier. Точные lifecycle,
read-back и transaction rules находятся в
[Task Manager flow](task-manager-flow.md).

Для каждого существенного implementation/rework, dispatch, integration/review,
blocker и terminal решения восстанови три роли: Strategic Outcome всего scope,
применимые Human Requirements и формальный Agent Plan/issue contract. После
compaction, interruption, resume или material scope change не полагайся на
усечённый локальный контекст. Для single issue без Goal прочитай её parent chain
и выведи общий ориентир из Task/Epic и current scope.

В существующем checkpoint удерживай короткую карту обязательной приёмки всего
run: ID каждого открытого критерия → полученное доказательство с
кандидатом/условиями либо «нет» → оставшийся пробел. После сжатия восстанови
её по current issue contracts и receipts; перед `Done` и остановкой сверь с
ней все открытые ID всего active scope, даже если последнее сообщение о другой
карточке. Если у критерия нет ID, сохрани точную ссылку на его текст; ID не
выдумывай. Ссылка на подробный журнал достаточна, отдельный tracker или новые
Tasks не создавай.

Strategic Outcome направляет trade-offs и помогает понять вклад issue, но не
создаёт Requirement, Task, acceptance, verification или blocker. Известный gap
между завершёнными issue и широким outcome попадёт в final report и может стать
входом отдельного planning flow; текущий run он не удерживает открытым.

Если это может быть первый turn новой Codex task и host показывает title
capability, после разрешения canonical scope прочитай
[title contract](thread-title.md). Только доказанный catalog placeholder
получает не более одной best-effort попытки `Issue Grinder · ...` до первой Task
Manager mutation. Meaningful title сохраняй; отсутствие, deferred или failure
capability не блокируют delivery.

## Goal lifecycle

Новый Goal создавай только когда run явно вызван и live scope содержит больше
одного issue. Сначала прочитай current Goal. Совместимый Goal продолжай лишь при
доказанной continuity этого run; несовместимый или чужой Goal не присваивай и
не завершай. Implicit run Goal не создаёт и autonomy из него не наследует.

Если scope вырос с одного issue до нескольких, создай Goal тогда; последующее
сокращение scope Goal не уничтожает. Перед `create_goal` изучи весь frontier и
сформулируй один стратегический outcome: какую общую проблему решает работа и
что должно стать истинным. Включи selector, обязательство вывести весь scope из
трёх рабочих статусов, ограничения и наблюдаемый done criterion. Issue являются
декомпозицией цели, а не её заменой.

Добавь в objective краткое обязательство: «При завершении, блокировке или
контрольной точке подготовить и проверить итог всего запуска по правилам
<доступный абсолютный путь final-report.md>. Смена статуса Goal не заменяет
отчёт». Разреши [правила отчёта](final-report.md) от фактически загруженного
skill и проверь доступность пути. Обязательство и ссылка сохраняются в continuity
после compaction и смены модели. Для уже существующего Goal без ссылки либо
недоступного старого пути восстанови их в checkpoint текущего run; не пересоздавай
Goal и не обещай редактирование objective неподдерживаемым интерфейсом.
Без Goal используй тот же контекст финализации.

В течение run сохраняй существенные outcomes с evidence, нерешённые дефекты,
ограничения и остаток в существующем checkpoint. Откат или обходной путь не
стирает дефект. Перед каждым видом финализации полностью перечитай
[правила отчёта](final-report.md) и проверь фактически исходящий текст всего run.

Goal завершается только после fresh full inventory с пустым active scope и
final reflection, не нашедшей обязательной доступной работы внутри current
issue contracts. Нетерминальный
checkpoint `Экономичного` режима сохраняет Goal активным. Terminal blocker
требует [Консультант routing](consultant.md) с исключением для ведущей Astra
и может изменить Goal только после publication и platform blocker gates из
[Strategic Explainer routing](strategic-explainer.md).

После первого принятого blocker-handoff сохрани checkpoint fingerprint из
selector/frontier, заблокированного effect, current причин, authority boundary и
resume signal. Автоматическое продолжение Goal без релевантного user signal или
изменения primary state засчитывает следующий platform audit turn по этому
checkpoint, но не повторяет browser/profile discovery, external proof, вопрос
пользователю, publication, Консультанта или Strategic Explainer. До порога закончи turn без
нового handoff; на пороге выполни только отсутствующий `update_goal(blocked)`.
Релевантный resume signal либо изменившийся primary state инвалидирует
checkpoint и возвращает run к live reconciliation.
