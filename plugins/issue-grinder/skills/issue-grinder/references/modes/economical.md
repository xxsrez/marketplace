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

Все содержательные решения и работа режима выполняются Luna Max: scope analysis,
repository research, decomposition, implementation, tests, preliminary и final
self-review, independent critique и reduction. Каждый child dispatch явно
задаёт `model="gpt-5.6-luna"`, `reasoning_effort="max"`, bounded `fork_turns` и
проходит routing guard. Имя и тип child не заменяют проверку его effective
profile.

Если top-level root уже запущен не на Luna, он является только
transport/authority оболочкой: разрешает live refs, хранит Goal и receipts,
выполняет Task Manager mutations, механический fan-in и publication. До любой
содержательной repository analysis, source mutation, test authoring/execution
или review он вызывает direct Luna Max supervisor. Supervisor возвращает
решения и готовые bounded packets, а root механически dispatch-ит их direct Luna
workers. Root не считается экономичным исполнителем и не переносит эту работу
на Sol при нехватке Luna. Для настоящего end-user run без Sol
top-level session также должна быть Luna Max.

При ambiguity, contract/context conflict, unexpected tool/environment state,
scope expansion или proof gap Luna сохраняет checkpoint/evidence и прекращает
одинаковые corrective retry. Режим может оставить resumable checkpoint.
Недоступная automatic Luna уменьшает capacity: используй доступную serial Luna
lane, а если её нет — сохрани достаточный resumable checkpoint. Не подменяй её
Sol/GPT-5.4 и не продолжай содержательную дорогую работу под именем
`Экономичный`. Явный пользовательский role profile не подменяй.

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
