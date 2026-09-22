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
initial_scope_size: small | medium | large | uncertain
mode_selection_reason: <explicit choice or concrete scope assessment>
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

До выбора нового режима получи actual model/effort через обязательный
bundled `scripts/main_profile.py`; для explicit balance/economical вызови сразу
с `--mode <mode>`. Перед effects нужен `allowed=true` именно этого mode в
текущем turn для balance/economical. После automatic выбора balance добавь
mode admission через `scripts/main_profile.py --mode balance`;
результат root не подменяется нормализованным worker profile. Helper читает
только CODEX_THREAD_ID, без config и чужих журналов. Нет receipt → нет запуска.

Для нового run явный выбор поддерживаемого режима имеет приоритет.
`single`, `сингл`, «одним агентом», «без субагентов» означают `solo` при
явном mode intent; случайное слово в описании продукта не является selector-ом.
Без явного выбора exact `gpt-6-luna/max` → `balance` независимо от объёма.
Неизвестный model/effort не подтверждает Luna Max. Для остальных профилей
оцени содержание live scope до стратегической декомпозиции:
небольшой/средний объём → `solo`, обоснованно крупный → `classic`.
Крупный объём требует нескольких содержательных направлений, каждое с отдельным
контекстом реализации и проверки, и существенной пользы распределения работы
или независимой проверки относительно затрат на координацию. Опирайся на
доступные contracts, подсистемы и зависимости; не запускай отдельное исследование
ради классификации. Несколько связанных исправлений, много однотипных правок
или одна сложная тесно связанная задача сами по себе не означают крупный объём.
Несколько самостоятельных функций или существенные изменения нескольких
подсистем с общей интеграцией могут обосновать `classic`.
В этой ветке число карточек, их декомпозиция, model/effort, квота и capacity не выбирают режим.
При недостаточных основаниях сохрани `uncertain` и выбери `solo`; это не отменяет
разрешение неизвестного contract, если он мешает самой delivery.
Сохрани размер и краткое обоснование конкретным составом работ в mode record;
для `large` назови также пользу многоагентного исполнения.
Для автоматического `balance` запиши причиной Luna Max; оценка объёма не нужна,
допустимо `uncertain`. `economical` только явно.
`По умолчанию` — этот resolver, не отдельный режим.

До Goal, mutations или child dispatch выполни `IG-MODE-20`: `balance` и
`economical` требуют exact effective current root `gpt-6-luna` / `max`.
Mismatch или неизвестный profile → откажись: «Для этого режима переключите
основную сессию на Luna Max». Не запускай работу, supervisor, дорогую оболочку,
автоматический другой режим и не меняй настройки пользователя. Это также
проверка на каждом resume и explicit switch, даже при сохранённом mode record.
Роль child и role override не обходят root admission. Чистая справка разрешена.
Новый balance record имеет `mode_contract_version=luna-coordinator-v1`.
Старый balance record без этой версии сохраняется, но требует explicit switch.
Удалённые `swarm`, `manager`, `roy`, `roi`, «Менеджер», «Рой» не подменяй default.

После выбора сохрани mode record, один раз назови режим и кратко объясни автовыбор.
В Соло единственный
execution profile — exact effective current top-level model и effort текущего
turn; смена профиля не меняет canonical mode и не разрешает execution-subagents.
Для остальных режимов до normalization и dispatch прочитай
[multi-agent routing](multi-agent-routing.md).
В следующих turns загружай только выбранный mode-файл.

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
соседнего файла. Общие resolver, invariants и switch barrier остаются в этом reference;
normalization, coordination и review packet — в условном
[multi-agent routing](multi-agent-routing.md); выбранный файл их не переопределяет.

## Explicit mode switch

Только явная команда пользователя меняет mode незавершённого run. Перед сменой:

1. прекрати новые dispatch;
2. доведи active writers до commit/checkpoint или безопасно останови их;
3. выполни `assert-unchanged` integration checkout и reconcile ownership;
4. сохрани exact candidates, checks и deferred effects;
5. проверь current root по `IG-MODE-20`; при mismatch откажи без смены режима; затем запиши новый canonical mode с `mode_origin=explicit`;
6. вне Соло прочитай [multi-agent routing](multi-agent-routing.md) до профильного
   решения; примени режим только к следующей wave/review decision; если новый mode —
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

## Итоговый отчёт о расходе

По `IG-FLOW-08` при completion, подтверждённом blocker или resumable checkpoint
добавь в ответ пользователю полный отчёт через `codex-token-usage-report`, если
он доступен в каталоге skills текущей среды. Прочитай его актуальный SKILL.md
и следуй его workflow и формату. Не сокращай таблицу: сохрани все ненулевые
строки моделей под каждой категорией и после Total, даже в кратком ответе и
после обработки текста Strategic Explainer. Сначала сообщи delivery outcome,
затем период и таблицу с оговорками provider-а. В Task Manager её не публикуй.

Измеряй только текущий run: сохраняй его начало и подтверждённые session IDs
root/descendants в continuity; перед отчётом зафиксируй telemetry cutoff.
Применяй штатный фильтр helper-а либо его поддерживаемый входной каталог с
временными копиями только этих logs (`--codex-home` с `sessions/` и
`archived_sessions/`, если поддерживается текущей версией). Не изменяй записи:
дедупликацией занимается provider. Копии держи вне Git и удали после отчёта.
Не считай все сессии машины за тот же период расходом run. Укажи, что ещё не
записанный финальный ответ не измерен; неполное покрытие обозначь явно.

Skill является необязательной интеграцией: без него продолжай без отчёта,
не устанавливай его и не копируй формат, цены или коэффициенты в Grinder.
Если точную телеметрию выделить нельзя или helper недоступен, кратко раскрой
ограничение вместо придуманного расчёта; это не blocker delivery. Для отдельно
запрошенного сравнения расхода используй тот же skill и заданный scope измерения.
