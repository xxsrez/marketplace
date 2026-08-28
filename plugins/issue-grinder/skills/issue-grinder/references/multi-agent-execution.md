# Multi-agent execution

Применяет `IG-MA-01..18` при двух независимых полезных пакетах либо при явном
topology rule пользователя. Делегация не расширяет scope, authority или
автономность неявного run.

## Topology и ownership

Явное natural-language правило пользователя о числе, ролях, условиях или
opt-out имеет приоритет. Root/coordinator не входит в названное число
субагентов. Без override делегируй только когда dependency-ready frontier
содержит минимум два полезных независимых пакета; число writers определяется
реальной шириной, capacity и стоимостью fan-in/verification.

Сначала построй dependency graph и карту write surfaces: modules, files,
schemas, generated artifacts и shared state. Пересекающиеся либо зависимые
packets не выполняй параллельно. Каждый packet получает exact scope, owned
surfaces, source facts, constraints, expected outcome и verification.

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

## Writer admission — hard gate

Каждый одновременно пишущий subagent до первой записи получает собственные
feature branch и Git worktree от подтверждённой integration base. Один writable
worktree принадлежит одному writer. Одного file ownership или пути в prompt
недостаточно. Используй bundled
[`writer_worktree_guard.py`](../scripts/writer_worktree_guard.py) по его
абсолютному resolved path; если Python/Git/isolated capacity
недоступны, выполняй write packets последовательно и не эмулируй parallel writer.

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

## Checkpoints и профили

Task-owned unfinished branch/worktree сохраняется. Перед takeover докажи, что
прежний writer остановлен и ownership quiescent; затем выполни `resume` exact
checkpoint и передай его новому exclusive writer через тот же admission. Для
доказанно task-owned dirty checkpoint `resume --allow-dirty` и
`admit --allow-dirty` допустимы только с exact expected branch/HEAD и после
quiescence. Active/unknown ownership не перехватывай, parallel replacement не
создавай и чужие изменения не очищай.

Без user profile override `gpt-5.6-luna` с `max` получает только пакет, который
одновременно bounded, self-contained, disjoint, объективно проверяем, не требует
material product/creative/architecture/cross-issue judgment и не зависит от
непредсказуемой среды или риска. Маленький diff сам по себе не simple.

При ambiguity, contract/context conflict, unexpected tool/environment state,
scope expansion или proof gap Luna прекращает corrective mutations, фиксирует
checkpoint/evidence и возвращает тот же packet coordinator-у. Coordinator
продолжает на current profile без Luna retry loop. Недоступная automatic Luna
не блокирует: packet выполняется current profile. Явный пользовательский profile
не подменяй. Strategic Explainer не входит в routing рабочих packets.
