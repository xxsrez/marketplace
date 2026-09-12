---
name: task-composer
description: "Сформулировать, создать или декомпозировать Task Manager работу через $issue-grinder:task-composer или planning intent. Запись только по просьбе создать; без delivery и status/audit."
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
  separate create-and-deliver contract.
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
attachment mapping, parent-child hierarchy, labels и relation graph всего
candidate scope.

Сохрани три различимые роли: **Strategic Outcome** объясняет общую проблему и
направляет решения; **Human Requirements** содержат только явно данные или
согласованные человеком обязательства; **Agent Plan** хранит изменяемые
decomposition, boundaries, dependencies, acceptance, evidence и technical
detail. Фиксированные заголовки не обязательны, но предположение агента,
рекомендуемый hardening или широкий стратегический ориентир никогда не становятся
Human Requirement из-за одной формулировки.

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
содержит problem и Strategic Outcome, отличимые Human Requirements/exact scope и
Agent Plan с достаточной technical конкретикой, objective acceptance criteria и
evidence. Не создавай Epic с одной формальной подзадачей.

Создай Epic/parent Task, когда outcome требует нескольких independently
deliverable частей, разных проверяемых результатов или настоящих dependencies.
Используй самую мелкую полезную hierarchy.

Для Epic/подзадач прочитай [Epic planning](references/epic-planning.md).
При Astra (`gpt-6-astra`) Strategic Explainer не вызывай: description native;
только не-Astra Epic требует связанный там publication contract.

Strategic Outcome помогает выбирать реализацию и проверять связность, но не
расширяет exact scope и не создаёт новую задолженность. Material работа вне
согласованного scope остаётся planning gap либо вопросом человеку; не маскируй
её под Human Requirement.

Если переданы пользовательские файлы, до первой mutation прочитай
[attachments](references/attachments.md): оцени уместность и native route,
сохрани обязательные attachments на соответствующих Tasks и проверь read-back.

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
intended strategic/technical split в descriptions, а также каждый обязательный
attachment на сопоставленной Task. Проверь, что title не дублирует
type/classification Label, кроме exact verbatim user title.
Отдельно проверь, что Strategic Outcome, Human Requirements и Agent Plan
различимы и ни одно agent-owned предположение не записано как требование
человека.

Финальный ответ перечисляет созданный scope, duplicate disposition, Project,
Release, status, hierarchy, labels/label gaps, relations, attachment disposition
и любой unreconciled outcome. Не представляй Task Manager planning projection
как implementation или delivery evidence.
