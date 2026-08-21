# Комментарии ShipTask

Native Task Manager comment — часть существенного lifecycle transition и
material blocker. Он предназначен человеку, а не внутреннему workflow.

## Когда комментарий обязателен

- готовый candidate перед `In Progress → In Review`;
- доказанный defect перед `In Review → In Progress`;
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
4. Создать native comment.
5. Перечитать comment.
6. Только затем выполнить status write и перечитать Task.

Не повторять create вслепую. Не использовать `description`, acceptance или
другой Task field как fallback. Ответ в Codex не является durable Task comment.

Если comment create/list/read отсутствует или сломан, сначала применить правило
восстановления необходимого инструмента. Пока channel не восстановлен,
существенный status transition запрещён. Это инфраструктурная незавершённость
ShipTask, а не повод ослабить lifecycle contract.

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

Объяснить воспроизводимое нарушение current acceptance, фактическое влияние,
почему Task возвращается и что будет исправлено. Не называть product defect то,
что является только отсутствием способа проверки.

### Приёмка заблокирована

Объяснить, что нельзя установить и почему, что уже проверено/восстановлено, а
также 2–4 различных способа получить доказательство. Для каждого кратко назвать
предпосылки, доказательную силу и trade-off; затем дать рекомендацию и observable
success signal.

### Завершено

Объяснить полученный outcome, его значение, ключевое evidence, exact result или
environment identity, если она нужна для проверки, и реальные ограничения.

### Canceled, Duplicate или reopen

Объяснить material причину и связь с пользовательским результатом. Не создавать
новый terminal status или reopen без comment/read-back.

## Task boundary

Комментарий относится к exact Task. Shared batch evidence можно кратко назвать,
но нельзя копировать общий журнал всем Tasks. Secrets, signed URLs и private raw
logs не публикуются.
