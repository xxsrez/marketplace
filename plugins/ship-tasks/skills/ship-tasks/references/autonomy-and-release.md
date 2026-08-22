# Автономность и release authority

Этот reference уточняет границы, но не задаёт пошаговый универсальный flow.

## Продолжение работы

- Самостоятельно выбирай обратимые локальные решения, которые не меняют
  согласованный observable outcome.
- Сам выбирай достаточный способ продвинуть Task и проверить результат.
- Task-local blocker не останавливает независимые Tasks. В
  `batch-implementation` перед ожиданием пользователя перечитай весь scope; при
  наличии runnable work продолжай.
- Не повторяй неизменившуюся операцию, проверку или poll. Повтор нужен после
  material change в result, tool, environment, access, authority или Task
  contract.
- Нет фиксированного числа попыток. Остановка допустима, когда следующий шаг
  требует новой authority, внешнего state change или небезопасного действия.

Global `TASK CONTEXT ALARM` нужен только при конфликте exact scope, connector,
Goal, shared integration state или общей authority, из-за которого любые
оставшиеся writes небезопасны.

## Свобода способа и качество evidence

Выбирай, меняй и сочетай инструменты по ситуации. Отказ одного средства не
создаёт обязанности чинить именно его и не является готовым выводом о Task.
Решение оценивается по тому, доказывает ли итоговый evidence current acceptance.

Не объявляй разные способы равноценными без основания и не снижай требование
ради удобного результата. Если достаточный способ не найден в текущем scope и
полномочиях, объясни границу знания, влияние и условие возобновления. Это итог
приёмки, а не политика использования отдельного инструмента.

Native Task comment нельзя заменить полем description или ответом в Codex.
Существенный lifecycle transition завершён только при фактическом comment и
read-back, но выбор технического пути к этому результату остаётся за агентом.

## Task-local defer

Сохраняй правдивый status:

- `To Do`, если execution не начинался;
- `In Progress`, если есть partial implementation или rework;
- `In Review`, если candidate предъявлен, но приёмка объективно заблокирована
  или current contract требует material решения.

Перед defer опубликуй и перечитай понятный comment. Он сообщает, что уже
установлено, что мешает продолжить, влияние и exact resume condition. Не
создавай provider status `Blocked` и не создавай follow-up Task без authority.
Task-local blocker оставляет применимый Goal массовой имплементации активным.

## Non-production

Обычный необходимый release в local/dev/test/QA/UAT/staging/preview/sandbox
входит в delivery authority после надёжной проверки target. Разрешены build,
publish, deploy/redeploy, smoke и bounded repair/rollback затронутого
non-production surface.

Это не разрешает permanent deletion, destructive durable-data reset без
recovery, secrets/privacy/access-policy changes, external-recipient action,
unbounded cost или unrelated cleanup.

## Production

Production workflow требует явного approval для exact target и candidate.
Без него оставь понятный Task comment, сохрани правдивый non-terminal status и
продолжай независимую работу. Approval не отменяет checks, comment/read-back и
terminal evidence.

Production approval меняет authority, но не Goal policy. Release-only не создаёт
Goal из-за production target, Project/Release selector или количества прочитанных
Tasks. Уже активный Goal допустим только если он возник из массовой имплементации
минимум двух Tasks и production release был его исходным done criterion.

## Resume

После нового evidence, authority или внешнего state перечитай Task, comments,
result identity и affected environment. Не считать старый handoff текущим
доказательством. В первом содержательном chat update назови unresolved acceptance
incidents выбранного scope и возобнови с первого безопасного действия,
указанного в comment.
