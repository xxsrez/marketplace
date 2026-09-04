# Баланс

Применяет `IG-MODE-04`, `IG-MODE-13..18` и mode-specific части
`IG-MA-01..19` только когда сохранённый `canonical_mode=balance`. Общие
resolver, authority, recovery и evidence rules бери из
[Execution modes](../execution-modes.md); остальные mode-файлы не читай.

Цель — получить терминальный результат сложного scope быстрее `Соло` за счёт
ограниченной параллельной помощи Luna, сохранив основную модель владельцем
архитектуры, интеграции и итоговой приёмки. Не оптимизируй число agents или долю
Luna саму по себе.

## Main-owned план и admission

1. Основной профиль одним full-scope pass изучает current state, requirements,
   dependency graph, write surfaces, acceptance, risks и integration points. Он
   выбирает общий способ решения и оставляет за собой downstream integration,
   CLI, общий harness, cross-cutting risk либо другой связующий пакет, который
   может выполняться по уже стабильным interfaces.
2. До dispatch допусти пакет Luna только когда одновременно доказаны:
   dependency readiness, понятный outcome, self-contained inputs, отдельная
   непересекающаяся write surface, стабильный interface, изолированный candidate,
   локальный oracle и отсутствие общего external/irreversible effect.
3. Параллельная wave требует минимум два подходящих пакета. Запусти не больше
   трёх Luna workers и одну одновременно активную execution wave. Число Tasks,
   файлов или свободных slots само по себе не создаёт пакет и не оправдывает
   делегацию. Но учитывай устранимое последовательное время tool-bound работы:
   два независимых пакета с долгими build/test/analyzer gates могут окупить
   dispatch даже при небольшом diff.
4. Если подходящей ширины нет, основной профиль продолжает scope последовательно.
   Это нормальный `Баланс`, а не повод искусственно дробить работу или менять
   режим.
5. Полностью заверши admission до первой source mutation и первого долгого
   packet-local gate. Если минимум два пакета прошли admission, dispatch всей
   волны обязателен до собственной реализации main profile. Не начинай один
   пригодный пакет последовательно и не откладывай Luna до следующей frontier.

## Одна активная Luna High wave

Сокращай время до dispatch: после admission не проектируй подробно реализацию
Luna-пакетов. Собери исходный inventory, подготовку изоляции и routing checks
независимых packets одним tool batch, когда это допускает среда. Используй
проверенный локальный helper, если он сохраняет ownership и candidate identity.
Не заменяй неизвестный результат проверки предположением ради быстрого запуска.

Каждый worker получает отдельный admitted worktree либо отдельный task-owned
shadow tree root и один self-contained packet:

```text
packet_id | purpose | exact base/root | owned surfaces | forbidden surfaces
source facts | architecture constraints | expected outcome | quick check
handoff: candidate identity | changed surfaces | checks | findings | unknowns
```

Default worker profile — exact `model="gpt-5.6-luna"`,
`reasoning_effort="high"` и bounded `fork_turns`. Перед dispatch обязателен
зелёный packet-bound routing receipt. Явный пользовательский profile override
сохраняется и проходит тот же guard.

Luna выполняет относящееся к пакету research, implementation и быструю локальную
проверку. Она не получает роли manager, supervisor, critic, reviewer, reducer,
integration owner или lifecycle owner; не создаёт descendants и не передаёт
packet дальше. Она не читает Issue Grinder `SKILL.md`, соседние mode references,
history или sibling modules, если нужные facts уже материализованы в packet-е.
После последнего quick check сразу возвращает compact handoff. Долгий процесс
нельзя оставить без доступного owner-а и результата.

Все packets текущей wave dispatch-ятся в одном окне. Основной профиль сразу
продолжает независимую полезную работу: выполняет critical/integration packet,
готовит общий test harness, исследует cross-cutting risk либо готовит fan-in. Он
не повторяет активную Luna work и не создаёт фиктивную занятость.

Если текущий checkout — отдельный clean task-owned run, Luna writers работают в
разных shadow roots общей exact base и их owned surfaces не пересекаются, main
profile владеет integration checkout и пишет свой downstream packet прямо в
него. До dispatch зафиксируй exact base, исходный status и main-owned surfaces;
после каждого handoff докажи, что observed integration diff ограничен этими
surfaces. Не создавай `.lanes/main` или другой main shadow и не копируй затем
весь main candidate обратно. Luna по-прежнему не пишет в integration checkout.
Для dirty пользовательского checkout, Git-worktree writers или пересекающихся
surfaces действует общий read-only integration rule.

После завершения собственной независимой работы ожидай всех оставшихся workers
одним collective event-driven wait до результата, внимания или общего stage
deadline. Не используй polling, повторные status lists, пустые nudges,
произвольный десятиминутный timeout или промежуточную synthesis неизменившегося
состояния. Технический timeout раньше deadline продолжает то же ожидание.

## Fan-in, tools и итоговая приёмка

После quiescence wave основной профиль последовательно:

1. сверяет candidate identity и owned surfaces каждого handoff;
2. механически переносит task-owned commits/bytes в integration candidate;
3. разрешает реальные conflicts и сохраняет пригодную часть неполного handoff;
4. не реализует заново успешно выполненный packet;
5. завершает сам недостающую часть дефектного packet-а без штатной цепочки
   replacement agents;
6. запускает независимые долгие build/test/analyzer gates точной интегрированной
   версии ровно одним parallel tool batch и сохраняет handles/results каждого
   процесса без повторного запуска;
7. перечитывает requirements и exact integrated diff одним итоговым проходом;
8. проверяет архитектурные и межпакетные решения, raw check results, known risks
   и remaining unknowns, исправляет найденное и принимает точный candidate.

Полный patch не пересказывай в model context: fan-in использует commit, copy,
sync либо автоматически построенный diff с отдельной identity-сверкой.
Identity, changed surfaces и diffs всех завершённых packets собирай одним
batch; изучи diff перед принятием и один раз сверь bytes после переноса.
Повторное полное чтение тех же файлов до и после копирования нужно только при
конфликте, дефекте или нехватке контекста; успешная identity-сверка не требует
нового содержательного обзора неизменившейся версии.
Package-local quick checks принадлежат Luna; общие, длительные и integration
checks — основной lane. Параллельность инструментов не создаёт новых agents.
Если exact-diff review потребовал material rework, повтори только затронутые
gates и минимальный общий acceptance. Не запускай полный уже зелёный batch снова
без утраты применимости evidence (cross-cutting изменение, среда/зависимости
или неизвестное влияние). Механический fan-in сам по себе не требует повтора.

Отдельный independent reviewer не является штатной ролью `Баланса`. Добавляй
его только по явному требованию пользователя, project policy или exact scope.
Даже тогда main profile сохраняет final acceptance и режим не превращается в
Manager Loop.

После fan-in и проверки версии новая независимая frontier может открыть
следующую волну после того же admission. Receipts и batches считай по волнам;
неполный прежний handoff завершает main без replacement workers. Terminal result требует exact
integrated checks и final acceptance main profile; нетерминальный checkpoint
`Экономичного` режима недоступен.

## Evidence и отклонения

Mode evidence содержит admission reasons, подтверждение admission до первой
source mutation, учёт tool-bound critical path, dependency/write-surface map,
routing receipts, worker profiles, dispatch window, peak Luna concurrency,
collective wait, candidate identities, ownership receipts, fan-in identity,
main-owned integration receipt, отсутствие main shadow/copy-back, число полных
integrated gate batches, targeted rework gates, integrated tool checks и final
acceptance.

Не являются доказательством `Баланса`:

- пересекающиеся writers либо запись worker-а в integration checkout;
- больше трёх Luna workers или две active waves;
- Luna manager/reviewer/reducer либо nested delegation;
- несколько конкурирующих реализаций одного packet-а или Best-of-N;
- повтор Luna work основным профилем без конкретного defect/conflict;
- последовательные dispatch/wait cycles там, где packets были одновременно
  готовы;
- source mutation или долгий packet-local gate до обязательного dispatch;
- отдельный main shadow или обратная копия main candidate в task-owned
  integration checkout при доступном безопасном main-owned path;
- новая wave до окончания и приёмки предыдущей либо без нового admission;
- повтор полного integrated gate batch без утраты применимости evidence;
- polling и coordination churn вместо одного collective wait;
- принятие worker self-report без exact integrated checks;
- отсутствие итогового exact-diff review основным профилем.

Локальный measurement noise или отклонение, не мешающее безопасно сохранить и
проверить candidate, записывается как nuance. Оно не обнуляет функциональный
результат, но topology deviation остаётся видимой в mode evidence.
