import { useEffect, useState } from "react";
import { api, type HeadcountRow, type MovementPoint } from "../lib/api";

export default function AnalyticsView() {
  const [headcount, setHeadcount] = useState<HeadcountRow[]>([]);
  const [movements, setMovements] = useState<MovementPoint[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    Promise.all([api.headcount(), api.movements("2018-01-01", "2026-12-31")])
      .then(([loadedHeadcount, loadedMovements]) => {
        if (cancelled) return;
        setHeadcount(loadedHeadcount);
        setMovements(loadedMovements);
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
    <div className="mt-8 grid gap-6 lg:grid-cols-2">
      <section className="rounded-lg border bg-white p-6">
        <h2 className="font-medium">Headcount by department</h2>
        <div className="mt-4 divide-y">
          {headcount.map((row) => <div key={row.department} className="flex justify-between py-2 text-sm"><span>{row.department}</span><strong>{row.count}</strong></div>)}
          {!headcount.length && <p className="py-3 text-sm text-zinc-500">No headcount data.</p>}
        </div>
      </section>
      <section className="rounded-lg border bg-white p-6">
        <h2 className="font-medium">Movements</h2>
        <div className="mt-4 divide-y">
          {movements.slice(-20).map((row) => <div key={row.date} className="grid grid-cols-4 gap-2 py-2 text-sm"><span>{row.date}</span><span>+{row.new_entries}</span><span>-{row.exits}</span><strong>{row.net_change >= 0 ? "+" : ""}{row.net_change}</strong></div>)}
          {!movements.length && <p className="py-3 text-sm text-zinc-500">No movements recorded.</p>}
        </div>
      </section>
    </div>
  );
}
