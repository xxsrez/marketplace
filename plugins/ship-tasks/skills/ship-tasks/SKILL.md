---
name: ship-tasks
description: "Доставлять однозначно выбранный Task Manager scope до проверенного результата, автономно создавая тестовые входы, выдавая причинный blocker-report и применяя строгий critical-codebase fallback, когда active frontier исчерпан и все оставшиеся In Review требуют существенной человеческой приёмки. Для batch scope применять лёгкий per-Task gate и периодический review-batch с одним exact UAT release, а не деплоить каждую Task. Использовать явно через $ship-tasks или неявно только при delivery intent с exact Task вроде TM-123 либо уже выбранным Task Manager Project/Release/current scope; одного delivery-глагола недостаточно. Create-and-deliver разрешён для явной просьбы создать и сразу выполнить одну Task. Bare invocation определяет mode по live inventory; Goal нужен только для implementation/rework минимум двух Tasks. Не использовать для чтения, статуса, аудита, объяснения, planning/backlog capture или обычной работы с кодом без Task Manager anchor."
---

# Ship Tasks

Доведи exact Task Manager scope до результата, подтверждённого фактами и правдиво
отражённого в Task. Сам выбирай реализацию и проверки, не подменяя outcome процедурой.

## 1. Выбери mode и exact scope

- Exact Task или create-and-deliver одной Task → `single`, без Goal.
- Реальная implementation/rework минимум двух concrete Tasks →
  `batch-implementation`, с Goal.
- Commit/push/publish/deploy/smoke/rollback уже подготовленного candidate →
  `release`, без Goal, включая production release.
- Чтение, аудит и объяснение → не delivery; planning/backlog capture с Task
  Manager intent принадлежит Task Composer.
- Обычная просьба исправить код/продукт без Task Manager anchor → не ShipTask.
Project, Release, current scope, несколько Tasks и bare `$ship-tasks` — selectors,
не modes. Mode определяет live inventory; чтение/проверка/reconciliation нескольких Tasks не создают Goal.
Exact Task и явно перечисленные refs — closed selectors; Project, Release и resolved
current scope — live selectors. Стартовый inventory — audit snapshot, не frozen refs/count cap. Перед frontier, ожиданием, blocker, Goal status и finish перечитай full inventory.
Для bare invocation прочитай [project memory](references/project-memory.md). Prompt selector
имеет приоритет, но не обновляет memory. Через Task Manager adapter разреши полный exact scope и перечитай Task state, acceptance, relations и comments.
Memory не заменяет live state. В первом update назови unresolved acceptance incidents.
Если Task входит в parent-child hierarchy, до первой implementation/rework mutation прочитай current parent chain до ближайшего relevant Epic целиком; установи общий outcome, вклад Task, применимые parent requirements, constraints/non-goals, exact child boundary, acceptance и dependencies и передай bounded context любому implementation/review packet.
Epic задаёт смысл и планку качества, но не расширяет selector/scope, не разрешает sibling work и не доказывает completion; material конфликт → `task-contract-conflict` до затронутой mutation, недоступный обязательный parent context → `TASK CONTEXT ALARM`, а независимую runnable работу продолжай.
После смены session и до новой implementation surface найди task-owned Git state. Если exact
unfinished worktree/branch существует, прежний writer остановлен и ownership
exclusive, прими тот же artifact и продолжай в нём вместо нового checkout или
повтора работы. Проверь scope, acceptance, diff/commits и partial effects; active/unknown ownership не перехватывай и artifact не очищай.
Canonical `Backlog` не входит в delivery inventory и не переводится в работу без отдельного
явного решения пользователя. Любая current matching non-Backlog Task из live selector входит
автоматически, включая появившуюся после старта или переведённую `Backlog → To Do`: это не scope expansion и не требует approval.
Вышедшую из selector/ушедшую в Backlog больше не реализуй; partial effects reconciliate.
В safe lane сначала `In Progress`/`In Review`, затем `To Do`; независимые lanes параллельны.
`Duplicate` отдельно не исполняй, но прочитай outgoing canonical relation и incoming duplicates.
`blocked by` открывает dependent implementation по readiness gate, а не связывает lifecycle statuses. Открой dependent
Task после доказанного fan-in нужного blocking Task contract в exact integration candidate; `Done` не требуется, а isolated code/comment/status gate не открывают. Pending upstream verification/effect сохраняет её non-terminal, но dependent Task уже runnable и может получить собственный `Done`.
Relation и Task/candidate/contract attribution сохраняй. Закрой gate повторно только при attributed upstream change/defect, который нарушает используемый contract; не инвалидируй независимые Tasks/evidence.
Если это может быть первый turn новой Codex task и host показывает title capability,
прочитай [title contract](references/thread-title.md): доказанный catalog placeholder получает не более одной best-effort попытки `ShipTask · ...`;
отсутствие/deferred/failure capability не блокируют workflow; meaningful title сохраняй.
Создавай Goal после подтверждения `batch-implementation` минимум двух Tasks и до implementation
mutation; `single`/`release` без Goal. Release-only run Goal не создаёт, не переиспользует, не ретаргетит и не завершает; compatible Goal продолжается лишь при исходном release done criterion; сам release новый Goal не создаёт.
Goal live selector хранит identity/predicate, не стартовые refs/count; новая matching Task сохраняет его active без approval. Goal не определяет Task outcome/попытки/status.
Неразрешимый scope/shared-state/authority конфликт → `TASK CONTEXT ALARM` до writes.
## 2. Соблюдай обязательные требования

### Исполняй topology rule пользователя, иначе выбирай автоматически
Сначала выведи effective topology rule из natural-language указаний пользователя.
Сохраняй их смысл и scope: exact/relative число, роли, общий/узкий opt-out,
условие по длительности, сложности или другому названному признаку. Совместимые
rules комбинируй; позднее более конкретное заменяет прежнее того же scope. Root
не входит в явно названное число субагентов. Exact count — обязательное число,
не ceiling; «побольше» materially увеличивает полезную delegation, а condition
проверяется по названной пользователем метрике.
Только без применимого rule сам решай, где delegation полезна и сколько
субагентов использовать. Параллельны лишь независимые packets. Safety,
authority, useful ownership, worktree isolation и проверяемый fan-in сильнее
topology preference. Не создавай фиктивные packets ради quota. Невыполнимое
rule не подменяй молча: назови exact конфликт, topology и влияние на result.
Основной агент — единственный integration owner и владелец Goal, Task Manager
comments/status/version writes, общего candidate и Task attribution. Каждый
одновременно пишущий implementation subagent до первой writable mutation
получает собственную feature branch и собственный Git worktree для своей exact
Task. Один writable worktree принадлежит одному writer: не пиши в integration
target, чужую branch или чужой worktree. Read-only scouts, reviewers и comment
Explainer отдельного worktree не требуют. Только integration owner делает fan-in
и проверяет exact объединённый result.
Явный user model/effort для subagents имеет приоритет. Без него только genuinely
simple packet запускай на `gpt-5.6-luna`/`max`: bounded self-contained scope,
ясные acceptance/evidence, disjoint ownership, без material creative/product/
architecture/authority/risk judgment и environment uncertainty. Остальные
packets и Strategic Explainer наследуют current model/effort; маленький diff сам
по себе не simple. При ambiguity, context/contract conflict, unexpected
environment/tool state, scope expansion или proof gap Luna прекращает packet без
corrective mutations, guess и ослабления acceptance; bounded read-only read-back
partial effects обязателен. Верни exact effects/evidence/unknown integration
owner на current profile без повторного cheap Luna loop. Если current primary —
Luna и broader context не снимает uncertainty, сообщи
`luna-escalation=not-available` без скрытой подмены Sol. Для override выбери
compatible bounded context. Genuine unavailable auto Luna → current profile;
явный unavailable user profile не подменяй и считай role capacity недоступной.
Общий no-subagent rule означает ноль субагентов во всём run. Он не отключает
доступный Strategic Explainer Fast, потому что Fast работает текущим агентом;
ordinary Explainer при этом недоступен, поэтому без Fast communication mode
становится native. Role-scoped rule меняет
только названную роль. Goal/lifecycle/authority не
меняются. При user rule сообщай его смысл и соблюдение либо material deviation;
внутренняя target/width accounting не требуется. Profile/handoff называй только
когда это помогает понять результат.
### Правдивый статус и обязательный комментарий
Перед любым существенным status transition опубликуй понятный native Task
Manager comment и перечитай его. Затем измени status и перечитай Task. Обычный
старт `To Do → In Progress` комментария не создаёт.
Комментарий обязателен перед `In Progress → In Review`, `In Review → In
Progress`, `In Review → Done`, reopen из terminal status, новым
`Canceled`/`Duplicate` и необычной lifecycle-корректировкой. Material blocker
также требует comment, даже если status не меняется. Ответ в Codex не заменяет
Task comment; `description` и другие fields не являются fallback.
Native comment create/list/read — гарантированная часть current Task Manager
adapter. Всегда создай и перечитай обязательный comment. Неизвестный write
outcome сначала разреши через native read-back; не повторяй write вслепую.
Каждый комментарий сформулируй в выбранном communication mode из
[reference](references/strategic-explainer.md): ordinary provider имеет
приоритет, затем доступный Fast, иначе native ShipTask writing. Opt-out также
выбирает native. Native mode сообщает обязательные lifecycle facts по
собственному truth contract ShipTask: не читает, не применяет и не имитирует
внутреннюю методику Explainer и не заявляет эквивалентное качество. Отсутствие
provider не блокирует comment/read-back/transition и не требует capability
warning. Обычный `To Do → In Progress` не запускает publication unit:
комментария для старта нет.
### Доказательство важнее выбранного способа

Сам выбирай инструменты, способы диагностики и приёмки. Сбой одного выбранного
способа не делает его обязательным, не доказывает defect и не создаёт blocker
сам по себе. Можно исправить инструмент, заменить его, объединить несколько
источников evidence или перестроить проверку — в зависимости от ситуации.
Browser/controller/session switch — диагностика, не repair. Достаточное evidence
exact candidate/server path остаётся product incident при сбое browser login/MFA. Пока
доступна безопасная in-scope работа с продуктом, browser logistics не stop condition.
Не снижай current acceptance ради удобства и не называй непроведённую проверку
проведённой. Единственное ослабление — строгий `critical-codebase-accepted` ниже.
Если в current scope/authority не доказать success или failure, это
`verification-blocked`: назови границу знания, влияние и resume condition.
Фиксированного числа попыток или обязательного порядка действий нет.
Обязательный Task comment — результат, не tool flow. Transition завершён, только
когда comment существует в Task и перечитан; как обеспечить это, решает агент.
### Не скрывай приёмочный инцидент

`verified-failure`, `verification-blocked` и `task-contract-conflict`, найденные
при проверке exact candidate или release scope, — material acceptance incidents.
Только `verified-failure` называй bug.
До repair, status write или blocking handoff немедленно сообщи в Codex chat:
exact Task/criterion, expected result, observed failure либо границу знания,
impact, установленный outcome и следующий шаг. Затем опубликуй и перечитай
opening Task comment. Для proven defect comment предшествует началу rework.
Не стирай opening после repair. При material state change сообщай progress в
chat; после исправления свяжи resolution/completion comment с тем же criterion,
cause/confidence, fix/result identity и retest evidence. Пока incident unresolved
и active run продолжается без material change, напоминай о нём в chat примерно
каждые 10 минут. Таймерные updates не дублируй в Task comments.
Ожидаемая red/green iteration, отказ одного способа проверки и batch failure без
task-level attribution не являются incident конкретной Task.

### Сохраняй authority boundary

Не расширяй scope и сохраняй unrelated пользовательские изменения, Tasks и
внешнее состояние; не используй blind rollback или destructive cleanup для
скрытия partial/unknown effect. Production, destructive durable-data changes, secrets,
privacy/access-policy changes, external recipients и unbounded cost требуют
явной authority. Обычные необходимые local/dev/test/QA/UAT/staging/preview/
sandbox effects разрешены после надёжного определения non-production target.
Подробности autonomy и release: [reference](references/autonomy-and-release.md).

## 3. Выполни и проверь result

Реализуй целостный in-scope result и получи evidence, достаточный для current
Task contract и обязательных runtime/external effects. Способ реализации,
диагностики и проверки выбирай сам. Недоступное называй `not-available`, а не
verified. Сначала свежо сопоставь finding с selector: новая matching non-Backlog
Task live selector уже in-scope и не требует approval. Только действительно
unmatched finding не исправляй и не превращай в новую Task без authority.

До `verification-blocked` исчерпай автономную test frontier: сам создай
synthetic fixtures/seed data (PDF, ZIP, PNG, изображения, Markdown) и проведи
их через поддерживаемые ingress, пробуя другой безопасный path после сбоя.
User-provided файл не blocker, если input можно создать самому; второй
principal/session, внешний account, provider evidence или access-policy effect —
authority boundary после проверки synthetic/ephemeral substitute, без создания
identity или ACL без authority.
Когда candidate готов, сформулируй problem-first comment о результате и
проведённой проверке, опубликуй и перечитай его, затем переведи
`In Progress → In Review`. Сразу проведи приёмку: `In Review` не является
ожиданием человека.
### Периодический UAT batch release

Перед review каждой Task выполняй лёгкий targeted gate. Совместимые ready
candidates группируй в разумный exact integrated batch и периодически, по
cadence/`batch_target`/review WIP, wave/frontier, общему UAT effect,
acceptance/Done, checkpoint или final flush, проводи один thorough batch gate и
деплой один exact candidate в verified UAT. Не деплой UAT после каждой bug/Task
по умолчанию. UAT — разрешённый non-production effect без approval после проверки
target; high-risk/coupled Task, explicit singleton или final flush могут быть
batch из одной Task. В report фиксируй members, SHA, checks, UAT read-back/smoke;
failure локализуй, а unaffected evidence не обнуляй без attribution.
## 4. Разбери приёмку по текущим фактам

История изменений acceptance и количество прошлых попыток сами по себе ничего
не доказывают. Те же outcomes применяй при release verification; lifecycle
effects выполняй только для точно attributed Tasks.
- `task-contract-conflict`: current mandatory requirements противоречат друг
  другу или не определяют observable result. Сначала chat update и opening
  comment. Однозначное исправление по accepted source выполни и перечитай;
  иначе сохрани `In Review` и назови recommended decision.
- `verified-failure`: exact candidate в совместимой среде прямо нарушает current
  acceptance. Сначала chat update и opening comment с expected/observed,
  evidence, impact и причиной возврата, затем `In Review → In Progress`;
  перечитай Task и продолжай rework в этом же run.
- `verification-blocked`: после автономной frontier success/failure всё ещё не
  доказаны в current scope/authority. До ожидания выпусти в chat и native Task
  comment report: Task/criterion, attempts, proven/unknown, cause, recommended
  path, prerequisites/authority, alternatives, success signal, safe continuation
  и exact resume condition. При сохранённом Explainer report проходит его
  независимый read-only проход; Task остаётся `In Review`, bug не объявляется.
- `verified-success`: acceptance, checks, result identity и required effects
  доказаны. Сначала понятный completion comment и read-back, затем
  `In Review → Done` и Task read-back. Если incident был открыт в этом run,
  comment явно закрывает его либо называет remaining risk.

После обычной матрицы перечитай fresh full inventory. Если `Backlog`/terminal
исключены, `To Do == 0`, `In Progress == 0`, `In Review > 0`, каждая оставшаяся
Task `verification-blocked` только потому, что нужен существенный human verifier,
и exact integrated candidate стабилен, примени
[critical codebase review](references/critical-codebase-review.md). Bounded
approval/MFA/invite/access unlock, доступный normal test path или tool inconvenience
gate не открывают. Ровно один read-only `critic` с `fork_turns="none"` независимо
проверяет current Tasks, code и tests. Grounded problem → обычный per-Task rework;
grounded approval → обязательный отдельно объяснённый
`critical-codebase-accepted` comment/read-back и `Done`; inconclusive → `In Review`.
Любое изменение candidate/contract/inventory до writes требует fresh gate.
Падение общего batch gate без task-level attribution не доказывает defect каждой
Task. Сначала получи separating evidence; до attribution сообщай scope-level
finding только в chat/final report и не возвращай весь batch в rework. Defect в
terminal Task сначала получает opening comment, затем exact Task reopen.

## 5. Обеспечь человеческое объяснение

В начале run выбери communication mode по live skill catalog и effective user
rule: ordinary `$strategic-explainer:strategic-explainer`, иначе Fast
`$strategic-explainer-fast:strategic-explainer-fast`, иначе native. При
установленных обоих приоритет у ordinary. Сохраняй mode до explicit rule change
или failure выбранного provider-а. Каждый Task Manager comment, отдельный
Task/scope report, blocker explanation и final — отдельная publication unit;
routine chat/progress unit не создают.

Следуй [client protocol](references/strategic-explainer.md). Ordinary provider
вызывай новым `default` read-only subagent с `fork_turns="none"`, exact terminal
role lock и opaque clean envelope; не читай provider-internal contract. Fast
выполняй текущим агентом без subagent: установи одну compact user-facing task,
exact scope и resolvable read-only anchors, затем загрузи Fast skill и его
reference. Inherited history не является evidence.

Provider result — готовый пользовательский text и отдельно обозначенный source
basis либо operational refusal. Публикуй только text; basis к нему не дописывай.
Проверь material facts по authoritative sources и не делай второй editorial
rewrite. Changed facts/scope/anchors получают новый pass того же provider.
Invalid ordinary invocation исправь одним новым clean subagent; старый не
продолжай. Повторный refusal или любая другая ошибка выбранного provider-а
переводит mode в native без вызова второго provider. Не заявляй для Fast
независимость или stateless isolation.

В native mode самостоятельно сформулируй unit только по собственному
truth/lifecycle/reporting contract ShipTask. Не загружай и не имитируй внутренний
метод provider и не заявляй эквивалентное качество. Отсутствие, opt-out или
failure Explainer не является capability gap: comment публикуется и
перечитывается, разрешённый transition выполняется. Immediate incident update
всё равно обязателен. [Comment](references/delivery-report.md).

До окончательного blocker claim используй provider explanation/source basis как
reflection input, если provider mode дал result; в native mode выполни ту же
повторную проверку safe frontier непосредственно по primary facts. Заново проверь
исходную цель, primary/cascade cause, applicable context и всю safe in-scope
diagnostic/repair/verification/reconciliation frontier. Достаточный путь отменяет
stale blocker; иначе создай новый blocker unit. Wording не является
evidence/authority. Срочный proven incident публикуй отдельным nonterminal unit
без преждевременного release-wide blocker claim.

Финальный ответ — новая scope-level unit с исходным вопросом, exact scope и
anchors всего run, не склейка Task comments или прежний draft. В provider mode
используй выбранный provider; в native mode сразу сообщи grounded result без
служебного предупреждения об отсутствии provider.

## 6. Продолжай автономно и финализируй

Task-local blocker не останавливает независимую runnable работу. В `batch-implementation` перед
ожиданием пользователя перечитай полный current inventory live selector, включая новые Tasks;
пока остаётся безопасная in-scope работа, продолжай. Goal active до завершения current membership.
Перед финальным ответом перечитай affected Tasks, comments, statuses, применимый
Goal и external effects. Собери authoritative anchors, обязательные factual
fields и краткий перечень существенных инцидентов по [run report](references/run-report.md), затем проведи publication unit
через раздел 5. Не составляй второй explanation candidate после provider
result. При user topology rule сохрани факт его соблюдения/deviation; без rule —
только material delegation/profile handoff.
Материальные `blocked by` gates отчитай отдельно от acceptance blocking Tasks: exact candidate/contract, downstream effect и attributed invalidation.
Product failure и доступную repair frontier сообщи раньше browser/OAuth/MFA
logistics; смену средства не называй repair или обязательным user action, пока
есть безопасная работа с продуктом. Unresolved incident не совместим с clean
success affected Task. Не заменяй Task comments этим ответом. Подробности: [run report](references/run-report.md).
Goal `batch-implementation` заверши после fresh full inventory без незавершённой current in-scope работы; release-only run его не финализирует. В final отдельно назови Tasks, закрытые через `critical-codebase-accepted`, непроведённую functional check и residual risk.
