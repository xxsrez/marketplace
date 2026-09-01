# Соло

Применяет `IG-MODE-11` и mode-specific части `IG-MA-*` только когда сохранённый
`canonical_mode=solo`. Общие resolver, authority, recovery и evidence rules
бери из [Execution modes](../execution-modes.md); остальные mode-файлы не читай.

Цель — терминально выполнить любой выбранный scope строго текущей моделью без
внутренней агентной topology.

1. Полностью прочитай live scope, существенный source context, dependencies,
   acceptance и risk surfaces текущей основной моделью.
2. Сохраняй одну execution lane: выбери один dependency-ready issue или пакет,
   доведи его текущую итерацию до проверенного status transition либо настоящего
   blocker gate и только затем переходи к следующему.
3. Не вызывай subagents ни для анализа, ни для реализации, ни для тестов,
   критики, проверки, supervisor routing, альтернативных candidates, reduction
   или публикационного текста. Не заменяй их отдельной Codex task/session.
   Число subagents и одновременных execution lanes всегда равно `0` и `1`.
4. Единственный execution profile — exact effective current top-level model и
   effort этого turn. Controller/worker normalization не используется для
   dispatch. Смена current profile не меняет сохранённый canonical mode.
5. Текущая модель сама реализует, интегрирует, запускает применимые проверки и
   проводит self-review exact result. Её уверенность и self-report не являются
   evidence; terminal acceptance опирается на наблюдаемые checks и факты.
6. Используй native communication mode: Strategic Explainer не вызывается,
   однако causal explanation и post-explanation reflection сохраняются.
7. Goal, Task Manager writes, recovery, authority и environment rules остаются
   общими. Несколько issue не включают Goal автоматически вне `IG-GOAL-01` и не
   разрешают делегацию.
8. Продолжай до empty active scope либо принятого terminal blocker handoff.
   Resumable checkpoint `Экономичного` режима недоступен.
