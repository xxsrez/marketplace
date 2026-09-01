# Классический

Применяет `IG-MODE-03` и mode-specific части `IG-MA-14..17` только когда
сохранённый `canonical_mode=classic`. Общие resolver, authority, recovery и
evidence rules бери из [Execution modes](../execution-modes.md); остальные
mode-файлы не читай.

Цель — полный терминальный результат с высокой уверенностью.

1. Controller/reviewer самостоятельно читает весь live scope и существенный
   source context, строит стратегию, dependency graph, acceptance и risk map.
2. Product, architecture, cross-issue, migration, concurrency, security,
   shared-state и иное weak-oracle judgment остаются этому profile.
3. Controller/reviewer остаётся основным исполнителем и делает почти всю
   implementation сам. Luna Max получает только действительно тривиальный
   strict-simple packet: bounded, self-contained, disjoint, objectively
   verifiable, low-risk, без material product/creative/architecture/cross-issue
   judgment и без зависимости от непредсказуемой среды или риска. Маленький diff
   сам по себе не simple.
4. Обычная delegation разделяет полезные независимые packets и read-only
   critique. Не создавай competing full implementations по умолчанию.
5. При ambiguity, contract/context conflict, unexpected tool/environment state,
   scope expansion или proof gap Luna сохраняет checkpoint/evidence и прекращает
   одинаковые corrective retry. Problem fallback получает controller/reviewer:
   он продолжает сам, уточняет contract либо формирует новый безопасный packet.
6. Если automatic Luna недоступна, используй controller profile; это само по
   себе не является blocker-ом. Явный пользовательский role profile не
   подменяй.
7. Один integration owner выполняет fan-in и aggregate checks.
8. Controller/reviewer читает exact integrated diff, material source и raw
   evidence, проводит final code/result review и только затем разрешает Done.
9. Продолжай до empty active scope либо принятого terminal blocker handoff.
