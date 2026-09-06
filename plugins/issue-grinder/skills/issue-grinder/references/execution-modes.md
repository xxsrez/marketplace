# Execution modes

Применяет общую часть `IG-MODE-01..11` и выбирает ровно один mode-specific
runtime contract. Прочитай этот reference после доказательства
нового/продолжающегося run и до стратегической декомпозиции. В продолжающемся
run с уже сохранённым mode record не выбирай режим повторно; восстанови record
и загрузи соответствующий ему mode-файл. Вернись сюда только для явного
переключения или общего profile/review решения.

## Mode record — hard gate

Один run имеет один record:

```text
canonical_mode: solo | classic | balance | economical
mode_origin: explicit | automatic
initial_main_model: <exact effective model>
initial_scope_task_count: <unique live tasks before decomposition>
initial_main_effort: <exact effective effort when available>
mode_contract_version: <balance: luna-coordinator-v1>
controller_profile: <effective profile>
worker_profile: <effective profile>
routing_guard_path: <once-resolved absolute bundled path>
routing_receipts: <pre-dispatch receipts and observed child profiles>
campaign_stages: <bounded direct or nested owner stages and deadlines>
review_session: <independent reviewer identity and current state when mode requires it>
expensive_work_ledger: <reason-coded controller/reviewer work when required by mode>
```

До первого implementation dispatch:

Сначала восстанови record доказанного продолжения: смена модели, числа
оставшихся задач или квоты не запускает автовыбор повторно.

Если exact actual profile неизвестен, используй один read-only вызов bundled
`scripts/main_profile.py` по CODEX_THREAD_ID; чужие sessions/config не являются
доказательством. Не повторяй helper в том же turn при сохранённом receipt.

Для нового run явный выбор поддерживаемого режима имеет приоритет.
`single`, `сингл`, «одним агентом», «без субагентов» означают `solo` при
явном mode intent; случайное слово в описании продукта не является selector-ом.
Без явного выбора exact `gpt-5.6-luna/max` → `balance`, exact
`gpt-5.6-sol/xhigh` → `classic` независимо от числа задач. Для прочих профилей
non-Luna и более одной live задачи → `classic`, иначе — `solo`.
`economical` только явно. `По умолчанию` — этот resolver, не отдельный режим.
Неизвестное число задач для последней ветки сначала разреши, не угадывай.

До Goal, mutations или child dispatch выполни `IG-MODE-20`: `balance` и
`economical` требуют exact effective current root `gpt-5.6-luna` / `max`.
Mismatch или неизвестный profile → откажись: «Для этого режима переключите
основную сессию на Luna Max». Не запускай работу, supervisor, дорогую оболочку,
автоматический другой режим и не меняй настройки пользователя. Это также
проверка на каждом resume и explicit switch, даже при сохранённом mode record.
Роль child и role override не обходят root admission. Чистая справка разрешена.
Новый balance record имеет `mode_contract_version=luna-coordinator-v1`.
Старый balance record без этой версии сохраняется, но требует explicit switch.
Удалённые `swarm`, `manager`, `roy`, `roi`, «Менеджер», «Рой» не подменяй default.

После выбора нормализуй profiles, сохрани mode record и один раз назови режим.
В следующих turns загружай только выбранный mode-файл.

## Profile normalization

Обычный economical baseline:

```text
model = "gpt-5.6-luna"
reasoning_effort = "max"
```

Явная команда пользователя в prompt о model/effort всех либо конкретных ролей
имеет приоритет. Обычный выбор top-level model в UI является входом resolver-а,
но не означает «все agents обязаны наследовать этот profile».

Без role override:

- main profile семейства Luna при любом effort до `max` включительно:
  `controller_profile = worker_profile = Luna Max`;
- иной main profile: сохраняй его для controller/reviewer, Luna Max используй
  как economical worker, если только repository evaluation не содержит
  отдельного надёжного правила, что main profile не сильнее Luna Max;
- неизвестное cross-family отношение не угадывай по цене, имени или одному
  прошлому результату.

В `Соло` не применяй эти profiles к исполнению: единственный execution profile
равен exact effective current top-level model и effort этого turn. Не создавай
Luna Max supervisor/worker для работы Issue Grinder и не подменяй current main
profile. При продолжении после смены top-level model/effort сохрани canonical
mode `Соло`, но следующую работу выполняй уже фактически текущим root profile.
Нормализованные поля mode record могут сохраняться только для безопасного
явного переключения в другой режим, но в `Соло` не дают права вызвать
соответствующего execution-agent. Внешний semantic provider не является такой
ролью и управляет своим profile через собственный interface.

Для `balance` controller/worker — фактический Luna Max root, а
specialist/reviewer — Sol Extra High. Для `economical` все роли Luna Max.
Правила схлопывания выше относятся к `classic`; для этих двух режимов
недопустимый root отклоняется, а не нормализуется через supervisor.

## Model routing — hard gate

Для каждого Issue Grinder execution-child dispatch заранее отдели смысловую роль от platform
`agent_type` и выполни bundled `scripts/model_routing_guard.py` по его
абсолютному resolved path. Разреши этот путь один раз при загрузке multi-agent
runtime и сохрани его в mode record; не ищи script и не читай `--help` перед
каждым dispatch без доказанного package/path mismatch. Receipt фиксирует unique `packet_id`, canonical
mode, semantic role, `agent_type`, exact requested model/effort, bounded
`fork_turns`, fingerprint этих dispatch args и, когда tool surface его
раскрывает, exact observed model/effort. До зелёного pre-dispatch receipt child
не создавай; точные поля receipt перенеси в один фактический spawn, а доступный
observed profile сверяй сразу. Receipt другого packet-а либо с другим
fingerprint не переиспользуй. При mismatch закрой wave до содержательной работы.
Если runtime surface не раскрывает actual profile, пометь его
`telemetry_pending`, не выдумывай подтверждение и сохрани exact spawn args для
внешней recursive telemetry проверки.

В многоагентных режимах mode-specific Luna-lane использует явные `model="gpt-5.6-luna"` и `reasoning_effort="max"`, bounded `fork_turns`. Не оставляй model/effort на
наследование root. Имя или тип агента не выбирает профиль режима. Используй
effective profile текущего dispatch и сверяй наблюдаемый профиль.

Пример pre-dispatch gate для содержательного writer-а:

```bash
python3 <installed-skill>/scripts/model_routing_guard.py \
  --packet-id TM-123-implementation-1 \
  --mode economical --semantic-role implementation --agent-type worker \
  --model gpt-5.6-luna --effort max --fork-turns none
```

Отсутствие либо отрицательный receipt, фактически другая model/effort или
substantive Sol/GPT-5.4 child там, где mode требует Luna, являются routing
failure. Не продолжай дорогую implementation wave под видом выбранного режима:
сохрани evidence и примени fallback выбранного mode-файла. Явный пользовательский
profile override сохраняется, но также проходит receipt с этим exact profile.

Bundled guard валидирует объявленный dispatch, но не перехватывает raw
platform `spawn_agent`. Пропуск guard-а, отсутствующий packet-bound receipt либо
расхождение fingerprint являются наблюдаемым protocol violation и делают
cost-aware run невалидным; не называй это физически непропускаемым platform
enforcement. Настоящий unskippable gate возможен только в dispatcher-е, который
одновременно валидирует и создаёт child.

## Общие инварианты

Во всех режимах:

- live scope, authority, Production prohibition и truthful lifecycle не
  меняются;
- startup recovery всегда предшествует fresh candidate/worktree;
- каждый concurrent writer проходит собственный worktree admission;
- один effect owner владеет Goal, Task Manager writes, fan-in и publication;
- один exact integrated candidate является предметом acceptance;
- worker self-report и isolated green tests не являются final evidence;
- independent reviewer обязателен только в mode contract, user instruction или
  project policy;
- intentional cheap failures не размножают deployment, shared mutable state,
  external recipients, paid external calls или terminal Task Manager effects;
- после material rework exact changed candidate проходит применимый review
  заново.
- routing receipt относится к реальному semantic packet, а не к названию child:
  текстовая роль `critic` не разрешает скрыть implementation, а
  `material_judgment` не является универсальным обходом Luna-lane.
- mode topology охватывает только агентов, которым Issue Grinder передал
  содержательную delivery-работу. Agent-backed Strategic Explainer, connector
  или другой semantic provider остаётся вне неё, пока исполняет только свой
  ограниченный interface; техническое положение в thread tree этого не меняет.

## Владение стадиями и событийная координация

Каждая многоагентная campaign разбивается на ограниченные mode-specific стадии с одним владельцем каждой волны.

Normal lifecycle каждого нового direct owner-а:

```text
resolved routing guard ×1 → spawn owner ×1 → event wait ×1
```

Owner не отправляет routine progress родителю. Event wait длится до результата,
запроса внимания или owner deadline, который короче остатка общего run ceiling.
Передай в wait весь остаток этого deadline, а не произвольные десять минут;
ранний технический timeout только продолжает то же ожидание без `list`, probe,
commentary или nudge.
Независимых owners одной стадии можно ждать общим event mechanism. При
неизменном состоянии запрещены discovery/`--help`, status polling, повторные
`list` и пустые nudges. На deadline owner возвращает candidate/ledger либо
точный partial handoff, а не исчезает без evidence.

Каждый owner получает конечный outcome и stopping condition. Один tool call может пакетировать связанные чтения или
проверки без потери evidence.

Coordinator материализует exact candidate root, прямой source manifest, owned
files, checks, stopping condition и compact output envelope. Candidate owner использует их
как готовую policy: не перечитывает Issue Grinder `SKILL.md`, references,
Architecture или routing guard, не ищет tool catalog/parent messaging и не
делает directory discovery уже переданных путей. Он один раз читает
contract/source, изменяет настоящий изолированный candidate абсолютными путями
под candidate root и запускает основной suite; известный заведомо красный
baseline suite не повторяет без диагностической ценности. Он не реконструирует
весь проект несколькими in-memory программами и не запускает open-ended
simulation/fuzzing. После последнего check сразу возвращает короткий handoff
`candidate | owned files | checks | findings<=3 | unknowns | next`. Те же
запреты на повторные чтения, messaging discovery и read-only cleanup действуют
для reviewer-а.

После material rework в режиме с independent reviewer переиспользуй того же
owner-а и ту же reviewer session:

```text
follow-up exact changed candidate ×1 → event wait ×1
```

Новый reviewer guard/spawn допустим только после доказанной недоступности
прежнего reviewer-а, потери независимости либо material изменения review
contract. На большом scope reviewer `Классического` может создавать внутренние
read-only lenses, когда capability доказана.

## Выбранный режим — обязательная загрузка

После сохранения mode record полностью прочитай ровно один соответствующий
файл. Он является единственным runtime-местом с mode-specific topology,
role/profile routing, fallback, review и stop promise:

| `canonical_mode` | Mode contract |
|---|---|
| `solo` | [Соло](modes/solo.md) |
| `classic` | [Классический](modes/classic.md) |
| `balance` | [Баланс](modes/balance.md) |
| `economical` | [Экономичный](modes/economical.md) |

Не загружай остальные mode-файлы «для сравнения» во время delivery и не
собирай mode-specific policy из `mode-help.md`, Architecture, старого run или
соседнего файла. Общие resolver, normalization, invariants, switch barrier и
review packet остаются в этом reference; выбранный файл их не переопределяет.

## Explicit mode switch

Только явная команда пользователя меняет mode незавершённого run. Перед сменой:

1. прекрати новые dispatch;
2. доведи active writers до commit/checkpoint или безопасно останови их;
3. выполни `assert-unchanged` integration checkout и reconcile ownership;
4. сохрани exact candidates, checks и deferred effects;
5. проверь current root по `IG-MODE-20`; при mismatch откажи без смены режима; затем запиши новый canonical mode с `mode_origin=explicit`;
6. примени его только к следующей wave/review decision; если новый mode —
   `Соло`, поставь сохранённые packets в последовательную очередь и не создавай
   новых Issue Grinder execution-subagents.

Переключение не меняет scope, target environment или authority и не позволяет
повторно реализовать уже сохранённую работу.

## Проверки и длительные инструменты — все режимы

Применяй `IG-MODE-09` и `IG-MODE-19`, включая self-review Solo. Сохраняй owner,
candidate/root, command, существенные входы code/tests/config/environment и
dependencies, handle, state, exit code и log/result location каждого долгого
процесса. Не отбрасывай session id при сокращении tool output. Receipt обновляй
рабочим вызовом, без отдельного bookkeeping turn. Доступный процесс продолжай
ждать, доступный результат прочитай вместо перезапуска. Unknown effect сначала
reconcile. После изменения переиспользуй лишь доказанно применимые результаты:
имени файла или candidate id недостаточно, учитывай зависимости и среду.
Повтори затронутые gates и необходимую интеграционную приёмку; неизвестное
влияние расширяет проверку. Смена фазы или reviewer-а не требует полного
повторения. Reviewer сам оценивает доказательства, исследует риски и при
необходимости перепроверяет их; summary не заменяет независимое суждение.

Пакетируй независимые чтения и проверки, сохраняя отдельные результаты и ошибки,
без конфликтующих writes, shared effects и ресурсных ограничений. В Solo это
остаётся внутри одного текущего пакета. Цена делегации включает постановку,
context transfer, ожидание, интеграцию и проверку; расчётный отчёт не нужен.
Свободные slots сами по себе не требуют agents.

## Review packet

Для controller/reviewer подготовь:

- exact selector, base и candidate identity;
- problem/acceptance и integrated diff с source anchors;
- deterministic и aggregate checks с raw evidence;
- material decisions и отклонённые alternatives;
- finding ledger, где каждый material finding имеет evidence и disposition
  `fixed | refuted_with_evidence | escalate`;
- unresolved objections, negative evidence, known defects и unknowns;
- reason-coded expensive-work ledger для режима, который ограничивает работу
  controller/reviewer profile;
- migrations, shared-state, permissions и внешние effects;
- один точный вопрос либо требуемое решение reviewer-а.

Summary помогает найти evidence, но не заменяет чтение существенного кода и
фактов reviewer-ом. Количество общих одобрений не перевешивает один
воспроизводимый material defect.
