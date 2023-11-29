## Create table
```
migrate create -seq -ext=.sql -dir=./migrations create_user_table
```

Or

```
make migration_new name=create_user_table
```