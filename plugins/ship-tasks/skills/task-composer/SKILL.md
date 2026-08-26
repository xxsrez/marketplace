---
name: task-composer
description: "Формулировать и по явному planning intent создавать Task Manager Tasks: оставлять один independently deliverable outcome одной Task, а составную работу превращать в Epic с problem-first описанием через Strategic Explainer, конкретными подзадачами, live labels, hierarchy и реальными relations. Использовать явно через $ship-tasks:task-composer и неявно для постановки, декомпозиции или backlog capture в Task Manager. Не использовать для implementation, delivery, release, status/audit-only запросов или изменения Label taxonomy."
---

# Task Composer

Создай связную и исполнимую planning-модель, которая сохраняет требования
человека и честно отражена в Task Manager. Draft без просьбы о записи остаётся
draft; явный create/add/backlog intent разрешает только planning mutations.

## 1. Сохрани границу workflow

- Работай только через Task Manager adapter. Не используй fallback tracker.
- Не реализуй, не тестируй, не выпускай и не принимай созданную работу; не
  создавай Goal и не выводи Tasks из `Backlog`.
- Одновременный запрос создать ровно одну Task и сразу выполнить её передай
  ShipTask create-and-deliver contract.
- Не создавай, не переименовывай и не архивируй Labels. Не помещай secret
  values, credentials или signed URLs в Task text.

## 2. Разреши live scope до writes

Установи current workspace, exact Project, write access, workflow statuses,
live active Labels и materially relevant duplicate candidates. Каждая Task
требует однозначного Project; memory и repository context не заменяют live
canonical refs.

Все новые элементы создавай в canonical status `Backlog`. Если его нет, не
подменяй default status и не начинай частичный create.

Назначай Release только когда пользователь выбрал его явно либо live Project
context однозначно определяет current unreleased Release. Не угадывай current
Release по максимальному номеру или дате. Если он неизвестен или неоднозначен,
создавай без `releaseRef` и сообщи об этом. Released Release требует отдельного
явного подтверждения.

До write выполни bounded duplicate search. Exact duplicate не создавай;
material overlap изучи до решения. Не изменяй и не reuse существующую Task без
соответствующего intent.

## 3. Скомпонуй целую модель

До первой mutation подготовь Project/Release/status, titles, descriptions,
parent-child hierarchy, labels и relation graph всего candidate scope.

Не превращай шаги исходного плана в Tasks механически. Строй outcome graph из
independently verifiable results; порядок задавай только настоящими
dependencies.

Не объединяй независимые desired outcomes в искусственный umbrella Epic:
сохраняй их отдельными Tasks/Epics и связывай только при реальной relation.
Title кратко называет ожидаемый результат и объект изменения: Epic — strategic
outcome, subtask — конкретный deliverable, без расплывчатой процессной формулы.
Type/classification выражай native Label и hierarchy, не title: не добавляй
`BUG:`, `EPIC:`, `[Bug]`, `Epic —`, `Feature:` или их эквиваленты. Missing Label
не заменяй textual prefix; исключение — exact verbatim title пользователя.
Legacy-prefixed и clean outcome title считай одним duplicate candidate.

Оставь одну Task, когда есть один independently deliverable outcome. Она
содержит problem, observable outcome, exact scope, material constraints,
достаточную technical конкретику, objective acceptance criteria и evidence.
Не создавай Epic с одной формальной подзадачей.

Создай Epic/parent Task, когда outcome требует нескольких independently
deliverable частей, разных проверяемых результатов или настоящих dependencies.
Используй самую мелкую полезную hierarchy.

Epic сохраняет problem, beneficiary, desired outcome, strategic intent, exact
scope, human requirements, cross-cutting acceptance/non-goals и целостный вклад
подзадач. Перед его созданием примени sibling
`$ship-tasks:strategic-explainer` как отдельный publication unit: новый built-in
`default` subagent с `fork_turns="none"`, одна compact task, exact planning scope
и resolvable read-only anchors без inherited turns/tool transcript/process
diary, твоего analysis/strategic view, требований к форме ответа или прежнего
candidate. Не читай provider-internal contract, не составляй explanation draft
и не применяй методику Explainer самостоятельно. Прими готовый description и
source basis либо refusal. Material facts проверь по authoritative planning
sources; factual/structural correction всегда передавай новому clean subagent,
старый не продолжай и текст самостоятельно не улучшай. Explainer не выбирает
decomposition, Project, status, labels, relations или write authority. Если
готовый grounded description недоступен, не создавай Epic; single Task, которой
Epic не нужен, от этого не блокируется.

Каждая подзадача получает один конкретный результат, exact change boundary,
свой вклад в desired outcome Epic, material technical details, применимые parent
requirements/dependencies, acceptance criteria и expected evidence.
Cross-cutting requirement остаётся в Epic и отражается в каждой применимой
подзадаче. Native parent link ведёт к полному strategic context, но одной ссылки
недостаточно: child description содержит компактную самодостаточную проекцию
вклада, применимых constraints/non-goals и качеств, которыми нельзя пожертвовать
ради локального упрощения. Исполнитель должен понять общий смысл, exact Task
boundary и планку качества без догадки; Epic context не расширяет scope child.
Material противоречие исправь до write. Для secrets указывай только имя
credential/secret store и target, никогда значение.

## 4. Назначь metadata по смыслу

Выбирай Labels только из live active catalog отдельно для каждой Task. Не
предполагай inheritance от Epic и не ставь нерелевантный label ради заполнения
поля. Если подходящего Label нет, создай Task без него и явно перечисли label
gap; taxonomy не расширяй и не дублируй classification в title.

Создавай native parent-child hierarchy. Не дублируй её `related` relation.
Добавляй relation только при реальном смысле: `blocks` для обязательной
dependency и `related` для полезной недирективной связи. `duplicate_of` не
создавай для нового planning set: exact duplicate вообще не создаётся, а
lifecycle mutation существующей Task требует отдельного intent. Не строй
последовательную цепочку по умолчанию; до write проверь direction `blocks`
человеческой фразой.

## 5. Выполни безопасные planning writes

Создай standalone Task или Epic с canonical Project, `Backlog`, confirmed
Release при его наличии и resolved label refs. Подзадачи создавай native
subtask operation с current parent version. После каждой mutation перечитывай
authoritative Task state. Relations создавай только после read-back обоих
endpoints со stable idempotency key.

Unknown write outcome сначала reconciles через reads и duplicate search; не
повторяй create вслепую. Если multi-Task create остановился частично, не скрывай
результат и не выполняй destructive cleanup без authority: перечисли created,
confirmed и not-created элементы и точное условие безопасного продолжения.

## 6. Проверь результат

Read-back должен подтвердить canonical identities, `Backlog`, Release либо его
честное отсутствие, hierarchy, labels/label gaps, relation type/direction и
intended strategic/technical split в descriptions. Проверь, что title не
дублирует type/classification Label, кроме exact verbatim user title.

Финальный ответ перечисляет созданный scope, duplicate disposition, Project,
Release, status, hierarchy, labels/label gaps, relations и любой unreconciled
outcome. Не представляй Task Manager planning projection как implementation
или delivery evidence.
