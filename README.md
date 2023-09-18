# testTaskEM

## Сервис работы с данными пользователя
Сервис получает ФИО, обогащает его данными из открытых апи и сохраняет в базу данных, после чего данными мощно кправлять описаными ниже методами.

### Доступные методы:
 - /user
 Метод: GET
 Возвращает всех пользователей
 Ответ: список или ошибка
 *Пример*:
 http://localhost:8080/user
 *Ответ*:
 ```
[
    {
        "ID": 1,
        "Name": "Dmitriy",
        "Surname": "Ushakov",
        "Patronymic": "Vasilevich",
        "Age": 42,
        "Gender": "male",
        "Country": "UA"
    },
    {
        "ID": 2,
        "Name": "Anton",
        "Surname": "Ivanov",
        "Patronymic": "Vitalyevich",
        "Age": 56,
        "Gender": "male",
        "Country": "RU"
    }
] 
```


 - /user
 Метод: POST
 Добавляет пользователя
 Принимает два обязательных поля ("Name" и "Surname") и одно не обязательное ("Patronymic")
 Ответ: пользователь или ошибка
 *Пример*:
 http://localhost:8080/user?name=Anton&surname=Ivanov&patronymic=Vitalyevich
 *Ответ*:
 ```
[
    {
        "ID": 3,
        "Name": "Anna",
        "Surname": "Dmitrieva",
        "Patronymic": "Pavlovna",
        "Age": 51,
        "Gender": "female",
        "Country": "PL"
    }
] 
```


 - /user
 Метод: DELETE
 Удаление пользователя. 
 Принимает id в качестве path параметра
 Ответ: Success или ошибка.


- /pagination
 Метод: GET
 Фильтрация и пагинация пользователей 
 Ответ: список или ошибка.
 *Пример*:
 http://localhost:8080/pagination?name=i&perPage=10&page=1
 *Ответ*:
``` 
"Users": [
        {
            "ID": 1,
            "Name": "Dmitriy",
            "Surname": "Ushakov",
            "Patronymic": "Vasilevich",
            "Age": 42,
            "Gender": "male",
            "Country": "UA"
        }
    ],
    "RespInfo": {
        "Total": 1,
        "PageCount": 1
    }
``` 


- /user
 Метод: PUT
 Изменение пользователя
 Принимает список изменяемых параметров в формате JSON
 Ответ: пользователь или ошибка
 *Пример*:
 http://localhost:8080/user
 Body:
 {
    "ID": 1,
    "Name": "Oleg"
}
 *Ответ*:
 ```
[
    {
        "ID": 1,
        "Name": "Oleg",
        "Surname": "Ushakov",
        "Patronymic": "Vasilevich",
        "Age": 42,
        "Gender": "male",
        "Country": "UA"
    }
] 
```