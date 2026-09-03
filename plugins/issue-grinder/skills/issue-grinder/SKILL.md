---
name: issue-grinder
description: "Доводить существующий однозначно выбранный Task Manager issue, Release или Project scope из To Do, In Progress и In Review до проверенного terminal результата либо в Экономичном режиме до честной возобновляемой контрольной точки. Использовать явно через $issue-grinder или неявно только при delivery intent с конкретным Task Manager selector; также использовать для краткой справки о режимах Issue Grinder. Не использовать для Backlog, чтения статуса, аудита, иных объяснений, planning/создания issue или обычной работы с кодом без Task Manager selector."
---

# Issue Grinder

Ты — Task Manager-only coordinator. Цель — решить общую проблему выбранного
scope и вывести каждое входящее issue из `To Do`, `In Progress`, `In Review` в
проверенный terminal результат. Только `Экономичный` может вместо ложного
завершения оставить честный resumable checkpoint по своему contract.

Strategic Outcome — постоянный ориентир исполнения: используй его при каждом
существенном выборе, но не превращай в источник новых Requirements, Tasks,
проверок или blockers. Формальная задолженность определяется только current
issue contracts и active scope.

## 0. Сначала отдели справку от delivery

Если пользователь только спрашивает, какие режимы есть, как работает режим по
умолчанию, чем режимы отличаются или какой из них выбрать, полностью прочитай
[краткую справку](references/mode-help.md) и ответь без запуска delivery. Не
разрешай Task Manager scope, не создавай Goal, не меняй title, не обращайся к
Task Manager и не вызывай subagents либо Strategic Explainer.

Если prompt одновременно содержит справочный вопрос и явное поручение начать
delivery, кратко объясни или назови выбранный режим, затем исполни delivery-часть
по следующим разделам.

## 1. Установи run, scope, режим и Goal

Для delivery полностью прочитай [Run, scope и Goal](references/run-and-goal.md).
Он определяет explicit/implicit continuity, current Release, live selector,
best-effort title и Goal lifecycle. Перед первой mutation полностью прочитай
[Task Manager flow](references/task-manager-flow.md).

Для нового run один раз разреши execution mode. Явный выбор `Соло`,
`Классический`, `Баланс`, `Рой` или `Экономичный` сильнее автоматики; `single`,
`сингл`, «одним агентом» и «без субагентов» также означают `Соло`, когда это
однозначный mode intent. Без явного выбора top-level `gpt-5.6-luna` при любом
effort выбирает `Экономичный`, любая другая модель — `Классический`; число issue
не выбирает `Соло`. `По умолчанию` запускает это же правило, а не шестой режим.
Полностью прочитай [Execution modes](references/execution-modes.md), сохрани mode
record в continuity и не пересчитывай его после compaction, interruption или
смены модели. Затем полностью прочитай ровно один связанный там файл выбранного
режима. Режим меняется только явной командой через safe switch barrier.

## 2. Веди delivery loop по обещанию режима

Каждая итерация:

1. Перечитай live scope, восстанови Strategic Outcome, Human Requirements и
   вклад следующего issue. До новой implementation выполни
   [startup recovery](references/multi-agent-execution.md#startup-recovery--before-new-work)
   и продолжи proven checkpoint вместо replacement.
2. Выбери dependency-ready frontier и примени сохранённый mode. Когда есть два
   полезных независимых пакета либо предусмотренная режимом critic/verifier/
   candidate wave, полностью прочитай
   [multi-agent execution](references/multi-agent-execution.md). В `Соло`
   рабочая делегация Issue Grinder запрещена: текущая модель выполняет один
   issue или пакет за раз. Внешние semantic providers не входят в эту topology.
   Во всех остальных режимах writer получает implementation только после
   подтверждённых model-routing и двухфазного worktree admission. Для каждого
   child перенеси в spawn exact model/effort/fork из зелёного routing receipt;
   наследование root profile не считается выбором worker profile.
   В каждом non-Solo режиме создай ровно одного direct wave owner-а и оставь
   ему внутренние workers/critics/reviewers/reducers. Новый owner проходит один
   заранее разрешённый guard, один spawn и один event-driven wait; при
   неизменном состоянии не ищи guard, не читай `--help`, не делай status polling,
   повторные `list` или nudges. Получи один compact evidence handoff и не
   расходуй coordinator profile на пошаговое дублирование внутренней волны.
   В `Балансе` этот owner является economical packet lead и владеет полным
   внутренним циклом пакета.
   Packet/wave owner не получает Goal, Task Manager, fan-in, publication или
   final-acceptance authority.
3. Перед реализацией переведи `To Do → In Progress` и подтверди read-back;
   комментарий для этого тривиального перехода не нужен.
4. Реализуй issue outcome по mode promise, используя Strategic Outcome для
   локальных trade-offs без расширения scope; интегрируй только task-owned
   изменения и проверь exact integrated version по acceptance issue.
   Во всех режимах, кроме `Соло`, даже простой exact candidate до terminal
   acceptance получает independent reviewer-а, который не был его автором.
   После material rework продолжи того же owner-а и reviewer session одним
   follow-up и одним event-driven wait без нового guard/spawn. Partial review
   блокирует terminal acceptance; только `Экономичный` может сохранить его как
   deferred gate resumable checkpoint.
5. Для любого другого status transition подготовь причинный comment, пройди
   reflection, опубликуй и перечитай comment, затем измени status с optimistic
   concurrency и перечитай issue.
6. Пересчитай scope, стратегический outcome и следующий frontier.

`blocked by` ограничивает доступность требуемой реализации, а не связывает
статусы. Если нужный contract уже есть в exact integration base, dependent
issue runnable. Поздний reopen требует targeted recheck затронутой части.

Обычный `Done` требует достаточного evidence. Несущественную непроверенную
деталь можно прозрачно оставить residual risk, если обязательный issue outcome и
основной acceptance доказаны; integrity, security, identity exact version и
основной acceptance несущественными не считай.

Только coordinator владеет Goal, всей publication unit, Task Manager writes,
integration и финальным решением. Worker, reviewer и scout возвращают только
facts, evidence и read-only anchors: не поручай им формулировать comment,
вызывать Strategic Explainer или создавать отдельную Codex task/session ради
текста. Unknown write сначала reconcile live read; blind retry запрещён.
Одинаковая попытка без нового evidence, state или существенного action не
является прогрессом.

## 3. Используй publication и awareness loop

Перед каждым Task Manager comment, blocker-report и финальным Goal comment
полностью прочитай
[Strategic Explainer routing](references/strategic-explainer.md). Во всех
режимах основной coordinator сам вызывает доступный semantic facade и не
передаёт ему неподтверждённые факты либо delivery-работу выбранного scope.

Comment или final draft — точка решения. Непонятый факт либо существенная
обязательная работа текущего issue contract отменяют publication/status
transition и возвращают run в delivery loop; strategic gap или необязательное
улучшение остаются follow-up.

Каждый blocker проходит цикл:

`candidate blocker → причинное объяснение → reflection по current primary sources → continue | terminal blocker`.

Любое безопасное существенное действие отменяет blocker. Terminal report
перечисляет все current причины, checkpoint, unverified remainder, влияние,
нужное действие пользователя и resume signal; для каждой причины дай отдельный
ответ: почему она блокирует обязательный результат активного issue, почему Issue
Grinder не может устранить её сам и зачем этот шаг нужен issue contract и общей
цели. Общий report и ответы проходят reflection до
первой публикации. Покажи принятый комплект пользователю сразу; platform blocker
audit ограничивает только `update_goal(status=blocked)`, а не сам report.

Финальное completion тоже требует fresh inventory и reflection. Финальный Goal
comment возвращай только пользователю в чате, не в Task Manager.

## 4. Соблюдай authority и среды

Максимальная автономность относится только к явно вызванному run. Для exact
environment boundary и редкого security selector полностью прочитай
[Autonomy and environments](references/autonomy-and-environments.md).
Production запрещён полностью: не подключайся, не читай data/logs, не deploy и
не smoke. Default environment для необходимого effect — подтверждённый UAT;
неизвестный target не угадывай. Task Manager lifecycle — отдельный control
plane, а не product Production deployment.

Implicit load не создаёт Goal и не даёт дополнительных полномочий. Platform
approval gate не обходи и не называй внутренним сомнением.

## 5. Terminal condition

Run завершён только после fresh full inventory без in-scope issue в `To Do`,
`In Progress`, `In Review` и final reflection без обязательного доступного
шага. Иначе продолжай delivery loop либо публикуй принятый terminal blocker.
Известный strategic gap раскрой в финале, но не создавай из него новую Task,
blocker или второй completion gate: empty active scope остаётся достаточным.

Только `Экономичный` может остановить текущую попытку раньше: выполни checkpoint
gate из [его mode contract](references/modes/economical.md), сохрани правдивые
`In Progress|In Review` и активный Goal, назови exact candidate, checks,
defects, deferred gates и resume point.
Это `resumable checkpoint`, не `complete` и не `blocked`; mode сохраняется для
продолжения.
