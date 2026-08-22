# Комментарии ShipTask

Native Task Manager comment — часть существенного lifecycle transition и
material blocker. Current Task Manager adapter гарантирует native comment
create/list/read. ShipTask всегда создаёт и перечитывает обязательный comment.
Он предназначен человеку, а не внутреннему workflow.

## Когда комментарий обязателен

- готовый candidate перед `In Progress → In Review`;
- доказанный defect перед `In Review → In Progress`;
- opening любого `verification-blocked` или `task-contract-conflict`;
- resolution доказанного defect после repair и повторной проверки;
- доказанный success перед `In Review → Done`;
- reopen из terminal status;
- новый `Canceled`/`Duplicate` и необычная lifecycle-корректировка;
- material blocker, даже если status остаётся прежним.

Обычный старт `To Do → In Progress` комментария не требует. Внутренние red/green
iterations также не создают noise.

## Порядок effects

1. Сформулировать comment с помощью Strategic Explainer contract.
2. Проверить, что текст соответствует current facts и выбранному outcome.
3. Найти equivalent comment, если outcome прошлого write неизвестен.
4. Создать native comment со stable idempotency key logical write.
5. Перечитать comment.
6. Только затем выполнить status write и перечитать Task.

Не повторять create вслепую. Не использовать `description`, acceptance или
другой Task field как fallback. Ответ в Codex не является durable Task comment.

Если outcome write неизвестен, сначала reconciliate через native reads. Если
обязательный comment не появился и не был перечитан, существенный status
transition не завершён. Это failure обязательной operation. Конституция не
предписывает конкретный tool flow, но требует фактический durable result и
запрещает подменять его `description` либо ответом в Codex.

## Содержание

Формального шаблона нет. В большинстве случаев достаточно пяти смыслов:

- результат или точная граница знания;
- значение для пользователя;
- основание вывода;
- причина status change либо blocker;
- следующий шаг или условие возобновления.

Текст начинается с результата. Он не содержит process diary, tool inventory,
reason-code dump, raw logs или внутреннюю orchestration, если они не нужны для
понимания. Technical detail получает человеческую роль: что именно он
доказывает и почему это важно.

### Готово к review

Объяснить, что реализовано, как это наблюдается и какая проверка пройдена.

### Нужна доработка

До repair немедленно сообщить incident в Codex chat. Opening Task comment
объясняет expected result, воспроизводимое нарушение current acceptance,
evidence, фактическое влияние, почему Task возвращается и что будет исправлено.
Не называть product defect то, что является только отсутствием способа проверки.

После repair новый resolution/completion comment связывает тот же Task и
acceptance criterion с cause/confidence, fix или exact result identity,
повторной проверкой, final state и remaining risk. Opening comment не
редактируется и не исчезает после успешного завершения.

### Приёмка заблокирована

Немедленно сообщить в chat, что приёмка не завершена и наличие product bug не
установлено. Task comment объясняет, что нельзя установить и почему, что уже
доказано, а также рекомендуемый feasible способ получить evidence, его
prerequisites и observable success signal. Несколько вариантов сравниваются
только при реальном material выборе.

### Завершено

Объяснить полученный outcome, его значение, ключевое evidence, exact result или
environment identity, если она нужна для проверки, и реальные ограничения. Если
в этом run был acceptance incident, явно назвать его final state и retest.

### Canceled, Duplicate или reopen

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
