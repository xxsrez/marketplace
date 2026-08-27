# Критическая приёмка по кодовой базе

Используй этот fallback только после обычной приёмки и автономной test frontier.
Он сознательно даёт `Done` более слабую доказательственную силу, но не разрешает
выдавать непроведённую функциональную проверку за проведённую.

## Eligibility gate

Перечитай свежий полный inventory текущего selector. `Backlog` и terminal
statuses не входят в active counts. Fallback разрешён, только если:

- `To Do == 0`, `In Progress == 0`, `In Review > 0`;
- каждая оставшаяся `In Review` Task уже `verification-blocked` после всех
  доступных безопасных способов обычной проверки;
- каждый blocker требует человека как содержательного verifier, а не как
  bounded unlocker;
- exact integrated candidate стабилен и никто его параллельно не изменяет.

Approval, MFA, invite, bounded access grant и другое малое действие, после
которого агент сам может завершить проверку, не открывают gate. Не открывают его
также tool inconvenience, первая неудачная попытка, отсутствующий создаваемый
fixture или ещё доступный supported path.

## Независимый critic

Запусти ровно одного read-only subagent role `critic` с
`fork_turns="none"`. Если effective user topology rule запрещает critic,
fallback недоступен и основной агент не симулирует независимость.

Reviewer получает нейтральный самодостаточный packet: canonical selector или
Task refs, exact candidate identity и поручение самостоятельно перечитать live
Task contracts, код и тесты. Не передавай inherited conversation, producer
rationale, прежний verdict или process diary. Critic самостоятельно повторяет
релевантные tests, проверяет acceptance surfaces и нужные cross-project связи,
ничего не исправляет и возвращает per-Task evidence map и grounded verdict.

## Disposition

- Доказанное нарушение, test failure или отсутствующая реализация — обычный
  `verified-failure` с точной attribution, обязательным comment и
  `In Review → In Progress`.
- Grounded approval по точному коду, самостоятельно повторённым tests и связям —
  `critical-codebase-accepted`, обязательный comment и `In Review → Done`.
- Mere absence of findings, speculation или непокрытый material criterion — не
  approval; Task остаётся `In Review`.
- Изменение candidate, Task contract или active inventory делает review stale;
  перечитай gate до lifecycle writes.

## Обязательный factual packet для completion comment

Перед `Done` availability-selected Strategic Explainer получает pass с anchors
на grounded packet. ShipTask не составляет второй explanation draft: ordinary
provider method остаётся opaque. При отсутствии или opt-out provider-а ShipTask
формулирует comment в native mode.
Factual anchors обязаны включать:

- точную непроведённую функциональную проверку;
- причину и объяснение, почему нужен существенный human verifier, а не approval
  либо иное малое разблокирующее действие;
- исчерпанные autonomous paths;
- exact candidate, criteria, code surfaces и tests;
- grounded verdict critic;
- явное указание, что Task закрывается по критической проверке кодовой базы, а
  не по полноценной функциональной приёмке;
- residual knowledge boundary и риск.

Опубликуй и перечитай comment, затем измени status и перечитай Task. Fallback не
создаёт отсутствующий external effect и не расширяет production, durable-data,
secrets/privacy/access-policy, external-recipient или unbounded-cost authority.
