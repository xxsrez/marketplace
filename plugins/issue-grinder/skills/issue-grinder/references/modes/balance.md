# Баланс

Применяет `IG-MODE-04`, mode-specific части `IG-MODE-08..09` и
`IG-MA-15..19` только когда сохранённый `canonical_mode=balance`. Общие
resolver, authority, recovery и evidence rules бери из
[Execution modes](../execution-modes.md); остальные mode-файлы не читай.

Цель — сохранить terminal outcome, автономность и сильную приёмку
`Классического`, но существенно уменьшить расход controller/reviewer profile.
Оптимизируй принятый результат на единицу дефицитной квоты, а не число Luna,
общее количество токенов или красивую параллельность.

## Дорогой control plane

1. Controller/reviewer одним целостным проходом выполняет первоначальный
   full-scope analysis: восстанавливает Strategic Outcome и Human Requirements,
   строит dependency/risk map, acceptance, package boundaries и integration
   points.
2. Сохрани `expensive-work ledger`. В нём допустимы только конкретные
   `material_judgment`, `integration_decision` и `final_review`. Для каждого
   другого содержательного действия controller-а назови, почему его нельзя
   отделить и передать Luna; статус root/coordinator сам по себе причиной не
   является.
3. Если сложное решение отделимо от исполнения, controller принимает только
   решение и передаёт Luna новый bounded contract. Обычные repository research,
   implementation, test authoring/execution, preliminary critique и rework на
   дорогой lane не переносятся.

## Luna-owned packet loop

4. До первой source mutation или targeted test создай первую Luna Max execution
   stage. Обычно один Luna packet owner получает self-contained contract и
   владеет research, implementation, tests и self-review одного bounded
   candidate. Material fork, слабый oracle либо высокая ожидаемая ценность могут
   оправдать несколько purpose-distinct candidate owners от общей exact base;
   свободные слоты сами по себе этого не оправдывают. Одновременно с execution
   stage запусти независимого Luna reviewer-а только для построения bounded
   review plan по requirements, risk map и исходной base; candidate до handoff
   он не оценивает. Если execution owner наблюдаемо умеет nested delegation, он
   может добавить mode-compatible
   внутренние роли. Если нет, это normal platform shape, а не proof gap:
   candidate завершается этим owner-ом, после чего coordinator запускает
   отдельную direct Luna review stage.
5. Каждый execution-child получает уникальный `packet_id`, зелёный
   `issue-grinder/model-routing/v2` receipt, точные
   `model="gpt-5.6-luna"`, `reasoning_effort="max"`, bounded `fork_turns` и
   совпадающий dispatch fingerprint. Packet lead может создавать только
   mode-compatible Luna execution-children с такими же receipts и не получает
   Task Manager, publication, integration либо final-acceptance authority.
6. Каждый exact candidate до дорогого final gate, включая простой или малый
   scope, передай независимому Luna verifier/critic. При недоступной nested
   delegation это отдельный direct read-only owner: его plan может готовиться
   параллельно, но exact review начинается только после завершения candidate
   stage. Если material fork дал несколько candidates, этот же независимый owner
   сначала сводит их по oracle и negative evidence, затем проверяет выбранный
   exact candidate.
   Дай ему current requirements, exact diff/source anchors и oracle, но не
   используй итоговый self-report автора как доказательство. Отсутствующая либо
   незавершённая проверка остаётся незакрытым gate и запрещает terminal result.
7. Packet result содержит один exact candidate либо checkpoint, owned surfaces,
   checks с raw results, source anchors, material decisions, known defects,
   unknowns и finding ledger. Для каждого material finding допустим только
   disposition `fixed | refuted_with_evidence | escalate`; несколько общих
   одобрений не перевешивают один воспроизводимый дефект.
8. Исправление возвращай тому же Luna candidate owner. Новая wave оправдана новым
   evidence, изменённым состоянием, иной гипотезой или сохраняющейся ожидаемой
   ценностью; одинаковый retry без новой информации запрещён.
   После material rework продолжи того же reviewer-а с exact changed candidate;
   replacement guard/spawn без доказанной причины не создавай.

Каждый новый direct stage owner проходит заранее разрешённый routing guard и
один spawn; `event-driven wait` длится до результата, запроса внимания или
переданного stage deadline. Не используй произвольные десятиминутные отсечки.
Технический timeout ожидания до stage deadline означает немедленно продолжить
то же событийное ожидание без `list`, status probe, commentary или nudge. Owner
не ищет messaging tool: его final response является одним consolidated handoff.

Execution packet заранее ограничивает полезные actions и checks. Review packet
задаёт один source/diff pass, один основной deterministic suite и только
оправданные targeted risk probes; для малого/среднего scope по умолчанию не
больше трёх. Open-ended fuzzing, повторное чтение неизменившихся файлов,
исследование parent messaging tools и cleanup из read-only reviewer-а запрещены.
После rework reviewer проверяет исходный reproducer, изменённые surfaces и один
основной suite; новый общий поиск допустим только при доказанно новой risk
surface, созданной изменением.

## Адаптивная избыточность

9. По умолчанию используй один implementation candidate и независимую дешёвую
   проверку. Дополнительный candidate создавай только при реальной развилке,
   проваленном или слабом oracle, неудаче прежнего подхода либо высокой
   ожидаемой ценности отдельной попытки.
10. Различай purpose и approach дополнительных candidates. Множество одинаковых
    prompts не является независимостью. Reducer сохраняет один рекомендуемый
    result и весь material negative evidence; решение принимается по findings и
    evidence, а не большинством голосов.

## Узкая эскалация и final gate

11. При material uncertainty, contract/context conflict, unexpected
    tool/environment state, scope expansion, proof gap или проблеме за
    границами packet-а Luna останавливает corrective mutations и возвращает
    compact evidence packet с одним узким вопросом либо требуемым решением.
12. Controller/reviewer разрешает именно неделимую неопределённость. После
    решения создай новый bounded Luna contract и верни отделимое исполнение или
    material rework в economical wave. Controller продолжает весь packet сам
    только когда judgment неотделим от исполнения; причину сохрани в
    `expensive-work ledger`.
13. Integration owner объединяет только task-owned commits, проверяет exact
    integrated candidate и формирует компактный review packet. Summary служит
    навигацией к source/check evidence; сырой transcript всех Luna waves не
    является обязательным reviewer input.
14. Только после завершения Luna candidate/review/rework stages
    controller/reviewer проводит final code/result gate точной интегрированной
    версии: один source/diff pass и один integrated deterministic suite. Он не
    повторяет exploratory probes Luna без противоречия в evidence.
    После material rework повтори exact-candidate gate. Без final gate не
    выдавай непроверенный result за terminal.
15. Первый содержательный Luna handoff является gate режима. Отсутствующий или
    отрицательный routing receipt, несовпавший observed profile, скрытая
    ordinary Sol/GPT-5.4 work либо основная содержательная работа controller-а
    делают прогон невалидным `Балансом`, даже если функциональный результат
    успешен. Не исправляй routing failure дорогой реализацией под прежним mode.
16. Недоступная automatic Luna уменьшает capacity. Используй совместимую serial
    Luna lane либо сохрани evidence и честно остановись/попроси явную смену
    режима. Самостоятельно режим не переключай.
