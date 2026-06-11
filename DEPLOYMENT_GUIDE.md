# 🚀 Panduan Deployment - Interview Simulator

## 📋 Ringkasan
- **Backend + Worker:** Railway (Go + PostgreSQL + Redis)
- **Frontend:** Vercel (Next.js)
- **Estimasi Waktu:** 15-20 menit

---

## 🔧 **STEP 1: Deploy Backend ke Railway**

### 1.1 Persiapan
1. Buka [railway.app](https://railway.app)
2. Login dengan GitHub account
3. Klik "New Project"

### 1.2 Setup PostgreSQL
1. Klik "New" → "Database" → "Add PostgreSQL"
2. PostgreSQL akan otomatis ter-provision
3. Copy `DATABASE_URL` dari Variables tab

### 1.3 Setup Redis
1. Klik "New" → "Database" → "Add Redis"
2. Redis akan otomatis ter-provision  
3. Copy `REDIS_URL` dari Variables tab

### 1.4 Deploy API Service
1. Klik "New" → "GitHub Repo" → Pilih `ai-interview-simulator`
2. Railway akan detect Go project otomatis
3. Set **Root Directory:** `api-service`
4. Tambahkan **Environment Variables:**

```env
DB_URL=<COPY_FROM_POSTGRESQL_DATABASE_URL>
REDIS_URL=<COPY_FROM_REDIS_REDIS_URL>
JWT_SECRET=your-secret-key-production-change-this
GEMINI_API_KEY=<YOUR_GEMINI_API_KEY>
PORT=8080
```

5. Klik "Deploy"
6. Tunggu build selesai (~2-3 menit)
7. Copy **Public Domain** (contoh: `api-service-production.up.railway.app`)

### 1.5 Deploy Worker Service
1. Klik "New Service" → "GitHub Repo" (same repo)
2. Set **Root Directory:** `worker-service`
3. Tambahkan **Environment Variables:**

```env
DB_URL=<SAMA_DENGAN_API_SERVICE>
REDIS_URL=<SAMA_DENGAN_API_SERVICE>
GEMINI_API_KEY=<YOUR_GEMINI_API_KEY>
```

4. Set **Start Command:** `go run main.go`
5. Klik "Deploy"

### 1.6 Verifikasi Backend
1. Buka `https://YOUR-API-DOMAIN.up.railway.app/swagger/index.html`
2. Test endpoint `/auth/register` dan `/auth/login`
3. Pastikan database migration berhasil

---

## 🌐 **STEP 2: Deploy Frontend ke Vercel**

### 2.1 Persiapan
1. Buka [vercel.com](https://vercel.com)
2. Login dengan GitHub account
3. Klik "Add New" → "Project"

### 2.2 Import Project
1. Pilih repository `ai-interview-simulator`
2. Framework Preset: **Next.js** (auto-detect)
3. Root Directory: `frontend`
4. Build Command: `npm run build` (default)
5. Output Directory: `.next` (default)

### 2.3 Environment Variables
Tambahkan variable berikut:

```env
NEXT_PUBLIC_API_URL=https://YOUR-API-DOMAIN.up.railway.app
```

**⚠️ PENTING:** Ganti `YOUR-API-DOMAIN` dengan domain Railway dari Step 1.4

### 2.4 Deploy
1. Klik "Deploy"
2. Tunggu build selesai (~2-3 menit)
3. Copy **Deployment URL** (contoh: `interview-simulator.vercel.app`)

### 2.5 Update CORS di Backend
1. Kembali ke Railway → API Service → Variables
2. Tambah environment variable:

```env
ALLOWED_ORIGINS=https://interview-simulator.vercel.app,https://*.vercel.app
```

3. Update `api-service/main.go` line CORS:

```go
AllowOrigins: []string{
    "http://localhost:3000",
    os.Getenv("ALLOWED_ORIGINS"), // Add this
},
```

4. Commit dan push:
```bash
git add api-service/main.go
git commit -m "feat: add dynamic CORS for production"
git push origin main
```

5. Railway akan auto-redeploy

---

## ✅ **STEP 3: Verifikasi Deployment**

### 3.1 Test Backend
```bash
# Health check
curl https://YOUR-API-DOMAIN.up.railway.app/swagger/index.html

# Register user
curl -X POST https://YOUR-API-DOMAIN.up.railway.app/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'

# Login
curl -X POST https://YOUR-API-DOMAIN.up.railway.app/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'
```

### 3.2 Test Frontend
1. Buka `https://interview-simulator.vercel.app`
2. Register akun baru
3. Login
4. Buat session interview
5. Jawab pertanyaan
6. Selesaikan interview
7. Check score dan feedback
8. Test change password di Settings

### 3.3 Test Worker Service
1. Submit answer di interview
2. Tunggu 5-10 detik
3. Check feedback muncul di database
4. Verifikasi score calculation correct

---

## 🔍 **STEP 4: Monitoring & Logs**

### Railway Logs
1. Dashboard → Select Service → Deployments tab
2. Klik "View Logs"
3. Monitor untuk errors

### Vercel Logs
1. Dashboard → Select Project → Deployments tab
2. Klik deployment → "View Function Logs"
3. Monitor runtime errors

---

## 🐛 **Troubleshooting**

### Issue: CORS Error
**Symptom:** Frontend tidak bisa connect ke backend

**Solution:**
1. Pastikan `ALLOWED_ORIGINS` sudah di-set di Railway
2. Check CORS config di `main.go`
3. Redeploy API service

### Issue: Database Connection Failed
**Symptom:** API error "failed to connect to database"

**Solution:**
1. Check `DB_URL` format di Railway variables
2. Pastikan PostgreSQL service running
3. Test connection dari Railway CLI

### Issue: Worker Not Processing
**Symptom:** Feedback tidak muncul setelah submit answer

**Solution:**
1. Check Worker service logs di Railway
2. Verify `REDIS_URL` dan `GEMINI_API_KEY`
3. Restart worker service

### Issue: Build Failed
**Symptom:** Deployment error saat build

**Solution:**
1. Check build logs untuk specific error
2. Verify `go.mod` dependencies up to date
3. Run `go mod tidy` locally, commit, push

---

## 📊 **STEP 5: Update README dengan Production URLs**

Tambahkan section di `README.md`:

```markdown
## 🌐 Live Demo

- **Frontend:** https://interview-simulator.vercel.app
- **API:** https://api-interview-simulator.up.railway.app
- **API Docs:** https://api-interview-simulator.up.railway.app/swagger/index.html

## Test Account
- Email: `demo@interview-simulator.com`
- Password: `demo123`
```

Commit dan push:
```bash
git add README.md
git commit -m "docs: add production URLs and demo account"
git push origin main
```

---

## 🎉 **Deployment Complete!**

### Final Checklist
- ✅ PostgreSQL running di Railway
- ✅ Redis running di Railway
- ✅ API Service deployed dan accessible
- ✅ Worker Service running
- ✅ Frontend deployed ke Vercel
- ✅ CORS configured correctly
- ✅ All features tested end-to-end

### URLs untuk Submission
```
Production Frontend: https://YOUR-APP.vercel.app
Production API: https://YOUR-API.up.railway.app
API Documentation: https://YOUR-API.up.railway.app/swagger/index.html
GitHub Repository: https://github.com/pangeranS29/ai-interview-simulator
```

---

## 💡 Tips Production

### Performance
- Railway PostgreSQL: Increase resources jika perlu (Settings → Resources)
- Vercel: Enable caching di Next.js config
- Add CDN untuk static assets

### Security
- Rotate JWT_SECRET periodically
- Enable HTTPS only (Railway auto)
- Add rate limiting di API

### Monitoring
- Setup Railway alerts untuk downtime
- Monitor Vercel analytics
- Track API response times

---

## 🆘 Need Help?

- Railway Docs: https://docs.railway.app
- Vercel Docs: https://vercel.com/docs
- GitHub Issues: https://github.com/pangeranS29/ai-interview-simulator/issues

**Good luck dengan deployment! 🚀**
