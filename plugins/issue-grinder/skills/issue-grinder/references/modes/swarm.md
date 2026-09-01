# Рой

Применяет `IG-MODE-05`, mode-specific части `IG-MODE-08..09` и
`IG-MA-15..17` только когда сохранённый `canonical_mode=swarm`. Общие resolver,
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

Intentional candidate получает purpose, base, identity, branch/worktree и
verification. Он может затрагивать те же files, что другой candidate, только в
полной изоляции. Existing checkpoint сначала обнаруживается; новый candidate
не называется его resume/replacement.

После каждой wave:

1. deterministic checks удаляют явно несостоятельные варианты;
2. economical critics/judges сравнивают оставшиеся, сохраняя provenance,
   dissent и negative evidence;
3. reducer оставляет один recommended candidate и максимум один runner-up при
   действительно material unresolved fork;
4. integration owner принимает только task-owned commit выбранного candidate;
5. reviewer получает compressed review packet, а не transcripts всех agents;
6. замечания создают новую bounded rework wave и повторный exact review.

При ambiguity, contract/context conflict, unexpected tool/environment state,
scope expansion или proof gap Luna сохраняет checkpoint/evidence и прекращает
одинаковые corrective retry. Следующая wave может запустить намеренно иной
candidate, но не параллельную копию того же подхода. Недоступная automatic Luna
не блокирует, пока существует разрешённый mode-compatible economical fallback;
её отсутствие не разрешает неограниченный расход дефицитного profile. Явный
пользовательский role profile не подменяй.

Останови новые attempts, когда candidate прошёл acceptance/final gate,
дальнейшие waves перестали добавлять независимые гипотезы/evidence, исчерпан
envelope или появился настоящий authority/environment blocker. Без final review
режим не завершает scope и не превращается в `Экономичный` автоматически.
