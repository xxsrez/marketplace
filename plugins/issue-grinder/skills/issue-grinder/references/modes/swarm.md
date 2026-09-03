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

Для всей bounded wave создай одного direct Luna Max swarm owner-а. Он является
единственным child, видимым Sol/controller-у, и внутри своей волны управляет
scouts, candidates, critics, test authors, economical judges, independent
reviewer-ом, evidence reducers и bounded rework. Sol/controller получает один
consolidated handoff, а не ведёт каждого участника напрямую.

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
implementation зарегистрируй хотя бы одну настоящую Best-of-M wave с `M >= 2`
независимыми Luna candidates от общей exact base. Один writer на каждую Task без
конкурирующей material wave не выполняет обещание `Роя`. Если candidate-friendly
развилки действительно нет, сохрани проверяемую причину и всё равно используй
Luna для scouts/critics; свободная capacity не создаёт искусственную развилку.

После каждой wave:

1. deterministic checks удаляют явно несостоятельные варианты;
2. economical critics/judges сравнивают оставшиеся, сохраняя provenance,
   dissent и negative evidence;
3. reducer оставляет один recommended candidate и максимум один runner-up при
   действительно material unresolved fork;
4. integration owner принимает только task-owned commit выбранного candidate;
5. независимый reviewer, который не был автором выбранного candidate, получает
   compressed review packet, а не transcripts всех agents; прежние critics,
   judges и большинство одобрений эту проверку не заменяют;
6. замечания создают bounded rework внутри того же owner-а; после material
   rework та же reviewer session получает exact changed candidate одним
   follow-up и проходит один новый event-driven wait без replacement
   guard/spawn.

Новый direct swarm owner проходит один заранее разрешённый routing guard, один
spawn и один event-driven wait до результата, запроса внимания или deadline.
Owner применяет тот же gate к своим children, не отправляет routine progress
родителю и возвращает полный либо частичный finding ledger. Sol/controller не
ищет guard, не читает его `--help`, не делает status polling, повторные `list`
или пустые nudges при неизменном состоянии.

При ambiguity, contract/context conflict, unexpected tool/environment state,
scope expansion или proof gap Luna сохраняет checkpoint/evidence и прекращает
одинаковые corrective retry. Следующая wave может запустить намеренно иной
candidate, но не параллельную копию того же подхода. Недоступная automatic Luna
не разрешает заменить wave Sol/GPT-5.4 children. Используй доступную serial Luna
lane либо сохрани candidates/evidence и честно остановись до следующей capacity;
явный пользовательский role profile не подменяй. Если swarm owner не может
создать обязательную независимую проверку, верни checkpoint/evidence и не
объявляй terminal acceptance.

Останови новые attempts, когда candidate прошёл acceptance/final gate,
дальнейшие waves перестали добавлять независимые гипотезы/evidence, исчерпан
envelope или появился настоящий authority/environment blocker. Без final review
режим не завершает scope и не превращается в `Экономичный` автоматически.
