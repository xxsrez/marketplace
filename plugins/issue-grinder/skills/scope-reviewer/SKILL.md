---
name: scope-reviewer
description: "Проверять однозначно выбранный Task Manager план или Release через независимые Luna Max оптики и возвращать один понятный человеку обзор. Автоматически использовать без явного имени skill-а для ревью плана перед запуском и для изучения активного долгого Issue Grinder / Task Manager delivery run. По явному intent улучшать только agent-owned planning model, не меняя Human Requirements. Не использовать для delivery, lifecycle mutations, ordinary single-Task lookup или произвольной Codex task без Task Manager scope."
---

# Scope Reviewer

Дай человеку целостную картину выбранного плана или Release без необходимости
читать все карточки и технические материалы. Глубоко анализируй scope через
несколько независимых уместных оптик, но публикуй только один лёгкий для чтения
отчёт с существенными выводами.

## 0. Автоматически входи в два класса review

Не требуй явного `$issue-grinder:scope-reviewer`, когда пользователь:

- просит проверить или объяснить выбранный Task Manager план перед запуском;
- просит понять, что происходит в текущем активном долгом Issue Grinder /
  Task Manager delivery run.

Если запрос относится к одному из этих классов, полностью прочитай
[автоматический вход и live-run handoff](references/invocation-and-live-run.md).
Короткое «посмотри на неё» является достаточным только при однозначной current
run continuity. Не срабатывай по таймеру, после каждого plan creation, на raw
single-Task lookup или на generic Codex task без Task Manager scope.

## 1. Разреши intent и exact scope

Выбери ровно один режим:

- **Plan review** — read-only анализ и отчёт о планируемой работе;
- **Plan improvement** — анализ, исправление только agent-owned planning model,
  повторная проверка и отчёт; требует однозначного intent «улучши», «почини»,
  «подготовь план к исполнению» или эквивалентного;
- **Release review** — всегда read-only отчёт о текущем delivery scope.

Просьба показать, объяснить, проверить план, дать status или рассказать, что
происходит с Release, не разрешает writes. Не расширяй разрешение на
implementation, testing, release, Goal, status transitions, production,
destructive/access/privacy/secret/external-recipient decisions.

Через Task Manager adapter разреши один exact Project, Release, Epic, Task или
явный связный набор Tasks. Для live-run review разрешай relative wording через
доказанную current run continuity. Не выбирай current объект по максимальному
номеру, последней дате или похожему имени. Если selector неоднозначен,
остановись до анализа и запроси exact выбор.

## 2. Построй согласованный snapshot

Полностью прочитай [scope и review contract](references/scope-and-review.md).
Пройди все страницы inventory, прочитай materially relevant details, hierarchy,
relations, statuses, Strategic Outcome, Human Requirements, Agent Plan, comments
и применимое evidence. Различай общий ориентир, защищённые обязательства
человека и изменяемый способ выполнения; Task Manager state доказывает только
собственную projection.

Зафиксируй snapshot manifest: exact selector, canonical refs, versions и source
anchors. Все оптики получают один логический snapshot. Перед repair и итоговым
report перечитай version vector; material change требует нового snapshot и
повторной проверки затронутых выводов.

## 3. Запусти независимые оптики

Основной agent сам выбирает минимальный набор materially useful оптик. Для plan
review обязательна отдельная Requirements integrity optic; остальные выводятся
из problem, Requirements, task graph, состояния и рисков. Не создавай две
оптики, если они проверят одно и то же решение.

В каждом review проверь стратегическую связность Tasks и происхождение
обязательств. Широкий Strategic Outcome, agent assumption, рекомендуемый
hardening или удобная практика не становятся Human Requirement, обязательной
проверкой либо blocker. Для этого можно выбрать отдельную оптику или включить
вопрос в другую независимую materially useful оптику.

Каждую оптику исполняет отдельный built-in `default` subagent через top-level
collaboration surface с точными параметрами:

- `fork_turns="none"`;
- `model="gpt-5.6-luna"`;
- `reasoning_effort="max"`;
- одна compact optic task с exact snapshot identity и resolvable read-only
  anchors;
- запрет mutations, orchestration и delegation.

Перед dispatch проверь requested profile bundled
`scripts/lens_routing_guard.py`; invalid receipt запрещает запуск. Не передавай
Luna caller reasoning, готовый общий вывод, ответы других оптик или process
transcript. Result каждой оптики должен быть bounded finding packet: понятный
вывод, source anchors, влияние, disposition, возможный repair/вопрос,
uncertainty и confidence boundary.

Если collaboration или exact Luna Max profile недоступны либо пользователь
запретил необходимые subagents, не имитируй независимый review текущей моделью и
не объявляй полное покрытие. Верни доступный factual snapshot с явным coverage
gap и условием полноценного повторного запуска.

## 4. Синтезируй и при необходимости исправь план

Проверь material claims оптик по primary snapshot sources. Отбрось
`no-material-finding`, объедини дубли, сохрани настоящие различия и disagreement.
Report subagent-а не является evidence. Requirements issue никогда не
переклассифицируй в auto-fix ради продолжения.

Строй две разные связи: `Strategic Outcome → вклад Tasks/известный gap` для
направления и `Human Requirement → plan → acceptance → evidence` для формальных
обязательств. Стратегический gap не создаёт скрытую задолженность.

В Plan improvement полностью прочитай
[plan-improvement contract](references/plan-improvement.md). Human Requirements
не изменяй ни при каких обстоятельствах. Если их границу нельзя однозначно
распознать и доказуемо сохранить, запрети auto-repair этой Task и сообщи
representation blocker. Каждый разрешённый write использует current version,
optimistic concurrency и read-back; unknown outcome сначала reconciles через
reads.

Формулировку Strategic Outcome можно уточнить, а Agent Plan — перестроить только
при сохранённом пользовательском смысле, exact scope и Requirements. Material
изменение problem или desired outcome верни человеку.

После repair построй новый snapshot и повтори применимые оптики. Plan readiness
не запускает Issue Grinder, не меняет `Backlog`/рабочие статусы и не доказывает
реализацию.

## 5. Верни один человекочитаемый отчёт

Полностью прочитай [reporting contract](references/reporting.md). Сначала собери
factual target из current snapshot и принятых findings, затем проверь active
model. При Astra (`gpt-6-astra`) не вызывай Strategic Explainer: coordinator сам
формулирует native report. В остальных случаях передай target отдельному
`$strategic-explainer:strategic-explainer` как явную editing task с exact scope,
языком и resolvable anchors. Не читай provider-internal contract и не передавай
Explainer-у planning, readiness, lifecycle или authority decisions.

Проверь готовый текст на factual conflict и reverse coverage. Потерянный факт
или ошибка требуют нового clean editing call, а не самостоятельной стилистической
переписи. Если Explainer operationally unavailable в не-Astra режиме, верни
собственный factual report и честно отметь отсутствие независимого editorial
pass. В Astra-режиме проверь native report тем же factual и reverse-coverage
gate, без вызова provider-а.

Публикация — один связный ответ, а не dump ответов Luna, task-by-task inventory
или process diary. Сначала дай главный причинный вывод и требуемое действие;
technical refs и служебные квитанции оставь в коротком source basis только если
они нужны для проверки.

В Release review отдельно назови формальную полноту Tasks и достижение
Strategic Outcome. Empty active scope достаточно для формального завершения
Goal; даже подтверждённый стратегический gap не разрешает создавать новую Task,
blocker или удерживать Goal открытым.
