# 📊 So Sánh Benchmark: Gin vs wx Framework

## 🎯 Test Configuration
- **Tool**: k6 load testing
- **Scenario**: 100 VUs (Virtual Users) trong 5 giây
- **Test Type**: File upload endpoint
- **Total Requests**: 500 requests mỗi framework

---

## ⚡ Kết Quả So Sánh

### 1. **Response Time (Latency)** 🏆 **wx THẮNG**

| Metric | Gin Framework | wx Framework | Cải thiện |
|--------|---------------|--------------|-----------|
| **Average** | 10.3ms | **3.58ms** | **65% nhanh hơn** ⚡ |
| **Median** | 5.66ms | **2.23ms** | **61% nhanh hơn** ⚡ |
| **Max** | 64.19ms | **17.28ms** | **73% nhanh hơn** ⚡ |
| **p(90)** | 29.04ms | **9.03ms** | **69% nhanh hơn** ⚡ |
| **p(95)** | 31.26ms | **11.46ms** | **63% nhanh hơn** ⚡ |

**Kết luận**: wx framework có latency thấp hơn đáng kể, đặc biệt ở các percentile cao (p90, p95).

### 2. **Throughput (Requests per Second)**

| Metric | Gin Framework | wx Framework | So sánh |
|--------|---------------|--------------|---------|
| **Requests/sec** | 98.14 req/s | 98.82 req/s | Tương đương (~0.7% nhanh hơn) |

**Kết luận**: Throughput gần như tương đương trong test này.

### 3. **Success Rate**

| Framework | Checks Total | Success Rate | Failed |
|-----------|--------------|--------------|--------|
| **Gin** | 1000 checks | 100.00% | 0 |
| **wx** | 500 checks | 100.00% | 0 |

**Kết luận**: Cả hai đều đạt 100% success rate. (Gin có 2 checks, wx có 1 check)

### 4. **Network Usage**

| Metric | Gin Framework | wx Framework | So sánh |
|--------|---------------|--------------|---------|
| **Data Received** | 79 kB | 62 kB | wx nhận ít hơn 21% |
| **Data Sent** | 512 kB | 520 kB | wx gửi nhiều hơn 1.6% |

**Kết luận**: wx nhận ít data hơn (có thể do response nhỏ hơn), gửi tương đương.

### 5. **Error Rate**

| Framework | Failed Requests | Error Rate |
|-----------|-----------------|------------|
| **Gin** | 0 / 500 | 0.00% ✅ |
| **wx** | 0 / 500 | 0.00% ✅ |

**Kết luận**: Cả hai đều không có lỗi.

---

## 📈 Biểu Đồ So Sánh

### Response Time Distribution

```
Gin Framework:
  avg: 10.3ms  ████████████
  med: 5.66ms  ██████
  p90: 29.04ms ████████████████████████████
  p95: 31.26ms ██████████████████████████████

wx Framework:
  avg: 3.58ms  ████
  med: 2.23ms  ██
  p90: 9.03ms  █████████
  p95: 11.46ms ███████████
```

---

## 🎯 Tổng Kết

### ✅ Ưu điểm của wx Framework

1. **Latency thấp hơn 65%** - Đặc biệt quan trọng cho real-time applications
2. **Consistency tốt hơn** - p90 và p95 thấp hơn đáng kể
3. **Max latency thấp hơn 73%** - Ít outliers hơn
4. **100% success rate** - Ổn định và đáng tin cậy

### ⚖️ Điểm tương đương

1. **Throughput** - Gần như tương đương
2. **Error rate** - Cả hai đều 0%
3. **Stability** - Cả hai đều ổn định

### 📊 Kết Luận Cuối Cùng

**wx Framework vượt trội về performance**, đặc biệt là:
- **Response time nhanh hơn 65%** ở average
- **Latency ổn định hơn** ở các percentile cao
- **Max latency thấp hơn 73%**

Điều này cho thấy wx framework được tối ưu tốt cho:
- ✅ High-performance APIs
- ✅ Low-latency requirements
- ✅ Real-time applications
- ✅ Microservices với SLA nghiêm ngặt

---

## 🔍 Chi Tiết Kỹ Thuật

### Test Environment
- **Load**: 100 concurrent users
- **Duration**: 5 seconds
- **Total Requests**: 500 requests
- **Test Type**: File upload với multipart/form-data

### Performance Metrics Explained

- **Average**: Thời gian trung bình xử lý request
- **Median**: Giá trị ở giữa (50% requests nhanh hơn, 50% chậm hơn)
- **p(90)**: 90% requests nhanh hơn giá trị này
- **p(95)**: 95% requests nhanh hơn giá trị này
- **Max**: Request chậm nhất

---

*Benchmark được thực hiện với k6 load testing tool*

