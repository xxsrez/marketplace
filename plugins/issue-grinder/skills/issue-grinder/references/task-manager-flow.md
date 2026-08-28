# Task Manager flow

Применяет `IG-FLOW-01..06`, `IG-SCOPE-01..03` и Task Manager часть
`IG-GOAL-01..06`.

## Live scope

- Current prompt задаёт selector сильнее project context. Default current
  Release допустим только для явного `$issue-grinder` без selector-а.
- Разрешай Project, Release, relations и statuses в canonical refs. Full
  inventory включает все страницы до подтверждённого конца.
- Active membership: только `To Do`, `In Progress`, `In Review`. Не меняй
  `Backlog`, не отменяй issue и не управляй terminal/foreign statuses.
- Храни selector identity/predicate и последний snapshot раздельно. Refresh
  нужен после замеченного membership change, результата packet-а, перед новой
  dispatch wave и перед terminal decision; постоянный polling не нужен.
- Исключённое issue больше не получает lifecycle writes. Сохрани его
  recoverable checkpoint; интегрируй уже сделанное только если оно независимо
  необходимо оставшемуся scope.

## Lifecycle и `blocked by`

- `To Do → In Progress`: status write с current version и read-back, без comment.
- `In Progress → In Review`: comment объясняет реализованный outcome и
  готовность к проверке.
- `In Review → In Progress`: comment называет обнаруженный defect/proof gap и
  требуемый rework.
- `In Review → Done`: comment называет exact result, проведённые проверки,
  evidence identity и residual risk.

Status сам по себе не доказывает code/evidence. Для `In Progress` и `In Review`
сначала восстанови exact target branch, task-owned diff/commits, checks и
partial effects.

`blocked by` закрывает frontier только пока требуемого upstream contract нет в
exact integration base зависимого issue. Код в чужом worker branch, Task
comment или status этого не доказывают. После fan-in contract dependent issue
runnable; поздний upstream defect инвалидирует только доказанно затронутый
downstream evidence.

## Publication transaction

Для любого нетривиального status transition порядок целостности:

1. собрать current facts и evidence;
2. получить publication-ready text через выбранный communication mode;
3. выполнить reflection и отменить переход при существенном пробеле;
4. создать comment;
5. подтвердить comment read-back;
6. изменить status с optimistic concurrency;
7. перечитать issue и итоговый status.

Comment описывает подтверждённые факты и причину готовности, но не утверждает,
что будущий status write уже успешен. Поэтому `comment committed / status
failed` остаётся правдивым checkpoint. Reconcile live status и повтори только
недостающий effect. При unknown comment outcome сначала прочитай thread; при
unknown status outcome — issue. Blind retry и duplicate comment запрещены.

После каждой terminal mutation fresh read определяет следующий scope. Run не
завершён, пока active membership не пуст.
