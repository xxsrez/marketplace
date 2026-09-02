# Scope и независимый review

Этот reference обязателен для Plan review, Plan improvement и Release review.
Он задаёт согласованный источник анализа, выбор оптик и формат findings, но не
разрешает ни одной mutation.

## Exact selector и полное чтение

Разреши workspace и один exact Task Manager selector. Project/Release/Epic/Task
ref, состав явно выбранного набора и current context должны давать единственное
толкование. Память, максимальный номер, последнее изменение или похожий title не
выбирают scope.

До review:

1. пройди pagination каждого inventory до terminal cursor;
2. прочитай full detail каждой materially relevant Task;
3. восстанови parent/child hierarchy, relations, statuses, Requirements и Agent
   Plan;
4. прочитай comments и внешние evidence anchors только когда они нужны для
   заявленного состояния;
5. запиши manifest из selector, canonical refs, versions и source anchors.

Сохраняй различие:

- **Human Requirements** — problem, desired outcome, пользовательские
  требования, constraints и non-goals;
- **Agent Plan** — decomposition, Task boundaries, acceptance, expected
  evidence, hierarchy, relations, sequencing и technical detail.

Отдельные native поля и versions предпочтительны. Структурированная граница в
общем description допустима только когда Human Requirements можно извлечь и
после write подтвердить byte-identical. Requirement projection дочерней Task не
становится новым source смысла и сверяется с canonical Requirements.

## Snapshot integrity

Оптики анализируют один логический snapshot. Перед каждым плановым write и
перед итоговым report перечитай version vector. Изменение одной Task
инвалидирует выводы, зависящие от неё; изменение Requirements или task graph
требует нового согласованного snapshot и полного materially applicable review.
Не соединяй старые findings с новым состоянием без явной переоценки.

Различай fact, plan, interpretation, risk, unknown и unverified result. Lifecycle
projection Task Manager не доказывает код, integration, tests, deployment или
product outcome без соответствующего primary evidence.

## Выбор materially useful оптик

Для review плана Requirements integrity — отдельная обязательная оптика. Затем
оцени необходимость классов:

- outcome/decomposition, coverage и ownership;
- dependencies, sequencing и critical path;
- acceptance, expected evidence и observable readiness;
- risks, authority, review и human-attention boundaries;
- применимые domain/product, UX, data/migration, compatibility,
  security/privacy, reliability, performance/cost или integration/release
  concerns.

Для Release review выбирай оптики по current goal, delivered/in-flight scope,
evidence, critical path и рискам. Формальный одинаковый набор для всех Releases
не нужен.

Оптика допустима, только если её вывод способен изменить понимание outcome,
план, Requirement question, acceptance/evidence, material risk/dependency,
human action, readiness или confidence. Разные названия не оправдывают
дублирующий анализ.

## Clean Luna packet

Каждая compact task содержит:

- точную строку role lock `SCOPE_REVIEWER_LENS_V1`;
- одну названную оптику и её exact question;
- snapshot identity и resolvable read-only anchors;
- запрет writes, orchestration, delegation, Strategic Explainer и решений за
  coordinator-а;
- требуемый result packet.

Не передавай conversation history, caller reasoning, прежний candidate,
соседние findings или желаемый итог. Luna самостоятельно читает только
необходимые anchors и возвращает:

- краткий человеческий вывод;
- material facts и source anchors;
- влияние на outcome/readiness;
- один disposition: `auto-fixable`, `human-requirement-decision`,
  `risk/observe` или `no-material-finding`;
- предлагаемый repair либо минимальный вопрос человеку, если они обоснованы;
- uncertainty и confidence boundary.

Lens не пишет Tasks, не меняет Requirements, не решает lifecycle и не вызывает
других agents. Failure или недоступность оптики — coverage gap, а не PASS.

## Синтез findings

Coordinator подтверждает каждое material утверждение primary sources. Дубли
объединяются с strongest evidence; factual conflict возвращается к sources;
настоящее disagreement остаётся видимым. Optional improvement не становится
blocker-ом, а Requirement issue — auto-fix-ом.

Для плана построй внутреннюю coverage map:

```text
Requirement → Task/plan element → acceptance → expected evidence
```

Карта выявляет скрытый остаток, duplicate ownership и обязательный outcome без
observable proof. Она служит основанием readiness, но не навязывает форму
пользовательского отчёта.
