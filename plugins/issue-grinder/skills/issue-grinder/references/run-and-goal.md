# Run, scope и Goal

Применяет `IG-GOAL-01..07`, `IG-SCOPE-01..03`, `IG-UI-01` и continuity-часть
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

Goal завершается только после fresh full inventory с пустым active scope и
final reflection, не нашедшей обязательной доступной работы. Нетерминальный
checkpoint `Экономичного` режима сохраняет Goal активным. Terminal blocker
может изменить Goal только после publication и platform blocker gates из
[Strategic Explainer routing](strategic-explainer.md).
