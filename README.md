# NickR

Requirements
------------

  + Download [Go](http://golang.org/dl)

  + Setup and [test your Go installation](https://golang.org/doc/install#testing)

API
---

  + These are the available endpoints:

    URL | Method  | DESC.
    --- | --- | ---
    `/login/github` | `GET` |  login page
    `/users` | `GET` | List all users
    `/user/NAME` | `GET` | List a user
    `/user` | `POST`  | Create a user
    `/user` | `PUT` | Update a user
    `/user` | `DELETE` | Delete a user

    All requests (except `login`) require token in the `Authorization` header i.e. (`Authorization: Bearer token_value`)

    For each "non" `GET` request, send the following parameters in the request `BODY`,

    ```json
    {
      "name": "the users name",
      "nicknames": ["the users nicknames"]
    }
    ```



