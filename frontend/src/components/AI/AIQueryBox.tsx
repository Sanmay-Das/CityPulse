import { useState } from "react";

const apiBase = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export default function AIQueryBox() {
  const [question, setQuestion] = useState("");
  const [answer, setAnswer] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [isOpen, setIsOpen] = useState(true);

  const ask = async () => {
    if (!question.trim()) return;
    setLoading(true);
    setAnswer("");
    setError("");

    try {
      const res = await fetch(`${apiBase}/api/ai/query`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ question }),
      });

      if (!res.ok) throw new Error("Failed to get answer");
      const data = await res.json();
      setAnswer(data.answer);
    } catch (e: any) {
      setError("Could not get an answer. Try again.");
    } finally {
      setLoading(false);
    }
  };

  if (!isOpen) return (
    <button
      onClick={() => setIsOpen(true)}
      className="absolute bottom-10 right-5 z-20 bg-chicago-red hover:bg-red-700 text-white text-xs px-3 py-2 rounded-xl shadow-lg transition-colors"
    >
      Ask AI
    </button>
  );

  return (
    <div className="absolute bottom-10 right-5 z-20 w-80 bg-slate-900/95 backdrop-blur border border-slate-700 rounded-xl shadow-xl p-4">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-white font-bold text-sm">Ask AI</h3>
        <button onClick={() => setIsOpen(false)} className="text-slate-400 hover:text-white text-lg leading-none">×</button>
      </div>

      <textarea
        value={question}
        onChange={(e) => setQuestion(e.target.value)}
        onKeyDown={(e) => e.key === "Enter" && !e.shiftKey && (e.preventDefault(), ask())}
        placeholder="e.g. Which ZIP in Chicago had the most robberies in 2020?"
        className="w-full bg-slate-800 border border-slate-600 text-slate-200 text-xs rounded-lg px-3 py-2 resize-none focus:outline-none focus:ring-1 focus:ring-chicago-red"
        rows={3}
      />

      <button
        onClick={ask}
        disabled={loading || !question.trim()}
        className="mt-2 w-full bg-chicago-red hover:bg-red-700 disabled:opacity-50 text-white text-xs font-semibold py-2 rounded-lg transition-colors"
      >
        {loading ? "Thinking…" : "Ask"}
      </button>

      {error && (
        <p className="text-red-400 text-xs mt-2">{error}</p>
      )}

      {answer && (
        <div className="mt-3 bg-slate-800 rounded-lg p-3">
          <p className="text-slate-300 text-xs leading-relaxed">{answer}</p>
        </div>
      )}
    </div>
  );
}
