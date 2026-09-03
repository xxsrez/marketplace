# Multi-agent execution

Применяет `IG-MA-01..19` при двух независимых полезных пакетах, при явном
topology rule пользователя либо когда сохранённый execution mode требует
critic/verifier/candidate wave. Сначала используй mode и profiles из
[Execution modes](execution-modes.md) и полностью прочитай выбранный там
mode-файл. Делегация не расширяет scope, authority или автономность неявного
run. Если сохранённый mode запрещает рабочую delegation Issue Grinder, не
применяй path ниже: startup recovery остаётся обязательным, но execution-
subagents и parallel writer admission запрещены. Agent-backed semantic
providers регулируются собственными interfaces и не входят в этот path.

## Topology и ownership

Явное natural-language правило пользователя о числе, ролях, условиях или
opt-out имеет приоритет. Root/coordinator не входит в названное число
субагентов. Допустимые topology, роли и envelope бери только из выбранного
mode-файла. Число writers определяется этим envelope, capacity, ожидаемой
ценностью и стоимостью fan-in/verification; свободные слоты сами по себе не
создают полезную роль.

Сначала построй dependency graph и карту write surfaces: modules, files,
schemas, generated artifacts и shared state. Пересекающиеся либо зависимые
обычные packets не выполняй параллельно. Только предусмотренные выбранным
mode-файлом intentional candidates могут затрагивать одинаковые surfaces, и то
только в разных branches/worktrees; они никогда не пишут в общий integration
target. Каждый packet/candidate получает exact scope, owned surfaces, source
facts, constraints, expected outcome, verification и, для конкурирующего
варианта, отдельные purpose/candidate identity.

Каждый packet также несёт Strategic Outcome всего scope, вклад конкретного issue
или пакета, применимые Human Requirements и отличимый Agent Plan. Outcome
направляет локальные решения, но не разрешает child расширять exact scope,
изобретать requirements, работу, verification либо stopping condition. После
resume или material scope change packet пересобирается из current sources.

Оркестрация строится как bounded campaign из стадий. Каждая direct стадия
состоит из одной или нескольких самостоятельных атомарных рабочих волн с одним
owner-ом и одним понятным outcome у каждой волны; последующая проверочная волна
также имеет одного независимого owner-а. Поэтому direct candidate authors не
считаются внутренними авторами общей скрытой волны. Не предполагай nested
delegation из имени, типа или обещания
execution-child: она считается доступной только когда фактическая child tool
surface позволяет создавать и ожидать mode-compatible agents. При доказанной
capability один direct owner может вести внутренние роли с отдельными routing
receipts. При её отсутствии operational coordinator использует только прямые
стадии, предусмотренные mode-файлом; отсутствие вложенности само по себе не
создаёт checkpoint и не снимает terminal promise.

`Баланс` без nested delegation использует последовательные direct stages:
обычно один Luna execution owner возвращает exact candidate, затем другой Luna
owner, не бывший автором, возвращает independent finding ledger. Material fork
может оправдать несколько атомарных candidate waves; тогда reviewer также
сводит варианты перед проверкой выбранного exact candidate. `Рой` использует
bounded direct candidate stage от общей exact base, затем отдельного Luna
reducer/reviewer. Прямое число owners ограничивается полезными ролями и
campaign envelope, а не свободными слотами. Packet/campaign owner не пишет Task
Manager, не вызывает Strategic Explainer, не делает fan-in в общий candidate и
не принимает final result.

## Startup recovery — before new work

До новой implementation, branch, worktree или dispatch прочитай весь текущий
Git inventory repository, а не только `cwd`: `git worktree list --porcelain -z`,
локальные branches/commits и их integration reachability, status каждого
доступного checkout со staged/unstaged/untracked files и известные worker/run
checkpoints. Сопоставляй candidate с live scope по нескольким evidence; похожее
имя branch само по себе ownership не доказывает.

- Уже интегрированный exact result переиспользуй и проверь без повторной
  реализации.
- Quiescent exact-scope worktree продолжай на месте. Для существующей branch без
  worktree вызови bundled guard `resume`, создав linked checkout той же branch;
  не создавай новую replacement branch. Dirty task-owned checkpoint сохраняй и
  допускай только явным `--allow-dirty` после проверки owner/quiescence.
- Active owner не перехватывай: продолжай через него либо дождись доказанной
  остановки. Ambiguous/unrelated changes не очищай, не присваивай и не вливай.
- Task-owned dirty root checkout основной агент может продолжить только как
  serial exclusive lane. Его нельзя выдавать subagent-у или использовать как
  integration checkout параллельной wave.

Startup recovery предшествует `prepare`. Отсутствие Goal, receipt прошлой
сессии или уже созданного worktree не позволяет проигнорировать однозначно
связанные со scope branch, commit либо локальные изменения.

Новый intentional candidate, если его разрешает выбранный mode-файл, создаётся
только после inventory. Зафиксируй, чем его purpose/approach отличается от
existing work, выдай отдельную candidate identity и общую exact base. Без этого
он является запрещённым parallel replacement, а не допустимым Best-of-N
вариантом.

## Model routing admission — hard gate

При первой загрузке этого reference один раз разреши абсолютный путь bundled
[`model_routing_guard.py`](../scripts/model_routing_guard.py), сохрани его в mode
record и используй стабильный interface без повторного discovery или `--help`.
До каждого `spawn_agent`, которому Issue Grinder передаёт delivery-работу,
назначь unique `packet_id`, сформулируй настоящий semantic role и получи зелёный receipt от bundled
[`model_routing_guard.py`](../scripts/model_routing_guard.py). Exact параметры
и dispatch fingerprint успешного receipt перенеси в один фактический spawn без
наследования или подмены: для Luna-lane явно укажи `model="gpt-5.6-luna"`,
`reasoning_effort="max"` и bounded `fork_turns`. Если полный history нужен в
packet-е, передай необходимые facts и anchors явно; `fork_turns="all"` либо
omitted fork нельзя совмещать с mode-controlled profile.

Platform `agent_type` не является semantic role и сам по себе не выбирает
model/effort. Для любого child используй effective profile текущего dispatch и
явно передавай профиль, который требует выбранный режим. Наблюдаемый child
profile сразу добавь в receipt; mismatch останавливает wave до FileChange,
fan-in, тестов и lifecycle effects. Если tool не раскрывает actual profile,
сохрани `telemetry_pending` и exact spawn args для внешней recursive проверки.
Receipt с другим `packet_id` либо fingerprint не подтверждает этот dispatch.
Сам Python guard не перехватывает прямой platform spawn: его пропуск является
protocol violation, которое должно остановить mode-valid run при обнаружении.

## Direct-stage lifecycle — bounded event path

Для каждого нового direct owner-а normal trace содержит один заранее
разрешённый routing guard, один spawn и событийное ожидание до результата,
запроса внимания либо переданного stage deadline. `wait` получает весь остаток
stage deadline, а не произвольные десять минут. Если platform технически вернул
timeout раньше deadline, сразу продолжи то же ожидание без `list_agents`, status
probe, commentary, сообщения owner-у или новой интерпретации состояния.
Несколько одновременно запущенных независимых owners можно ждать одним
доступным multi-target/event mechanism; событие одного owner-а разрешает
обработать только его handoff, затем снова ждать оставшихся до того же deadline.

Каждая стадия возвращает один consolidated handoff на owner через его final
response: exact candidate
либо ledger, task-owned identity, checks, known defects, negative evidence и
следующий разрешённый переход. На малом scope direct reviewer читает exact
candidate сам. На большом он может использовать nested read-only lenses только
при доказанной capability; иначе coordinator создаёт ограниченные direct lenses
по независимым поверхностям риска и передаёт их результаты одному final reviewer,
не заставляя дорогой профиль читать их сырые transcripts. Owner не исследует
tool catalog ради отдельного parent messaging: platform сам возвращает final
response родителю.

Control brief задаёт каждому reviewer/reducer-у `review_envelope`: exact
candidate identity, risk surfaces, action budget и stopping condition. Для
малого/среднего scope defaults: candidate owner — десять tool calls, review plan
— три, exact review — пять, recheck — три, controller final gate — три.
Budget увеличивается до dispatch только по независимым risk surfaces. Candidate
owner делает один source pass, работает с настоящим изолированным candidate и
запускает основной suite; он не повторяет заведомо красный baseline без
диагностической пользы, не реконструирует проект несколькими in-memory
симуляциями и не ищет parent messaging. Reviewer получает один source/diff pass,
один suite и максимум три targeted probes. После rework он проверяет reproducer,
changed surfaces и suite. Controller читает compact ledger и делает один exact
pass с integrated suite, не повторяя exploration без противоречащего evidence.
Параллельный review-plan turn заканчивается final response-ом в пределах трёх
calls; это handoff и quiescence turn, а не потеря session. Exact candidate
приходит в ту же session follow-up-ом с отдельным лимитом в пять calls. Reviewer
не ищет messaging tool и не пытается оставаться активным между этими turns.

В `Балансе` normal order — `control brief → parallel Luna execution candidate(s)
+ independent review planning → Luna exact reduction/review/rework → integrated
checks → controller final review`. Первый дорогой
review exact candidate начинается только после завершения Luna stages. В `Рое`
normal order — `control brief → Luna candidate campaign → Luna reduction and
independent review → integration → controller final review`. Coordinator может
механически вести routing, ownership и fan-in между стадиями, но не выполняет за
Luna их research, implementation, tests или critique.

После material rework продолжи существующего candidate owner-а и ту же
независимую reviewer session:

```text
follow-up exact changed candidate ×1 → event wait ×1
```

Не запускай новый guard/spawn, пока не доказаны недоступность прежнего
reviewer-а, потеря независимости либо material изменение review contract.
Незавершённый review запрещает terminal acceptance; `Экономичный` может
сохранить partial ledger как deferred gate resumable checkpoint.

## Writer admission — hard gate

Каждый одновременно пишущий subagent до первой записи получает собственные
feature branch и Git worktree от подтверждённой integration base. Один writable
worktree принадлежит одному writer. Одного file ownership или пути в prompt
недостаточно. Используй bundled
[`writer_worktree_guard.py`](../scripts/writer_worktree_guard.py) по его
абсолютному resolved path; если Python/Git/isolated capacity
недоступны, выполняй write packets последовательно и не эмулируй parallel writer.

До чтения implementation guard source или попытки `prepare` сделай дешёвый
preflight: writable ли фактический Git common dir и разрешено ли создать linked
worktree. Если Git metadata read-only, не запускай заведомо невозможный flow.
При одной последовательной writing lane coordinator создаёт внутри разрешённого
sandbox task-owned shadow tree точной base без `.git`, фиксирует его identity и
передаёт одному Luna owner-у исключительное право записи. Owner изменяет файлы и
запускает tests прямо в shadow tree; coordinator затем строит точный diff к base,
механически применяет его в своём checkout и удаляет только созданный им shadow
artifact после проверки. Read-only review planning может идти параллельно,
поскольку второго автора нет. Если даже shadow tree недоступен, используй один
read-only patch-return packet. Controller не реализует candidate заново.

Shadow path задаётся напрямую из packet identity; не исследуй reference tree,
`.gitignore` или cache-каталоги ради его выбора. После review coordinator одним
механическим действием переносит exact bytes только declared owned files
(`rsync`/copy либо один автоматически построенный patch по local policy), затем
отдельно сверяет identity. Не печатай полный diff в model context и не набирай
его вручную по chunks; повторный integration action допустим только после
конкретного conflict или tool failure.

Перед каждой writer wave:

1. Выбери clean exact integration worktree/base. Dirty пользовательский checkout
   не stash/reset; создай отдельный task-owned integration worktree либо перейди
   к последовательному выполнению.
2. Сохрани `snapshot` integration worktree в receipt вне всех Git worktree. Для
   каждого нового writer выполни `prepare` с уникальными owner, branch и
   worktree path; любой `--output` также держи вне всех worktree.
   Команда использует `git worktree add --lock` без force и возвращает
   `issue-grinder/writer-worktree/v1`.
3. Первый turn writer-а является admission-only: запрети FileChange и попроси
   первым Git-действием запустить `admit` из exact worktree `cwd` с обязательными
   ожидаемыми worktree, branch, HEAD и common dir. Только после проверки receipt вызови
   отдельный follow-up с implementation packet.
4. Пока writer активен, integration checkout read-only. Параллельная запись
   coordinator-а требует его собственной admitted writer lane. После каждого
   interaction/return и перед fan-in выполняй `assert-unchanged`.
5. Изменившийся integration checkout, admission mismatch или отсутствующий
   receipt немедленно закрывает write wave: не интегрируй, не меняй lifecycle,
   останови дальнейшие mutations, прерви активных writers и сохрани неизвестный
   diff для reconciliation без reset/cleanup.

Worker завершает успешный packet task-owned commit-ом в ожидаемой branch и
возвращает commit/evidence. Uncommitted checkpoint можно сохранить для resume,
но нельзя fan-in. Coordinator проверяет branch/worktree/commit и только сам
объединяет изменения. Read-only scout/reviewer/Strategic Explainer отдельный
worktree не получает, пока фактически не становится writer-ом.

Только coordinator владеет Goal, Task Manager comments/status/version writes,
fan-in и exact integrated acceptance. Worker report и green checks в его
worktree — evidence input, не общий result. После каждого return/blocker,
scope change или готовой dependency пересчитай frontier.

Это владение включает всю publication unit и semantic facade. Worker, reviewer
или scout не формулирует готовый comment и не вызывает Strategic Explainer:
он возвращает coordinator-у только facts, evidence и resolvable read-only
anchors. Уже созданный built-in child может не иметь top-level namespace
`collaboration`, поэтому попытка породить provider из worker-а не переносится в
другой transport. Coordinator после evidence handoff сам вызывает facade;
`create_thread`, отдельная projectless task, fork или новая session ради текста
запрещены.

## Checkpoints и профили

Task-owned unfinished branch/worktree сохраняется. Перед takeover докажи, что
прежний writer остановлен и ownership quiescent; затем выполни `resume` exact
checkpoint и передай его новому exclusive writer через тот же admission. Для
доказанно task-owned dirty checkpoint `resume --allow-dirty` и
`admit --allow-dirty` допустимы только с exact expected branch/HEAD и после
quiescence. Active/unknown ownership не перехватывай, parallel replacement не
создавай и чужие изменения не очищай.

Profiles всегда бери из сохранённого mode record по
[Execution modes](execution-modes.md), а допустимые роли и packets — только из
выбранного mode-файла. Не наследуй current model/effort механически и не
пересчитывай automatic mode после смены top-level модели. Profile normalization
может сделать controller и workers одинаковыми; это не меняет topology и
delivery promise выбранного режима.

Mode record хранит все pre-dispatch и observed routing receipts текущей wave.
Название child, его self-report или намерение coordinator-а не доказывают
маршрутизацию без exact model/effort evidence.

При ambiguity, contract/context conflict, unexpected tool/environment state,
scope expansion или proof gap Luna сохраняет checkpoint/evidence и прекращает
неограниченные corrective mutations. Следующее действие выбирай по fallback
выбранного mode-файла. Luna retry loop с тем же подходом/evidence запрещён.

Недоступная automatic Luna не блокирует, пока существует mode-compatible путь.
Fallback и допустимый расход дефицитного profile бери из выбранного mode-файла.
Явный пользовательский role profile не подменяй. Strategic Explainer не входит
в routing рабочих packets и не считается execution-agent режима, пока получает
только свой semantic request. Его communication topology задаёт собственный
interface, а не mode-файл Issue Grinder.
