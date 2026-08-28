# Strategic Explainer routing

Применяет `IG-FLOW-03..04`, `IG-GOAL-04..06` и publication часть остальных
переходов. Strategic Explainer — optional communication interface, не источник
scope, authority, evidence, repair или status decision.

## Выбор mode

В начале run установи communication mode:

- `ordinary`, если доступен и разрешён
  `$strategic-explainer:strategic-explainer`;
- иначе `native` без capability warning и без блокировки lifecycle.

Каждый Task Manager comment, отдельный blocker-report и финальный Goal comment
является новой publication unit. Обычный `To Do → In Progress` comment не
создаёт. Routine chat progress не отправляй в Explainer.

Для ordinary unit вызови semantic facade только с назначением, исходной
ситуацией/вопросом, exact scope, языком, material constraints и разрешимыми
read-only anchors. Не передавай provider topology, profile, internal protocol
или editing method. Caller error исправь по возвращённой причине и создай
корректный semantic call. При genuine technical/unavailability error перейди в
native и не имитируй provider method; остановка допустима только когда
обязательный безопасный текст невозможно сформировать либо ошибка одновременно
сломала более широкую необходимую capability.

## Reflection gate

Publication text и source basis читаются как reflection input. Проверь их по
current primary sources; сам Explainer result ничего не доказывает. Ready text
не переписывай для вкуса. Эта повторная проверка — обязательная
`post-explanation reflection`.

Для Task comment отмени publication/status transition, если обнаружены
неподтверждённый факт, неверная causal связь или существенная обязательная
работа issue. Для final отмени completion при новом active member, существенном
остатке или доступном обязательном шаге. Optional improvement остаётся
follow-up.

Для candidate blocker:

1. до explanation перечитай safe frontier;
2. подготовь причинный draft ordinary/native;
3. извлеки из draft/source basis возможные пути продолжения;
4. проверь их по current primary sources и authority;
5. любой подтверждённый существенный in-scope action отменяет blocker;
6. terminal report публикуется только после доказанного отсутствия такого пути.

Terminal report связывает exact scope с попыткой, primary cause, checkpoint,
unverified remainder, impact, нужным user action и observable resume condition.
Reason code, raw error или tool diary этого не заменяют. Goal status `blocked`
может следовать только после принятого отчёта и platform blocker audit. Сам
принятый report показывается пользователю и завершает текущую попытку даже если
audit ещё не разрешает Goal mutation; в таком случае Goal остаётся активным.

Финальный Goal comment возвращается только пользователю в чате. Task Manager
получает только issue-level comments, требуемые его lifecycle.
