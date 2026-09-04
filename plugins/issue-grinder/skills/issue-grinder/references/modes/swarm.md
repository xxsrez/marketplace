# Рой

Применяет `IG-MODE-05`, mode-specific части `IG-MODE-08..09` и
`IG-MA-15..19` только когда сохранённый `canonical_mode=swarm`. Внутреннее имя
`swarm` сохраняется для совместимости continuation, но current topology —
Manager Loop, а не Best-of-N. Общие resolver, authority, recovery и evidence
rules бери из [Execution modes](../execution-modes.md); остальные mode-файлы не
читай.

Цель — заменить длительное пошаговое исполнение Sol/controller-а постоянными
Luna Max manager и implementer sessions, сохранив терминальный результат,
независимую проверку и один точный кандидат.

## Control brief и постоянные роли

До первой implementation Sol/controller одним целостным проходом изучает весь
live scope и материализует control brief: Strategic Outcome, Human Requirements,
Agent Plan, exact base/candidate, dependency и risk maps, acceptance/oracles,
external-effect gates, конечный phase envelope и final gate. Он не превращает
это в пошаговую реализацию.

После startup recovery создай три Luna Max роли с отдельными routing receipts:

1. одну постоянную manager session без repository paths, shell, source access,
   FileChange и Task Manager authority;
2. одну постоянную implementer session с одним admitted writer worktree либо
   task-owned shadow tree и одним exact candidate на весь loop;
3. одну независимую reviewer session только после manager `complete`.

Каждый первоначальный dispatch явно задаёт `model="gpt-5.6-luna"`,
`reasoning_effort="max"`, bounded `fork_turns` и проходит routing guard. Имя или
`agent_type` не заменяет observed profile. Sol/controller сохраняет Goal,
Task Manager, external effects, fan-in, integration и final acceptance, но не
выполняет routine research, implementation, tests, phase decisions или cheap
review.

## Manager Loop

Manager получает control brief без source tree и возвращает ровно один пакет:

```text
phase:<stable id> | goal | acceptance | material dependencies | evidence expected
rework:<same phase id> | reproduced defect | exact acceptance delta
complete | accepted phases | remaining unknowns | final-review packet
escalate | one material question | evidence | blocked decision
```

Manager поддерживает общий plan state, выбирает только следующую
dependency-ready фазу и после evidence принимает `accept`, `targeted rework`,
`complete` либо одну material escalation. Фаза объединяет связный проверяемый
outcome; отдельный checklist item, файл или мелкий тест не создаёт manager turn.
Одновременно активна ровно одна phase/rework wave.

Coordinator механически пересылает phase packet в ту же implementer session, а
её compact evidence — обратно той же manager session. Он не пересказывает и не
расширяет пакеты, не выбирает следующую фазу и не дублирует Luna work. Если
доказан прямой совместимый transport между sessions, его можно использовать с
теми же identities и packet contract; sibling messaging не является
предусловием режима.

Implementer применяет текущий packet к тому же exact candidate, делает один
целевой source pass, реализует фазу, запускает относящиеся к ней checks и
self-review и сразу возвращает:

```text
phase id | changed surfaces | checks | findings<=3 | unknowns | next
```

Он не перечитывает Issue Grinder `SKILL.md`, references/Architecture, routing
guard или tool catalog, не ищет parent messaging и не меняет общий phase plan.
Все mutations используют absolute candidate path. Новая фаза не начинается до
решения manager-а по предыдущей. Повтор без нового reproducer, evidence или
materially иного способа запрещён.

Большой Teams-подобный scope дели на небольшое число крупных фаз по реальным
dependency/integration boundaries, например contract/data model,
authorization/API, UI/synchronization и сквозная verification. Это пример, а не
фиксированный шаблон; не дели механически по числу Tasks, файлов или tests.
Phase action budget конечен и соразмерен её outcome, но универсальный потолок
малого Balance packet-а не применяется к большой фазе `Роя`.

Best-of-N, параллельные реализации, массовые scouts/critics, reducer и
использование свободных слотов запрещены normal path. После доказанного тупика
manager может назначить один materially иной восстановительный путь, сохранив
старое negative evidence и не запуская варианты для голосования.

## Независимая проверка и завершение

Reviewer запускается только после manager `complete` и сам читает exact final
candidate, requirements, diff/source anchors и checks. Он сохраняет одну
постоянную session, конечный review plan и один finding ledger. Даже на большом
scope reviewer последовательно проходит крупные risk sections сам и не
делегирует descendants. Open-ended fuzzing, cleanup из read-only роли, поиск
messaging tool и повторное чтение mode policy запрещены.

Каждый material finding содержит reproducer/evidence и disposition
`fixed | refuted_with_evidence | escalate`. Targeted rework получает прежняя
implementer session; exact changed candidate одним follow-up получает прежняя
reviewer session. Новый reviewer или implementer допустим только после
доказанной недоступности прежнего, потери независимости либо material изменения
контракта, с сохранением уже собранного evidence.

Новый direct owner проходит `guard ×1 → spawn ×1 → event-driven wait` до
результата, внимания или переданного stage deadline. Follow-up существующей
session использует один новый event-driven wait без повторного guard/spawn.
Технический timeout до deadline продолжает то же ожидание без status polling,
`list`, commentary или nudge. Final response session является её handoff.

После reviewer pass coordinator механически переносит task-owned bytes/commit в
integration target, проверяет identity и exact integrated candidate, запускает
основной suite и передаёт Sol/controller-у сжатый final packet. Final gate имеет
один source/evidence pass и не воспроизводит manager/implementer exploration.
Без independent review и final gate `Рой` не завершает scope и не превращается
автоматически в `Экономичный`.

Mode-invalid, но потенциально функционально полезным считается прогон, где
manager читает source или использует рабочие tools, implementer заменён без
причины, одновременно выполняются две фазы, reviewer стартует до `complete`,
reviewer создаёт собственный рой либо Sol/controller выполняет routine
implementation. Сохрани такие отклонения в evidence и исправь topology до
terminal acceptance.
