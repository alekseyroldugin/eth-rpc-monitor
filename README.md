# 🚀 Ethereum RPC Monitor

This project monitors multiple Ethereum RPC endpoints, measuring **latency** and **block height consistency**, and exposes metrics for **Prometheus & Grafana**.

---

## **📌 Features**

✅ Monitors multiple RPC endpoints from `config.yaml`\
✅ Measures **RPC response time & block number**\
✅ Exposes **Prometheus metrics** at `/metrics`\
✅ Supports **Grafana dashboards** for visualization\
✅ Runs in **Docker Compose** with **Prometheus & Grafana**

---

## **🛠️ Installation & Setup**

### **1️⃣ Clone the Repository**

```sh
git clone https://github.com/alekseyroldugin/eth-rpc-monitor
cd ethereum-rpc-monitor
```

### **2️⃣ Configure Endpoints**

Edit `config.yaml` to specify **Ethereum RPC endpoints**:

```yaml
rpc_endpoints:
  - name: "Ankr"
    url: "https://rpc.ankr.com/eth"
  - name: "LlamaRPC"
    url: "https://eth.llamarpc.com"
  - name: "Tenderly"
    url: "https://mainnet.gateway.tenderly.co"
  - name: "Blast"
    url: "https://eth.rpc.blxrbdn.com"
```

💡 **Initially only 1st endpoint is active and 3 other are commented**

---

### **3️⃣ Run the service Locally**

Ensure you have **Go** installed (version 1.18+):

```sh
go run main.go
```

🔍 Visit Prometheus metrics: [http://localhost:9090/metrics](http://localhost:9090/metrics)

---

### **4️⃣ Run Prometheus and Grafana with Docker Compose**

```sh
docker-compose up -d --build
```

This starts:

- 📡 **Prometheus** (Port `9091`)
- 📊 **Grafana** (Port `3000`)

---

### **5️⃣ Open Grafana Dashboard**

👉 **Grafana UI:** [http://localhost:3000](http://localhost:3000)\
✍️ **Default Login:**

- **Username:** `admin`
- **Password:** `admin`

---

## **📊 Prometheus & Grafana Setup**

### **1️⃣ Add Prometheus Data Source in Grafana**

1. Open Grafana → `Configuration` → `Data Sources`
2. Click `Add data source` → Select **Prometheus**
3. Set **URL**: `http://host.docker.internal:9091/`
4. Click `Save & Test`

### **2️⃣ Import Dashboard**

1. Go to `Dashboards` → `Import`
2. Upload `grafana-dashboard.json`
3. Click `Import`

🎉 Now you can visualize RPC latency and block updates! 🚀

---

## **🔄 Updating the Project**

### **1️⃣ Pull Latest Changes**

```sh
git pull origin main
```

### **2️⃣ Restart Docker Services**

```sh
docker-compose down && docker-compose up -d --build
```

---

## **❓ Troubleshooting**

### **Prometheus Shows No Data**

- **Check if Go service is running:**
  ```sh
  docker ps | grep eth-rpc-monitor
  ```
- **Manually test **``** endpoint:**
  ```sh
  curl http://localhost:9090/metrics
  ```
- **Restart Prometheus & Grafana:**
  ```sh
  docker-compose restart prometheus grafana
  ```

### **Can't Access Grafana?**

- Try **default credentials:** `admin / admin`
- If login fails, reset with:
  ```sh
  docker-compose down -v && docker-compose up -d --build
  ```

---

## **📜 License**

This project is open-source and licensed under **MIT**.
