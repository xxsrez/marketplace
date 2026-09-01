# Баланс

Применяет `IG-MODE-04`, mode-specific части `IG-MODE-08..09` и
`IG-MA-15..17` только когда сохранённый `canonical_mode=balance`. Общие
resolver, authority, recovery и evidence rules бери из
[Execution modes](../execution-modes.md); остальные mode-файлы не читай.

Цель — terminal result с меньшим расходом controller/reviewer profile.

1. Controller/reviewer выполняет первоначальный full-scope analysis,
   декомпозицию, acceptance и risk classification. Это анализ задачи и рисков,
   а не разрешение самому проводить обычное repository research, писать код или
   прогонять основную verification вместо Luna.
2. Сразу оставляет себе только packets, где само решение требует material
   product/architecture judgment, integration decision или final review. Метка
   `material_judgment` должна называть конкретное необратимое либо существенно
   неоднозначное решение; ею нельзя обозначить обычную implementation ради
   обхода Luna. Если сложное решение отделимо от исполнения, передай Luna уже
   принятое решение и bounded implementation contract.
3. До первой source mutation или targeted test создай первую Luna Max wave.
   Luna выполняет основную массу implementation, repository research, test
   authoring, test execution, preliminary verification, critique и bounded
   rework для лёгких и средних packets. Каждый child получает явные
   `model="gpt-5.6-luna"`, `reasoning_effort="max"`, bounded `fork_turns` и
   зелёный routing receipt. Встроенные `critic`/`reviewer` не используй;
   semantic critics/test authors запускаются как `default`, `explorer` или
   `worker` с Luna profile.
4. Первый Luna handoff является gate режима: пока нет наблюдаемого Luna child
   profile и содержательного evidence по packet-у, controller не начинает
   ordinary implementation/research/test lane сам. Несовпадение profile закрывает
   wave как routing failure, а не молча превращает `Баланс` в Sol-конвейер.
5. Дополнительный candidate создавай только при реальной развилке, проваленном
   oracle или высокой ожидаемой ценности независимой дешёвой проверки.
6. Integration owner собирает содержательную batch и формирует review packet.
7. При существенной неопределённости, contract/context conflict, слабом oracle,
   unexpected tool/environment state, scope expansion, proof gap или проблеме
   за границами packet-а Luna сохраняет изменения, checks, evidence и причину
   остановки и возвращает fallback controller/reviewer-у без одинакового retry.
   Тот уточняет material decision или contract и возвращает bounded execution в
   новую Luna wave; сам продолжает только неделимый material-judgment packet.
8. Reviewer может вернуть bounded material rework в новую economical wave.
   После rework повтори exact-candidate gate.
9. Недоступная automatic Luna уменьшает capacity. Controller может выполнить
   только неделимое material judgment; обычную Luna-работу не переносит на
   дорогих children. Если mode-compatible serial Luna lane также недоступна,
   сохрани evidence и честно остановись либо попроси пользователя сменить режим,
   не называя такой прогон валидным `Балансом`.
10. Без final gate не выдавай непроверенный result за terminal; режим сам не
   переключай.
