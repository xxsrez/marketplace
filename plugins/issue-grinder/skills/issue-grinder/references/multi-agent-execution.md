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

Каждый одновременно пишущий subagent до первой записи получает собственные
feature branch и Git worktree от подтверждённой integration base. Один writable
worktree принадлежит одному writer. Read-only scout/reviewer отдельный worktree
не получает.

Только coordinator владеет Goal, Task Manager comments/status/version writes,
fan-in и exact integrated acceptance. Worker report и green checks в его
worktree — evidence input, не общий result. После каждого return/blocker,
scope change или готовой dependency пересчитай frontier.

## Checkpoints и профили

Task-owned unfinished branch/worktree сохраняется. Перед takeover докажи, что
прежний writer остановлен и ownership quiescent; затем передай тот же checkpoint
новому exclusive writer. Active/unknown ownership не перехватывай, parallel
replacement не создавай и чужие изменения не очищай.

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
