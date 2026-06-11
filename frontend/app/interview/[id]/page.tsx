"use client";

import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import api from "@/lib/api";

interface Question {
  id: number;
  category: string;
  content: string;
  difficulty: string;
}

interface Session {
  id: number;
  category: string;
  status: string;
  version: number;
}

export default function InterviewPage() {
  const router = useRouter();
  const params = useParams();
  const sessionId = params.id;

  const [session, setSession] = useState<Session | null>(null);
  const [questions, setQuestions] = useState<Question[]>([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [answer, setAnswer] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [submitted, setSubmitted] = useState<number[]>([]);
  const [finishing, setFinishing] = useState(false);
  const [feedbackStatus, setFeedbackStatus] = useState<{total: number, ready: number}>({total: 0, ready: 0});
  const [errorMessage, setErrorMessage] = useState<string>("");
  const [showSuccessAnimation, setShowSuccessAnimation] = useState(false);

  useEffect(() => {
    fetchSession();
  }, []);

  useEffect(() => {
    if (session) {
      fetchQuestions();
    }
  }, [session]);

  // Auto-check feedback status setiap 3 detik jika ada jawaban yang belum ada feedback
  useEffect(() => {
    if (submitted.length > 0 && feedbackStatus.ready < feedbackStatus.total) {
      const interval = setInterval(() => {
        fetchSession();
      }, 3000);
      return () => clearInterval(interval);
    }
  }, [submitted, feedbackStatus]);

  const fetchSession = async () => {
    try {
      const res = await api.get(`/sessions/${sessionId}`);
      setSession(res.data.session);
      
      // Update feedback status
      const answers = res.data.answers || [];
      const totalAnswers = answers.length;
      const readyFeedbacks = answers.filter((a: any) => a.feedback !== null).length;
      setFeedbackStatus({ total: totalAnswers, ready: readyFeedbacks });
    } catch (err) {
      console.error(err);
    }
  };

  const fetchQuestions = async () => {
    try {
      if (!session) return;
      // Fetch hanya 1 soal untuk testing
      const res = await api.get(`/questions?category=${session.category}`);
      // Backend sudah return 5 soal random, kita ambil 1 saja untuk testing
      const limitedQuestions = (res.data || []).slice(0, 1);
      setQuestions(limitedQuestions);
    } catch (err) {
      console.error(err);
    }
  };

  const submitAnswer = async () => {
    if (!answer.trim() || !questions[currentIndex]) return;
    setSubmitting(true);
    setErrorMessage("");
    try {
      await api.post(`/sessions/${sessionId}/answers`, {
        question_id: questions[currentIndex].id,
        answer_text: answer,
      });
      setSubmitted([...submitted, questions[currentIndex].id]);
      setAnswer("");
      
      // Show success animation
      setShowSuccessAnimation(true);
      setTimeout(() => setShowSuccessAnimation(false), 2000);
      
      if (currentIndex < questions.length - 1) {
        setCurrentIndex(currentIndex + 1);
      }
    } catch (err) {
      setErrorMessage("Gagal mengirim jawaban. Silakan coba lagi.");
      console.error(err);
    } finally {
      setSubmitting(false);
    }
  };

  const finishSession = async () => {
    if (!session) return;
    
    // Double check feedback status
    if (feedbackStatus.ready < feedbackStatus.total) {
      setErrorMessage("Mohon tunggu, AI masih menganalisis jawaban Anda...");
      return;
    }
    
    setFinishing(true);
    setErrorMessage("");
    try {
      await api.put(`/sessions/${sessionId}/finish`, {
        version: session.version,
      });
      // Success! Redirect to results
      router.push(`/session/${sessionId}`);
    } catch (err: any) {
      if (err.response?.status === 409) {
        setErrorMessage("Terjadi konflik data. Halaman akan di-refresh...");
        setTimeout(() => {
          fetchSession();
          setErrorMessage("");
        }, 2000);
      } else if (err.response?.status === 425) {
        setErrorMessage("AI masih menganalisis jawaban Anda. Mohon tunggu sebentar...");
      } else {
        setErrorMessage("Terjadi kesalahan. Silakan coba lagi.");
      }
    } finally {
      setFinishing(false);
    }
  };

  const currentQuestion = questions[currentIndex];
  const progress = questions.length > 0 ? ((currentIndex) / questions.length) * 100 : 0;

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-purple-50">
      {/* Modern Navbar with Glass Effect */}
      <nav className="bg-white/80 backdrop-blur-md shadow-sm border-b border-indigo-100 px-6 py-4 sticky top-0 z-50">
        <div className="max-w-4xl mx-auto flex justify-between items-center">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-gradient-to-br from-indigo-600 to-purple-600 rounded-xl flex items-center justify-center shadow-lg">
              <span className="text-white text-xl font-bold">🎯</span>
            </div>
            <div>
              <h1 className="text-lg font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
                Interview Simulator
              </h1>
              <p className="text-xs text-gray-500">Powered by AI</p>
            </div>
          </div>
          <button
            onClick={() => router.push("/dashboard")}
            className="flex items-center gap-2 text-sm text-gray-600 hover:text-indigo-600 transition-colors font-medium"
          >
            <span>←</span> Dashboard
          </button>
        </div>
      </nav>

      <div className="max-w-4xl mx-auto p-6 pb-12">
        {/* Success Animation Modal */}
        {showSuccessAnimation && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/20 backdrop-blur-sm animate-fadeIn">
            <div className="bg-white rounded-2xl p-8 shadow-2xl animate-scaleIn">
              <div className="text-center">
                <div className="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4 animate-bounce">
                  <span className="text-4xl">✓</span>
                </div>
                <h3 className="text-xl font-bold text-gray-800 mb-2">Jawaban Terkirim!</h3>
                <p className="text-sm text-gray-600">AI sedang menganalisis jawaban Anda...</p>
              </div>
            </div>
          </div>
        )}

        {/* Error Message Banner */}
        {errorMessage && (
          <div className="mb-6 bg-red-50 border-l-4 border-red-500 rounded-lg p-4 flex items-start gap-3 shadow-sm animate-slideDown">
            <span className="text-red-500 text-xl">⚠️</span>
            <div className="flex-1">
              <p className="text-sm font-medium text-red-800">{errorMessage}</p>
            </div>
            <button
              onClick={() => setErrorMessage("")}
              className="text-red-400 hover:text-red-600 transition-colors"
            >
              ✕
            </button>
          </div>
        )}

        {/* Modern Info Card with Gradient */}
        <div className="relative overflow-hidden bg-gradient-to-br from-blue-500 to-indigo-600 rounded-2xl p-6 mb-6 shadow-xl">
          <div className="absolute top-0 right-0 w-64 h-64 bg-white/10 rounded-full -mr-32 -mt-32"></div>
          <div className="absolute bottom-0 left-0 w-48 h-48 bg-white/10 rounded-full -ml-24 -mb-24"></div>
          
          <div className="relative z-10">
            <div className="flex items-start gap-4 mb-4">
            
            </div>
            
            {/* Feedback Status with Modern Design */}
            {submitted.length > 0 && (
              <div className="mt-4 pt-4 border-t border-white/20">
                <div className="flex items-center justify-between bg-white/10 backdrop-blur-sm rounded-xl p-3">
                  <div className="flex items-center gap-3">
                    <span className="text-sm font-semibold text-white">
                      📊 Status Analisis AI:
                    </span>
                    {feedbackStatus.ready < feedbackStatus.total ? (
                      <>
                        <span className="text-sm text-blue-100">
                          {feedbackStatus.ready}/{feedbackStatus.total} selesai
                        </span>
                        <div className="flex gap-1">
                          <div className="w-2 h-2 bg-white rounded-full animate-bounce" style={{ animationDelay: '0ms' }}></div>
                          <div className="w-2 h-2 bg-white rounded-full animate-bounce" style={{ animationDelay: '150ms' }}></div>
                          <div className="w-2 h-2 bg-white rounded-full animate-bounce" style={{ animationDelay: '300ms' }}></div>
                        </div>
                      </>
                    ) : (
                      <div className="flex items-center gap-2 bg-green-500/30 px-3 py-1 rounded-full">
                        <span className="text-white">✓</span>
                        <span className="text-sm text-white font-semibold">Analisis Selesai!</span>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Progress Card with Modern Design */}
        {/* Progress Card with Modern Design */}
        <div className="bg-white/80 backdrop-blur-sm rounded-2xl shadow-lg border border-gray-100 p-6 mb-6">
          <div className="flex justify-between text-sm mb-4">
            <div className="flex items-center gap-2">
              <div className="w-8 h-8 bg-gradient-to-br from-indigo-500 to-purple-500 rounded-lg flex items-center justify-center text-white font-bold text-xs">
                {currentIndex + 1}
              </div>
              <span className="text-gray-600 font-medium">dari {questions.length} pertanyaan</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="px-3 py-1 bg-gradient-to-r from-indigo-100 to-purple-100 text-indigo-700 rounded-full text-xs font-semibold">
                {session?.category}
              </span>
              <span className="px-3 py-1 bg-green-100 text-green-700 rounded-full text-xs font-semibold">
                {submitted.length}/{questions.length} ✓
              </span>
            </div>
          </div>
          
          {/* Question Navigation Dots */}
          <div className="flex gap-3 mb-4 justify-center">
            {questions.map((q, idx) => (
              <button
                key={q.id}
                onClick={() => setCurrentIndex(idx)}
                className={`relative group transition-all duration-300 ${
                  submitted.includes(q.id)
                    ? 'w-12 h-12'
                    : idx === currentIndex
                    ? 'w-12 h-12'
                    : 'w-10 h-10'
                }`}
              >
                <div className={`w-full h-full rounded-xl flex items-center justify-center font-bold text-sm transition-all shadow-lg ${
                  submitted.includes(q.id)
                    ? 'bg-gradient-to-br from-green-400 to-green-600 text-white scale-100'
                    : idx === currentIndex
                    ? 'bg-gradient-to-br from-indigo-500 to-purple-500 text-white scale-100'
                    : 'bg-gray-200 text-gray-600 hover:bg-gray-300 scale-95'
                }`}>
                  {submitted.includes(q.id) ? '✓' : idx + 1}
                </div>
                
                {/* Tooltip */}
                <div className="absolute -top-10 left-1/2 transform -translate-x-1/2 bg-gray-800 text-white text-xs px-2 py-1 rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap pointer-events-none">
                  {submitted.includes(q.id) ? 'Sudah dijawab' : idx === currentIndex ? 'Sedang aktif' : 'Belum dijawab'}
                </div>
              </button>
            ))}
          </div>

          {/* Progress Bar */}
          <div className="relative w-full h-3 bg-gray-200 rounded-full overflow-hidden">
            <div
              className="absolute inset-y-0 left-0 bg-gradient-to-r from-indigo-500 to-purple-500 rounded-full transition-all duration-500 ease-out"
              style={{ width: `${progress}%` }}
            >
              <div className="absolute inset-0 bg-white/30 animate-pulse"></div>
            </div>
          </div>
        </div>

        {/* Modern Question Card */}
        {currentQuestion && (
          <div className="bg-white rounded-2xl shadow-xl border border-gray-100 p-8 mb-6 transform transition-all hover:shadow-2xl">
            <div className="flex gap-3 mb-6">
              <span className="px-4 py-1.5 bg-gradient-to-r from-purple-100 to-pink-100 text-purple-700 rounded-full text-xs font-semibold shadow-sm">
                {currentQuestion.category}
              </span>
              <span className="px-4 py-1.5 bg-gray-100 text-gray-600 rounded-full text-xs font-semibold">
                {currentQuestion.difficulty}
              </span>
            </div>
            <div className="relative">
              <div className="absolute -left-4 top-0 w-1 h-full bg-gradient-to-b from-indigo-500 to-purple-500 rounded-full"></div>
              <p className="text-gray-800 text-lg font-medium leading-relaxed pl-4">
                {currentQuestion.content}
              </p>
            </div>
          </div>
        )}

        {/* Modern Answer Box */}
        <div className="bg-white rounded-2xl shadow-xl border border-gray-100 p-8 mb-6">
          <label className="flex items-center gap-2 text-sm font-semibold text-gray-700 mb-4">
            <span className="w-8 h-8 bg-gradient-to-br from-indigo-500 to-purple-500 rounded-lg flex items-center justify-center text-white">
              ✍️
            </span>
            <span>Jawaban Anda:</span>
          </label>
          
          <textarea
            value={answer}
            onChange={(e) => setAnswer(e.target.value)}
            rows={8}
            className="w-full border-2 border-gray-200 rounded-xl px-4 py-3 focus:outline-none focus:border-indigo-500 focus:ring-4 focus:ring-indigo-100 resize-none transition-all text-gray-700 leading-relaxed"
            placeholder="Tulis jawaban Anda dengan detail di sini... Semakin lengkap jawaban Anda, semakin akurat feedback dari AI."
            disabled={submitted.includes(currentQuestion?.id)}
          />

          {submitted.includes(currentQuestion?.id) ? (
            <div className="mt-4 p-4 bg-gradient-to-r from-green-50 to-emerald-50 border-l-4 border-green-500 rounded-lg">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 bg-green-500 rounded-full flex items-center justify-center text-white animate-bounce">
                  ✓
                </div>
                <div>
                  <p className="text-green-800 font-semibold text-sm">Jawaban Terkirim!</p>
                  <p className="text-green-600 text-xs mt-0.5">AI sedang menganalisis jawaban Anda...</p>
                </div>
              </div>
            </div>
          ) : (
            <button
              onClick={submitAnswer}
              disabled={submitting || !answer.trim()}
              className="mt-4 w-full bg-gradient-to-r from-indigo-600 to-purple-600 text-white px-6 py-4 rounded-xl hover:from-indigo-700 hover:to-purple-700 disabled:opacity-50 disabled:cursor-not-allowed font-semibold shadow-lg hover:shadow-xl transition-all transform hover:scale-[1.02] active:scale-[0.98] flex items-center justify-center gap-2"
            >
              {submitting ? (
                <>
                  <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                  <span>Mengirim...</span>
                </>
              ) : (
                <>
                  <span>📤</span>
                  <span>Kirim Jawaban</span>
                </>
              )}
            </button>
          )}
        </div>

        {/* Modern Navigation */}
        <div className="flex justify-between items-center gap-4">
          <button
            onClick={() => setCurrentIndex(Math.max(0, currentIndex - 1))}
            disabled={currentIndex === 0}
            className="flex items-center gap-2 px-6 py-3 border-2 border-gray-300 rounded-xl text-sm font-medium text-gray-700 hover:border-indigo-500 hover:text-indigo-600 disabled:opacity-30 disabled:cursor-not-allowed transition-all hover:shadow-md"
          >
            <span>←</span>
            <span>Sebelumnya</span>
          </button>

          <div className="flex-1 flex flex-col items-end gap-2">
            {/* Warning Messages */}
            {currentIndex === questions.length - 1 && submitted.length < questions.length && (
              <div className="bg-orange-50 border-l-4 border-orange-500 rounded-lg px-4 py-2 animate-pulse">
                <p className="text-xs text-orange-700 font-semibold flex items-center gap-2">
                  <span>⚠️</span>
                  <span>Jawab pertanyaan ini untuk melanjutkan</span>
                </p>
              </div>
            )}
            
            {currentIndex === questions.length - 1 && submitted.length === questions.length && feedbackStatus.ready < feedbackStatus.total && (
              <div className="bg-blue-50 border-l-4 border-blue-500 rounded-lg px-4 py-3 animate-pulse">
                <div className="flex items-center gap-3">
                  <div className="flex gap-1">
                    <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce" style={{ animationDelay: '0ms' }}></div>
                    <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce" style={{ animationDelay: '150ms' }}></div>
                    <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce" style={{ animationDelay: '300ms' }}></div>
                  </div>
                  <p className="text-xs text-blue-700 font-semibold">
                    Tunggu, AI sedang menganalisis jawaban Anda...
                  </p>
                </div>
              </div>
            )}

            {/* Action Buttons */}
            {currentIndex < questions.length - 1 ? (
              <button
                onClick={() => setCurrentIndex(currentIndex + 1)}
                className="flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-xl text-sm font-semibold hover:from-indigo-700 hover:to-purple-700 shadow-lg hover:shadow-xl transition-all transform hover:scale-105"
              >
                <span>Selanjutnya</span>
                <span>→</span>
              </button>
            ) : (
              <button
                onClick={finishSession}
                disabled={finishing || submitted.length < questions.length || feedbackStatus.ready < feedbackStatus.total}
                className="relative group px-8 py-4 bg-gradient-to-r from-green-500 to-emerald-600 text-white rounded-xl text-sm font-bold hover:from-green-600 hover:to-emerald-700 disabled:from-gray-300 disabled:to-gray-400 disabled:cursor-not-allowed shadow-xl hover:shadow-2xl transition-all transform hover:scale-105 disabled:scale-100 disabled:shadow-none flex items-center gap-3 overflow-hidden"
                title={
                  submitted.length < questions.length 
                    ? 'Jawab semua soal terlebih dahulu' 
                    : feedbackStatus.ready < feedbackStatus.total
                    ? 'Tunggu AI selesai menganalisis jawaban'
                    : 'Selesaikan interview dan lihat hasil'
                }
              >
                {/* Shimmer effect */}
                {feedbackStatus.ready >= feedbackStatus.total && (
                  <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent animate-shimmer"></div>
                )}
                
                {finishing ? (
                  <>
                    <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                    <span>Menyelesaikan...</span>
                  </>
                ) : feedbackStatus.ready < feedbackStatus.total ? (
                  <>
                    <span>⏳</span>
                    <span>Tunggu Analisis AI...</span>
                  </>
                ) : (
                  <>
                    <span className="text-xl">🎉</span>
                    <span>Selesaikan Interview</span>
                    <span>→</span>
                  </>
                )}
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}