# Автономность, восстановление и release authority

Этот reference уточняет границы, но не задаёт пошаговый универсальный flow.

## Продолжение работы

- Самостоятельно выбирай обратимые локальные решения, которые не меняют
  согласованный observable outcome.
- Если безопасный in-scope repair способен продвинуть Task, выполни его,
  перечитай state и продолжай.
- Task-local blocker не останавливает независимые Tasks. В batch перед ожиданием
  пользователя перечитай весь scope; при наличии runnable work продолжай.
- Не повторяй неизменившуюся операцию, проверку или poll. Повтор нужен после
  material change в result, tool, environment, access, authority или Task
  contract.
- Нет фиксированного числа попыток. Остановка допустима, когда следующий шаг
  требует новой authority, внешнего state change или небезопасного действия.

Global `TASK CONTEXT ALARM` нужен только при конфликте exact scope, connector,
Goal, shared integration state или общей authority, из-за которого любые
оставшиеся writes небезопасны.

## Необходимый инструмент

Сбой нужного инструмента является частью текущей проблемы:

1. назови требуемую операцию и наблюдаемый отказ;
2. проверь current tool/catalog/configuration и ближайшую причину;
3. выполни безопасный bounded repair в существующей authority;
4. повтори исходную операцию после изменения;
5. только затем оцени равноценную альтернативу.

Альтернатива должна доказывать тот же requirement. Нельзя заменить runtime
acceptance компиляцией, authenticated scenario — публичным endpoint, а native
Task comment — полем description или ответом в Codex.

Если восстановление невозможно, обязательный lifecycle/blocker comment
объясняет причину, impact, выполненную диагностику и условие возобновления. Если
сломался сам comment channel, не выполнять существенный status transition и не
начинать новую delivery mutation, которую нельзя будет объяснить в Task.

## Task-local defer

Сохраняй правдивый status:

- `To Do`, если execution не начинался;
- `In Progress`, если есть partial implementation или rework;
- `In Review`, если candidate предъявлен, но приёмка объективно заблокирована
  или current contract требует material решения.

Перед defer опубликуй и перечитай понятный comment. Он сообщает, что уже
установлено, что мешает продолжить, влияние и exact resume condition. Не
создавай provider status `Blocked` и не создавай follow-up Task без authority.
Task-local blocker оставляет batch Goal активным.

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

## Resume

После нового evidence, authority или внешнего state перечитай Task, comments,
result identity и affected environment. Не считать старый handoff текущим
доказательством. Возобнови с первого безопасного действия, указанного в comment.
