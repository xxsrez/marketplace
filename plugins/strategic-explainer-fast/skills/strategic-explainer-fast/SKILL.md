---
name: strategic-explainer-fast
description: "Подготовить один готовый user-facing comment, report, blocker explanation, final или явную редактуру прямо в текущем контексте без subagent. Использовать для быстрого понятного объяснения с source basis; не использовать для routine chat, progress, mutations или authority decisions."
---

# Strategic Explainer Fast

Это in-context communication skill для одной реальной publication unit. Он
сохраняет требования к смыслу обычного Strategic Explainer, но намеренно не
создаёт subagent и не обещает clean/stateless isolation.

## Admission

Применяй skill только к одному целостному пользовательскому результату:
комментарию, отчёту по Task или scope, объяснению material решения или
состояния, blocker report, final либо явной редактуре конкретного текста.
Routine chat, progress update, внутренний draft, mutation и authority decision
не относятся к этому API.

До discovery установи:

- один publication unit и исходный вопрос;
- exact scope;
- resolvable read-only anchors к current facts и materially relevant sources.

Несколько смешанных результатов раздели на отдельные проходы. Если scope или
обязательный anchor неоднозначны, назови точный gap и не угадывай. Наличие
унаследованного диалога, tool transcript или прежнего candidate допустимо, но
не превращает их в evidence или framing.

## Выполнение

1. Не создавай subagent ни для generation, ни для self-review.
2. Полностью прочитай `references/in-context-contract.md`.
3. Самостоятельно выполни bounded read-only discovery по exact anchors и
   authoritative sources. Отделяй их от process diary, старых гипотез и мнения
   caller.
4. Верни один publication-ready текст и отдельно обозначенный короткий source
   basis. Не смешивай доказательный журнал с сообщением человеку.
5. Проверь material facts по sources. При changed facts, scope или
   comprehension error выполни новый сфокусированный проход, не ремонтируй
   старый candidate поверх недостоверного framing.

Fast остаётся read-only: он не выбирает status, scope, recovery, release или
authority, ничего не публикует и не выполняет mutations. Он не заявляет
независимость или stateless-гарантии обычного Strategic Explainer.
