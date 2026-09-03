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
canonical_mode: solo | classic | balance | swarm | economical
mode_origin: explicit | automatic
initial_main_model: <exact effective model>
initial_main_effort: <exact effective effort when available>
controller_profile: <effective profile>
worker_profile: <effective profile>
routing_guard_path: <once-resolved absolute bundled path>
routing_receipts: <pre-dispatch receipts and observed child profiles>
campaign_stages: <bounded direct or nested owner stages and deadlines>
review_session: <independent reviewer identity and current state>
expensive_work_ledger: <reason-coded controller/reviewer work when required by mode>
```

До первого implementation dispatch:

Порядок разрешения фиксирован: `explicit user mode → automatic model rule`.

1. Если это доказанное продолжение, восстанови record из run continuity и не
   применяй automatic resolver снова. Смена current model/effort не меняет уже
   выбранный mode.
2. Для нового run найди в текущем prompt явное намерение выбрать режим. Прими
   каноническое имя и свободные формы вроде «Соло», `single`, `сингл`, «одним
   агентом», «без субагентов», «используй профиль Balance», «запусти Рой» или
   «работай экономично», только когда речь явно о режиме Issue Grinder.
   Случайное слово «баланс» или «один» в описании продукта не является
   selector-ом.
3. При явном selector-е выбери его. Иначе exact top-level model семейства
   `gpt-5.6-luna` при любом effort выбирает `economical`; любая другая модель —
   `classic`.
4. `По умолчанию` запускает то же automatic правило и не создаёт шестой mode.
5. Выполни profile normalization ниже, сохрани record в доступной continuity
   текущего run и один раз коротко назови выбранный канонический режим в
   progress update.

Не читай config, прошлый run или выбранные child profiles вместо effective
top-level model текущего нового run. Не пересчитывай automatic mode из процента
квоты, доступности reviewer-а или модели, на которой позже продолжили работу.

## Profile normalization

Economical baseline:

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

Если root уже запущен на Luna ниже Max и mode не `Соло`, оставь его
единственной transport/authority оболочкой. Он вызывает direct built-in Luna
Max supervisor с bounded packet для scope analysis, dispatch decision или
review и сам сохраняет Goal, Task Manager writes, fan-in и publication unit.
Supervisor не пишет comment, не вызывает Strategic Explainer и не создаёт
отдельную Codex task/session. Его результат — decisions, facts, evidence и
resolvable anchors для root.

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

В `Балансе`, `Рое` и `Экономичном` Luna-lane всегда запускается с явно
переданными `model="gpt-5.6-luna"`, `reasoning_effort="max"` и
`fork_turns="none"` либо положительным bounded числом. Не оставляй model/effort
на наследование root. Имя или тип агента не выбирает профиль режима и не
является основанием для запрета. Используй effective profile текущего dispatch,
передавай требуемые режимом model/effort явно и сверяй наблюдаемый профиль.

Пример pre-dispatch gate для содержательного writer-а:

```bash
python3 <installed-skill>/scripts/model_routing_guard.py \
  --packet-id TM-123-implementation-1 \
  --mode balance --semantic-role implementation --agent-type worker \
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
- любой non-Solo exact candidate до terminal acceptance получает independent
  reviewer-а, который не был его автором; малый или простой scope review не
  отменяет;
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

Каждая многоагентная campaign разбивается на ограниченные стадии. Каждая direct
рабочая или проверочная волна внутри стадии имеет ровно одного owner-а и один
понятный outcome. Internal roles остаются за одним owner-ом только при
наблюдаемой nested delegation; неизвестная child tool surface не считается
такой capability. Когда вложенность недоступна, `Баланс` использует отдельные
direct Luna execution и review stages, а `Рой` — bounded batch самостоятельных
direct candidate waves и последующий direct Luna reducer/reviewer. Это
mode-compatible fallback, а не
основание остановить terminal run.

Normal lifecycle каждого нового direct owner-а:

```text
resolved routing guard ×1 → spawn owner ×1 → event wait ×1
```

Owner не отправляет routine progress родителю. Event wait длится до результата,
запроса внимания или owner deadline, который короче остатка общего run ceiling.
Независимых owners одной стадии можно ждать общим event mechanism. При
неизменном состоянии запрещены discovery/`--help`, status polling, повторные
`list` и пустые nudges. На deadline owner возвращает candidate/ledger либо
точный partial handoff, а не исчезает без evidence.

После material rework переиспользуй того же owner-а и ту же reviewer session:

```text
follow-up exact changed candidate ×1 → event wait ×1
```

Новый reviewer guard/spawn допустим только после доказанной недоступности
прежнего reviewer-а, потери независимости либо material изменения review
contract. На большом scope доказанно способный owner может создавать внутренние
read-only lenses; иначе direct lenses ограничиваются независимыми risk surfaces,
а один final reviewer получает их compact findings и exact candidate.

## Выбранный режим — обязательная загрузка

После сохранения mode record полностью прочитай ровно один соответствующий
файл. Он является единственным runtime-местом с mode-specific topology,
role/profile routing, fallback, review и stop promise:

| `canonical_mode` | Mode contract |
|---|---|
| `solo` | [Соло](modes/solo.md) |
| `classic` | [Классический](modes/classic.md) |
| `balance` | [Баланс](modes/balance.md) |
| `swarm` | [Рой](modes/swarm.md) |
| `economical` | [Экономичный](modes/economical.md) |

Не загружай остальные четыре mode-файла «для сравнения» во время delivery и не
собирай mode-specific policy из `mode-help.md`, Architecture, старого run или
соседнего файла. Общие resolver, normalization, invariants, switch barrier и
review packet остаются в этом reference; выбранный файл их не переопределяет.

## Explicit mode switch

Только явная команда пользователя меняет mode незавершённого run. Перед сменой:

1. прекрати новые dispatch;
2. доведи active writers до commit/checkpoint или безопасно останови их;
3. выполни `assert-unchanged` integration checkout и reconcile ownership;
4. сохрани exact candidates, checks и deferred effects;
5. запиши новый canonical mode с `mode_origin=explicit`;
6. примени его только к следующей wave/review decision; если новый mode —
   `Соло`, поставь сохранённые packets в последовательную очередь и не создавай
   новых Issue Grinder execution-subagents.

Переключение не меняет scope, target environment или authority и не позволяет
повторно реализовать уже сохранённую работу.

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
