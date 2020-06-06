# Sample Application
## Prerequsites
Before starting, you need send to a metric to Newrelic.
Simply change the account id and the api key in the following command and execute it
```bash
docker run -d --restart unless-stopped   --name newrelic-statsd -h $(hostname) -e NR_ACCOUNT_ID=1234 -e NR_API_KEY=api_key -p 8125:8125/udp newrelic/nri-statsd:latest
```

## Send The query to NewRelic
```bash 
echo "test.k8s.num:1|g" | nc  -w 1 -u localhost 8125
```

## Define your query
In sample/externalmetric.yaml you can decide on your query

## Deploy
Deploy The smaple application
```bash
kubectl apply -f sample/
```

## Scale
1. Check the number of pods
```bash
kubectl get po
```
```
NAME                                 READY   STATUS    RESTARTS   AGE
sample-application-757df6d98-g6rkq   1/1     Running   0          5s
```
2. Increate your metric to 5
```bash 
echo "test.k8s.num:5|g" | nc  -w 1 -u localhost 8125
```
3. Wait a few seconds and see that The deployment scaled
```bash
kubectl get po
```
```bash
NAME                                 READY   STATUS    RESTARTS   AGE
sample-application-757df6d98-4722t   1/1     Running   0          20s
sample-application-757df6d98-5j779   1/1     Running   0          20s
sample-application-757df6d98-6jcrr   1/1     Running   0          3s
sample-application-757df6d98-g6rkq   1/1     Running   0          2m17s
sample-application-757df6d98-v8rqk   1/1     Running   0          20s
```