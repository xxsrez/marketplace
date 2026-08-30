# Execution modes

Применяет `IG-MODE-01..11` и mode-specific часть `IG-MA-*`. Прочитай этот
reference после доказательства нового/продолжающегося run и до стратегической
декомпозиции. В продолжающемся run с уже сохранённым mode record не выбирай
режим повторно; вернись сюда только для явного переключения или mode-specific
dispatch/review решения.

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

## Соло

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
4. Текущая модель сама реализует, интегрирует, запускает применимые проверки и
   проводит self-review exact result. Её уверенность и self-report не являются
   evidence; terminal acceptance опирается на наблюдаемые checks и факты.
5. Используй native communication mode: Strategic Explainer не вызывается,
   однако causal explanation и post-explanation reflection сохраняются.
6. Goal, Task Manager writes, recovery, authority и environment rules остаются
   общими. Несколько issue не включают Goal автоматически вне `IG-GOAL-01` и не
   разрешают делегацию.
7. Продолжай до empty active scope либо принятого terminal blocker handoff.
   Resumable checkpoint `Экономичного` режима недоступен.

## Классический

Цель — полный терминальный результат с высокой уверенностью.

1. Controller/reviewer самостоятельно читает весь live scope и существенный
   source context, строит стратегию, dependency graph, acceptance и risk map.
2. Product, architecture, cross-issue, migration, concurrency, security,
   shared-state и иное weak-oracle judgment остаются этому profile.
3. Controller/reviewer остаётся основным исполнителем и делает почти всю
   implementation сам. Luna Max получает только действительно тривиальный
   strict-simple packet: bounded, self-contained, disjoint, objectively
   verifiable, low-risk и без material judgment.
4. Обычная delegation разделяет полезные независимые packets и read-only
   critique. Не создавай competing full implementations по умолчанию.
5. Один integration owner выполняет fan-in и aggregate checks.
6. Controller/reviewer читает exact integrated diff, material source и raw
   evidence, проводит final code/result review и только затем разрешает Done.
7. Продолжай до empty active scope либо принятого terminal blocker handoff.

## Баланс

Цель — terminal result с меньшим расходом controller/reviewer profile.

1. Controller/reviewer выполняет первоначальный full-scope analysis,
   декомпозицию, acceptance и risk classification.
2. Сразу оставляет себе packets, где само решение требует material judgment.
   Если сложное решение отделимо от исполнения, передай economical worker-у
   уже принятое решение и bounded implementation contract.
3. Luna Max является предпочтительным исполнителем лёгких и средних bounded
   packets и выполняет основную массу implementation, research, tests и
   preliminary verification. Средний packet может требовать содержательной
   работы, но не material product/architecture judgment и должен иметь ясные
   границы и проверяемый contract. Независимые economical critics/test authors
   атакуют candidate до дорогого review.
4. Дополнительный candidate создавай только при реальной развилке, проваленном
   oracle или высокой ожидаемой ценности независимой дешёвой проверки.
5. Integration owner собирает содержательную batch и формирует review packet.
6. При существенной неопределённости, contract/context conflict, слабом oracle
   или проблеме за границами packet-а Luna сохраняет изменения, checks,
   evidence и причину остановки и возвращает fallback controller/reviewer-у без
   одинакового retry. Тот уточняет contract, продолжает сам либо создаёт новый
   безопасный packet. Reviewer также может вернуть bounded material rework в
   новую economical wave. После rework повтори exact-candidate gate.
7. Без final gate не выдавай непроверенный result за terminal; режим сам не
   переключай.

## Рой

Цель — использовать большой объём дешёвого поиска, реализации и критики с
одним проверяемым выходом.

До wave задай конечный envelope: явный user budget либо bounded число волн с
условиями остановки. Свободная capacity не является бесконечным budget.

Допустимые роли: scouts, design candidates, implementation candidates,
minimal-diff/reliability alternatives, critics, test authors, economical
judges и evidence reducers. Различай подходы; множество одинаковых prompts
создаёт коррелированные ошибки и не является достаточным разнообразием.

Intentional candidate получает purpose, base, identity, branch/worktree и
verification. Он может затрагивать те же files, что другой candidate, только в
полной изоляции. Existing checkpoint сначала обнаруживается; новый candidate
не называется его resume/replacement.

После каждой wave:

1. deterministic checks удаляют явно несостоятельные варианты;
2. economical critics/judges сравнивают оставшиеся, сохраняя provenance,
   dissent и negative evidence;
3. reducer оставляет один recommended candidate и максимум один runner-up при
   действительно material unresolved fork;
4. integration owner принимает только task-owned commit выбранного candidate;
5. reviewer получает compressed review packet, а не transcripts всех agents;
6. замечания создают новую bounded rework wave и повторный exact review.

Останови новые attempts, когда candidate прошёл acceptance/final gate,
дальнейшие waves перестали добавлять независимые гипотезы/evidence, исчерпан
envelope или появился настоящий authority/environment blocker. Без final review
режим не завершает scope и не превращается в `Экономичный` автоматически.

## Экономичный

Цель — максимальный безопасный прогресс почти без расхода более дефицитного
profile. Economical controller/supervisor и workers могут анализировать,
реализовывать, тестировать, проводить self-review, independent critique и
bounded Best-of-N. Не сохраняй бесконтрольное множество вариантов: своди его к
одному recommended candidate.

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
