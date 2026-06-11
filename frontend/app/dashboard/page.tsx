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
  const { token, logout } = useAuthStore();
  const [sessions, setSessions] = useState<Session[]>([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState("behavioral");

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

  const statusColor: Record<string, string> = {
    in_progress: "bg-yellow-100 text-yellow-700",
    completed: "bg-green-100 text-green-700",
  };

  const categoryColor: Record<string, string> = {
    behavioral: "bg-blue-100 text-blue-700",
    technical: "bg-purple-100 text-purple-700",
    situational: "bg-orange-100 text-orange-700",
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Navbar */}
      <nav className="bg-white shadow-sm px-6 py-4 flex justify-between items-center">
        <h1 className="text-xl font-bold text-indigo-700">🎯 Interview Simulator</h1>
        <div className="flex gap-4">
          <button
            onClick={() => router.push("/analytics")}
            className="text-sm text-indigo-600 hover:underline font-medium"
          >
            📊 Analytics
          </button>
          <button
            onClick={() => { logout(); router.push("/login"); }}
            className="text-sm text-red-500 hover:underline"
          >
            Logout
          </button>
        </div>
      </nav>

      <div className="max-w-4xl mx-auto p-6">
        {/* Start Interview */}
        <div className="bg-white rounded-xl shadow p-6 mb-6">
          
          
          <h2 className="text-lg font-semibold mb-4">Mulai Sesi Interview</h2>
          <div className="flex gap-3 flex-wrap mb-4">
            {["behavioral", "technical", "situational"].map((cat) => (
              <button
                key={cat}
                onClick={() => setSelectedCategory(cat)}
                className={`px-4 py-2 rounded-lg text-sm font-medium border transition-all ${
                  selectedCategory === cat
                    ? "bg-indigo-600 text-white border-indigo-600"
                    : "bg-white text-gray-600 border-gray-300 hover:border-indigo-400"
                }`}
              >
                {cat.charAt(0).toUpperCase() + cat.slice(1)}
              </button>
            ))}
          </div>
          <button
            onClick={createSession}
            disabled={creating}
            className="bg-indigo-600 text-white px-6 py-2 rounded-lg hover:bg-indigo-700 disabled:opacity-50 font-medium"
          >
            {creating ? "Memulai..." : "🚀 Mulai Interview"}
          </button>
        </div>

        {/* Session History */}
        <div className="bg-white rounded-xl shadow p-6">
          <h2 className="text-lg font-semibold mb-4">Riwayat Sesi</h2>

          {loading && <p className="text-gray-400 text-center">Memuat...</p>}
          {!loading && sessions.length === 0 && (
            <p className="text-gray-400 text-center">Belum ada sesi. Mulai interview pertamamu!</p>
          )}

          <div className="space-y-3">
            {sessions.map((session) => (
              <div
                key={session.id}
                className="border rounded-lg p-4 flex justify-between items-center hover:bg-gray-50 cursor-pointer"
                onClick={() => router.push(`/session/${session.id}`)}
              >
                <div className="flex gap-2 items-center">
                  <span className={`text-xs px-2 py-1 rounded-full font-medium ${categoryColor[session.category] || "bg-gray-100 text-gray-600"}`}>
                    {session.category}
                  </span>
                  <span className={`text-xs px-2 py-1 rounded-full font-medium ${statusColor[session.status] || "bg-gray-100"}`}>
                    {session.status}
                  </span>
                  {session.status === "completed" && (
                    <span className="text-sm font-semibold text-gray-700">
                      Score: {session.score}/100
                    </span>
                  )}
                </div>
                <span className="text-xs text-gray-400">
                  {new Date(session.created_at).toLocaleDateString("id-ID")}
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>  
    </div>
  );
}