sh "echo 'Running Kubernetes Configuration'"

sh "kubectl apply -f ./Database.yaml"

sh "echo 'Database Manifest has been Applied Successfully! Applying Application Manifest'"

sh "kubectl apply -f ./Application.yaml"

sh "echo 'Application has been Deployed to Kubernetes Cluster, Go Check Results! \n
NOTE: If you are on local machine, you can Access Application via command: \n
`kubectl port-forward $(kubectl get pods --namespace=app-namespace)[0] 3000:3000 --namespace=app-namespace` 
and go to the Browser and Type `localhost:3000` \n 
Or If you are in the Production Cluster At the Cloud...'"