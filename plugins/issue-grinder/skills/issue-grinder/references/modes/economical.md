# Экономичный

Применяет `IG-MODE-06`, mode-specific части `IG-MODE-08..09` и
`IG-MA-15..19` только когда сохранённый `canonical_mode=economical`. Общие
resolver, authority, recovery и evidence rules бери из
[Execution modes](../execution-modes.md); остальные mode-файлы не читай.

Обязателен `IG-MODE-20`: exact текущий основной `gpt-5.6-luna/max`.
При другом или неизвестном profile откажись до effects/dispatch; дорогая
оболочка и supervisor вместо подходящего root запрещены.

Цель — максимальный безопасный прогресс почти без расхода более дефицитного
profile. Economical controller/supervisor и workers могут анализировать,
реализовывать, тестировать, проводить self-review, independent critique и
bounded Best-of-N. Не сохраняй бесконтрольное множество вариантов: своди его к
одному recommended candidate.

При Luna top-level текущая сессия совмещает координацию, последовательную
реализацию и интеграцию. Не создавай supervisor или writer только для формального
разделения ролей; они нужны при конкретной пользе независимой работы либо
изоляции контекста с учётом стоимости передачи и проверки. Экономь также Luna
tokens. Независимый итоговый reviewer остаётся отдельным неавтором кандидата;
Max profile, правила изоляции и checkpoint gate не меняются.

Все содержательные решения и работа режима выполняются Luna Max: scope analysis,
repository research, decomposition, implementation, tests, preliminary и final
self-review, independent critique и reduction. Каждый child dispatch явно
задаёт `model="gpt-5.6-luna"`, `reasoning_effort="max"`, bounded `fork_turns` и
проходит routing guard. Имя и тип child не заменяют проверку его effective
profile.

Каждый exact candidate, включая простой или малый scope, до terminal acceptance
получает independent Luna Max review owner, который не был его автором. Для
Luna top-level это один direct child проверочной волны. Reviewer возвращает
один finding ledger; его отсутствие запрещает terminal result, но может
сохраняться как deferred gate resumable checkpoint.

Новый direct owner проходит один заранее разрешённый routing guard, один spawn и
один event-driven wait до результата, запроса внимания или deadline. Не ищи
guard, не читай его `--help`, не делай status polling, повторные `list` или
пустые nudges при неизменном состоянии. После material rework продолжи того же
owner-а и ту же reviewer session одним follow-up и одним новым event-driven wait
без replacement guard/spawn.

При ambiguity, contract/context conflict, unexpected tool/environment state,
scope expansion или proof gap Luna сохраняет checkpoint/evidence и прекращает
одинаковые corrective retry. Режим может оставить resumable checkpoint. Явный
пользовательский role profile не подменяй.

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
