# Репетиторы

### Получить всех репетиторов

`GET` http://localhost:8080/admin/tutors/

<details>
<summary><b>Пример ответа</b></summary>

```json
{
  "tutors": [
    {
      "id": 1,
      "full_name": "Нечепорк Максим Алексеевич",
      "tg": "https://t.me/maxim_jordan",
      "has_balance_negative": true,
      "has_only_trial": false,
      "has_newbie": false
    }
  ]
}
```

</details>

### Получить репетитора по ID

`GET` http://localhost:8080/admin/tutors/1


<details>
<summary><b>Пример ответа</b></summary>

```json
{
  "tutor": {
    "id": 1,
    "full_name": "Нечепорк Максим Алексеевич",
    "phone": "89826588317",
    "tg": "https://t.me/maxim_jordan",
    "cost_per_hour": "1,500.00",
    "subject_id": 1
  }
}
```

</details>

### Поиск по репетиторам

`GET` http://localhost:8080/admin/tutors/search?search=Нечепорук


<details>
<summary><b>Пример ответа</b></summary>

```json
{
  "tutors": [
    {
      "id": 1,
      "full_name": "Нечепорук Максим Алексеевич"
    },
    {
      "id": 2,
      "full_name": "Узянов Даниил Евгеньевич"
    }
  ]
}
```

</details>

### Создание репетитора

`POST` http://localhost:8080/admin/tutors


<details>
<summary><b>Пример запроса</b></summary>

```json
{
  "full_name": "Узянов Даниил Евгеньевич",
  "phone": "+7 995 677 8781",
  "tg": "https://t.me/danzelVash",
  "cost_per_hour": "1500",
  "subject_id": 1
}
```

</details>

Ответы: `200`, `500`

### Удаление репетитора

`DELETE` http://localhost:8080/admin/tutors/1

Ответы: `200`, `500`

### Получение информации по финансам репетитора

`POST` http://localhost:8080/admin/tutors/1/finance

<details>
<summary><b>Пример запроса</b></summary>

```json
{
  "from": "2023-01-01 00:00:00",
  "to": "2023-03-31 23:59:59"
}
```

</details>


<details>
<summary><b>Пример ответа</b></summary>

```json
{
  "data": {
    "conversion": 30,
    "count": 23,
    "amount": 23
  }
}
```

</details>

### Репетитор проводит демо урок

`POST` http://localhost:8080/admin/tutors/trial_lesson

<details>
<summary><b>Пример запроса</b></summary>

```json
{
  "student_id": 2
}
```

</details>

Ответы: `200`, `500`

### Репетитор проводит обычный урок

`POST` http://localhost:8080/admin/tutors/conduct_lesson

<details>
<summary><b>Пример запроса</b></summary>

```json
{
  "student_id": 2,
  "duration": 1
}
```

</details>

Ответы: `200`, `500`

# Студенты

### Получить всех студентов

`GET` http://localhost:8080/admin/students


<details>
<summary><b>Пример ответа</b></summary>

```json
{
  "students": [
    {
      "id": 1,
      "first_name": "Максим",
      "last_name": "Нечепорук",
      "middle_name": "Алексеевич",
      "tg": "https://t.me/maxim_jordan",
      "is_only_trial_finished": true,
      "is_balance_negative": false,
      "is_newbie": false
    }
  ]
}
```

</details>

### Получить студентов репетитора

`GET` http://localhost:8080/admin/students?tutor_id=1


<details>
<summary><b>Пример ответа</b></summary>

```json
{
  "students": [
    {
      "id": 1,
      "first_name": "Максим",
      "last_name": "Нечепорук",
      "middle_name": "Алексеевич",
      "tg": "https://t.me/maxim_jordan",
      "is_only_trial_finished": true,
      "is_balance_negative": false,
      "is_newbie": false
    }
  ]
}
```

</details>

### Получить информацию о студенте

`GET` http://localhost:8080/admin/students/1


<details>
<summary><b>Пример ответа</b></summary>

```json
{
  "student": {
    "id": 1,
    "first_name": "Максим",
    "last_name": "Нечепорук",
    "middle_name": "Алексеевич",
    "phone": "89826588317",
    "tg": "https://t.me/maxim_jordan",
    "cost_per_hour": "$1,500.00",
    "subject_id": 1,
    "tutor_id": 1,
    "parent_full_name": "Нечепорук Алексей Владимирович",
    "parent_phone": "89826588317",
    "parent_tg": "https://t.me/maxim_jordan",
    "balance": "$1,000.00",
    "has_buttons": false,
    "is_only_trial_finished": true,
    "is_balance_negative": false,
    "is_newbie": false
  }
}
```

</details>

### Поиск по студентам

`GET` http://localhost:8080/admin/students/search?search="Нечепорук Максим Алексеевич"


<details>
<summary><b>Пример ответа</b></summary>

```json
{
  "students": [
    {
      "id": 1,
      "first_name": "Максим",
      "last_name": "Нечепорук",
      "middle_name": "Алексеевич",
      "parent_full_name": "Нечепорук Алексей Владимирович"
    }
  ]
}
```

</details>

### Создание студента

`POST` http://localhost:8080/admin/students


<details>
<summary><b>Пример запроса</b></summary>

```json
{
  "first_name": "Узянов",
  "last_name": "Даниил",
  "middle_name": "Евгеньевич",
  "phone": "+7 995 677 8781",
  "tg": "https://t.me/danzelVash",
  "cost_per_hour": "1500",
  "subject_id": 1,
  "tutor_id": 1,
  "parent_full_name": "Узянов Даниил Евгеньевич",
  "parent_phone": "+7 995 677 8781",
  "parent_tg": "https://t.me/danzelVash"
}
```

</details>

Ответы: `200`, `500`

### Удаление студента

`DELETE` http://localhost:8080/admin/students/1

Ответы: `200`, `500`

### Получение информации по финансам студента

`POST` http://localhost:8080/admin/students/1/finance

<details>
<summary><b>Пример запроса</b></summary>

```json
{
  "from": "2023-01-01 00:00:00",
  "to": "2023-03-31 23:59:59"
}
```

</details>

<details>
<summary><b>Пример ответа</b></summary>

```json
{
  "data": {
    "count": 10,
    "amount": 23
  }
}
```

</details>



