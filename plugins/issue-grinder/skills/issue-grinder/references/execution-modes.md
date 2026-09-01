# Execution modes

Применяет общую часть `IG-MODE-01..11` и выбирает ровно один mode-specific
runtime contract. Прочитай этот reference после доказательства
нового/продолжающегося run и до стратегической декомпозиции. В продолжающемся
run с уже сохранённым mode record не выбирай режим повторно; восстанови record
и загрузи соответствующий ему mode-файл. Вернись сюда только для явного
переключения или общего profile/review решения.

## Mode record — hard gate

Один run имеет один record:

```text
canonical_mode: solo | classic | balance | swarm | economical
mode_origin: explicit | automatic
initial_main_model: <exact effective model>
initial_main_effort: <exact effective effort when available>
controller_profile: <effective profile>
worker_profile: <effective profile>
```

До первого implementation dispatch:

Порядок разрешения фиксирован: `explicit user mode → automatic model rule`.

1. Если это доказанное продолжение, восстанови record из run continuity и не
   применяй automatic resolver снова. Смена current model/effort не меняет уже
   выбранный mode.
2. Для нового run найди в текущем prompt явное намерение выбрать режим. Прими
   каноническое имя и свободные формы вроде «Соло», `single`, `сингл`, «одним
   агентом», «без субагентов», «используй профиль Balance», «запусти Рой» или
   «работай экономично», только когда речь явно о режиме Issue Grinder.
   Случайное слово «баланс» или «один» в описании продукта не является
   selector-ом.
3. При явном selector-е выбери его. Иначе exact top-level model семейства
   `gpt-5.6-luna` при любом effort выбирает `economical`; любая другая модель —
   `classic`.
4. `По умолчанию` запускает то же automatic правило и не создаёт шестой mode.
5. Выполни profile normalization ниже, сохрани record в доступной continuity
   текущего run и один раз коротко назови выбранный канонический режим в
   progress update.

Не читай config, прошлый run или выбранные child profiles вместо effective
top-level model текущего нового run. Не пересчитывай automatic mode из процента
квоты, доступности reviewer-а или модели, на которой позже продолжили работу.

## Profile normalization

Economical baseline:

```text
model = "gpt-5.6-luna"
reasoning_effort = "max"
```

Явная команда пользователя в prompt о model/effort всех либо конкретных ролей
имеет приоритет. Обычный выбор top-level model в UI является входом resolver-а,
но не означает «все agents обязаны наследовать этот profile».

Без role override:

- main profile семейства Luna при любом effort до `max` включительно:
  `controller_profile = worker_profile = Luna Max`;
- иной main profile: сохраняй его для controller/reviewer, Luna Max используй
  как economical worker, если только repository evaluation не содержит
  отдельного надёжного правила, что main profile не сильнее Luna Max;
- неизвестное cross-family отношение не угадывай по цене, имени или одному
  прошлому результату.

В `Соло` не применяй эти profiles к исполнению: единственный execution profile
равен exact effective current top-level model и effort этого turn. Не создавай
Luna Max supervisor/worker и не подменяй current main profile. При продолжении
после смены top-level model/effort сохрани canonical mode `Соло`, но следующую
работу выполняй уже фактически текущим root profile. Нормализованные поля mode
record могут сохраняться только для безопасного явного переключения в другой
режим, но в `Соло` не дают права вызвать соответствующего агента.

Если root уже запущен на Luna ниже Max и mode не `Соло`, оставь его
единственной transport/authority оболочкой. Он вызывает direct built-in Luna
Max supervisor с bounded packet для scope analysis, dispatch decision или
review и сам сохраняет Goal, Task Manager writes, fan-in и publication unit.
Supervisor не пишет comment, не вызывает Strategic Explainer и не создаёт
отдельную Codex task/session. Его результат — decisions, facts, evidence и
resolvable anchors для root.

## Общие инварианты

Во всех режимах:

- live scope, authority, Production prohibition и truthful lifecycle не
  меняются;
- startup recovery всегда предшествует fresh candidate/worktree;
- каждый concurrent writer проходит собственный worktree admission;
- один effect owner владеет Goal, Task Manager writes, fan-in и publication;
- один exact integrated candidate является предметом acceptance;
- worker self-report и isolated green tests не являются final evidence;
- intentional cheap failures не размножают deployment, shared mutable state,
  external recipients, paid external calls или terminal Task Manager effects;
- после material rework exact changed candidate проходит применимый review
  заново.

## Выбранный режим — обязательная загрузка

После сохранения mode record полностью прочитай ровно один соответствующий
файл. Он является единственным runtime-местом с mode-specific topology,
role/profile routing, fallback, review и stop promise:

| `canonical_mode` | Mode contract |
|---|---|
| `solo` | [Соло](modes/solo.md) |
| `classic` | [Классический](modes/classic.md) |
| `balance` | [Баланс](modes/balance.md) |
| `swarm` | [Рой](modes/swarm.md) |
| `economical` | [Экономичный](modes/economical.md) |

Не загружай остальные четыре mode-файла «для сравнения» во время delivery и не
собирай mode-specific policy из `mode-help.md`, Architecture, старого run или
соседнего файла. Общие resolver, normalization, invariants, switch barrier и
review packet остаются в этом reference; выбранный файл их не переопределяет.

## Explicit mode switch

Только явная команда пользователя меняет mode незавершённого run. Перед сменой:

1. прекрати новые dispatch;
2. доведи active writers до commit/checkpoint или безопасно останови их;
3. выполни `assert-unchanged` integration checkout и reconcile ownership;
4. сохрани exact candidates, checks и deferred effects;
5. запиши новый canonical mode с `mode_origin=explicit`;
6. примени его только к следующей wave/review decision; если новый mode —
   `Соло`, поставь сохранённые packets в последовательную очередь и не создавай
   новых subagents.

Переключение не меняет scope, target environment или authority и не позволяет
повторно реализовать уже сохранённую работу.

## Review packet

Для controller/reviewer подготовь:

- exact selector, base и candidate identity;
- problem/acceptance и integrated diff с source anchors;
- deterministic и aggregate checks с raw evidence;
- material decisions и отклонённые alternatives;
- unresolved objections, negative evidence, known defects и unknowns;
- migrations, shared-state, permissions и внешние effects;
- один точный вопрос либо требуемое решение reviewer-а.

Summary помогает найти evidence, но не заменяет чтение существенного кода и
фактов reviewer-ом.
