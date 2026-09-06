# Баланс

Применяет `IG-MODE-04`, `IG-MODE-20`, `IG-MODE-08..09`, `IG-MA-15..19`
при `canonical_mode=balance`, `mode_contract_version=luna-coordinator-v1`.
Общие resolver, authority и evidence — [Execution modes](../execution-modes.md).

Основной агент обязательно exact `gpt-5.6-luna/max`; при unknown/mismatch
откажись до Goal, mutations и child dispatch. Не создавай вместо него supervisor.

Luna сама ведёт scope, исследует, реализует обычную работу, тестирует,
исправляет и интегрирует. Не разделяй эти роли формально. Экономь дорогие Sol
токены и общий расход при неизменной полноте результата. Не решай задачу
полностью только ради подробной постановки другому исполнителю.

Sol Extra High (`gpt-5.6-sol`, `reasoning_effort=xhigh`) получает конкретное
решение или целый сложный участок до проверенного результата. Неопределённый
общий контракт либо две разные безуспешные содержательные попытки — повод
передать evidence и точный вопрос Sol. Обычные локализованные исправления
Luna доводит сама. Sol не ведёт статусы, ожидание процессов и управление Luna.

Перед terminal acceptance обязателен один независимый Sol xhigh reviewer
точного интегрированного кандидата по исходным требованиям, коду и рискам,
включая малый scope. Он не автор; если Sol ранее писал код, нужен другой
неавтор. Роль `final_review`; для сложной работы — `specialist`.
Обе роли проходят bundled model_routing_guard.py с exact Sol/xhigh,
`fork_turns=none`. Это предусмотренный режимом маршрут, не user override.
Другие обычные execution-child roles по умолчанию Luna Max; не создавай их
без полезной самостоятельной работы.

Handoff: цель, requirements anchors, owned scope, exact diff, checks/raw
results и неизвестное. Не передавай полную историю. Reviewer самостоятельно
изучает первичные источники, закрывает пробелы собственными проверками и
переиспользует применимые длительные results. После material rework продолжи
того же reviewer по изменённому diff. Самоотчёт автора не доказывает качество.

Один guard → один spawn → одно событийное ожидание до результата или deadline.
Без status polling и nudges; поддерживай процессные handles. Для каждого
child прочитай [multi-agent mechanics](../multi-agent-execution.md) один раз
перед первым admission. Каждый Sol writer получает task-owned worktree через существующий admission;
при недоступной записи Git применяй предусмотренный read-only patch handoff.
Одновременная запись в общий checkout запрещена.

Luna исправляет понятные замечания и проверяет итоговую интеграцию. Перед
завершением обязательны project checks и независимое accepted review.
Не снижай качество ради стоимости. Прогон продолжается до terminal результата
или настоящего blocker; особый resumable checkpoint Экономичного недоступен.
