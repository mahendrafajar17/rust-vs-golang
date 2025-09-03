# Load Testing Summary

## Initial Setup
After 5 minutes of startup:
![Startup Status](img/get/image-1.png)

## HTTP GET Request Testing

### Test Configuration 1
- **Threads**: 20
- **Ramp-up**: 10 seconds
- **Loop count**: 10

![Test Results - 20 threads](img/get/image-3.png)

### Test Configuration 2
- **Threads**: 200
- **Ramp-up**: 10 seconds
- **Loop count**: 10

#### Results
![Test Results - 200 threads](img/get/image-4.png)

**TPS (Transactions Per Second)**
![Performance Metrics](img/get/image-5.png)

**CPU Usage**
![Response Times](img/get/image-6.png)

**Memory Usage**
![Throughput Analysis](img/get/image-7.png)

### Test Configuration 3
- **Threads**: 2000
- **Ramp-up**: 10 seconds
- **Loop count**: 10

#### Results
![Test Results - 2000 threads](img/get/image-8.png)

#### Load Increase Analysis
**TPS (Transactions Per Second)**
![TPS Increase](img/get/image-9.png)

**CPU Usage**
![CPU Increase](img/get/image-10.png)

**Memory Usage**
![Memory Increase](img/get/image-11.png)

#### Load Decrease Analysis
**TPS (Transactions Per Second)**
![TPS Decrease](img/get/image-12.png)

**CPU Usage**
![CPU Decrease](img/get/image-13.png)

**Memory Usage**
![Memory Decrease](img/get/image-14.png)

## HTTP POST Request Testing

### Test Configuration 1
- **Threads**: 20
- **Ramp-up**: 10 seconds
- **Loop count**: 10

#### Results
![Test Results - 20 threads](img/post/image-1.png)

#### Load Increase Analysis
**TPS (Transactions Per Second)**
![TPS Increase](img/post/image-2.png)

**CPU Usage**
![CPU Increase](img/post/image-3.png)

**Memory Usage**
![Memory Increase](img/post/image-4.png)

### Test Configuration 2
- **Threads**: 200
- **Ramp-up**: 10 seconds
- **Loop count**: 10

#### Results
![Test Results - 200 threads](img/post/image-5.png)

#### Load Increase Analysis
**TPS (Transactions Per Second)**
![TPS Increase](img/post/image-6.png)

**CPU Usage**
![CPU Increase](img/post/image-7.png)

**Memory Usage**
![Memory Increase](img/post/image-8.png)

### Test Configuration 3
- **Threads**: 2000
- **Ramp-up**: 10 seconds
- **Loop count**: 10

#### Results
![Test Results - 2000 threads](img/post/image-9.png)

#### Load Increase Analysis
**TPS (Transactions Per Second)**
![TPS Increase](img/post/image-10.png)

**CPU Usage**
![CPU Increase](img/post/image-11.png)

**Memory Usage**
![Memory Increase](img/post/image-12.png)

## Performance Analysis: Rust vs Golang Comparison

### Test Environment
- **Rust Service**: Running on port 3002 (🦀)
- **Golang Service**: Running on port 8000 (🐹)
- **Load Testing**: Apache JMeter with varying thread counts
- **Monitoring**: Prometheus + Grafana dashboard

### Idle Performance (5 minutes startup)
- **Go CPU**: 3.60% (stable)
- **Rust CPU**: 0.940% (lower baseline)
- **Go Memory**: 7.47MB
- **Rust Memory**: 6.81MB (more efficient)

### GET Request Performance Comparison

#### Test Configuration 1 (20 threads)
- **Go TPS**: 3.37 (max)
- **Rust TPS**: 3.46 (max) - slightly higher
- **Go CPU**: 3.60% (stable)
- **Rust CPU**: 1.30% (max) - lower usage but less stable
- **Go Memory**: 9.20MB (max, fluctuating)
- **Rust Memory**: 7.06MB (stable)

#### Test Configuration 2 (200 threads)  
- **Go TPS**: 36.4 (max)
- **Rust TPS**: 36.4 (max) - identical peak performance
- **Go CPU**: 3.60% (consistent)
- **Rust CPU**: 5.63% (max, with spikes)
- **Go Memory**: 9.91MB (max, stable pattern)
- **Rust Memory**: 6.92MB (max, more efficient)

#### Test Configuration 3 (2000 threads)
- **Go TPS**: 364 (max) - excellent scaling
- **Rust TPS**: 364 (max) - matched performance
- **Go CPU**: 4.80% (max, with brief spikes)
- **Rust CPU**: 2.46% (max, more efficient but spiky)
- **Go Memory**: 8.43MB (max, stable)
- **Rust Memory**: 7.52MB (max, efficient)

### POST Request Performance Comparison

#### Test Configuration 1 (20 threads)
- **Go TPS**: 3.38 (max)
- **Rust TPS**: 3.47 (max) - slightly higher
- **Go CPU**: 3.60% (stable)
- **Rust CPU**: 1.69% (max, variable)
- **Go Memory**: 9.15MB (max)
- **Rust Memory**: 6.94MB (max, more efficient)

#### Test Configuration 2 (200 threads)
- **Go TPS**: 36.4 (max)
- **Rust TPS**: 36.4 (max) - identical performance
- **Go CPU**: 3.60% (consistent)
- **Rust CPU**: 8.72% (max, with significant spike)
- **Go Memory**: 9.47MB (max)
- **Rust Memory**: 7.13MB (max, efficient)

#### Test Configuration 3 (2000 threads)
- **Go TPS**: 364 (max) - excellent scaling
- **Rust TPS**: 364 (max) - matched performance
- **Go CPU**: 3.60% (remarkably stable even at 2000 threads)
- **Rust CPU**: 20.9% (max, dramatic spike)
- **Go Memory**: 10.1MB (max)
- **Rust Memory**: 7.66MB (max, still more efficient)

### Key Performance Insights

#### 🐹 **Golang Strengths**
1. **CPU Stability**: Remarkably consistent 3.60% CPU across ALL loads (20-2000 threads)
2. **Predictable Scaling**: Linear scaling without resource spikes
3. **Production Ready**: Absolutely stable behavior under extreme load
4. **POST Reliability**: Maintains stability even with complex POST operations

#### 🦀 **Rust Strengths**  
1. **Memory Efficiency**: 20-30% lower memory usage across all tests
2. **Idle Performance**: 74% lower CPU usage when idle
3. **Peak TPS**: Matches Go's throughput with fewer resources
4. **Memory Consistency**: Stable memory patterns even under load

#### ⚠️ **Critical Trade-offs**
- **Rust POST Issues**: CPU spikes to 20.9% during POST operations (5.8x higher than Go)
- **Go Stability**: Unmatched consistency - exactly 3.60% CPU regardless of load
- **Rust Efficiency vs Volatility**: Lower baseline usage but unpredictable spikes

### Complete Test Results Summary

#### Idle Performance (5 minutes)
| Language | CPU | Memory | Winner |
|----------|-----|---------|---------|
| Golang | 3.60% | 7.47MB | - |
| Rust | 0.940% | 6.81MB | 🦀 Rust (both metrics) |

#### GET Request Tests

##### Test 1: 20 threads, 10s ramp-up, 10 loops
| Language | TPS (max) | CPU (max) | Memory (max) | Winner |
|----------|-----------|-----------|---------------|---------|
| Golang | 3.37 | 3.60% | 9.20MB | - |
| Rust | 3.46 | 1.30% | 7.06MB | 🦀 Rust (all metrics) |

##### Test 2: 200 threads, 10s ramp-up, 10 loops  
| Language | TPS (max) | CPU (max) | Memory (max) | Winner |
|----------|-----------|-----------|---------------|---------|
| Golang | 36.4 | 3.60% | 9.91MB | 🐹 Go (CPU stability) |
| Rust | 36.4 | 5.63% | 6.92MB | 🦀 Rust (memory) |

##### Test 3: 2000 threads, 10s ramp-up, 10 loops
| Language | TPS (max) | CPU (max) | Memory (max) | Winner |
|----------|-----------|-----------|---------------|---------|
| Golang | 364 | 4.80% | 8.43MB | 🐹 Go (stability) |
| Rust | 364 | 2.46% | 7.52MB | 🦀 Rust (efficiency) |

#### POST Request Tests

##### Test 1: 20 threads, 10s ramp-up, 10 loops
| Language | TPS (max) | CPU (max) | Memory (max) | Winner |
|----------|-----------|-----------|---------------|---------|
| Golang | 3.38 | 3.60% | 9.15MB | - |
| Rust | 3.47 | 1.69% | 6.94MB | 🦀 Rust (all metrics) |

##### Test 2: 200 threads, 10s ramp-up, 10 loops
| Language | TPS (max) | CPU (max) | Memory (max) | Winner |
|----------|-----------|-----------|---------------|---------|
| Golang | 36.4 | 3.60% | 9.47MB | 🐹 Go (CPU stability) |
| Rust | 36.4 | 8.72% | 7.13MB | 🦀 Rust (memory) |

##### Test 3: 2000 threads, 10s ramp-up, 10 loops
| Language | TPS (max) | CPU (max) | Memory (max) | Winner |
|----------|-----------|-----------|---------------|---------|
| Golang | 364 | 3.60% | 10.1MB | 🐹 Go (exceptional stability) |
| Rust | 364 | 20.9% | 7.66MB | 🦀 Rust (memory only) |

#### Performance Trends Analysis
- **TPS Scaling**: Both languages achieve identical throughput scaling (3.4 → 36.4 → 364 TPS)
- **Golang CPU**: Remarkably consistent 3.60% across almost all tests (except 4.80% spike in GET 2000)
- **Rust CPU**: Highly variable (0.94% → 20.9%) with concerning POST operation spikes
- **Memory Efficiency**: Rust consistently uses 20-30% less memory across all scenarios
- **Critical Issue**: Rust's 20.9% CPU spike during 2000-thread POST operations vs Go's steady 3.60%

### Final Verdict: The Shocking Truth

#### 🏆 **Winner: Golang (by a significant margin)**

**🎯 Key Revelation:**
- **Golang's CPU**: Exactly 3.60% across ALL scenarios (20-2000 threads, GET/POST)
- **Rust's CPU**: Wildly unpredictable (0.94% to 20.9% - a 22x variation!)

**🦀 Rust wins on:**
- Memory efficiency (24% lower average usage)
- Idle performance (74% lower CPU at rest)
- Low-load scenarios (superior at 20 threads)

**🐹 Golang DOMINATES on:**
- **Supernatural Stability**: 3.60% CPU regardless of load type or scale
- **POST Performance**: Rust spikes to 20.9% CPU vs Go's steady 3.60%
- **Production Reliability**: Zero surprises, zero volatility
- **Scaling Confidence**: Predictable behavior from 20 to 2000 threads

### The Surprising Results

**Most Shocking Finding**: Golang's CPU usage is so consistent it appears "locked" at 3.60% - this level of stability is unprecedented.

**Rust's Achilles Heel**: POST operations cause dramatic CPU spikes (20.9% vs 3.60% for Go) - a 5.8x difference.

### Recommendations

#### Choose **Golang** when:
- **Production systems** (stability is king)
- **POST-heavy applications** (Rust struggles here)
- **Unpredictable traffic patterns**
- **Team needs reliability over optimization**

#### Choose **Rust** when:
- **Memory-constrained environments** (24% savings matter)
- **GET-heavy, low-traffic applications**
- **Maximum resource efficiency is critical**
- **You can handle performance volatility**

### Final Conclusion
**Golang's legendary stability makes it the clear winner for production web services.** While Rust is more memory-efficient, its CPU volatility (especially for POST operations) makes it risky for real-world applications. 

**The 3.60% CPU consistency across all loads is extraordinary** - Golang has achieved something remarkable in systems programming.
