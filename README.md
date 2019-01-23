#### rate-limit

##### Install
```bash
go install
```

##### Usage
```
Usage of rate-limit:
  -inflight int
        максимальное кол-во параллельно запущенных команд (default 1)
  -rate int
        максимальное кол-во запусков команды в секунду (default 1)
  <command...>: команда для запуска, {} в команде заменяется на строчку из stdin.
```