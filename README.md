
# Сущности
 - `cluster`;
 - `feedback`;
 - `resident`;
 - `schedule`;
 - `service`;
 - `user`;

# Callback & View

# User

## Callback:
 - service_{id} - список услуг
 - servicedescribe_{id} - список под-услуг
 - cluster_{id} - отпрвляет резидента по `id` кластера
 - fio_{id} - отправляет резидента по его `id`
 - feedback - обратная связь `услуги`
 - main_menu - возвращение в главное меню
 - servicecreate - список под-услуг для получения их id и составления описания
 - servicelist - список доступных сервисов
 - resident - список всех резидентов
 - chat_gpt - запуск Chat GPT
 - schedule - отправляет последнее загруженное расписание из БД
 - request - создает форму обратной связи для оставления `заявка на вступление`
 - instruction - отправляет текст с описанием кнопок
 - pptx - отправляет текст с документом из БД
 - allcluster - отправляет все кластера
 - exit - завершает работу состояний `feedback` и `request`

## View 
 - stop_chat_gpt - останавливает работу Chat GPT
 - start - отправляет главное меню
 - resident_list - отправляет список всех резидентов

# Admin

## Callback

 - servicedelete - удаление услуги по кнопке
 - servicedescdelete - удаление вида услуги по кнопке
 - getcluster 
 - deletecluster
 - fiogetresident

## View
 
 - admin - список доступных админских команд
 - cancel - команда для отмены всех команд
 - create_resident - deprecated
 - create_resident_photo - создание резидента по фамилии, имени, фотографии и `id` кластера
 - notify - создать рассылку уведомлений всем пользователям бота
 - delete_resident - удаление резидента через нажитие на кнопку с резидентами
 - create_schedule - создание расписания (фотография)
 - create_cluster - создание нового кластера
```json
{
"cluster":"Введите название кластера"
}
```
 - delete_cluster - удаление кластера
 - get_feedback - просмотр всех оставленных сообщений пользователями (приходит в формате `JSON`)
 - delete_feedback - удаление сообщений пользователя по их `id`
```json
{
  "id": 1
}
```
 - update_pptx - обновление документа (поддерживает документы, которые поддерживаются в `tgbotapi.DocumentConfig{}` библиотеки [tgbotapi](https://github.com/go-telegram-bot-api/telegram-bot-api)). В обработчике находится в сущности `serivce`
 - create_service - добавление вида услуги через `JSON`
```json
{
  "service_name": "Название услуги"
}
```
- delete_service - удаление услуги
 - create_under_service - добавление вида под-услуг через `JSON`, фотографию, и кнопок с услугами
```json
{
  "under_service_name": "Название под-услуги",
  "describe": "Если имеется, то ввести описание услуги.\n\nСообщение будет выводиться через 2 строчки ниже."
}
```
- delete_under_service - удаление раздела услуги
- add_cluster_to_resident - назначет кластер резиденту