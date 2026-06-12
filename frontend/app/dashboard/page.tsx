"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import api from "@/lib/api";
import { useAuthStore } from "@/lib/store";

interface Session {
  id: number;
  category: string;
  status: string;
  score: number;
  version: number;
  created_at: string;
}

export default function DashboardPage() {
  const router = useRouter();
  const { token, logout, user } = useAuthStore();
  const [sessions, setSessions] = useState<Session[]>([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState("frontend");

  useEffect(() => {
    if (!token) {
      router.push("/login");
      return;
    }
    fetchSessions();
  }, [token]);

  const fetchSessions = async () => {
    try {
      const res = await api.get("/sessions");
      setSessions(res.data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const createSession = async () => {
    setCreating(true);
    try {
      const res = await api.post("/sessions", { category: selectedCategory });
      router.push(`/interview/${res.data.id}`);
    } catch (err) {
      console.error(err);
    } finally {
      setCreating(false);
    }
  };

  const statusConfig: Record<string, { bg: string; text: string; icon: string }> = {
    in_progress: { bg: "from-yellow-400 to-orange-500", text: "text-yellow-700", icon: "⏳" },
    completed: { bg: "from-green-400 to-emerald-500", text: "text-green-700", icon: "✓" },
  };

  const categoryConfig: Record<string, { color: string; gradient: string; icon: string; description: string }> = {
    frontend: {
      color: "blue",
      gradient: "from-blue-500 to-cyan-600",
      icon: "⚛️",
      description: "React, Next.js, State Management, CSS, Performance"
    },
    backend: {
      color: "purple",
      gradient: "from-purple-500 to-pink-600",
      icon: "�",
      description: "API Design, Database, Authentication, Scalability"
    },
  };

  // Stats calculation
  const completedSessions = sessions.filter(s => s.status === "completed");
  const avgScore = completedSessions.length > 0
    ? Math.round(completedSessions.reduce((sum, s) => sum + s.score, 0) / completedSessions.length)
    : 0;
  const inProgressCount = sessions.filter(s => s.status === "in_progress").length;

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-purple-50">
      {/* Modern Navbar with Glass Effect */}
      <nav className="bg-white/80 backdrop-blur-md shadow-sm border-b border-indigo-100 px-6 py-4 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto flex justify-between items-center">
          <div className="flex items-center gap-4">
            <div className="w-12 h-12 bg-gradient-to-br from-indigo-600 to-purple-600 rounded-2xl flex items-center justify-center shadow-lg transform hover:scale-110 transition-transform">
              <span className="text-white text-2xl font-bold">🎯</span>
            </div>
            <div>
              <h1 className="text-xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
                Interview Simulator
              </h1>
              <p className="text-xs text-gray-500">Powered by AI - InaAI Competition</p>
            </div>
          </div>
          
          <div className="flex items-center gap-4">
            <button
              onClick={() => router.push("/analytics")}
              className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-indigo-600 hover:bg-indigo-50 rounded-xl transition-all"
            >
              <span>📊</span>
              <span>Analytics</span>
            </button>
            <button
              onClick={() => router.push("/settings")}
              className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 rounded-xl transition-all"
            >
              <span>⚙️</span>
              <span>Settings</span>
            </button>
            {/* Settings navigation button */}
            <div className="w-px h-6 bg-gray-300"></div>
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 bg-gradient-to-br from-indigo-400 to-purple-500 rounded-full flex items-center justify-center text-white font-bold text-sm">
                {user?.email?.charAt(0).toUpperCase() || "U"}
              </div>
              <button
                onClick={() => { logout(); router.push("/login"); }}
                className="text-sm text-red-500 hover:text-red-600 font-medium transition-colors"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto p-6 pb-12">
        {/* Hero Section with Stats */}
        <div className="mb-8">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h2 className="text-3xl font-bold text-gray-800 mb-2">
                Selamat Datang! 👋
              </h2>
              <p className="text-gray-600">Siap untuk latihan interview dengan AI?</p>
            </div>
          </div>

          {/* Stats Cards */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            {/* Total Sessions */}
            <div className="bg-white rounded-2xl shadow-lg border border-gray-100 p-6 transform hover:scale-105 transition-all hover:shadow-xl">
              <div className="flex items-center justify-between mb-4">
                <div className="w-12 h-12 bg-gradient-to-br from-blue-400 to-blue-600 rounded-xl flex items-center justify-center">
                  <span className="text-2xl">📝</span>
                </div>
                <span className="text-xs font-semibold text-gray-500 uppercase">Total</span>
              </div>
              <div className="text-3xl font-bold text-gray-800 mb-1">{sessions.length}</div>
              <p className="text-sm text-gray-600">Interview Sessions</p>
            </div>

            {/* Average Score */}
            <div className="bg-white rounded-2xl shadow-lg border border-gray-100 p-6 transform hover:scale-105 transition-all hover:shadow-xl">
              <div className="flex items-center justify-between mb-4">
                <div className="w-12 h-12 bg-gradient-to-br from-green-400 to-emerald-600 rounded-xl flex items-center justify-center">
                  <span className="text-2xl">⭐</span>
                </div>
                <span className="text-xs font-semibold text-gray-500 uppercase">Average</span>
              </div>
              <div className="text-3xl font-bold text-gray-800 mb-1">{avgScore}</div>
              <p className="text-sm text-gray-600">Score / 100</p>
            </div>

            {/* In Progress */}
            <div className="bg-white rounded-2xl shadow-lg border border-gray-100 p-6 transform hover:scale-105 transition-all hover:shadow-xl">
              <div className="flex items-center justify-between mb-4">
                <div className="w-12 h-12 bg-gradient-to-br from-orange-400 to-red-600 rounded-xl flex items-center justify-center">
                  <span className="text-2xl">⏳</span>
                </div>
                <span className="text-xs font-semibold text-gray-500 uppercase">Active</span>
              </div>
              <div className="text-3xl font-bold text-gray-800 mb-1">{inProgressCount}</div>
              <p className="text-sm text-gray-600">In Progress</p>
            </div>
          </div>
        </div>

        {/* Start New Interview Section */}
        <div className="relative overflow-hidden bg-gradient-to-br from-indigo-600 to-purple-600 rounded-3xl shadow-2xl p-8 mb-8">
          {/* Background Decorations */}
          <div className="absolute top-0 right-0 w-96 h-96 bg-white/10 rounded-full -mr-48 -mt-48"></div>
          <div className="absolute bottom-0 left-0 w-64 h-64 bg-white/10 rounded-full -ml-32 -mb-32"></div>
          
          <div className="relative z-10">
            {/* Testing Mode Banner */}
            <div className="inline-flex items-center gap-2 bg-yellow-400/20 backdrop-blur-sm px-4 py-2 rounded-full mb-6">
              <span className="text-xl">⚡</span>
              <span className="text-sm font-semibold text-yellow-100">
                Mode Testing: 1 Soal per Interview
              </span>
            </div>

            <h3 className="text-2xl font-bold text-white mb-3">
              Mulai Interview Baru
            </h3>
            <p className="text-indigo-100 mb-6 max-w-2xl">
              Pilih kategori interview dan mulai latihan dengan feedback AI yang detail dan personal
            </p>

            {/* Category Selection */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
              {Object.entries(categoryConfig).map(([key, config]) => (
                <button
                  key={key}
                  onClick={() => setSelectedCategory(key)}
                  className={`relative overflow-hidden group transition-all duration-300 ${
                    selectedCategory === key
                      ? 'scale-105 shadow-2xl'
                      : 'scale-100 hover:scale-105'
                  }`}
                >
                  <div className={`bg-white/95 backdrop-blur-sm rounded-2xl p-6 border-2 transition-all ${
                    selectedCategory === key
                      ? 'border-white shadow-xl'
                      : 'border-white/20 hover:border-white/40'
                  }`}>
                    {/* Selection Indicator */}
                    {selectedCategory === key && (
                      <div className="absolute top-3 right-3 w-6 h-6 bg-green-500 rounded-full flex items-center justify-center animate-scaleIn">
                        <span className="text-white text-sm">✓</span>
                      </div>
                    )}
                    
                    <div className={`w-14 h-14 bg-gradient-to-br ${config.gradient} rounded-2xl flex items-center justify-center mb-4 shadow-lg`}>
                      <span className="text-3xl">{config.icon}</span>
                    </div>
                    
                    <h4 className="text-lg font-bold text-gray-800 mb-2 capitalize">
                      {key}
                    </h4>
                    <p className="text-xs text-gray-600 leading-relaxed">
                      {config.description}
                    </p>
                  </div>
                </button>
              ))}
            </div>

            {/* Start Button */}
            <button
              onClick={createSession}
              disabled={creating}
              className="group relative px-8 py-4 bg-white text-indigo-600 rounded-2xl font-bold text-lg hover:bg-indigo-50 disabled:opacity-50 disabled:cursor-not-allowed shadow-2xl hover:shadow-3xl transition-all transform hover:scale-105 disabled:scale-100 flex items-center gap-3 overflow-hidden"
            >
              {/* Shimmer effect */}
              <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/40 to-transparent group-hover:animate-shimmer"></div>
              
              {creating ? (
                <>
                  <div className="w-6 h-6 border-3 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
                  <span>Mempersiapkan Interview...</span>
                </>
              ) : (
                <>
                  <span className="text-2xl">🚀</span>
                  <span>Mulai Interview Sekarang</span>
                  <span className="text-xl">→</span>
                </>
              )}
            </button>
          </div>
        </div>

        {/* Session History */}
        <div className="bg-white/80 backdrop-blur-sm rounded-3xl shadow-xl border border-gray-100 p-8">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h3 className="text-2xl font-bold text-gray-800 mb-1">Riwayat Interview</h3>
              <p className="text-sm text-gray-600">Lihat semua sesi interview Anda</p>
            </div>
            <div className="w-12 h-12 bg-gradient-to-br from-indigo-100 to-purple-100 rounded-xl flex items-center justify-center">
              <span className="text-2xl">📚</span>
            </div>
          </div>

          {loading && (
            <div className="flex flex-col items-center justify-center py-12">
              <div className="w-16 h-16 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
              <p className="text-gray-400">Memuat data...</p>
            </div>
          )}

          {!loading && sessions.length === 0 && (
            <div className="text-center py-16">
              <div className="w-24 h-24 bg-gradient-to-br from-gray-100 to-gray-200 rounded-full flex items-center justify-center mx-auto mb-6">
                <span className="text-5xl">📭</span>
              </div>
              <h4 className="text-xl font-bold text-gray-700 mb-2">Belum Ada Interview</h4>
              <p className="text-gray-500 mb-6">Mulai interview pertama Anda sekarang!</p>
              <button
                onClick={createSession}
                className="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-xl font-semibold hover:from-indigo-700 hover:to-purple-700 shadow-lg hover:shadow-xl transition-all transform hover:scale-105"
              >
                <span>🚀</span>
                <span>Mulai Sekarang</span>
              </button>
            </div>
          )}

          <div className="space-y-4">
            {sessions.map((session) => {
              const categoryData = categoryConfig[session.category];
              const statusData = statusConfig[session.status];
              
              return (
                <button
                  key={session.id}
                  onClick={() => router.push(`/session/${session.id}`)}
                  className="w-full group"
                >
                  <div className="bg-white border-2 border-gray-100 rounded-2xl p-6 hover:border-indigo-300 hover:shadow-xl transition-all transform hover:scale-[1.02]">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-4">
                        {/* Category Icon */}
                        <div className={`w-14 h-14 bg-gradient-to-br ${categoryData.gradient} rounded-xl flex items-center justify-center shadow-lg`}>
                          <span className="text-2xl">{categoryData.icon}</span>
                        </div>
                        
                        <div className="text-left">
                          <div className="flex items-center gap-3 mb-2">
                            <span className="px-3 py-1 bg-gradient-to-r from-indigo-100 to-purple-100 text-indigo-700 rounded-full text-xs font-bold uppercase">
                              {session.category}
                            </span>
                            
                            <div className={`flex items-center gap-1.5 px-3 py-1 bg-gradient-to-r ${statusData.bg} text-white rounded-full text-xs font-bold`}>
                              <span>{statusData.icon}</span>
                              <span className="uppercase">{session.status}</span>
                            </div>
                            
                            {session.status === "completed" && (
                              <div className="flex items-center gap-1.5 px-3 py-1 bg-gradient-to-r from-yellow-400 to-orange-500 text-white rounded-full text-xs font-bold">
                                <span>⭐</span>
                                <span>{session.score}/100</span>
                              </div>
                            )}
                          </div>
                          
                          <p className="text-sm text-gray-600">
                            {new Date(session.created_at).toLocaleDateString("id-ID", {
                              weekday: 'long',
                              year: 'numeric',
                              month: 'long',
                              day: 'numeric'
                            })}
                          </p>
                        </div>
                      </div>

                      {/* Arrow */}
                      <div className="text-gray-400 group-hover:text-indigo-600 group-hover:translate-x-2 transition-all">
                        <span className="text-2xl">→</span>
                      </div>
                    </div>
                  </div>
                </button>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}
