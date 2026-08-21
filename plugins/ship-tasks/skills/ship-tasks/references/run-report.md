# Финальный отчёт ShipTask

Финальный ответ сообщает человеку состояние всего выбранного run. Он не
заменяет Task comments и не доказывает status write.

## Перед ответом

- Перечитать affected Tasks, comments, statuses, versions и relations.
- Перечитать Goal только в `batch-implementation` либо когда release продолжает
  уже активный совместимый Goal массовой имплементации.
- Проверить exact result identity и обязательные external effects.
- Сопоставить обещанный результат с фактическим.
- Убедиться, что больше нет безопасного in-scope пути, который агент считает
  достаточным для продолжения.

## Содержание

Начать с результата. Затем кратко сообщить:

- что выполнено и текущий Task/Goal status;
- что доказано и что осталось `not-available`;
- primary cause незавершённости, отдельно от последствий;
- что уже сделано и почему этого недостаточно для завершения;
- user action или внешнее условие, без которого продолжение невозможно;
- первый safe resume step.

Не писать process diary, полные URL/UUID/raw errors и внутренние названия
оркестрации, если они не нужны пользователю. Если run успешен, объяснить
полученный outcome, важное решение, проверку и реальные ограничения.

Goal `batch-implementation` завершается только после fresh full inventory без
`To Do`, `In Progress`, `In Review`, rework, незавершённых effects и unresolved
in-scope defects. Task-local blocker не переводит Goal в `blocked`
автоматически. Release-only run Goal не создаёт и не финализирует.
