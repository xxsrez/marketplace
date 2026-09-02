# Человекочитаемый итоговый отчёт

Этот reference обязателен для каждого результата Scope Reviewer. Итог — одна
связная композиция, а не последовательность ответов оптик, task inventory или
технический журнал.

## Factual synthesis

Сначала coordinator строит factual target только из current snapshot и
подтверждённых findings. Удали `no-material-finding`, повторы и подробности,
которые не меняют понимание, решение, действие, риск или confidence. Сохрани
material disagreement, uncertainty и границу evidence.

Для плана читатель должен суметь понять:

- какую проблему решает scope и какой будущий результат ожидается;
- где система находится сейчас и что существенно изменится;
- как Task-ы совместно дают outcome без пересказа каждой карточки;
- что было исправлено автоматически и что подтверждено read-back;
- какие Requirements требуют решения человека;
- какие риски, dependencies, blockers и unknowns остаются;
- где и когда реально нужен человек, а где его участие не требуется;
- готов ли план к исполнению и на каком основании.

Это смысловая полнота, а не обязательные заголовки. Структура и длина следуют
конкретному scope.

## Release review

Release review всегда read-only. Разделяй:

- фактически доказанный product/release outcome;
- current Task Manager lifecycle projection;
- известную незавершённую или in-flight работу;
- critical path и следующий observable state;
- current blockers и лишь возможные будущие risks;
- действие человека, необходимое сейчас, будущую review point и отсутствие
  human dependency.

Task count и status percentage могут быть supporting context, но не заменяют
оценку относительно общей цели. Worker narrative без primary evidence не
становится завершённым результатом. Scope Reviewer не пишет comments, не меняет
status/Goal и не ремонтирует delivery scope.

## Независимый редакторский проход

Передай factual target отдельному
`$strategic-explainer:strategic-explainer` как explicit editing task. Semantic
request содержит ровно одну publication unit, назначение отчёта, исходный
вопрос, exact scope, язык, factual target, material constraints и resolvable
read-only anchors. Не передавай внутренний invocation recipe, методику provider-а,
caller reasoning или authority decisions.

Strategic Explainer отвечает только за publication-ready реконструкцию. Он не
выбирает findings, repair, readiness, lifecycle или действие от имени
coordinator-а. После возврата проверь:

- material claims не противоречат current snapshot;
- Requirements, risks, disagreement и next action не потеряны;
- главное причинное сообщение видно раньше supporting detail;
- человек без технического погружения может пересказать outcome, границы знания
  и нужное действие;
- служебные identifiers и process details не несут основную мысль.

Factual defect или потеря reverse coverage исправляются новым clean editing
call с корректным target/source. Не переписывай готовый текст вручную ради
вкуса. Operational unavailability Explainer-а выбирает собственный factual
report, а не блокирует review; честно обозначь отсутствие независимого
editorial pass.

Верни один report и при необходимости короткий source basis с точными refs,
versions и evidence anchors. Source basis подтверждает текст, но не заменяет
понятное объяснение.
