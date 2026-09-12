# Epic и подзадачи

Читай только когда outcome требует Epic. Общие planning boundaries и live
refs остаются в [SKILL.md](../SKILL.md).

Epic сохраняет problem, beneficiary, Strategic Outcome, отличимые Human
Requirements/exact scope, Agent Plan, cross-cutting acceptance/non-goals и
целостный вклад подзадач. Перед его созданием сначала проверь active model. При
Astra (`gpt-6-astra`) Strategic Explainer не вызывай: problem-first description
формулирует сам coordinator в native mode. В остальных случаях перед формулировкой Epic прочитай
[publication contract](epic-publication.md). Native description проходит
тот же factual и coverage gate.

Каждая подзадача получает один конкретный результат, exact change boundary,
свой вклад в Strategic Outcome Epic, material technical details, применимые
Human Requirements/non-goals и Agent Plan с dependencies, acceptance criteria и
expected evidence.
Cross-cutting requirement остаётся в Epic и отражается в каждой применимой
подзадаче. Native parent link ведёт к полному strategic context, но одной ссылки
недостаточно: child description содержит компактную самодостаточную проекцию
вклада, применимых constraints/non-goals и качеств, которыми нельзя пожертвовать
ради локального упрощения. Исполнитель должен понять общий смысл, exact Task
boundary и планку качества без догадки; Epic context не расширяет scope child.
Material противоречие исправь до write. Для secrets указывай только имя
credential/secret store и target, никогда значение.
