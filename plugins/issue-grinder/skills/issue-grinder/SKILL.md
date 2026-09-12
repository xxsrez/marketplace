---
name: issue-grinder
description: "Выполнить выбранный Task Manager scope через $issue-grinder или явный delivery intent; также справка о режимах. Не использовать для Backlog, status/audit, planning или кода без Task Manager selector."
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

Для явно изолированного local model-forward запуска прочитай
[local evaluation](references/local-evaluation.md); обычный delivery его не загружает.

## Проверка основной сессии до эффектов

Применяет `IG-MODE-20`. До resolver/исполнения (кроме чистой справки) один раз
в текущем turn обязательно выполни `python3 <this-skill>/scripts/main_profile.py`
для actual profile. Для выбранного `balance`/`economical` обязательно вызови
его с `--mode balance` либо `--mode economical` и получи `allowed=true` до
Goal, writes и children. Если mode уже явно известен, достаточно одного вызова
с mode. Self-report модели, желаемый профиль и нормализация не заменяют receipt.
При `allowed=false`, unknown или ошибке helper выбранный Luna-only режим
откажется работать; не пытайся продолжить PLAN. После model change/resume
receipt получают заново, в одном turn без изменений его не повторяют.
Если выбран или восстановлен `Баланс`/`Экономичный`, actual root обязан быть `gpt-5.6-luna` с
`reasoning_effort=max`. Unknown/mismatch → отказ с просьбой переключить основную
сессию на Luna Max, до Goal, любых writes и execution-subagents. Не подменяй
режим, не вызывай Luna-supervisor вместо root, не меняй настройки. Сохранённый
режим не отменяет проверку после смены модели. Чистая справка выше не блокируется.

## 1. Установи run, scope, режим и Goal

Для delivery полностью прочитай [Run, scope и Goal](references/run-and-goal.md).
Он определяет explicit/implicit continuity, current Release, live selector,
best-effort title и Goal lifecycle. Перед первой mutation полностью прочитай
[Task Manager flow](references/task-manager-flow.md).

Для нового run явные `Соло`, `Классический`, `Баланс`, `Экономичный` важнее
default; `single`, `сингл`, «одним агентом», «без субагентов» означают `Соло`.
Без явного выбора Luna Max → `Баланс`, Sol Extra High → `Классический`
независимо от числа задач. Для прочих профилей: non-Luna и больше одной live
задачи → `Классический`, иначе `Соло`. `Экономичный` только явно.
Удалённый mode в запросе либо сохранённом record требует явного выбора
поддерживаемого режима без потери выполненной работы и без скрытой подмены.
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
   полезных независимых пакета либо предусмотренные режимом critic/verifier, полностью прочитай
   [multi-agent execution](references/multi-agent-execution.md). В `Соло`
   рабочая делегация Issue Grinder запрещена: текущая модель выполняет один
   issue или пакет за раз. Внешние semantic providers не входят в эту topology.
   Во всех остальных режимах writer получает implementation только после
   подтверждённых routing и двухфазного worktree admission либо разрешённого
   shadow tree admission. Перенеси exact model/effort/fork из routing receipt;
   наследование root profile не считается выбором worker profile.
   Не предполагай, что execution-child умеет создавать собственных agents:
   nested delegation допустима только при наблюдаемой такой capability. Если её
   нет, используй mode-specific direct stages без остановки terminal-режима. Новый direct owner проходит guard и один
   spawn; ожидание событийное до stage deadline без произвольной отсечки,
   polling, повторных `list` и nudges. Final response — handoff; coordinator не
   дублирует Luna work. Ни один child не получает Goal, Task
   Manager, publication или final-acceptance authority.
3. Перед реализацией переведи `To Do → In Progress` и подтверди read-back;
   комментарий для этого тривиального перехода не нужен.
4. Реализуй issue outcome по mode promise, используя Strategic Outcome для
   локальных trade-offs без расширения scope; интегрируй только task-owned
   изменения и проверь exact integrated version по acceptance issue.
   Independent reviewer запускается там, где его требует выбранный mode,
   пользователь, project policy или exact scope. Child не перечитывает
   policy/tool catalog, пишет по absolute path и возвращает handoff. После
   rework прежний reviewer проверяет changed candidate. Partial review
   блокирует acceptance; только `Экономичный` сохраняет deferred gate checkpoint.
5. Для любого другого status transition подготовь причинный comment, пройди
   reflection, опубликуй и перечитай comment, затем измени status с optimistic
   concurrency и перечитай issue.
6. Пересчитай scope, стратегический outcome и следующий frontier.

`blocked by` ограничивает доступность требуемой реализации в exact integration base, а не связывает статусы; поздний reopen требует targeted recheck. Неизвестный contract или препятствующий дефект обходить нельзя.
Готовое issue с внешней приёмкой оставь в `In Review` после всех доступных проверок; сохрани остаток и продолжай доступный scope, включая dependents и режим Соло. Это не deferred internal review.
Объединённый blocker-handoff допустим после исчерпания самостоятельной работы scope; [Task Manager flow](references/task-manager-flow.md) задаёт учёт и resume. Экономичный checkpoint остаётся отдельным выходом.

При выборе проверок или адаптации проекта прочитай [соразмерные проверки](references/verification.md).
Найди правила, оцени влияние и прежние результаты, выполни достаточный набор;
после ошибки пересмотри диагностику, сохрани итог и оставшиеся обязательства. Release/UAT сами по себе не требуют полного Quality Gate: он нужен по запросу,
широкому/высокорисковому изменению или для production-ready доказательств.
Применимый прежний результат используй повторно. Production запрещён.

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

При затруднении перед остановкой или просьбой о вмешательстве примени [Консультант routing](references/consultant.md)
по `IG-FLOW-07`; перед остановкой Goal действуют отдельные обязательные условия ниже.

## 3. Используй publication и awareness loop

Перед обычным Task Manager comment, blocker-report без Goal и успешным финальным Goal comment
полностью прочитай [Strategic Explainer routing](references/strategic-explainer.md)
и сначала проверь active model. При Astra (`gpt-6-astra`) основной coordinator
выбирает native writing и не вызывает Strategic Explainer; этот guard имеет
приоритет над availability и сохранённым mode. В остальных случаях coordinator
сам вызывает доступный semantic facade и не передаёт ему неподтверждённые факты
либо delivery-работу выбранного scope.

Comment или final draft — точка решения. Непонятый факт либо существенная обязательная
работа current issue contract отменяют publication/status transition и возвращают run в
delivery loop; strategic gap или необязательное улучшение остаются follow-up.

Перед остановкой Goal прочитай [Консультант routing](references/consultant.md):
ведущая Astra разбирает ситуацию сама, остальные, включая unknown, вызывают
Консультанта. Его проверенный вывод возвращает к работе либо становится отчётом;
Strategic Explainer здесь не вызывается. Недоступный советчик не подтверждает blocker. Не ограничивай разбор одной консультацией: проверяй альтернативы, уточняй существенные пробелы и новые факты; до запроса разрешения проверь весь предлагаемый путь.

Blocker: `candidate blocker → причинное объяснение → reflection по current primary sources → continue | terminal blocker`.

Любое безопасное существенное действие отменяет blocker. Terminal report
перечисляет все current причины, checkpoint, unverified remainder, влияние,
нужное действие пользователя и resume signal; для каждой причины дай отдельный
ответ: почему она блокирует обязательный результат активного issue, почему Issue
Grinder не может устранить её сам и зачем этот шаг нужен issue contract и общей
цели. Общий report и ответы проходят reflection до
первой публикации. Покажи принятый комплект пользователю сразу; platform blocker
audit ограничивает только `update_goal(status=blocked)`, а не сам report. На
автоматическом продолжении без нового релевантного сигнала переиспользуй blocker
fingerprint: не повторяй проверку, facade, handoff или user request; на достигнутом пороге выполни только Goal mutation.

Финальное completion требует fresh inventory и reflection; Goal comment возвращай только пользователю в чате.
При completion, blocker или checkpoint выполни [итоговый отчёт о расходе](references/execution-modes.md#итоговый-отчёт-о-расходе) через доступный `codex-token-usage-report`: полный формат со всеми моделями текущего run и субагентов; без skill продолжай без установки.

## 4. Соблюдай authority и среды

Максимальная автономность относится только к явно вызванному run. Для exact
environment boundary и полного разрешённого цикла UAT полностью прочитай
[Autonomy and environments](references/autonomy-and-environments.md).
Production запрещён полностью: не подключайся, не читай data/logs, не deploy и
не smoke. Default environment для необходимого effect — подтверждённый UAT;
неизвестный target не угадывай. Task Manager lifecycle — отдельный control
plane, а не product Production deployment. Generic connector/tool label
`production` не переименовывает подтверждённый project UAT и не разрешает
спрашивать у пользователя подтверждение обычного UAT deployment. В scope весь цикл UAT — изменения, тестовый доступ, сбои, восстановление и очистка — заранее разрешён; по техническим шагам не переспрашивай. На Production послабление не распространяется.

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
Это `resumable checkpoint`, не `complete` и не `blocked`; mode сохраняется для продолжения.
