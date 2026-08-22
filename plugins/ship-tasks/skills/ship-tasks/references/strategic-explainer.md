# Strategic Explainer quality contract для ShipTask

Каждый обязательный lifecycle/blocker comment и final report должен выражать
смысл на уровне проблемы, facts, impact, evidence/unknown и next state.
Strategic Explainer помогает добиться этого качества, но не выбирает facts,
status, scope, authority, repair или terminal outcome.

## Какие основания нужны

Независимо от способа работы агенту нужны:

- для кого предназначен result, какой observable outcome нужен и какой exact
  Task scope рассматривается;
- какой outcome уже установлен ShipTask, что доказано/failed/unverified и на
  каком evidence;
- user impact, confidence, реальные constraints и исходная authority;
- ближайшие relevant Project/Release/Epic/spec/design sources, если они меняют
  смысл;
- следующий state и подтверждённый action, если он нужен.

Форма context свободна. Не подменять facts готовым желаемым выводом и не
переносить полный tool transcript или process diary. Discovery остаётся
bounded/read-only и прекращается, когда дополнительный source уже не меняет
смысл.

## Что должно получиться

ShipTask проверяет, что material claims имеют source basis, а
decision-relevant facts не потеряны. Итоговый Task comment сохраняет:

- проблему и outcome;
- impact;
- причинную границу и confidence;
- constraint;
- next state/action.

Внутренняя provenance остаётся основанием для проверки. В comment не
перечисляется orchestration, если она не имеет пользовательского значения.

Обязательный результат — понятный grounded comment. `$ship-tasks:strategic-explainer`
можно вызвать напрямую, делегировать ему adaptation или не вызывать отдельно,
если основной workflow уже удовлетворяет quality contract. Конституция не
задаёт agent topology, invocation или recovery flow. Красивый текст не
компенсирует missing problem, неверный state или потерянный факт.
