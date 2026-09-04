# Соло

Применяет `IG-MODE-11` и mode-specific части `IG-MA-*` только когда сохранённый
`canonical_mode=solo`. Общие resolver, authority, recovery и evidence rules
бери из [Execution modes](../execution-modes.md); остальные mode-файлы не читай.

Цель — терминально выполнить любой выбранный scope одним исполнителем Issue
Grinder на текущей модели без рабочей делегации его delivery-семантики.

1. Полностью прочитай live scope, существенный source context, dependencies,
   acceptance и risk surfaces текущей основной моделью.
2. Сохраняй одну execution lane: выбери один dependency-ready issue или пакет,
   доведи его текущую итерацию до проверенного status transition либо настоящего
   blocker gate и только затем переходи к следующему.
3. Не передавай другим агентам анализ scope, source research, реализацию, тесты,
   техническую критику, проверку, supervisor routing, альтернативных candidates,
   reduction или integration decision. Не заменяй такую делегацию отдельной
   Codex task/session. Число Issue Grinder execution-subagents и одновременных
   содержательных execution lanes всегда равно `0` и `1`.
4. Единственный execution profile — exact effective current top-level model и
   effort этого turn. Controller/worker normalization не используется для
   dispatch. Смена current profile не меняет сохранённый canonical mode.
5. Текущая модель сама реализует, интегрирует, запускает применимые проверки и
   проводит self-review exact result. Её уверенность и self-report не являются
   evidence; terminal acceptance опирается на наблюдаемые checks и факты.
6. Для publication используй общий Strategic Explainer routing. При активной
   Astra (`gpt-6-astra`) routing выбирает native writing без вызова provider-а;
   в остальных случаях его отдельный provider-agent не входит в execution
   topology `Соло`, пока получает только publication request и не выполняет
   анализ, реализацию, тесты или review delivery scope. Coordinator сохраняет
   causal reflection и решение.
7. Goal, Task Manager writes, recovery, authority и environment rules остаются
   общими. Несколько issue не включают Goal автоматически вне `IG-GOAL-01` и не
   разрешают рабочую делегацию Issue Grinder. Внешний controller, создавший эту
   сессию, и ограниченные service/provider agents не считаются исполнителями
   режима; если им передана delivery-работа, это ограничение нарушено.
8. Продолжай до empty active scope либо принятого terminal blocker handoff.
   Resumable checkpoint `Экономичного` режима недоступен.

Независимые инструменты внутри текущего пакета можно выполнять одновременно по
`IG-MODE-19`, не создавая вторую execution lane. Самопроверка переиспользует
применимые checks по `IG-MODE-09`, не снижая доказательность.
