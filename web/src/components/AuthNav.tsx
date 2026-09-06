import { useEffect, useState } from "react";

export default function AuthNav() {
  const [hasToken, setHasToken] = useState(false);

  useEffect(() => {
    setHasToken(!!localStorage.getItem("token"));
  }, []);

  function handleLogout() {
    localStorage.removeItem("token");
    setHasToken(false);
    window.location.href = "/login";
  }

  if (hasToken) {
    return (
      <button
        onClick={handleLogout}
        type="button"
        className="text-xs font-medium px-2.5 py-1.5 rounded-md border border-[#E6DBC5] dark:border-[#2b303c] hover:bg-red-500/10 hover:border-red-500/30 text-red-600 dark:text-red-400 transition-colors cursor-pointer"
      >
        Sign out
      </button>
    );
  }

  return (
    <div className="flex items-center gap-2 text-xs">
      <a
        href="/login"
        className="px-2.5 py-1.5 rounded-md text-[#141E46] dark:text-slate-300 hover:text-[#41B06E] dark:hover:text-[#1DCD9F] transition-colors"
      >
        Sign in
      </a>
      <a
        href="/register"
        className="px-2.5 py-1.5 rounded-md bg-[#41B06E] dark:bg-[#1DCD9F] text-white dark:text-[#131519] font-semibold hover:opacity-90 transition-opacity"
      >
        Register
      </a>
    </div>
  );
}
