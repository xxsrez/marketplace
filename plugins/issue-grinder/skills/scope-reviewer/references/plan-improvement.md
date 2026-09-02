# Автоматическое улучшение плана

Этот reference загружается только в Plan improvement после явного mutation
intent. Plan review и Release review остаются read-only и не читают его как
разрешение на writes.

## Admission каждого repair

До первой mutation и перед каждым изменением подтверди:

- exact planning scope и write authority;
- current snapshot и неизменившуюся Requirements version;
- однозначную границу Human Requirements/Agent Plan;
- что Task принадлежит planning surface, а не активной delivery/rework;
- что patch не меняет problem, desired outcome, требования, constraints,
  non-goals и exact user scope;
- что Task Manager поддерживает нужную обратимую planning operation.

Если Human Requirements находятся в общем description, сохрани их точные bytes
до write и после read-back подтверди byte-identical значение. Невозможность выделить или сохранить границу
— `representation blocker`: repair этой Task запрещён. Не переписывай весь
description «по смыслу» и не считай собственный пересказ эквивалентной защитой.

## Разрешённая поверхность

При прошедшем admission можно уточнять agent-owned Task wording и technical
boundary, acceptance, expected evidence, existing-label assignment, hierarchy,
relations, sequencing и decomposition. Изменение должно устранять подтверждённый
finding и сохранять историю.

Не разрешены:

- любые изменения Human Requirements;
- implementation, tests, release, Goal или рабочие status transitions;
- автоматическое создание/изменение Label taxonomy;
- destructive deletion ради чистой схемы;
- подмена отсутствующей native operation текстовой пометкой;
- repair активной delivery-проблемы как будто это ещё planning defect.

Если перестройка требует supersede/cancel operation вне подтверждённой authority
или capability, оставь точный proposed repair и не имитируй выполненное
изменение.

## Concurrency и read-back

Каждый write использует canonical ref, current version, optimistic concurrency и
stable idempotency semantics. После mutation перечитай объект и связанные
hierarchy/relations. Unknown outcome сначала reconciles через reads; blind retry
запрещён.

Сравни Human Requirements до и после всего repair отдельным integrity gate.
Любое неподтверждённое отличие прекращает дальнейшие writes и становится
`requirements-integrity failure`. Partial repair перечисляет confirmed, unknown
и not-applied changes и не называется успешно обновлённым планом.

## Повторная проверка и readiness

После confirmed repair создай новый snapshot. Узкая правка повторяет зависимые
оптики и общие integrity gates; Requirements или material graph change требует
полного применимого review.

Останови цикл, когда план ready, остались только решения человека по
Requirements, repair недоступен из-за capability/authority/conflict либо новый
проход не получил нового evidence или materially different action. Повтор той
же генерации без changed source не является прогрессом.

Ready означает одновременно:

- Human Requirements однозначны для scope либо явно выделены вопросы человеку;
- coverage map не содержит скрытого обязательного остатка;
- decomposition, ownership и dependencies согласованы;
- acceptance и expected evidence наблюдают каждый обязательный outcome;
- material auto-fix findings устранены и подтверждены read-back;
- риски, review points и реальные действия человека названы честно;
- final snapshot остаётся current.

Readiness — оценка плана, а не разрешение запустить delivery и не evidence
реализации.
