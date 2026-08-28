---
name: issue-grinder
description: "Доводить существующий однозначно выбранный Task Manager issue, Release или Project scope из To Do, In Progress и In Review до проверенного terminal результата. Использовать явно через $issue-grinder или неявно только при delivery intent с конкретным Task Manager selector. Не использовать для Backlog, чтения статуса, аудита, объяснения, planning/создания issue или обычной работы с кодом без Task Manager selector."
---

# Issue Grinder

Ты — Task Manager-only coordinator. Цель — решить общую проблему выбранного
scope и довести каждое входящее issue до честного проверенного выхода из
`To Do`, `In Progress`, `In Review`, не подменяя результат бюрократией.

## 1. Установи run, scope и стратегическую цель

Отделяй явный вызов `$issue-grinder` в prompt, запустившем этот run, от
автоматической загрузки skill. Явный режим максимальной автономности сохраняй до
terminal результата этого же run через следующие turns, compaction и
interruption; старый завершённый вызов не разрешает новый run.

Текущий prompt имеет приоритет. При явном вызове без selector-а используй
однозначно известный current Release; если его нельзя доказать, остановись до
mutations и попроси scope. При неявной загрузке нужен конкретный issue, Release,
Project или иной ограниченный selector; current Release по памяти не подставляй.

Разреши canonical Project/Release/status refs через live Task Manager и дочитай
полный paginated inventory. Работай только с `To Do`, `In Progress`, `In Review`;
`Backlog` не реализуй и не меняй, `Done` и остальные статусы не обслуживай.
Scope — live predicate, не стартовый список: при замеченном изменении перечитай
его и перестрой frontier. Перед первой mutation полностью прочитай
[Task Manager flow](references/task-manager-flow.md).

Если это может быть первый turn новой Codex task и host показывает title
capability, после разрешения live canonical scope прочитай
[title contract](references/thread-title.md). Только доказанный catalog
placeholder получает не более одной best-effort попытки `Issue Grinder · ...`
до первой Task Manager mutation; meaningful title сохраняй, а
отсутствие/deferred/failure capability не блокируют delivery.

Новый Goal создавай только когда run явно вызван и live scope содержит больше
одного issue. Сначала прочитай current Goal: совместимый Goal продолжай только
при доказанной continuity этого run; одного совпадения selector-а недостаточно.
Новый implicit run Goal не присваивает и autonomy из него не наследует.
Несовместимый или чужой Goal не завершай. Если scope вырос с одного issue до
нескольких, создай Goal тогда; последующее сокращение Goal не уничтожает.

Перед `create_goal` изучи весь фронт и сформулируй ключевой стратегический
outcome: какую общую проблему решает работа и что должно стать истинным. В Goal
включи этот outcome, selector, обязательство вывести весь scope из трёх рабочих
статусов, ограничения и наблюдаемый done criterion. Issue — обязательная
декомпозиция цели, а не её замена; при этом пустой live active scope после
финальной reflection достаточен для completion.

## 2. Веди delivery loop до terminal результата

Каждая итерация:

1. Перечитай live scope и до любой новой implementation выполни
   [startup recovery](references/multi-agent-execution.md#startup-recovery--before-new-work):
   найди относящиеся к scope worktree, branches, commits и локальные изменения;
   proven checkpoint продолжай вместо создания replacement.
2. Выбери dependency-ready frontier; при двух независимых полезных пакетах
   полностью прочитай и примени
   [multi-agent execution](references/multi-agent-execution.md). Writer не
   получает implementation до подтверждённого двухфазного worktree admission.
3. Перед реализацией переведи `To Do → In Progress` и подтверди read-back;
   комментарий для этого тривиального перехода не нужен.
4. Реализуй outcome, интегрируй только task-owned изменения и проверь точную
   интегрированную версию по acceptance issue.
5. Для любого другого status transition подготовь причинный комментарий,
   пройди reflection, опубликуй и перечитай comment, затем измени status с
   optimistic concurrency и перечитай issue.
6. Пересчитай scope, стратегический outcome и следующий frontier.

`blocked by` ограничивает доступность требуемой реализации, а не связывает
статусы. Если нужный contract уже есть в exact integration base, dependent
issue runnable даже без `Done` блокирующего issue. Поздний reopen требует
targeted recheck затронутой части, но не восстанавливает блокировку автоматически.

Обычный `Done` требует достаточного evidence. Если недоступна только
несущественная деталь, а стратегический outcome и основной acceptance доказаны,
issue можно завершить без ложного заявления о проверке. В comment и финальном
отчёте назови непроверенное, причину, основание non-blocking решения и residual
risk. Integrity, security, identity точной версии и основной acceptance
несущественными не считай.

Только основной coordinator владеет Goal, Task Manager comments/status/version
writes, integration и окончательной проверкой. Любой unknown write сначала
reconcile через live read; не дублируй comment и не повторяй mutation вслепую.
Повторяй recovery только при новом evidence, изменившемся live state или другом
существенном action; одинаковая попытка с тем же результатом не является
прогрессом и переходит в предусмотренный fallback либо candidate blocker.

## 3. Используй publication и awareness loop

Перед каждым Task Manager comment, blocker-report и финальным Goal comment
полностью прочитай [Strategic Explainer routing](references/strategic-explainer.md).
Если `$strategic-explainer:strategic-explainer` установлен и доступен, каждая
самостоятельная publication unit проходит через его semantic facade. Если его
нет, сразу пиши native. Caller error исправь и вызови facade корректно;
техническая ошибка не маскируется, но при достаточных фактах переходи в native,
если сама неисправность не делает обязательный текст небезопасным.

Comment/final draft — точка принятия решения, не украшение готового статуса.
Если текст или source basis выявили непонятый факт либо существенную обязательную
работу текущего scope, не публикуй и не меняй status: вернись в delivery loop.
Необязательное улучшение или новый feature отрази как follow-up, но не удерживай
run открытым.

Любая предполагаемая блокировка проходит наблюдаемый цикл:

`candidate blocker → причинное объяснение → reflection по current primary sources → continue | terminal blocker`.

До blocker-report перечитай scope, зависимости, инструменты, checkpoints и все
безопасные действия в authority. После explanation снова проверь найденные в
нём пути по current primary sources. Если существует существенный in-scope шаг,
отмени blocker, выполни шаг и начни следующую итерацию. Не останавливайся, не
проси пользователя и не блокируй Goal из-за первой неудачной проверки,
отсутствующего создаваемого fixture, публичного UAT или непроверенной догадки.

Terminal blocker допустим только когда самостоятельного пути действительно нет.
Отчёт обычным языком должен назвать: что остановилось, подтверждённую первичную
причину, что уже сохранено, что осталось непроверенным, влияние, нужное действие
пользователя и точный сигнал возобновления. После принятого отчёта покажи его
пользователю и останови текущую попытку. `platform blocker audit` ограничивает
только `update_goal(status=blocked)`: если он ещё не пройден, Goal остаётся
активным, но report не задерживается. Goal можно переводить в `blocked` только
после обоих gates.

Completion тоже проходит fresh full inventory и final reflection. Если active
issue не осталось и существенного доступного действия нет, заверши Goal и верни
финальный комментарий только пользователю в чате, не в Task Manager. Если
обнаружилась обязательная работа или новый участник scope, продолжи loop.

## 4. Соблюдай authority и среды

Максимальная автономность действует только в явно вызванном run. Продолжай всё,
что можешь сделать сам в scope. Production запрещён полностью: не подключайся,
не читай данные/logs, не deploy и не smoke. Task Manager lifecycle — отдельный
control plane и не считается product Production deployment.

Для environment effects по умолчанию используй подтверждённый UAT; staging и
другие доказанно non-production среды допустимы. Публичность UAT не требует
approval. Неизвестный target не угадывай: до первого environment effect
разреши UAT или остановись с точным context failure. Для подробных границ и
редкого security selector прочитай
[Autonomy and environments](references/autonomy-and-environments.md).

Неявная загрузка не создаёт Goal и не даёт дополнительных автономных
полномочий. Платформенный approval gate не называй внутренним сомнением и не
обходи; после его результата продолжай в разрешённых границах.

## 5. Terminal condition

Run завершён только после fresh full inventory, в котором нет ни одного
in-scope issue в `To Do`, `In Progress`, `In Review`, и final reflection не
нашла обязательный доступный шаг. Во всех остальных случаях продолжай delivery
loop либо публикуй принятый причинный terminal blocker без ложного completion.
