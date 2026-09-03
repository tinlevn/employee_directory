import React, { useState } from "react";

export default function RegisterForm() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [personId, setPersonId] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleRegister(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const BASE = import.meta.env.PUBLIC_API_URL || "http://localhost:8080";
      const res = await fetch(`${BASE}/api/v1/auth/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            username,
            password,
            person_id: personId,
            role: "staff"
        }),
      });

      if (!res.ok) {
        const body = await res.json().catch(() => ({}));
        throw new Error(body.detail || body.message || "Registration failed");
      }

      const data = await res.json();
      localStorage.setItem("token", data.token);
      window.location.href = "/";
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="mx-auto mt-10 max-w-sm rounded-xl border border-[#E6DBC5] dark:border-[#2b303c] bg-white dark:bg-[#1c1f26] p-6 shadow-sm">
      <h2 className="text-xl font-semibold text-[#141E46] dark:text-white">Register</h2>
      <p className="mt-1 text-sm text-[#5A6578] dark:text-slate-400">Create an account using your Person ID.</p>

      {error && <div className="mt-4 rounded-md bg-red-500/10 border border-red-500/20 p-3 text-sm text-red-600 dark:text-red-400">{error}</div>}

      <form onSubmit={handleRegister} className="mt-6 space-y-4">
        <div>
          <label htmlFor="reg-username" className="block text-sm font-medium text-[#141E46] dark:text-slate-300">Username</label>
          <input
            id="reg-username"
            type="text"
            required
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="mt-1 block w-full rounded-md border border-[#E6DBC5] dark:border-[#2b303c] bg-white dark:bg-[#131519] px-3 py-2 text-[#141E46] dark:text-slate-100 shadow-sm focus:border-[#41B06E] dark:focus:border-[#1DCD9F] focus:outline-none focus:ring-1 focus:ring-[#41B06E] dark:focus:ring-[#1DCD9F] sm:text-sm transition-colors"
          />
        </div>
        <div>
          <label htmlFor="reg-password" className="block text-sm font-medium text-[#141E46] dark:text-slate-300">Password</label>
          <input
            id="reg-password"
            type="password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="mt-1 block w-full rounded-md border border-[#E6DBC5] dark:border-[#2b303c] bg-white dark:bg-[#131519] px-3 py-2 text-[#141E46] dark:text-slate-100 shadow-sm focus:border-[#41B06E] dark:focus:border-[#1DCD9F] focus:outline-none focus:ring-1 focus:ring-[#41B06E] dark:focus:ring-[#1DCD9F] sm:text-sm transition-colors"
          />
        </div>
        <div>
          <label htmlFor="reg-personid" className="block text-sm font-medium text-[#141E46] dark:text-slate-300">Person ID (UUID)</label>
          <input
            id="reg-personid"
            type="text"
            required
            value={personId}
            onChange={(e) => setPersonId(e.target.value)}
            className="mt-1 block w-full rounded-md border border-[#E6DBC5] dark:border-[#2b303c] bg-white dark:bg-[#131519] px-3 py-2 text-[#141E46] dark:text-slate-100 shadow-sm focus:border-[#41B06E] dark:focus:border-[#1DCD9F] focus:outline-none focus:ring-1 focus:ring-[#41B06E] dark:focus:ring-[#1DCD9F] sm:text-sm font-mono placeholder-[#7A869A] dark:placeholder-slate-600 transition-colors"
            placeholder="00000000-0000-0000-0000-000000000000"
          />
        </div>
        <button
          type="submit"
          disabled={loading}
          className="w-full flex justify-center rounded-md bg-[#41B06E] hover:bg-[#329057] text-white dark:bg-[#1DCD9F] dark:hover:bg-[#169976] dark:text-slate-950 px-4 py-2 text-sm font-semibold shadow-sm focus:outline-none focus:ring-2 focus:ring-[#41B06E] dark:focus:ring-[#1DCD9F] disabled:opacity-50 transition-colors"
        >
          {loading ? "Registering..." : "Register"}
        </button>
      </form>
      <div className="mt-4 text-center text-sm">
        <a href="/login" className="text-[#5A6578] dark:text-slate-400 hover:text-[#41B06E] dark:hover:text-[#1DCD9F] transition-colors">Already have an account? Sign in</a>
      </div>
    </div>
  );
}
