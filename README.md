# grafana dashboard importer

Did you also make a change to your kubernetes-mixin dashboards? Want to reimport the entire folder to grafana quickly?
```bash
./grafana-dashboard-importer \
    -token="Bearer xxxx" \
    -url="https://grafana.xxxxx.com" \
    -dir="/home/bwilsonhunt/kubernetes-mixin/dashboards_out"
```

Its not pretty code, so don't judge. It got the job done in 5 mins =)