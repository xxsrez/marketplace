# Финальный отчёт ShipTask

Финальный ответ сообщает человеку состояние всего выбранного run. Он не
заменяет Task comments и не доказывает status write.

## Evidence для ответа

Final report опирается на current affected Tasks/comments/statuses, применимый
Goal массовой имплементации, exact result identity и обязательные external
effects. Он сверяет обещанный результат с фактическим, а acceptance incident
openings — с resolution comments и current Task state. Успешный repair не
стирает найденный defect. Незавершённость можно передавать пользователю только
когда в текущем scope и authority не осталось достаточного безопасного способа
продолжить. Как собрать эти основания, решает агент.

## Содержание

Начать с результата. Затем кратко сообщить:

- что выполнено и текущий Task/Goal status;
- что доказано и что осталось `not-available`;
- compact ledger всех material acceptance incidents run, включая resolved;
- primary cause незавершённости, отдельно от последствий;
- что уже сделано и почему этого недостаточно для завершения;
- user action или внешнее условие, без которого продолжение невозможно;
- первый safe resume step.

Для каждого incident назвать exact Task/criterion, observed failure или границу
знания, cause/confidence, fix, retest evidence и final state. `verified-failure`
называется найденным defect; `verification-blocked` и contract conflict — нет.
Resolved incident остаётся видимым как `found and resolved`. Unresolved incident
располагается рядом с общим result и не совместим с clean success affected Task.

Не писать process diary, полные URL/UUID/raw errors и внутренние названия
оркестрации, если они не нужны пользователю. Incident ledger — audit trail
material outcomes, а не журнал каждой попытки. Если run успешен, объяснить
полученный outcome, важное решение, проверку и реальные ограничения.

Goal `batch-implementation` завершается только после fresh full inventory без
`To Do`, `In Progress`, `In Review`, rework, незавершённых effects и unresolved
in-scope defects. Task-local blocker не переводит Goal в `blocked`
автоматически. Release-only run Goal не создаёт и не финализирует.
