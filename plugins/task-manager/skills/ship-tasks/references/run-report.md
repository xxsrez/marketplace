# Причинный анализ и финальный ShipTask run report

Прочитать перед blocking input, terminal Goal decision или финальным ответом.
Task comment объясняет одну Task; этот report объясняет человеку весь run и не
заменяет Task comments, checks, status read-back или external evidence.

## Перед blocked

Нельзя завершать run на первом симптоме. До утверждения «продолжить невозможно»:

1. Сформулировать симптом и impact без protocol jargon.
2. Восстановить causal timeline: expected step, observed state, last successful
   step и первое расхождение.
3. Проверить вероятную root cause текущими authoritative reads/tools. Для
   capability issue отдельно проверить фактически доступные operations,
   authority, attempted outcome и current Task state; не выводить отсутствие
   capability только из отсутствующего результата.
4. Перечислить safe in-scope recovery actions и выполнить их, если они уже
   разрешены. Обычный retry после установленной transient cause, missing
   reconciliation read, доступный comment/status write, refresh или другой
   предусмотренный project recovery step не требуют нового user decision.
5. После recovery перечитать affected Task/Goal/source/effect state и повторить
   terminal decision.

Если доступное действие устраняет blocker, `blocked` запрещён: продолжить run.
Если действия исчерпаны и требуется внешнее решение/state change, показать
анализ пользователю до допустимого `update_goal(status="blocked")`. Повторный
poll без нового evidence не считается отдельной recovery attempt или новым
blocker occurrence.

## Human-readable explanation

Начинать с 2–5 предложений простым языком:

- что агент пытался завершить;
- что фактически произошло;
- почему это произошло — root cause, а не только symptom;
- смог ли агент исправить это сам и что именно сделал;
- если не смог — почему требуется человек или внешний state change.

После этого дать проверяемые детали. Reason code разрешён только как appendix,
не как замена объяснению. При недоказанной причине писать `вероятная причина`
и confidence, а не категоричное утверждение.

## Финальный формат

Каждый terminal exit skill заканчивается компактным отчётом:

```text
SHIPTASK RUN REPORT
Outcome: COMPLETED | PARTIAL | BLOCKED | NO WORK

Итог
<Что получил пользователь или почему результат пока неполный.>

Что произошло и почему
<Понятная causal narrative; evidence и inference разделены.>

Что агент проверил и сделал
<Ключевые diagnosis/recovery/actions и их результат.>

Состояние scope
<Task counts/identifiers/statuses, result/release identity, comments и Goal.>

Что осталось
<Zero remaining work либо один точный unresolved blocker/decision.>

Следующий шаг
<Одно конкретное действие или not-applicable.>

Technical evidence
<Exact refs, SHAs, checks, deploy/read-back identities; только полезные детали.>
```

Для completed/no-work не выдумывать incident section: коротко объяснить outcome
и доказательства. Для partial/blocked обязательно показать impact, root cause с
confidence, recovery attempts, почему дальнейшее self-recovery невозможно и
exact resume step.

Перед Goal `blocked` этот report должен уже доказать, что:

- тот же blocker удовлетворяет строгому tool threshold;
- complete inventory не содержит доступной runnable или recovery work;
- агент не может устранить причину имеющимися operations и authority;
- требуемое внешнее изменение или решение названо точно;
- Task comments имеют truthful disposition;
- после status write финальный ответ показывает фактический Goal state.
