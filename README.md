
## Payload examples

Start read-write workload
```json
{
    "plan_name": "read-write",
    "insert_workload": {
        "scale_factor": 1
    },
    "update_workload": {
        "scale_factor": 1
    },
    "read_workload": {
        "scale_factor": 1
    }
}
```

Start write workload
```json
{
    "plan_name": "write-only",
    "insert_workload": {
        "scale_factor": 1
    },
    "update_workload": {
        "scale_factor": 1
    },
}
```

Start read workload
```json
{
    "plan_name": "read-only",
    "read_workload": {
        "scale_factor": 1
    }
}
```