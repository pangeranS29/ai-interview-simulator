"use client";

import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import api from "@/lib/api";

interface Feedback {
  id: number;
  score: number;
  strengths: string;
  weaknesses: string;
  suggestion: string;
}

interface Question {
  id: number;
  content: string;
  category: string;
  difficulty: string;
}

interface Answer {
  id: number;
  answer_text: string;
}

interface AnswerWithFeedback {
  answer: Answer;
  question: Question;
  feedback: Feedback | null;
}

interface Session {
  id: number;
  category: string;
  status: string;
  score: number;
}

interface SessionDetail {
  session: Session;
  answers: AnswerWithFeedback[];
}

export default function SessionDetailPage() {
  const router = useRouter();
  const params = useParams();
  const sessionId = params.id;

  const [detail, setDetail] = useState<SessionDetail | null>(null);
  const [loading, setLoading] = useState(true);

  async function fetchDetail() {
    try {
      const res = await api.get(`/sessions/${sessionId}`);
      setDetail(res.data);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    fetchDetail();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const scoreColor = (score: number) => {
    if (score >= 80) return "text-green-600";
    if (score >= 60) return "text-yellow-600";
    return "text-red-600";
  };

  if (loading) return (
    <div className="min-h-screen flex items-center justify-center">
      <p className="text-gray-400">Memuat hasil...</p>
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

      <div className="max-w-3xl mx-auto p-6">
        {/* Score Summary */}
        {detail?.session.status === "completed" && (
          <div className="bg-white rounded-xl shadow p-6 mb-6 text-center">
            <p className="text-gray-500 mb-2">Skor Akhir</p>
            <p className={`text-6xl font-bold ${scoreColor(detail.session.score)}`}>
              {detail.session.score}
            </p>
            <p className="text-gray-400 text-sm mt-1">dari 100</p>
            <span className="inline-block mt-3 px-3 py-1 bg-indigo-100 text-indigo-700 rounded-full text-sm font-medium">
              {detail.session.category} interview
            </span>
          </div>
        )}

        {/* Answers & Feedback */}
        <div className="space-y-6">
          {detail?.answers.map((awf, index) => (
            <div key={awf.answer.id} className="bg-white rounded-xl shadow p-6">
              <div className="flex gap-2 mb-3">
                <span className="text-xs px-2 py-1 bg-gray-100 text-gray-600 rounded-full">
                  Pertanyaan {index + 1}
                </span>
                <span className="text-xs px-2 py-1 bg-gray-100 text-gray-600 rounded-full">
                  {awf.question.difficulty}
                </span>
              </div>

              <p className="font-medium text-gray-800 mb-3">{awf.question.content}</p>

              <div className="bg-gray-50 rounded-lg p-3 mb-4">
                <p className="text-xs text-gray-500 mb-1">Jawaban kamu:</p>
                <p className="text-sm text-gray-700">{awf.answer.answer_text}</p>
              </div>

              {awf.feedback ? (
                <div className="border-l-4 border-indigo-500 pl-4">
                  <div className="flex justify-between items-center mb-3">
                    <p className="text-sm font-semibold text-indigo-700">🤖 AI Feedback</p>
                    <span className={`text-2xl font-bold ${scoreColor(awf.feedback.score)}`}>
                      {awf.feedback.score}/100
                    </span>
                  </div>
                  <div className="space-y-2">
                    <div className="bg-green-50 rounded p-2">
                      <p className="text-xs font-medium text-green-700 mb-1">✅ Kekuatan</p>
                      <p className="text-sm text-gray-700">{awf.feedback.strengths}</p>
                    </div>
                    <div className="bg-red-50 rounded p-2">
                      <p className="text-xs font-medium text-red-700 mb-1">⚠️ Kelemahan</p>
                      <p className="text-sm text-gray-700">{awf.feedback.weaknesses}</p>
                    </div>
                    <div className="bg-blue-50 rounded p-2">
                      <p className="text-xs font-medium text-blue-700 mb-1">💡 Saran</p>
                      <p className="text-sm text-gray-700">{awf.feedback.suggestion}</p>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="bg-yellow-50 rounded-lg p-3 text-sm text-yellow-700">
                  ⏳ AI sedang menganalisis jawaban...
                </div>
              )}
            </div>
          ))}
        </div>

        <div className="mt-6 flex gap-3">
          <button
            onClick={() => router.push("/dashboard")}
            className="flex-1 border border-indigo-600 text-indigo-600 py-2 rounded-lg hover:bg-indigo-50 font-medium"
          >
            Kembali ke Dashboard
          </button>
          <button
            onClick={() => router.push("/analytics")}
            className="flex-1 bg-indigo-600 text-white py-2 rounded-lg hover:bg-indigo-700 font-medium"
          >
            📊 Lihat Analytics
          </button>
        </div>
      </div>
    </div>
  );
}