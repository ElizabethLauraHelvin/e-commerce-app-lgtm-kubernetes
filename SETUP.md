# 🚀 Setup & Deployment Guide

## 📋 Prerequisites

- Git
- Docker & Docker Hub account
- Go 1.21+
- Node.js 18+
- kubectl & Kubernetes cluster (optional, for K8s deployment)

---

## 1️⃣ Persiapan (Setup Initial)

### Clone Repository
```bash
git clone https://github.com/your-username/ecommerce-app-go-lgtm-fp.git
cd ecommerce-app-go-lgtm-fp
```

### Update Docker Hub Username
Ubah `your-dockerhub-username` di file berikut:
- `.github/workflows/build.yml` - semua reference ke Docker image
- `k8s/**/*.yaml` - semua image spec
- `README.md` - documentation examples

---

## 2️⃣ Local Development

### Run Backend Services

#### API Gateway
```bash
cd backend/api-gateway
go run main.go
# Server running on http://localhost:8080
```

#### Product Service (new terminal)
```bash
cd backend/services/product
go run main.go
# Server running on http://localhost:8080
```

#### Order Service (new terminal)
```bash
cd backend/services/order
go run main.go
```

#### User Service (new terminal)
```bash
cd backend/services/user
go run main.go
```

#### Payment Service (new terminal)
```bash
cd backend/services/payment
go run main.go
```

### Run Frontend

```bash
cd frontend
npm install
npm run dev
# Frontend running on http://localhost:5173
```

---

## 3️⃣ Docker Build & Push

### Login to Docker Hub
```bash
docker login
```

### Build All Images
```bash
# API Gateway
docker build -t your-dockerhub-username/ecommerce-api-gateway:latest ./backend/api-gateway
docker push your-dockerhub-username/ecommerce-api-gateway:latest

# Product Service
docker build -t your-dockerhub-username/ecommerce-product:latest ./backend/services/product
docker push your-dockerhub-username/ecommerce-product:latest

# Order Service
docker build -t your-dockerhub-username/ecommerce-order:latest ./backend/services/order
docker push your-dockerhub-username/ecommerce-order:latest

# User Service
docker build -t your-dockerhub-username/ecommerce-user:latest ./backend/services/user
docker push your-dockerhub-username/ecommerce-user:latest

# Payment Service
docker build -t your-dockerhub-username/ecommerce-payment:latest ./backend/services/payment
docker push your-dockerhub-username/ecommerce-payment:latest

# Frontend
docker build -t your-dockerhub-username/ecommerce-frontend:latest ./frontend
docker push your-dockerhub-username/ecommerce-frontend:latest
```

### Test Images Locally
```bash
# Run API Gateway
docker run -p 8080:8080 your-dockerhub-username/ecommerce-api-gateway:latest

# Test health check
curl http://localhost:8080/health
```

---

## 4️⃣ Kubernetes Deployment

### Prerequisites
```bash
# Check kubectl setup
kubectl cluster-info
kubectl get nodes
```

### Update Manifests
Replace `your-dockerhub-username` in all `k8s/**/*.yaml` files:

```bash
# Quick replace using sed
sed -i 's/your-dockerhub-username/YOUR_ACTUAL_USERNAME/g' k8s/**/*.yaml
```

### Deploy Services
```bash
# Apply all manifests
kubectl apply -f k8s/api-gateway/api-gateway.yaml
kubectl apply -f k8s/product-service/product.yaml
kubectl apply -f k8s/order-service/order.yaml
kubectl apply -f k8s/user-service/user.yaml
kubectl apply -f k8s/payment-service/payment.yaml
kubectl apply -f k8s/frontend/frontend.yaml
```

### Verify Deployment
```bash
# Check namespace
kubectl get namespace | grep ecommerce

# Check all resources
kubectl get all -n ecommerce

# Check specific deployments
kubectl get deployments -n ecommerce
kubectl get services -n ecommerce
kubectl get pods -n ecommerce

# Check pod logs
kubectl logs -f deployment/api-gateway -n ecommerce
```

### Access Application

#### Port Forward
```bash
# Terminal 1: API Gateway
kubectl port-forward svc/api-gateway 8080:8080 -n ecommerce

# Terminal 2: Frontend
kubectl port-forward svc/frontend 3000:3000 -n ecommerce

# Terminal 3: Product Service
kubectl port-forward svc/product-service 8081:8080 -n ecommerce
```

#### Get LoadBalancer IP
```bash
# If using cloud provider (AWS, GCP, Azure)
kubectl get svc api-gateway -n ecommerce
kubectl get svc frontend -n ecommerce

# External IP will appear in EXTERNAL-IP column after a moment
```

### Test Services
```bash
# Health checks
curl http://localhost:8080/health
curl http://localhost:8081/health

# API calls
curl http://localhost:8080/api/products
curl http://localhost:8080/api/orders
curl http://localhost:8080/api/users
```

---

## 5️⃣ GitHub Actions CI/CD

### Setup GitHub Secrets

Go to Repository Settings → Secrets and Variables → Actions

Add these secrets:
```
DOCKERHUB_USERNAME: your-dockerhub-username
DOCKERHUB_TOKEN: your-dockerhub-token
GITHUB_TOKEN: (automatically created)
```

### How CI/CD Works

1. Push code to `main` branch
2. GitHub Actions triggers automatically
3. Builds Docker images for changed services
4. Pushes images to Docker Hub
5. Uses cache layers for faster builds
6. Shows build summary

### Check Workflow Status
```
GitHub → Actions tab → Find latest workflow run
```

---

## 6️⃣ Troubleshooting

### Service Won't Start
```bash
# Check logs
kubectl logs deployment/product-service -n ecommerce

# Describe pod to see events
kubectl describe pod <pod-name> -n ecommerce
```

### Image Pull Errors
```bash
# Check image exists on Docker Hub
docker pull your-dockerhub-username/ecommerce-product:latest

# Check ImagePullSecrets if using private registry
kubectl get secret -n ecommerce
```

### Services Can't Communicate
```bash
# Test DNS resolution
kubectl exec -it <pod-name> -n ecommerce -- nslookup product-service.ecommerce.svc.cluster.local

# Check network policy
kubectl get networkpolicy -n ecommerce
```

### Port Already in Use (Local)
```bash
# Kill process on port
# macOS/Linux
lsof -i :8080
kill -9 <PID>

# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

---

## 7️⃣ Common Commands Reference

### Docker
```bash
# Build image
docker build -t name:tag .

# Run container
docker run -p 8080:8080 name:tag

# View images
docker images

# Remove image
docker rmi <image-id>

# Push to registry
docker push name:tag
```

### Kubernetes
```bash
# Deploy
kubectl apply -f manifest.yaml

# View resources
kubectl get pods -n ecommerce
kubectl get svc -n ecommerce
kubectl get deployments -n ecommerce

# Delete resources
kubectl delete pod <name> -n ecommerce
kubectl delete -f manifest.yaml

# Logs
kubectl logs <pod-name> -n ecommerce
kubectl logs -f deployment/api-gateway -n ecommerce

# Port forward
kubectl port-forward svc/api-gateway 8080:8080 -n ecommerce

# Exec into pod
kubectl exec -it <pod-name> -n ecommerce -- bash
```

### Go
```bash
# Run locally
go run main.go

# Build binary
go build -o binary-name

# Format code
go fmt ./...

# Run tests
go test ./...

# Check dependencies
go mod tidy
```

### Frontend
```bash
# Install dependencies
npm install

# Development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

---

## 📝 Git Workflow

### First Time Push
```bash
git add .
git commit -m "Initial commit: microservices setup"
git push origin main
```

### Regular Updates
```bash
# Make changes
git add <changed-files>
git commit -m "feat: add new feature"
git push origin main

# GitHub Actions will automatically build & push Docker images
```

### Create Feature Branch
```bash
git checkout -b feature/new-feature
# Make changes
git add .
git commit -m "Add new feature"
git push origin feature/new-feature
# Create Pull Request on GitHub
```

---

## 🎯 Next Steps

1. ✅ Update Docker Hub username
2. ✅ Build and push Docker images
3. ✅ Setup Kubernetes cluster
4. ✅ Deploy to Kubernetes
5. ✅ Configure CI/CD secrets
6. ✅ Test application
7. ✅ Monitor services
8. ✅ Scale services as needed

---

## 📞 Support

For issues or questions:
- Check logs: `kubectl logs <pod> -n ecommerce`
- Review GitHub Actions: Repository → Actions tab
- Check Docker Hub: Profile → Repositories

---

**Happy deploying! 🚀**
