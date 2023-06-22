# Replayer

## Setting new regions

Getting the region list:
```
REGIONS=$(fly platform regions | tail -n +2 | awk '{ print $1 }' | tr '\n' ' ' | awk '{$1=$1};1')
COUNT=$(wc -w <<< $REGIONS | awk '{$1=$1};1')
```

Deploying to all regions:
```
fly scale count $COUNT --max-per-region 1 --region $(echo $REGIONS | tr ' ' ,)
```

Possible gotchas:
- `dev` is not a true fly region, it is (who would guess) a development region. Remove it from the list, if present.