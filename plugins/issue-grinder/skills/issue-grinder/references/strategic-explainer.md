# Strategic Explainer routing

Применяет `IG-FLOW-03..04`, `IG-GOAL-04..07` и publication часть остальных
переходов. Strategic Explainer — optional communication interface, не источник
scope, authority, evidence, repair или status decision.

## Выбор mode

В начале run сначала проверь активную модель. При Astra (`gpt-6-astra`)
автоматическая publication unit всегда получает `native` mode: Strategic
Explainer не вызывается, потому что Astra сама формулирует текст. Этот guard
имеет приоритет над availability provider-а и сохранённым mode и применяется к
каждой новой unit. Явный отдельный запрос пользователя к
`$strategic-explainer:strategic-explainer` остаётся самостоятельным direct call,
но не включается автоматически этим workflow.

Во всех остальных случаях в начале run установи communication mode:

- `ordinary`, если доступен и разрешён
  `$strategic-explainer:strategic-explainer`;
- иначе `native` без capability warning и без блокировки lifecycle.

Это правило одинаково для всех трёх execution modes. Provider-agent Strategic
Explainer не входит в Issue Grinder execution topology: он получает только
publication request и не анализирует, не реализует, не тестирует и не проверяет
delivery scope. Если пользователь отдельно запретил вообще любых subagents,
используй `native`.

Каждый Task Manager comment, общий blocker-report, отдельный ответ по каждой
причине блокировки и финальный Goal comment являются отдельными publication units.
Обычный `To Do → In Progress` comment не создаёт. Routine chat progress не отправляй в
Explainer.

Всю ordinary publication unit исполняет основной coordinator. Рабочие
subagents возвращают ему только проверенные facts, evidence и resolvable
read-only anchors; им нельзя поручать готовый comment или вызов facade. После
такого handoff coordinator сам вызывает semantic facade из top-level
collaboration surface. Если direct child transport там недоступен, применяется
обычный native fallback ниже; отдельная Codex task/session не создаётся.

Для ordinary unit вызови semantic facade только с назначением, исходной
ситуацией/вопросом, exact scope, языком, material constraints и разрешимыми
read-only anchors. Не передавай provider topology, profile, internal protocol
или editing method. Caller error исправь по возвращённой причине и создай
корректный semantic call. При genuine technical/unavailability error перейди в
native и не имитируй provider method; остановка допустима только когда
обязательный безопасный текст невозможно сформировать либо ошибка одновременно
сломала более широкую необходимую capability.

Каждый завершившийся facade call закрыт навсегда. Не сохраняй его provider-child
как канал и не используй `followup_task` или `send_message` для исправления либо
следующей publication unit: caller repair и новая unit получают новый top-level
semantic call. Автоматическое продолжение Goal с уже принятым неизменным
blocker-handoff новой unit не создаёт и вообще не вызывает facade.

## Reflection gate

Publication text и source basis читаются как reflection input. Проверь их по
current primary sources; сам Explainer result ничего не доказывает. Ready text
не переписывай для вкуса. Эта повторная проверка — обязательная
`post-explanation reflection`.

Для Task comment отмени publication/status transition, если обнаружены
неподтверждённый факт, неверная causal связь или существенная обязательная
работа issue. Для final отмени completion при новом active member, существенном
остатке внутри current issue contracts или доступном обязательном шаге. Новый
strategic gap и optional improvement остаются follow-up и не создают второй
completion gate.

Для candidate blocker:

1. до explanation перечитай safe frontier;
2. подготовь причинный draft ordinary/native;
3. извлеки из draft/source basis возможные пути продолжения;
4. проверь их по current primary sources и authority;
5. любой подтверждённый существенный in-scope action отменяет blocker;
6. terminal report публикуется только после доказанного отсутствия такого пути.

Terminal report связывает exact scope с попыткой, primary cause, checkpoint,
unverified remainder, impact, нужным user action и observable resume condition.
Reason code, raw error или tool diary этого не заменяют. Общий report перечисляет все
подтверждённые current причины. Для каждой причины отдельный answer должен объяснить:
почему она блокирует обязательный результат активного issue, почему Issue Grinder
не может устранить её сам и что заблокированный шаг даст issue contract и общей
цели. Весь комплект готовится и проходит reflection до первой
публикации; stale reason или найденный safe path отменяют весь blocker candidate.

Goal status `blocked` может следовать только после публикации общего report, всех
отдельных ответов и platform blocker audit. Принятый комплект показывается
пользователю и завершает текущую попытку даже если audit ещё не разрешает Goal
mutation; в таком случае Goal остаётся активным. Следующий автоматический turn
с тем же blocker fingerprint не повторяет report, reason answers или user
request.

Финальный Goal comment возвращается только пользователю в чате. Task Manager
получает только issue-level comments, требуемые его lifecycle.
