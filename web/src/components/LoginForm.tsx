import React, { useState } from "react";

export default function LoginForm() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleLogin(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const BASE = import.meta.env.PUBLIC_API_URL || "http://localhost:8080";
      const res = await fetch(`${BASE}/api/v1/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      if (!res.ok) {
        const body = await res.json().catch(() => ({}));
        throw new Error(body.detail || body.message || "Login failed");
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
    <div className="mx-auto mt-10 max-w-sm rounded-xl border bg-white p-6 shadow-sm">
      <h2 className="text-xl font-semibold text-zinc-900">Sign in</h2>
      <p className="mt-1 text-sm text-zinc-500">Welcome back to the employee directory.</p>

      {error && <div className="mt-4 rounded bg-red-50 p-3 text-sm text-red-700">{error}</div>}

      <form onSubmit={handleLogin} className="mt-6 space-y-4">
        <div>
          <label htmlFor="login-username" className="block text-sm font-medium text-zinc-700">Username</label>
          <input
            id="login-username"
            type="text"
            required
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="mt-1 block w-full rounded-md border border-zinc-300 px-3 py-2 shadow-sm focus:border-zinc-500 focus:outline-none focus:ring-1 focus:ring-zinc-500 sm:text-sm"
          />
        </div>
        <div>
          <label htmlFor="login-password" className="block text-sm font-medium text-zinc-700">Password</label>
          <input
            id="login-password"
            type="password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="mt-1 block w-full rounded-md border border-zinc-300 px-3 py-2 shadow-sm focus:border-zinc-500 focus:outline-none focus:ring-1 focus:ring-zinc-500 sm:text-sm"
          />
        </div>
        <button
          type="submit"
          disabled={loading}
          className="w-full flex justify-center rounded-md border border-transparent bg-zinc-900 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-zinc-800 focus:outline-none focus:ring-2 focus:ring-zinc-900 focus:ring-offset-2 disabled:opacity-50"
        >
          {loading ? "Signing in..." : "Sign in"}
        </button>
      </form>
      <div className="mt-4 text-center text-sm">
        <a href="/register" className="text-zinc-600 hover:underline">Don't have an account? Register</a>
      </div>
    </div>
  );
}
