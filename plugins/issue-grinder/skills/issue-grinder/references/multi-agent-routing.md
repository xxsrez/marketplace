# Multi-agent routing

Читай только для `classic`, `balance` или `economical` после общего
[resolver](execution-modes.md), до normalization и первого dispatch.
В Соло reference не нужен. Это общая механика; topology, fallback и stop
promise задаёт выбранный mode-файл.

## Profile normalization

Обычный economical baseline:

```text
model = "gpt-6-luna"
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

В многоагентных режимах mode-specific Luna-lane использует явные `model="gpt-6-luna"` и `reasoning_effort="max"`, bounded `fork_turns`. Не оставляй model/effort на
наследование root. Имя или тип агента не выбирает профиль режима. Используй
effective profile текущего dispatch и сверяй наблюдаемый профиль.

Пример pre-dispatch gate для содержательного writer-а:

```bash
python3 <installed-skill>/scripts/model_routing_guard.py \
  --packet-id TM-123-implementation-1 \
  --mode economical --semantic-role implementation --agent-type worker \
  --model gpt-6-luna --effort max --fork-turns none
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
