# AI Usage Log - INaAI Competition 2026

**Peserta:** Pangeran Silaen  
**Track:** Full-stack Developer  
**Periode Development:** Juni 2026

---

## 🤖 Tools & Pattern Penggunaan

### AI Tools yang Digunakan
- **Kiro (powered by Claude Sonnet 4.5)** - Primary development assistant
- **Mode:** Agentic mode dengan autopilot untuk iterative development
- **Use Cases:**
  - Code generation & architecture design
  - Debugging & troubleshooting
  - Testing & TDD implementation
  - Deployment configuration
  - Documentation

### Pattern Penggunaan
1. **Agentic Mode (80%)** - AI mengeksekusi multi-step tasks secara autonomous
2. **Interactive Debugging (15%)** - Troubleshooting production issues
3. **Code Review & Refactoring (5%)** - Quality improvements

---

## 📝 History Prompt Utama (Fase 2)

### 1. **Initial Setup & Architecture**
```
"tolong perbaiki folder api-service saya banyak import yang error"
```
**Context:** Missing logger package menyebabkan compilation errors  
**AI Output:** Created `internal/logger/logger.go` with zerolog implementation  
**Decision:** ✅ **DITERIMA** - Structure dan implementation sudah sesuai dengan best practices Go logging

---

### 2. **Testing Mode Configuration**
```
"buatkan 1 soal ajah tiap tes, karna ini masih tes jadi buatkan 1 saja dulu ya"
```
**Context:** Menghemat token Gemini API selama development  
**AI Output:** Modified question limit dari 3 ke 1  
**Decision:** ✅ **DITERIMA** - Pragmatic approach untuk development testing

---

### 3. **Race Condition Bug Fix** ⭐ **MOST IMPACTFUL**
```
"tolong diperbaiki apakah errornya dibackend atau ui pada saat saya test dengan jawaban yang kompleks nilainya masuk ke database itu 92/100 tetapi ditampilan webnya menunjukkan score akhir 0"
```

**Context:** Score yang benar tersimpan di database (92/100) tapi UI menampilkan 0

**Root Cause Analysis (AI):**
- Worker service memproses feedback secara asynchronous (Redis Pub/Sub)
- User bisa klik "Selesai" sebelum worker selesai generate feedback
- Frontend fetch session data yang masih belum punya feedback → score = 0

**AI Solution:**
1. **Backend Fix:**
   - Tambah kolom `feedback_ready` di tabel `answers`
   - Handler `FinishSession` check apakah semua feedback ready
   - Return HTTP 425 (Too Early) jika belum siap
   
2. **Frontend Fix:**
   - Real-time tracking feedback status dengan polling 3 detik
   - Button "Selesai" disabled sampai semua feedback ready
   - Visual indicator bouncing dots untuk loading state
   - Replace alert() dengan toast notifications

**Decision: ✅ DITERIMA dengan MODIFIKASI**

**Alasan Terima:**
- ✅ Root cause analysis benar: classic race condition
- ✅ HTTP 425 (Too Early) adalah status code yang tepat semantically
- ✅ Polling mechanism simple dan reliable
- ✅ UX improvement dengan visual feedback clear

**Modifikasi yang Saya Lakukan:**
- Adjust polling interval dari 2s ke 3s (mengurangi load ke server)
- Tambah error handling untuk edge case network failure
- Refine animation CSS untuk smoother experience

**Why This Is Most Impactful:**
1. **Production Bug** - Data inconsistency adalah critical issue
2. **System Design Lesson** - Memahami trade-off asynchronous architecture
3. **User Experience** - Dari confusing (score 0) menjadi clear feedback state
4. **Proper HTTP Semantics** - Menggunakan status code 425 yang jarang dipakai tapi tepat

**Proof in Code:**
- Backend: `api-service/internal/handlers/session.go` (line ~180-200)
- Frontend: `frontend/app/interview/[id]/page.tsx` (line ~120-150)
- Git Commit: Dapat ditunjukkan di live defense

---

### 4. **UI/UX Redesign**
```
"perbaiki seluruh tampilan yang inovatif kreatif dan adapatif"
```
**Context:** Meningkatkan visual appeal untuk kompetisi  
**AI Output:** Modern design dengan gradient, glass morphism, animations  
**Decision:** ✅ **DITERIMA** - Sesuai dengan kriteria "innovative & creative" kompetisi

---

### 5. **TDD Implementation**
```
"nah saya ingin lolos kriteria penilaian yang diberikan tapi kita kerjakan langkah demi langkah ya"
```
**Context:** Must Have criteria - TDD dengan visible git history  
**AI Output:** 
- Created `change_password_test.go` FIRST (RED phase)
- Then implemented `ChangePassword()` function (GREEN phase)
- Separate commits untuk menunjukkan TDD workflow

**Decision:** ✅ **DITERIMA** - Clear TDD demonstration dengan git history sebagai bukti

**Git Commits Proof:**
1. `83f5136` - Test first (RED)
2. `a1479fa` - Implementation (GREEN)
3. `efa6bc0` - Regression tests

---

### 6. **Docker Containerization**
```
"loh kenapa ini tidak ada? CRITICAL: ❌ Docker/Containerization BELUM ADA!"
```
**Context:** Self-review menemukan gap - docker-compose.yml ada tapi Dockerfile tidak ada  
**AI Output:** Created 3 Dockerfiles (api-service, worker-service, frontend) dengan multi-stage build  
**Decision:** ✅ **DITERIMA dengan ADJUSTMENT**

**Alasan Terima:**
- ✅ Multi-stage build untuk optimasi image size
- ✅ Alpine base untuk minimal footprint
- ✅ Proper CGO_ENABLED=0 untuk static binary

**Adjustment yang Saya Lakukan:**
- Update Go version requirement (conflict resolution go.mod)
- Downgrade redis library version untuk compatibility
- Run `go mod tidy` untuk clean dependencies

**Why This Shows Good AI Judgment:**
- AI correctly identify gap antara "cloud deployment" vs "Docker containerization"
- Pragmatic approach: Railway/Vercel sudah jalan (Nixpacks), tapi kompetisi butuh Dockerfile
- Quick adaptation: dari identify gap sampai working solution dalam 1 session

---

### 7. **Production Deployment Issues**
```
"ada error CI/CD digithub saya"
```
**Context:** GitHub Actions timeout pulling large Docker images  
**AI Output:** Switch ke Alpine images, add caching, increase timeout  
**Decision:** ✅ **DITERIMA** - Semua suggestions valid untuk CI/CD optimization

---

### 8. **Worker Service Retry Mechanism**
```
"ini kenapa worker-service kita yang dirailway itu error... Failed to generate feedback: gemini API error: 503 UNAVAILABLE"
```
**Context:** Gemini API occasionally returns 503 during high demand  
**AI Output:** Implemented exponential backoff retry (3 attempts)  
**Decision:** ✅ **DITERIMA** - Standard practice untuk handling transient failures

---

### 9. **Database Seed Update**
```
"saya ingin ganti isi database yang railway supaya pertanyaan nya itu lebih ke pertanyaan yang sering ditanyakan untuk role frontend dan backend developer saja"
```
**Context:** Original questions terlalu generic (behavioral/technical/situational)  
**AI Output:** 16 realistic interview questions untuk frontend & backend roles  
**Decision:** ✅ **DITERIMA** - Questions relevant dan technically accurate

---

### 10. **Branch Cleanup**
```
"saya melihat di github nya branchnya ada banyak ini kalau tidak berguna bantu saya hapus ya"
```
**Context:** 9 stale branches (dependabot + testing branches)  
**AI Output:** Delete unused branches, remove dependabot.yml  
**Decision:** ✅ **DITERIMA** - Keep repository clean dan professional

---

## � Statistics

- **Total Interactions:** ~50+ across all development phases
- **Code Acceptance Rate:** ~85%
- **Code Rejection Rate:** ~5%
- **Code Modified Rate:** ~10%

### Breakdown by Category:
- **Architecture & Design:** 15%
- **Bug Fixes:** 30%
- **Feature Implementation:** 35%
- **Deployment & DevOps:** 15%
- **Documentation:** 5%

---

## 🎯 Key Learnings

### 1. **AI as Pair Programmer, Not Autopilot**
- AI excellent untuk identify root cause (race condition bug)
- Human judgment tetap crucial untuk production trade-offs (polling interval)

### 2. **Verification is Essential**
- Selalu test AI output di local environment
- Understand setiap line of code sebelum commit
- AI bisa salah (contoh: Go version conflict resolution butuh manual adjustment)

### 3. **Iterative Approach Works Best**
- Break complex problems into smaller prompts
- Review AI output per step, tidak langsung bulk accept
- Example: TDD implementation dilakukan step-by-step dengan clear git commits

### 4. **Domain Knowledge Matters**
- AI suggestion HTTP 425 (Too Early) - saya verify ini correct semantic untuk use case
- Redis Pub/Sub pattern - saya understand trade-offs asynchronous architecture
- Docker multi-stage build - saya verify image size optimization

---

## 🛡️ Production Readiness Checklist

- ✅ **Error Handling:** All critical paths have proper error handling
- ✅ **Logging:** Structured logging (zerolog) untuk production debugging
- ✅ **Testing:** TDD approach dengan visible git history
- ✅ **Security:** JWT authentication, bcrypt password hashing, input validation
- ✅ **Scalability:** Stateless API, Redis caching, async processing
- ✅ **Monitoring:** Structured logs ready untuk centralized logging (ELK/Loki)
- ✅ **Deployment:** Docker containerization + cloud deployment (Railway/Vercel)
- ✅ **Documentation:** README, API docs (Swagger), deployment guide

---

## 📌 Notes untuk Live Defense

**Saya siap explain setiap baris code yang di-generate AI karena:**

1. **Race Condition Fix** - Saya understand database transaction flow, HTTP status semantics, dan async processing patterns
2. **Docker Setup** - Saya understand multi-stage builds, image optimization, dan Go build flags
3. **TDD Implementation** - Saya follow RED-GREEN-REFACTOR cycle dengan conscious decisions
4. **Security** - Saya verify bcrypt rounds, JWT signing algorithms, dan CORS configuration

**AI membantu accelerate development, tapi architectural decisions dan production quality checks tetap saya yang decide.**

---

**Generated:** Juni 2026  
**Last Updated:** Juni 12, 2026
