#!/bin/bash

kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml
kubectl apply -f mysql-pvc.yaml
kubectl apply -f mysql-deployment.yaml
kubectl apply -f mysql-service.yaml
kubectl apply -f app-deployment.yaml
kubectl apply -f app-service.yaml
kubectl apply -f app-hpa.yaml

echo "Aguardando pods ficarem prontos..."
kubectl wait --for=condition=ready pod -l app=mysql -n soat-fiap --timeout=120s
kubectl wait --for=condition=ready pod -l app=soat-fiap -n soat-fiap --timeout=120s

echo "Deploy concluído!"
echo ""
echo "Para acessar a aplicação:"
echo "kubectl port-forward service/soat-fiap-service 8080:8080 -n soat-fiap"
echo ""
echo "Para verificar status:"
echo "kubectl get pods -n soat-fiap"
echo "kubectl get hpa -n soat-fiap"