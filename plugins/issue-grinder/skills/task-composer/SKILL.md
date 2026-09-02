---
name: task-composer
description: "Формулировать и по явному planning intent создавать Task Manager Tasks: оставлять один independently deliverable outcome одной Task, а составную работу превращать в Epic с problem-first описанием через Strategic Explainer, конкретными подзадачами, live labels, hierarchy и реальными relations. Использовать явно через $issue-grinder:task-composer и неявно для постановки, декомпозиции или backlog capture в Task Manager. Не использовать для implementation, delivery, release, status/audit-only запросов или изменения Label taxonomy."
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

Epic сохраняет problem, beneficiary, Strategic Outcome, отличимые Human
Requirements/exact scope, Agent Plan, cross-cutting acceptance/non-goals и
целостный вклад подзадач. Перед его созданием примени sibling
`$strategic-explainer:strategic-explainer` как semantic facade отдельной
publication unit. Передай только назначение description, исходный вопрос, exact
planning scope, язык, material constraints и resolvable read-only anchors. Не
выбирай и не передавай никакие другие invocation parameters или provider
instructions: внутренним исполнением полностью владеет facade.
Не читай provider-internal contract, не составляй explanation draft и не
применяй методику Explainer самостоятельно. Прими готовый description и
отдельно обозначенный source basis либо operational unavailability; в Epic
записывай только description. Material facts проверь по authoritative planning
sources; factual correction передавай новым semantic call, а текст
самостоятельно не улучшай. Explainer не выбирает
decomposition, Project, status, labels, relations или write authority. Если
готовый grounded description недоступен, не создавай Epic; single Task, которой
Epic не нужен, от этого не блокируется.

Каждая подзадача получает один конкретный результат, exact change boundary,
свой вклад в Strategic Outcome Epic, material technical details, применимые
Human Requirements/non-goals и Agent Plan с dependencies, acceptance criteria и
expected evidence.
Cross-cutting requirement остаётся в Epic и отражается в каждой применимой
подзадаче. Native parent link ведёт к полному strategic context, но одной ссылки
недостаточно: child description содержит компактную самодостаточную проекцию
вклада, применимых constraints/non-goals и качеств, которыми нельзя пожертвовать
ради локального упрощения. Исполнитель должен понять общий смысл, exact Task
boundary и планку качества без догадки; Epic context не расширяет scope child.
Material противоречие исправь до write. Для secrets указывай только имя
credential/secret store и target, никогда значение.

Strategic Outcome помогает выбирать реализацию и проверять связность, но не
расширяет exact scope и не создаёт новую задолженность. Material работа вне
согласованного scope остаётся planning gap либо вопросом человеку; не маскируй
её под Human Requirement.

Если пользователь создаёт Task по bug report и передал attachment, оцени его
уместность по содержанию, связи с самой конкретной создаваемой Task и пользе
исполнителю: материал должен помогать увидеть проявление, воспроизвести,
локализовать, понять релевантный контекст либо проверить исправление. Применяй
один критерий к screenshot, документу, логу, записи и любому другому файлу;
формат сам по себе не создаёт презумпцию уместности. Уместный материал сохрани
как native attachment самой конкретной создаваемой Task, для которой это
evidence. Не заменяй обязательный attachment пересказом, local path, base64
либо временной или protected URL и не пропускай его молча. Явно нерелевантный,
избыточный или нарушающий secret-safe boundary материал не добавляй и сообщи
его disposition.

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

Для каждого обязательного attachment до create подтверди доступный native source
route. После существования target Task свяжи с ней verified file identity со
stable independent bind key и перечитай attachment metadata. Не начинай create,
если native transport заведомо недоступен. Если bind остановился после создания
Task, сохрани это как partial result с exact missing attachment и безопасным
условием продолжения, а не как success; unknown upload/bind сначала reconcile
через reads и не повторяй с новой identity.

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
