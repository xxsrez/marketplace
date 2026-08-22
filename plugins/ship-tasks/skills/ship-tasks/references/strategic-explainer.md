# Strategic Explainer handoff для ShipTask

Использовать для каждого обязательного lifecycle/blocker comment. Explainer
помогает понять и выразить смысл; он не выбирает facts, status, scope, authority,
repair или terminal outcome.

## Что передать

`$ship-tasks:strategic-explainer` передаётся короткий handoff:

```text
Problem to solve
- Для кого результат и какой наблюдаемый outcome нужен.
- Exact Task scope.

Current-State Brief
- Какой исход уже установлен ShipTask и на каких фактах.
- Что работает, не работает, не проверено или не относится.
- User impact, confidence и реальные ограничения.
- Какой status/action уже разрешён исходной authority.

Strategic discovery anchors
- Exact Task/ref и ближайшие Project/Release/Epic/spec/design sources.

Decision support request
- Объяснить результат человеческим языком.
- Для verification blocker: рекомендовать strongest feasible способ получить
  evidence; сравнить alternatives только при реальном material выборе.
```

Не передавать готовый желаемый вывод, полный tool transcript или process diary.
Discovery только bounded/read-only и прекращается, когда дополнительный source
уже не меняет смысл.

## Что сделать с ответом

ShipTask проверяет forward trace к current facts и reverse coverage всех
decision-relevant входов. Затем сам пишет окончательный Task comment, сохраняя:

- проблему и outcome;
- impact;
- причинную границу и confidence;
- constraint;
- next state/action.

Source note остаётся внутренним основанием. В comment не упоминаются subagent,
Strategic Explainer, handoff или orchestration.

Обязательный результат — применённый Strategic Explainer и понятный comment.
Конституция не задаёт конкретный invocation или recovery flow. Если handoff не
содержит содержательной проблемы либо противоречит current facts, красивый текст
не компенсирует неверное состояние.
