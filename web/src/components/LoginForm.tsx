import { useState, type SyntheticEvent } from "react";
import { api } from "../lib/api";

export default function LoginForm() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleLogin(e: SyntheticEvent<HTMLFormElement>) {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const data = await api.login({ username, password });
      localStorage.setItem("token", data.token);
      window.location.href = "/";
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Login failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="mx-auto mt-10 max-w-sm rounded-xl border border-cream-border dark:border-dark-border bg-cream-card dark:bg-dark-card p-6 shadow-sm">
      <h2 className="text-xl font-semibold text-navy dark:text-white">Sign in</h2>
      <p className="mt-1 text-sm text-slate-500 dark:text-slate-400">Welcome back to the employee directory.</p>

      {error && <div className="mt-4 rounded-md bg-red-500/10 border border-red-500/20 p-3 text-sm text-red-600 dark:text-red-400">{error}</div>}

      <form onSubmit={handleLogin} className="mt-6 space-y-4">
        <div>
          <label htmlFor="login-username" className="block text-sm font-medium text-navy dark:text-slate-300">Username</label>
          <input
            id="login-username"
            type="text"
            required
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="mt-1 block w-full rounded-md border border-cream-border dark:border-dark-border bg-cream-card dark:bg-dark-bg px-3 py-2 text-navy dark:text-slate-100 shadow-sm focus:border-forest dark:focus:border-mint focus:outline-none focus:ring-1 focus:ring-forest dark:focus:ring-mint sm:text-sm transition-colors"
          />
        </div>
        <div>
          <label htmlFor="login-password" className="block text-sm font-medium text-navy dark:text-slate-300">Password</label>
          <input
            id="login-password"
            type="password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="mt-1 block w-full rounded-md border border-cream-border dark:border-dark-border bg-cream-card dark:bg-dark-bg px-3 py-2 text-navy dark:text-slate-100 shadow-sm focus:border-forest dark:focus:border-mint focus:outline-none focus:ring-1 focus:ring-forest dark:focus:ring-mint sm:text-sm transition-colors"
          />
        </div>
        <button
          type="submit"
          disabled={loading}
          className="w-full flex justify-center rounded-md bg-forest hover:bg-forest-hover text-white dark:bg-mint dark:hover:bg-mint-hover dark:text-slate-950 px-4 py-2 text-sm font-semibold shadow-sm focus:outline-none focus:ring-2 focus:ring-forest dark:focus:ring-mint disabled:opacity-50 transition-colors cursor-pointer"
        >
          {loading ? "Signing in..." : "Sign in"}
        </button>
      </form>
      <div className="mt-4 text-center text-sm">
        <a href="/register" className="text-slate-500 dark:text-slate-400 hover:text-forest dark:hover:text-mint transition-colors">Don't have an account? Register</a>
      </div>
    </div>
  );
}
