# Рой

Применяет `IG-MODE-05`, mode-specific части `IG-MODE-08..09` и
`IG-MA-15..19` только когда сохранённый `canonical_mode=swarm`. Общие resolver,
authority, recovery и evidence rules бери из
[Execution modes](../execution-modes.md); остальные mode-файлы не читай.

Цель — использовать большой объём дешёвого поиска, реализации и критики с
одним проверяемым выходом.

До wave задай конечный envelope: явный user budget либо bounded число волн с
условиями остановки. Свободная capacity не является бесконечным budget.

Допустимые роли: scouts, design candidates, implementation candidates,
minimal-diff/reliability alternatives, critics, test authors, economical
judges и evidence reducers. Различай подходы; множество одинаковых prompts
создаёт коррелированные ошибки и не является достаточным разнообразием.

До запуска задай одну bounded campaign: material fork, общую exact base,
различимые candidate purposes, максимальное число попыток и условия остановки.
Не рассчитывай на nested delegation. Если она доказанно доступна, один Luna Max
swarm owner может вести внутренние роли. Иначе Sol/controller одним dispatch
stage запускает напрямую только оправданные Luna candidates, а после их
завершения — одного direct Luna reducer/reviewer. Это штатная форма `Роя`, а не
checkpoint. Controller не ведёт исследование каждого кандидата и получает по
одному compact handoff на candidate и один итоговый ledger.

Все содержательные child roles `Роя` по умолчанию являются Luna Max: scouts,
design/implementation candidates, minimal-diff/reliability alternatives,
critics, test authors, economical judges, evidence reducers и bounded rework.
Каждый dispatch явно задаёт `model="gpt-5.6-luna"`,
`reasoning_effort="max"`, bounded `fork_turns` и проходит routing guard.
Имя и тип child не заменяют проверку его effective profile. Sol/controller
сохраняет effect ownership, fan-in, material integration decision и final
review, но не подменяет собой search, candidate writing или cheap critique.

Intentional candidate получает purpose, base, identity, branch/worktree и
verification. Он может затрагивать те же files, что другой candidate, только в
полной изоляции. Existing checkpoint сначала обнаруживается; новый candidate
не называется его resume/replacement.

Если scope содержит material candidate-friendly развилку, до ordinary
implementation зарегистрируй хотя бы одну настоящую Best-of-M stage с `M >= 2`
самостоятельными одно-ownerными Luna candidate waves от общей exact base. Один
writer на каждую Task без конкурирующей material stage не выполняет обещание
`Роя`. Если candidate-friendly развилки действительно нет, сохрани проверяемую
причину и всё равно используй Luna для одного candidate и независимых
scouts/critics; свободная capacity не создаёт искусственную развилку.

После candidate stage:

1. deterministic checks удаляют явно несостоятельные варианты;
2. direct Luna reducer/reviewer сравнивает оставшиеся, при необходимости
   используя доказанно доступные internal либо ограниченные direct read-only
   critics по независимым рискам, и сохраняет provenance, dissent и negative
   evidence;
3. reducer оставляет один recommended candidate и максимум один runner-up при
   действительно material unresolved fork;
4. integration owner принимает только task-owned commit выбранного candidate;
5. reducer/reviewer, который не был автором выбранного candidate, применяет
   compressed review packet к exact integrated candidate; прежние self-review и
   большинство одобрений эту проверку не заменяют;
6. замечания создают bounded rework внутри того же owner-а; после material
   rework та же reviewer session получает exact changed candidate одним
   follow-up и проходит один новый event-driven wait без replacement
   guard/spawn.

Campaign до dispatch задаёт reducer/reviewer action budget и единые критерии
сравнения. Reducer не перечитывает сырые transcripts и не воспроизводит research
каждого writer-а: он использует compact candidate handoffs, один source/diff
pass выбранного exact candidate, один основной deterministic suite и только
оправданные targeted risk probes; для малого/среднего scope по умолчанию не
больше трёх. Open-ended fuzzing, tool discovery ради parent messaging, cleanup
из read-only роли и повторный общий поиск после rework запрещены.

Каждый новый direct campaign owner проходит заранее разрешённый routing guard и
один spawn. Независимых candidate owners запускай как один bounded stage и
используй общий `event-driven wait` до переданного stage deadline, а не
произвольную десятиминутную отсечку или отдельный polling loop на каждого.
Технический timeout ожидания до deadline только продолжает то же ожидание без
`list`, probe, commentary или nudge. Reducer/reviewer запускается после
quiescence candidates. Final response каждого owner-а является его handoff;
отдельный messaging tool он не ищет. Sol/controller не
ищет guard, не читает его `--help`, не делает status polling, повторные `list`
или пустые nudges при неизменном состоянии.

При ambiguity, contract/context conflict, unexpected tool/environment state,
scope expansion или proof gap Luna сохраняет checkpoint/evidence и прекращает
одинаковые corrective retry. Следующая wave может запустить намеренно иной
candidate, но не параллельную копию того же подхода. Недоступная automatic Luna
не разрешает заменить wave Sol/GPT-5.4 children. Используй доступную serial Luna
lane либо сохрани candidates/evidence и честно остановись до следующей capacity;
явный пользовательский role profile не подменяй. Недоступность nested
delegation не является причиной checkpoint, пока coordinator может запустить
предусмотренные direct candidate и reducer/reviewer stages.

Останови новые attempts, когда candidate прошёл acceptance/final gate,
дальнейшие waves перестали добавлять независимые гипотезы/evidence, исчерпан
envelope или появился настоящий authority/environment blocker. Без final review
режим не завершает scope и не превращается в `Экономичный` автоматически.
