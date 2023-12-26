
# User

## Callback:
 - feedback - обратная связь
 - service_{id} - список услуг
 - servicedescribe_{id} - список под-услуг
 - servicecreate - список под-услуг для получения их id и составления описания

## View 
 - servicelist - список доступных сервисов

# Admin

## Callback

## View

 - create_service - добавление вида услуги через JSON
```json
{
  "service_name": "Название услуги"
}
```
 - create_under_service - добавление вида под-услуг через JSON
```json
{
  "under_service_name": "Название под-услуги",
  "describe": "Если имеется, то ввести описание услуги"
}
```