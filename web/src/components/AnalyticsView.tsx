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

  if (error) return <p className="rounded bg-red-50 p-4 text-sm text-red-700">{error}</p>;

  return (
    <section className="rounded-lg border bg-white p-6">
      <h2 className="font-medium">Headcount by department</h2>
      <div className="mt-4 divide-y">
        {headcount.map((row) => <div key={row.department} className="flex justify-between py-2 text-sm"><span>{row.department}</span><strong>{row.count}</strong></div>)}
        {!headcount.length && <p className="py-3 text-sm text-zinc-500">No headcount data.</p>}
      </div>
    </section>
  );
}