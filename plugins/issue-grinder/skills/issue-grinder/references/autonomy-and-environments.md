# Autonomy and environments

Применяет `IG-AUTO-01..05` только к явно вызванному `$issue-grinder` run.

## Режим

Явный marker сохраняется для того же незавершённого run через turns,
compaction и interruption. Terminal result, явная отмена или явная замена run
закрывают его. При сомнении в непрерывности не восстанавливай autonomy по
старому упоминанию skill.

Пока есть безопасное существенное действие текущего scope, выполняй его сам.
Не превращай выбор инструмента, обычный UAT effect, синтетические fixtures или
публичный UAT в permission loop.

## Environment boundary

- Production запрещён: no connection, reads, logs/data, mutation, deploy или
  smoke. Разрешено только локально прочитать metadata, чтобы распознать target
  как Production и отказаться.
- Default environment для необходимого effect — подтверждённый UAT. Staging и
  другая доказанно non-production среда допустимы.
- Target разрешай перед первым environment effect, а не как бессмысленный
  preflight для локальной работы. Неизвестный UAT — ранний context failure до
  environment mutation; Production не fallback.
- Task Manager writes остаются exact control-plane operations и не наследуют
  product Production classification.

## Исключительный security selector

Запрашивай решение только при действительно высоком риске: secrets,
access-policy expansion, необратимое изменение shared durable data, реальные
чувствительные данные, реальные внешние получатели, реальный платёж или
потенциально неограниченная стоимость. Сначала выбери synthetic/ephemeral
замену, если она доказывает acceptance.

Если prompt неизбежен и selector capability доступна, предложи `Да`, `Нет`,
`Да всегда`:

- `Да` разрешает только exact action;
- `Нет` запрещает action и эквивалентный обход, после чего ищется безопасный
  альтернативный путь;
- `Да всегда` сохраняет узкую категорию из action type, data class и
  environment в project memory и подтверждает read-back.

Не удалось подтвердить persistence — текущее `Да` сохраняет exact authority,
но постоянное разрешение не считается записанным. Ни один вариант не разрешает
Production и не расширяется на соседние категории или новый run.
