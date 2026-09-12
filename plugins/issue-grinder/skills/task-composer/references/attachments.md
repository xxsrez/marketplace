# Пользовательские attachments

Читай при наличии переданных файлов до первой mutation, включая draft planning.

Если пользователь создаёт Task по bug report и передал attachment, оцени его
уместность по содержанию, связи с самой конкретной создаваемой Task и пользе
исполнителю: материал должен помогать увидеть проявление, воспроизвести,
локализовать, понять релевантный контекст либо проверить исправление. Применяй
один критерий к screenshot, документу, логу, записи и любому другому файлу;
формат сам по себе не создаёт презумпцию уместности. Уместный материал сохрани
как native attachment самой конкретной создаваемой Task, для которой это
evidence. Не заменяй обязательный attachment пересказом, local path, base64
либо временной или protected URL и не пропускай его молча. Явно нерелевантный,
избыточный или нарушающий secret-safe boundary материал не добавляй и сообщи
его disposition.

## Bind и recovery

Для каждого обязательного attachment до create подтверди доступный native source
route. После существования target Task свяжи с ней verified file identity со
stable independent bind key и перечитай attachment metadata. Не начинай create,
если native transport заведомо недоступен. Если bind остановился после создания
Task, сохрани это как partial result с exact missing attachment и безопасным
условием продолжения, а не как success; unknown upload/bind сначала reconcile
через reads и не повторяй с новой identity.
