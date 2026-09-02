# Автоматический вход и live-run handoff

Этот reference обязателен, когда Scope Reviewer выбран явно или автоматически
для предзапускового ревью либо изучения активного долгого delivery run. Он
применяет `SR-13` и не расширяет Task Manager-only scope.

## Два автоматических trigger-класса

Выбирай Scope Reviewer без обязательного skill mention только когда запрос по
смыслу относится к одному из классов:

1. **Pre-launch review.** Человек просит проверить, объяснить или оценить
   готовность выбранного Task Manager плана до передачи в delivery: например,
   «проверь план перед запуском», «готовы ли мы начинать» или «объясни, что
   будем делать».
2. **Active long-run review.** Человек просит целостно понять текущий активный
   Issue Grinder / Task Manager delivery run: например, «посмотри, что сейчас
   происходит», «как идёт этот долгий запуск» или «где мы находимся и нужна ли
   моя помощь».

Implicit invocation требует пользовательского запроса. Не запускай review по
таймеру, из-за длительности run как таковой, после каждого Task Composer result
или перед каждым Issue Grinder start. Raw status одной карточки, ordinary Task
lookup и произвольная Codex task без однозначного Task Manager delivery scope
остаются вне Scope Reviewer.

Relative формулировка, включая «посмотри на неё», допустима только когда current
conversation и сохранённая continuity доказывают ровно один active run и один
Task Manager selector. При нуле или нескольких кандидатах запроси exact scope;
не выбирай последний, самый новый или наиболее похожий объект.

## Pre-launch review

Предзапусковый запрос выбирает `Plan review`. Он read-only, если пользователь
отдельно и однозначно не попросил улучшить план. Readiness report не запускает
Issue Grinder, не меняет planning или working statuses и не создаёт Goal.

## Active long-run review

Рассматривай live review как read-only sidecar owning Issue Grinder run, а не
как новый delivery run:

1. На ближайшей безопасной coordination boundary сохрани continuity владельца:
   run identity, exact selector, Goal ref/status при наличии, execution mode и
   его origin, task-owned checkpoints, active worker ownership и незавершённые
   effects.
2. Не создавай отдельную пользовательскую Codex task/session и не сбрасывай,
   не пересчитывай и не присваивай себе Goal, mode, scope, worker branches или
   lifecycle authority.
3. Построй current Scope Reviewer snapshot через live Task Manager reread.
   Factual run snapshot и resolvable candidate/check anchors могут дополнить
   его; worker narrative остаётся `in-flight`, пока primary evidence не
   подтверждает результат.
4. Не dispatch-ь новую delivery работу ради review и не отменяй уже работающих
   workers только для получения снимка. Перед итогом перечитай version vector;
   material change инвалидирует затронутые findings.
5. Верни один read-only report. Если пользователь добавил обзор к продолжающемуся
   run, передай управление owning Issue Grinder с сохранённой continuity, чтобы
   он продолжил с доказанного checkpoint. Явные stop, pause, cancel или
   replacement instructions сильнее автоматического возврата.

Отчёт не становится acceptance evidence, Task Manager comment, status
transition или новым blocker-ом. Scope Reviewer сообщает состояние; Issue
Grinder остаётся единственным владельцем delivery и дальнейшего исполнения.
