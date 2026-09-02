# Баланс

Применяет `IG-MODE-04`, mode-specific части `IG-MODE-08..09` и
`IG-MA-15..17` только когда сохранённый `canonical_mode=balance`. Общие
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

4. До первой source mutation или targeted test создай первую Luna Max wave.
   Обычная форма — один Luna packet lead, который получает self-contained
   contract и владеет полным внутренним циклом пакета: research,
   implementation, tests, economical critique и rework. Если nested delegation
   недоступна, coordinator сохраняет тот же дешёвый цикл через direct Luna
   lanes; он не выполняет их работу сам.
5. Каждый execution-child получает уникальный `packet_id`, зелёный
   `issue-grinder/model-routing/v2` receipt, точные
   `model="gpt-5.6-luna"`, `reasoning_effort="max"`, bounded `fork_turns` и
   совпадающий dispatch fingerprint. Packet lead может создавать только
   mode-compatible Luna execution-children с такими же receipts и не получает
   Task Manager, publication, integration либо final-acceptance authority.
6. Существенно изменившийся candidate до дорогого final gate передай
   независимому Luna verifier/critic, когда такая проверка возможна. Дай ему
   current requirements, exact diff/source anchors и oracle, но не используй
   итоговый self-report автора как доказательство. Отсутствие независимой
   проверки сохрани как явный proof gap.
7. Packet result содержит один exact candidate либо checkpoint, owned surfaces,
   checks с raw results, source anchors, material decisions, known defects,
   unknowns и finding ledger. Для каждого material finding допустим только
   disposition `fixed | refuted_with_evidence | escalate`; несколько общих
   одобрений не перевешивают один воспроизводимый дефект.
8. Исправление возвращай в Luna-owned loop. Новая wave оправдана новым
   evidence, изменённым состоянием, иной гипотезой или сохраняющейся ожидаемой
   ценностью; одинаковый retry без новой информации запрещён.

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
14. Reviewer проводит final code/result gate точной интегрированной версии.
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
