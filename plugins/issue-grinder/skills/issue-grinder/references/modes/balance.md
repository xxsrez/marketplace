# Баланс

Применяет `IG-MODE-04`, mode-specific части `IG-MODE-08..09` и
`IG-MA-15..17` только когда сохранённый `canonical_mode=balance`. Общие
resolver, authority, recovery и evidence rules бери из
[Execution modes](../execution-modes.md); остальные mode-файлы не читай.

Цель — terminal result с меньшим расходом controller/reviewer profile.

1. Controller/reviewer выполняет первоначальный full-scope analysis,
   декомпозицию, acceptance и risk classification.
2. Сразу оставляет себе packets, где само решение требует material judgment.
   Если сложное решение отделимо от исполнения, передай economical worker-у
   уже принятое решение и bounded implementation contract.
3. Luna Max является предпочтительным исполнителем лёгких и средних bounded
   packets и выполняет основную массу implementation, research, tests и
   preliminary verification. Средний packet может требовать содержательной
   работы, но не material product/architecture judgment и должен иметь ясные
   границы и проверяемый contract. Независимые economical critics/test authors
   атакуют candidate до дорогого review.
4. Дополнительный candidate создавай только при реальной развилке, проваленном
   oracle или высокой ожидаемой ценности независимой дешёвой проверки.
5. Integration owner собирает содержательную batch и формирует review packet.
6. При существенной неопределённости, contract/context conflict, слабом oracle,
   unexpected tool/environment state, scope expansion, proof gap или проблеме
   за границами packet-а Luna сохраняет изменения, checks, evidence и причину
   остановки и возвращает fallback controller/reviewer-у без одинакового retry.
   Тот уточняет contract, продолжает сам либо создаёт новый безопасный packet.
7. Reviewer может вернуть bounded material rework в новую economical wave.
   После rework повтори exact-candidate gate.
8. Недоступная automatic Luna не блокирует, пока существует разрешённый
   mode-compatible economical fallback. Не превращай её отсутствие в
   неограниченный расход дефицитного profile и не подменяй явный role override.
9. Без final gate не выдавай непроверенный result за terminal; режим сам не
   переключай.
