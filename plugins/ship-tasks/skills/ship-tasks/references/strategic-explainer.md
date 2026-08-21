# Strategic Explainer handoff для ShipTask

Прочитать перед каждым новым ShipTask Task report comment и перед terminal
user-facing handoff. Общий behavioral contract находится в sibling-skill
[`$strategic-explainer`](../../strategic-explainer/SKILL.md).

## Когда вызывать

Вызвать свежий Strategic Explainer обязательно:

- перед каждым новым Task report comment: `COMPLETED`, `REWORK REQUIRED`,
  `BLOCKED` или `CANCELED`;
- перед scope-level `PARTIAL` или `BLOCKED` user-facing handoff;
- когда для продолжения нужен user action, decision, новая authority, другой
  человек/аккаунт или external state change;
- для aggregate batch result, который не представлен одним task-scoped handoff.

Простой success не освобождает Task comment от Explainer pipeline. Обычная
промежуточная red/green iteration не создаёт comment и не является trigger.

## Выполнить parent finalization

Перед spawn перечитать exact current Task/result/effects и самостоятельно
зафиксировать report state, evidence, verified/unverified scenarios, current
constraints, safe self-recovery и допустимый next action. Explainer не выбирает
state, blocker, authority или terminal transition.

Отдельно сформулировать содержательную задачу, которую решает exact scope. Она
должна следовать из canonical Task description/acceptance и явного user/project
context. Task ref или технический title сами по себе недостаточны. ShipTask не
должен передавать готовый strategic view: его Explainer найдёт самостоятельно.

## Сформировать bounded Strategic Handoff

Передать один self-contained handoff:

```text
Target surface: TASK_COMMENT | RUN_REPORT
Artifact scope: <exact Task ref or aggregate run identity>
Authoritative report state: <state already chosen by ShipTask>

Problem to solve:
Beneficiary: <для кого предназначен результат>
Desired outcome: <наблюдаемое изменение или capability>
Exact scope: <canonical Task/run ref>
Why now/current gap: <только если materially relevant>

Reader purpose:

Current-State Brief:
Confirmed outcome:
Unfinished or unknown:
Observed user impact:
Evidence and confidence:
Current capability and attempts:
Constraints:
Candidate user dependency:
Next-state contract:
Output language and channel:
Required authoritative envelope:

Strategic discovery anchors:
Canonical Task/project/release refs:
Repository or project context:
Direct parent/design/spec links already attached to scope:

Scenario:
Expected user-visible behavior:
State: VERIFIED | FAILED | UNVERIFIED | NOT_APPLICABLE
Evidence basis:
Observed user impact:
Needed input:
```

Это input contract, не схема результата. Удалить неприменимые optional fields и
изложить факты связным текстом, если так яснее. Обязательные части —
`Problem to solve`, authoritative current-state facts и starting anchors.

`Candidate user dependency` включать только когда исходное evidence уже
показывает, что основной агент не может продолжить без человека или external
state. `Next-state contract` называет actor, минимальное действие, причину,
observable success signal и что ShipTask сможет продолжить после него.

Для material или multi-scenario ситуации заполнять отдельный scenario ledger.
`Needed input` не является разрешением просить пользователя. Не смешивать
independent scenarios, full conversation, process diary, raw logs, desired
wording или готовый strategic conclusion.

`Strategic discovery anchors` задают точку старта и scope, но не делают research
за Explainer. ShipTask передаёт exact refs и direct links, которые уже известны;
не пересказывает Epic/design/vision вместо самостоятельного discovery.

## Запустить свежий субагент правильно

Использовать built-in `default` agent с уникальным task name
`strategic_explainer_<surface>_<scope>` и без унаследованной истории. Передать
точный `fork_turns="none"`; `fork_turns="all"`, положительное число fork turns и
продолжение старого thread запрещены. Не задавать model override.

Initial task содержит только инструкцию применить установленный catalog skill
`$ship-tasks:strategic-explainer` и bounded `Strategic Handoff`. Full
conversation, inherited tool transcript и process diary не передаются.

Явно потребовать следующий порядок:

1. до анализа проверить context integrity;
2. затем до любого tool call проверить содержательный `Problem to solve`;
3. при чистом input самостоятельно выполнить bounded strategic discovery через
   доступные read-only tools;
4. не выполнять writes, recovery, status/Goal decisions или external actions;
5. вернуть свободное стратегическое объяснение и короткий parent-facing source
   basis, не structured result и не copy-ready comment.

Для ShipTask discovery разрешено только read-only:

- перечитать canonical Task и доступные relations/external context;
- пройти к exact Project/Release, parent/Epic или linked higher-level scope;
- в repository/project context искать product vision, high-level design,
  current specification и accepted ADR через bounded file search;
- отличать `current/accepted`, `proposed` и `historical` documents;
- остановиться, когда найденный уровень больше не меняет problem framing.

Не разрешать Task Manager mutations, broad log reading, source-code archaeology,
unrelated filesystem search или web research вне declared scope. Technical files
читать только когда без них нельзя понять причинную модель, impact, action или
confidence.

## Обработать исправляющие отказы

### Context integrity

При `CONTEXT_INTEGRITY_ERROR` не выполнять comment/status/Goal writes. Исправить
orchestration и один раз создать новый default subagent с `fork_turns="none"`.
Повторный отказ останавливает report workflow как внутреннюю orchestration
failure, а не blocker доставляемой Task.

### Missing problem

При `PROBLEM_CONTEXT_ERROR` Explainer не должен вызывать tools или продолжать
analysis. ShipTask перечитывает canonical Task description/acceptance и explicit
user/project context, исправляет handoff и один раз запускает новый fresh
subagent. Нельзя заменять semantic problem одним Task ref/title.

Если ShipTask всё ещё не может честно сформулировать beneficiary и desired
outcome, он останавливает report workflow до comment/status/Goal writes и прямо
запрашивает недостающий problem context. Не объявлять это product defect или
terminal blocker Task и не переходить к local wording fallback.

### Strategic source gap

Если Explainer не нашёл дополнительный strategic source, это не error: source
note фиксирует отсутствие и explanation опирается на явно переданную проблему.
Если current/accepted sources materially противоречат problem/current-state
brief либо обязательный source недоступен, ShipTask возвращается к finalization
и разрешает factual/context conflict до user-visible write. Не выбирать удобную
версию и не сглаживать conflict текстом.

## Использовать результат безопасно

Explainer возвращает два смысловых слоя: свободное problem-first объяснение и
короткий source basis для parent. Он не создаёт facts, authority, lifecycle
status или решение. ShipTask проверяет:

- forward trace каждого material утверждения к `Problem to solve`,
  `Current-State Brief` или exact discovered source;
- reverse coverage каждого decision-relevant current-state fact;
- source state: proposed/historical context не выдан за current;
- current execution evidence не переопределён design document;
- technical details остались только там, где меняют causal model, impact/risk,
  action или confidence.

После этого ShipTask самостоятельно пишет Task comment или user-facing handoff
своими словами. Он сохраняет problem, strategic meaning, outcome, impact,
confidence, подтверждённую dependency и next state. Authoritative envelope,
exact status и evidence refs добавляются по исходным фактам. Parent-facing source
note не копируется механически; exact source оставить пользователю только для
полезной навигации или проверки.

В user-visible тексте не упоминать субагента, Strategic Explainer, Strategic
Handoff или внутреннюю orchestration. Нельзя копировать ответ механически,
противоречить ему, добавлять неподтверждённый смысл или возвращаться к process
diary.

## Freshness, reuse и degraded adaptation

Переиспользовать explanation между `TASK_COMMENT` и `RUN_REPORT` можно только
когда совпадают problem, audience, artifact scope, facts, state, strategic source
basis, user dependency и next action. Planned comment/read-back и terminal
status reconciliation допустимы, если они совпали с next-state contract и не
выявили drift. Recovery, divergent write, changed problem/strategy или другое
material meaning change делают explanation stale.

Для aggregate batch, нескольких blocked Tasks или другого audience запустить
новый scope-level Explainer. Task-level context одного member не представляет
стратегический смысл всего batch.

Если subagent tools или `$ship-tasks:strategic-explainer` действительно
недоступны, ShipTask сам применяет тот же problem gate, bounded read-only
discovery, source-state classification и fidelity checks. Зафиксировать
`degraded-adaptation` только во внутреннем evidence. Tool unavailability не
разрешает raw technical summary; `CONTEXT_INTEGRITY_ERROR` и
`PROBLEM_CONTEXT_ERROR` не являются недоступностью и требуют исправленного fresh
invocation.
