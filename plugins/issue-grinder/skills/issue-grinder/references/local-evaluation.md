# Изолированная local evaluation

Применяется только при явном evaluation intent с локальным PLAN и запретом
Task Manager, Goal, сети и external effects. Это тот же runtime и те же
mode contracts; evaluation не разрешает live delivery.

Прочитай [Execution modes](execution-modes.md), ровно один выбранный там mode
и предоставленный local scope/PLAN. Вне Соло также прочитай связанный
multi-agent-routing; перед первым child прочитай
[multi-agent execution](multi-agent-execution.md) для admission.
Не загружай run/Goal, Task Manager, autonomy и publication references.
Routing, review и final gate обязательны. Сохраняй main profile receipt и
Luna-only root admission. Полноту чтения обеспечь доступным способом без
обрезанного вывода; уже прочитанные неизменные инструкции используй повторно.
