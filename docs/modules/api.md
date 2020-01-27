# HTTP API

### `GET /api/v1/alerts`

Get an alerts list

Response example:
```
[
    {
        "active": true,
        "name": "alert1",
        "script_name": "test1.lua"
    },
    {
        "active": false,
        "name": "alert2",
        "script_name": "test1.lua"
    }
]
```