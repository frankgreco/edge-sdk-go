# edge-sdk-go

> golang sdk for ubiquiti edgeos

## Quickstart
```
client, err := edge.Login("https://192.168.1.1", true, "ubnt", "ubnt")
if err != nil {
    os.Exit(1)
}
ruleset, err := client.Firewall.GetRuleset(context.Background(), "NO_SSH")
if err != nil {
    os.Exit(1)
}
log.Println(ruleset)
```
