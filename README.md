
# RESTful api

This project is a server side developed in Golang programming language using Gin framework. The project interacts with PostgreSQL database, uses RESTful architecture, CRUD system, and asynchronous programming using channels.




## Technologies
- [gin framework](https://gin-gonic.com/)
- postgresql
- [Insomnia](https://insomnia.rest/)

### Implemented:
- CRUD
- http server
- asynchronous programming using channels


   


## API Reference
#### Get all albums

```http
  GET /albums
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `album`   | `JSON`   | Returns all albums in --pretty JSON|

#### Get item

```http
  GET /albums/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Id of item to fetch |

#### Post Slice of albums 
**implement's asynchronous programming using channels** 


```http
  POST /albums
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `album`      | `JSON` | **Required**. JSON albums array |

#### DELETE album

```http
  DELETE /albums/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Id of item to delete |

#### PATCH album

```http
  PATCH /albums/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Id of item to PATCH |




