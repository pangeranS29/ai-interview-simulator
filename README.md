"# AI-Powered Interview Simulator

[![CI/CD Pipeline](https://github.com/pangeranS29/ai-interview-simulator/actions/workflows/ci.yml/badge.svg)](https://github.com/pangeranS29/ai-interview-simulator/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://go.dev/)
[![Next.js](https://img.shields.io/badge/Next.js-16-000000?logo=next.js)](https://nextjs.org/)

> Platform interview simulator berbasis AI dengan feedback real-time untuk persiapan karir Anda

## 🎯 Fitur Utama

- **Interview AI-Powered**: Latihan interview dengan 3 kategori (Behavioral, Technical, Situational)
- **Feedback Real-time**: Analisis jawaban dengan AI dan skor otomatis
- **Dashboard Analytics**: Tracking progress dan statistik interview
- **Modern UI/UX**: Design yang engaging dan responsive
- **Mode Testing**: 1 soal per session untuk hemat token

## 🏗️ Arsitektur

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Frontend  │────▶│ API Service │────▶│  PostgreSQL │
│  (Next.js)  │     │    (Go)     │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │    Redis    │
                    └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │Worker Service│
                    │   (AI API)   │
                    └─────────────┘
```

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.23+
- Node.js 20+

### Installation

1. Clone repository:
```bash
git clone https://github.com/pangeranS29/ai-interview-simulator.git
cd ai-interview-simulator
```

2. Setup environment variables:
```bash
cp .env.example .env
# Edit .env dengan konfigurasi Anda
```

3. Run dengan Docker Compose:
```bash
docker-compose up -d
```

4. Access aplikasi:
- Frontend: http://localhost:3000
- API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html

## 📦 Tech Stack

### Backend
- **Go 1.23** - API Service
- **Gin** - Web Framework
- **PostgreSQL** - Database
- **Redis** - Caching & Queue
- **Swagger** - API Documentation

### Frontend
- **Next.js 16** - React Framework
- **TypeScript** - Type Safety
- **Tailwind CSS** - Styling
- **Zustand** - State Management
- **Axios** - HTTP Client

## 🧪 Testing

```bash
# Backend tests
cd api-service
go test ./... -v

# Frontend tests
cd frontend
npm run test
```

## 📝 API Documentation

Swagger UI tersedia di: `http://localhost:8080/swagger/index.html`

### Main Endpoints:
- `POST /auth/register` - Register user baru
- `POST /auth/login` - Login user
- `POST /sessions` - Buat sesi interview
- `GET /sessions` - Get semua sesi
- `POST /sessions/:id/answers` - Submit jawaban
- `PUT /sessions/:id/finish` - Selesaikan interview
- `GET /analytics` - Dashboard analytics

## 🎨 Screenshots

### Dashboard
![Dashboard](docs/screenshots/dashboard.png)

### Interview
![Interview](docs/screenshots/interview.png)

### Results
![Results](docs/screenshots/results.png)

## 🏆 Kompetisi InaAI

Project ini dibuat untuk kompetisi InaAI dengan fokus pada:
- ✅ UX yang engaging dan modern
- ✅ AI integration untuk feedback
- ✅ Real-time progress tracking
- ✅ Optimasi performa dengan caching

## 📄 License

MIT License - lihat [LICENSE](LICENSE) untuk detail

## 👥 Contributors

- [pangeranS29](https://github.com/pangeranS29)

## 🙏 Acknowledgments

- InaAI untuk kompetisi dan dukungan
- OpenAI untuk AI API
- Community Go & Next.js

---

Made with ❤️ for InaAI Competition
" 
