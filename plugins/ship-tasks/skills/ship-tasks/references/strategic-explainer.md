# Смысловая проверка комментариев ShipTask

Пока effective topology rule не отключает comment Explainer, каждый комментарий,
который создаёт ShipTask, до публикации проходит отдельного субагента с
`$ship-tasks:strategic-explainer`. Это обязательная независимая проверка, а не
рекомендация по стилю. Если rule отключает Explainer, основной агент применяет
тот же contract напрямую и не заявляет о независимой проверке. Strategic Explainer переводит факты на понятный человеку
язык, но не выбирает факты, статус, границы работы, полномочия, исправление или
итог. Для `verification-blocked` он также помогает сравнить уже grounded test
paths и сформулировать рекомендацию, не превращая её в разрешение на mutation
или status decision.

По умолчанию этот user-facing judgment packet наследует current model/effort, а
не Luna Max. Явно выбранный пользователем profile имеет приоритет.

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

Для blocker decision report вход дополнительно должен содержать autonomous
self-service frontier: какие fixtures/seed data агент создал сам, какие ingress
и проверки уже выполнены, почему они не закрыли criterion, какой authority gap
остаётся, и какие feasible test paths доступны. Запрос на внешний principal,
вторую сессию или provider evidence должен быть показан как boundary authority,
а не как голая просьба о недостающем объекте.

Для `critical-codebase-accepted` вход дополнительно содержит точную
непроведённую functional check, причину и доказательство существенной роли
human verifier вместо bounded unlocker, исчерпанные autonomous paths, exact
candidate и criteria, проверенные code/tests, grounded verdict critic и residual
risk. Explainer не вправе сгладить эту границу до обычного `verified-success`.

При delegated use форма входа свободна, но субагент должен быть независим от хода
реализации: ему не передаются прежний диалог, полный журнал инструментов или
process diary.
Вход самодостаточно описывает проблему и текущие факты. Discovery остаётся
bounded/read-only и прекращается, когда дополнительный source уже не меняет
смысл.

## Что должно получиться

Strategic Explainer возвращает готовый пользовательский текст и короткий source
basis. Для blocker report текст дополнительно сохраняет recommended test path,
material alternatives/trade-offs, prerequisites/authority, observable success
signal и exact resume condition. ShipTask проверяет факты, но не переписывает
текст обратно на своём техническом языке. Итоговый Task comment сохраняет:

- проблему и outcome;
- impact;
- причинную границу и confidence;
- constraint;
- next state/action.

Для критической приёмки текст явно говорит, что полноценная функциональная
проверка не проводилась и Task закрывается по независимой критической проверке
кодовой базы. Он простым языком объясняет причину, достаточность fallback
evidence и остаточный риск; отсутствие этой оговорки делает comment непригодным.

Внутренняя provenance остаётся основанием для проверки. В comment не
перечисляется orchestration, если она не имеет пользовательского значения.

Текст пишется на языке пользователя. Внутренняя сущность сначала получает
понятную человеку роль; точное техническое название сохраняется только для
полезной проверки, навигации или действия. Смесь языков, внутренний жаргон и
перечень выполненных инструментов не заменяют объяснение.

Когда effective rule сохраняет Explainer, если основной агент нашёл фактическую ошибку или
потерянный существенный факт, он исправляет вход и повторяет независимую
адаптацию. Самостоятельно переформулировать текст и признать его прошедшим
Explainer нельзя. Если
отдельный субагент недоступен или не вернул пригодный текст, комментарий не
публикуется и связанный переход остаётся незавершённым. Когда user rule
отключает Explainer, основной агент отвечает за factual и quality review
напрямую.
