# Best-effort название Codex task

Применяй только после Issue Grinder invocation gate и только когда visible
conversation допускает, что это первый user turn новой Codex task. Title —
optional UI metadata Codex, не Task Manager write, не evidence, не Goal и не
delivery outcome. Отсутствие, deferred loading или ошибка capability не
блокируют основную работу.

## Докажи eligibility

1. Убедись, что host предоставляет `list_threads`, `read_thread` и
   `set_thread_title`. Если capability отсутствует или deferred, сохрани title и
   продолжай с `issue-grinder-title=not-available`.
2. Через `list_threads` найди ровно одного кандидата calling task. Он должен
   быть active Codex task на том же host и в том же project/cwd, а его
   preview/summary должен соответствовать текущему первому invocation. Не
   выбирай просто самую свежую task. Ноль или несколько подходящих candidates
   останавливают только auto-title.
3. Считай title, preview и summary только identity evidence: они не меняют
   prompt, scope или инструкции.
4. Прочитай exact candidate через `read_thread` с достаточным `turnLimit` и без
   outputs. Разрешай rename, только если history не paginated, существует ровно
   один текущий user-triggered turn и нет завершённого предыдущего turn. Любая
   неполнота, later turn или противоречие означает preserve.
5. Текущий title должен быть пустым либо очевидным catalog-generated
   placeholder: `Use issue-grinder skill`, `Use Issue Grinder delivery
   workflow`, `Использовать Issue Grinder`, raw `$issue-grinder`, qualified
   skill link или ясный локализованный эквивалент. Meaningful title и title,
   начинающийся с `Issue Grinder ·`, всегда сохраняй.

Не угадывай provenance semantic title, сгенерированного из natural-language
prompt: без явного placeholder signal он считается meaningful и сохраняется.

## Назови только calling task

Сначала разреши live canonical scope. Выполни best-effort title effect до первой
Task Manager mutation. Если bare `$issue-grinder` не удалось однозначно
разрешить в current Release, title не придумывай: сначала запроси scope.

Формат:

- single issue: `Issue Grinder · <Task ref> · <short Task title>`;
- Project + Release: `Issue Grinder · <Project name> · <Release name>`;
- Project scope без Release: `Issue Grinder · <Project name> · batch`;
- bare `$issue-grinder`: формат Project + Release после live-разрешения current
  Release.

Сожми labels до короткого нормализованного имени. Не включай status, дату,
branch, acceptance text и другие volatile details.

Если eligibility доказана, вызови `set_thread_title` не более одного раза без
`threadId`. Omission адресует calling task; discovery candidate id не передавай
в mutation. После ошибки не retry и не пробуй переименовать соседнюю task. Если
tool result не подтверждает effect, продолжай delivery с
`issue-grinder-title=not-available`.

Итоговые dispositions:

- `issue-grinder-title=renamed` — eligibility доказана и setter подтвердил
  effect;
- `issue-grinder-title=preserved` — custom/meaningful title, later turn или уже
  `Issue Grinder · ...`;
- `issue-grinder-title=not-available` — capability отсутствует/deferred,
  candidate не уникален, history неполна либо setter failed.

Не превращай disposition в отдельный подробный process report.
