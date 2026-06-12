"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import api from "@/lib/api";
import { useAuthStore } from "@/lib/store";

interface Summary {
  total_sessions: number;
  completed_sessions: number;
  average_score: number;
  best_score: number;
}

interface CategoryStat {
  category: string;
  total_sessions: number;
  average_score: number;
}

interface ScoreTrend {
  session_id: number;
  category: string;
  score: number;
  created_at: string;
}

interface Analytics {
  summary: Summary;
  category_stats: CategoryStat[];
  score_trends: ScoreTrend[];
}

export default function AnalyticsPage() {
  const router = useRouter();
  const { token } = useAuthStore();
  const [analytics, setAnalytics] = useState<Analytics | null>(null);
  const [loading, setLoading] = useState(true);

  async function fetchAnalytics() {
    try {
      const res = await api.get("/analytics");
      setAnalytics(res.data);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    if (!token) {
      router.push("/login");
      return;
    }
    fetchAnalytics();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [token]);

  const scoreColor = (score: number) => {
    if (score >= 80) return "text-green-600";
    if (score >= 60) return "text-yellow-600";
    return "text-red-600";
  };

  const scoreBarColor = (score: number) => {
    if (score >= 80) return "bg-green-500";
    if (score >= 60) return "bg-yellow-500";
    return "bg-red-500";
  };

  if (loading) return (
    <div className="min-h-screen flex items-center justify-center">
      <p className="text-gray-400">Memuat analytics...</p>
    </div>
  );

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm px-6 py-4 flex justify-between items-center">
        <h1 className="text-xl font-bold text-indigo-700">🎯 Interview Simulator</h1>
        <button
          onClick={() => router.push("/dashboard")}
          className="text-sm text-gray-500 hover:underline"
        >
          ← Dashboard
        </button>
      </nav>

      <div className="max-w-4xl mx-auto p-6">
        <h2 className="text-2xl font-bold text-gray-800 mb-6">📊 Analytics Dashboard</h2>

        {/* Summary Cards */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
          {[
            { label: "Total Sesi", value: analytics?.summary.total_sessions, suffix: "" },
            { label: "Sesi Selesai", value: analytics?.summary.completed_sessions, suffix: "" },
            { label: "Rata-rata Skor", value: analytics?.summary.average_score?.toFixed(1), suffix: "" },
            { label: "Skor Terbaik", value: analytics?.summary.best_score, suffix: "" },
          ].map((card) => (
            <div key={card.label} className="bg-white rounded-xl shadow p-4 text-center">
              <p className="text-gray-500 text-sm">{card.label}</p>
              <p className="text-3xl font-bold text-indigo-600 mt-1">{card.value ?? 0}</p>
            </div>
          ))}
        </div>

        {/* Category Stats */}
        <div className="bg-white rounded-xl shadow p-6 mb-6">
          <h3 className="font-semibold text-gray-800 mb-4">Performa per Kategori</h3>
          {analytics?.category_stats.length === 0 ? (
            <p className="text-gray-400 text-sm">Belum ada data.</p>
          ) : (
            <div className="space-y-4">
              {analytics?.category_stats.map((cat) => (
                <div key={cat.category}>
                  <div className="flex justify-between text-sm mb-1">
                    <span className="font-medium capitalize">{cat.category}</span>
                    <span className={scoreColor(cat.average_score)}>
                      {cat.average_score.toFixed(1)}/100
                    </span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-3">
                    <div
                      className={`h-3 rounded-full transition-all ${scoreBarColor(cat.average_score)}`}
                      style={{ width: `${cat.average_score}%` }}
                    />
                  </div>
                  <p className="text-xs text-gray-400 mt-1">{cat.total_sessions} sesi</p>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Score Trends */}
        <div className="bg-white rounded-xl shadow p-6">
          <h3 className="font-semibold text-gray-800 mb-4">Tren Skor (10 Sesi Terakhir)</h3>
          {analytics?.score_trends.length === 0 ? (
            <p className="text-gray-400 text-sm">Belum ada sesi yang selesai.</p>
          ) : (
            <div className="space-y-3">
              {analytics?.score_trends.map((trend, index) => (
                <div key={trend.session_id} className="flex items-center gap-3">
                  <span className="text-xs text-gray-400 w-4">{index + 1}</span>
                  <div className="flex-1 bg-gray-200 rounded-full h-4">
                    <div
                      className={`h-4 rounded-full transition-all ${scoreBarColor(trend.score)}`}
                      style={{ width: `${trend.score}%` }}
                    />
                  </div>
                  <span className={`text-sm font-semibold w-12 text-right ${scoreColor(trend.score)}`}>
                    {trend.score}
                  </span>
                  <span className="text-xs text-gray-400 w-20 capitalize">{trend.category}</span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}