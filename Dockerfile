apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    resources:
      limits:
        cpu: "1"
      requests:
        cpu: "0.5"      
    ports:
    - containerPort: 80
