# Экономичный

Применяет `IG-MODE-06`, mode-specific части `IG-MODE-08..09` и
`IG-MA-15..17` только когда сохранённый `canonical_mode=economical`. Общие
resolver, authority, recovery и evidence rules бери из
[Execution modes](../execution-modes.md); остальные mode-файлы не читай.

Цель — максимальный безопасный прогресс почти без расхода более дефицитного
profile. Economical controller/supervisor и workers могут анализировать,
реализовывать, тестировать, проводить self-review, independent critique и
bounded Best-of-N. Не сохраняй бесконтрольное множество вариантов: своди его к
одному recommended candidate.

При ambiguity, contract/context conflict, unexpected tool/environment state,
scope expansion или proof gap Luna сохраняет checkpoint/evidence и прекращает
одинаковые corrective retry. Режим может оставить resumable checkpoint.
Недоступная automatic Luna не блокирует, пока существует разрешённый
mode-compatible economical fallback; её отсутствие не разрешает
неограниченный расход дефицитного profile. Явный пользовательский role profile
не подменяй.

Выполняй все доступные deterministic и aggregate checks. Сохраняй raw results,
known defects, unknowns, rejected candidates и deferred gates. `In Review`
допустим только для действительно review-ready candidate; иначе оставь
`In Progress`.

Если terminal acceptance доказан обычным evidence, заверши scope. Иначе текущая
попытка может закончиться только после checkpoint gate:

- один exact recommended candidate;
- task-owned commit либо честный owned dirty checkpoint;
- base, branch/worktree и integration identity;
- выполненные checks с raw results;
- known defects, unknowns и deferred review/authority gates;
- точный следующий шаг и resume condition;
- правдивые Task Manager statuses и активный Goal.

Этот выход называется `resumable checkpoint`, не `complete` и не `blocked`.
Не вызывай `update_goal(complete|blocked)` только из-за него. Покажи
пользователю, что готово, что не принято и откуда продолжать. Дополнительные
cheap attempts после достаточного checkpoint запускай только при material
expected gain, а не ради занятости слотов.
