# Комментарии ShipTask

Native Task Manager comment — часть существенного lifecycle transition и
material blocker. Current Task Manager adapter гарантирует native comment
create/list/read. ShipTask всегда создаёт и перечитывает обязательный comment.
Он предназначен человеку, а не внутреннему workflow.

Пока effective rule не отключает comment Explainer, каждый комментарий до
публикации проходит один provider по
[availability protocol](strategic-explainer.md). Доступный Fast выполняется
текущим агентом без subagent; ordinary provider используется только когда Fast
отсутствует. Это относится не только к обязательным переходам.
Когда user rule отключает Explainer, основной агент сообщает обязательные
lifecycle facts по собственному truth contract и не читает, не применяет и не
имитирует provider method.

## Когда комментарий обязателен

- готовый candidate перед `In Progress → In Review`;
- доказанный defect перед `In Review → In Progress`;
- opening любого `verification-blocked` или `task-contract-conflict`;
- resolution доказанного defect после repair и повторной проверки;
- доказанный success либо `critical-codebase-accepted` перед `In Review → Done`;
- reopen из terminal status;
- новый `Canceled`/`Duplicate` и необычная lifecycle-корректировка;
- material blocker, даже если status остаётся прежним.

Обычный старт `To Do → In Progress` комментария не создаёт и Strategic
Explainer не запускает. Внутренние red/green iterations также не создают noise.
Если при старте возникло самостоятельное существенное событие, комментарий
относится к этому событию, а не к самому переходу.

## Инвариант effects

До связанного существенного status transition понятный comment должен
соответствовать current facts, фактически существовать в Task Manager и быть
перечитан. После transition Task также перечитывается. Как агент формулирует
text и организует create/reconciliation, конституция не предписывает.

Unknown write outcome reconciles через native reads; create не повторяется
вслепую. `description`, acceptance, другой Task field и ответ в Codex не
являются durable Task comment.

Если outcome write неизвестен, сначала reconciliate через native reads. Если
обязательный comment не появился и не был перечитан, существенный status
transition не завершён. Это failure обязательной operation. Конституция не
предписывает конкретный tool flow, но требует фактический durable result и
запрещает подменять его `description` либо ответом в Codex.

## Factual packet и publication

Пока effective rule сохраняет Explainer, основной агент устанавливает одну
короткую user-facing задачу, exact Task scope и resolvable anchors к current
facts, evidence, knowledge/authority boundary и next lifecycle state. Это не
готовый brief, explanation draft или требования к структуре текста. Выбранный
provider возвращает готовый text и отдельно обозначенный source basis. ShipTask
публикует только text, а basis использует для проверки material factual conflict;
исправление facts/anchors требует нового pass того же provider, а не второго
самостоятельного rewrite.

Если effective rule сохраняет Explainer, но выбранный provider недоступен или не
дал пригодный текст, комментарий не публикуется. Связанный существенный переход
остаётся незавершённым. Это capability failure, а не opt-out. Если user rule
отключает Explainer, основной агент публикует только необходимые grounded
lifecycle facts без provider methodology и без claim эквивалентного качества.
Ответ в Codex может немедленно показать incident, но не заменяет durable comment.

### Factual anchors: готово к review

Exact result, observable evidence и проведённая проверка.

### Factual anchors: нужна доработка

До repair немедленно сообщить incident в Codex chat. Opening Task comment
объясняет expected result, воспроизводимое нарушение current acceptance,
evidence, фактическое влияние, почему Task возвращается и что будет исправлено.
Не называть product defect то, что является только отсутствием способа проверки.

После repair новый resolution/completion comment связывает тот же Task и
acceptance criterion с cause/confidence, fix или exact result identity,
повторной проверкой, final state и remaining risk. Opening comment не
редактируется и не исчезает после успешного завершения.

### Factual anchors: Приёмка заблокирована

Немедленно сообщить в chat, что приёмка не завершена и наличие product bug не
установлено. Task comment объясняет, что нельзя установить и почему, что уже
доказано, а также рекомендуемый feasible способ получить evidence, его
prerequisites и observable success signal. Несколько вариантов сравниваются
только при реальном material выборе.

До этого report должен показать автономную test frontier: какие synthetic
fixtures/seed data агент создал сам, какие поддерживаемые ingress и проверки
прошёл, и почему этого всё ещё недостаточно. Стандартный PDF, ZIP, PNG,
изображение или Markdown, который агент может безопасно создать, не описывается
как user blocker. Если остаётся genuine authority boundary (например,
independent principal или вторая authenticated session), report отдельно
называет primary/cascade cause, сравнивает grounded test paths, рекомендует
лучший feasible путь, указывает prerequisites/authority, observable success
signal, продолжаемую safe работу и exact resume condition. Этот текст проходит
отдельного Strategic Explainer, пока effective rule его сохраняет; голое «нужен
файл/principal» без причинной рекомендации недостаточно.

### Factual anchors: завершено

Объяснить полученный outcome, его значение, ключевое evidence, exact result или
environment identity, если она нужна для проверки, и реальные ограничения. Если
в этом run был acceptance incident, явно назвать его final state и retest.

Для `critical-codebase-accepted` комментарий обязан прямо назвать
непроведённую функциональную проверку, почему она требует существенного human
verifier, а не bounded approval/unlock, исчерпанные autonomous paths, exact
candidate, проверенные criteria/code/tests, grounded verdict независимого critic
и residual risk. Человек должен без чтения сессии понять, что Task закрыта по
критической проверке кодовой базы, а не по полноценной функциональной приёмке.

### Factual anchors: Canceled, Duplicate или reopen

Объяснить material причину и связь с пользовательским результатом. Не создавать
новый terminal status или reopen без comment/read-back.

## Task boundary

Комментарий относится к exact Task. Shared batch evidence можно кратко назвать,
но нельзя копировать общий журнал всем Tasks. Secrets, signed URLs и private raw
logs не публикуются.

Batch failure без task-level attribution сообщается как scope-level finding в
chat и final report. Task comment появляется только после separating evidence.

## Progress channel

Приёмочный incident сначала виден в chat, затем durable в Task. Пока он
unresolved, chat получает короткий update при material state change и примерно
каждые 10 минут активной работы без такого изменения. Timer update не создаёт
новый Task comment; Task history получает opening, material blocker/change и
resolution.
