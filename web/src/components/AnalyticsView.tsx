import { useEffect, useState } from "react";
import { api, type HeadcountRow } from "../lib/api";

export default function AnalyticsView() {
  const [headcount, setHeadcount] = useState<HeadcountRow[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    api.headcount()
      .then((data) => {
        if (!cancelled) setHeadcount(data);
      })
      .catch((reason: Error) => {
        if (!cancelled) setError(reason.message);
      });
    return () => {
      cancelled = true;
    };
  }, []);

  if (error) return <p className="rounded-md bg-red-500/10 border border-red-500/20 p-4 text-sm text-red-600 dark:text-red-400">{error}</p>;

  return (
    <section className="rounded-lg border border-[#E6DBC5] dark:border-[#2b303c] bg-white dark:bg-[#1c1f26] p-6 shadow-sm">
      <h2 className="font-semibold text-[#141E46] dark:text-slate-100 flex items-center gap-2">
        <span className="h-2.5 w-2.5 rounded-full bg-[#41B06E] dark:bg-[#1DCD9F]"></span>
        Headcount by department
      </h2>
      <div className="mt-4 divide-y divide-[#E6DBC5] dark:divide-[#2b303c]">
        {headcount.map((row) => (
          <div key={row.department} className="flex justify-between py-2.5 text-sm hover:bg-[#8DECB4]/20 dark:hover:bg-[#1DCD9F]/5 px-2 rounded transition-colors">
            <span className="text-[#141E46]/90 dark:text-slate-300 font-medium">{row.department}</span>
            <strong className="font-mono text-[#41B06E] dark:text-[#1DCD9F]">{row.count}</strong>
          </div>
        ))}
        {!headcount.length && <p className="py-3 text-sm text-[#5A6578] dark:text-slate-400">No headcount data.</p>}
      </div>
    </section>
  );
}