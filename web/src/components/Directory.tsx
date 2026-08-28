import { useEffect, useRef, useState } from "react";
import { api, type Person } from "../lib/api";
import PersonHoverCard from "./PersonHoverCard";

const HOVER_OPEN_MS = 150;
const HOVER_CLOSE_MS = 150;

function readURL() {
  if (typeof window === "undefined") return { q: "", department: "", page: 1, pageSize: 20 };
  const sp = new URLSearchParams(window.location.search);
  return {
    q: sp.get("q") || "",
    department: sp.get("department") || "",
    page: Math.max(1, parseInt(sp.get("page") || "1", 10)),
    pageSize: Math.min(100, Math.max(1, parseInt(sp.get("page_size") || "20", 10))),
  };
}

function pushURL(q: string, department: string, page: number, pageSize: number) {
  const sp = new URLSearchParams();
  if (q) sp.set("q", q);
  if (department) sp.set("department", department);
  if (page !== 1) sp.set("page", String(page));
  if (pageSize !== 20) sp.set("page_size", String(pageSize));
  const qs = sp.toString();
  const url = `${window.location.pathname}${qs ? `?${qs}` : ""}`;
  window.history.pushState(null, "", url);
}

export default function Directory() {
  const [q, setQ] = useState("");
  const [department, setDepartment] = useState("");
  const [persons, setPersons] = useState<Person[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [totalPages, setTotalPages] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const requestVersion = useRef(0);
  const [hover, setHover] = useState<{ id: string; rect: DOMRect } | null>(null);
  const openTimer = useRef<number | null>(null);
  const closeTimer = useRef<number | null>(null);

  function clearHoverTimers() {
    if (openTimer.current !== null) {
      clearTimeout(openTimer.current);
      openTimer.current = null;
    }
    if (closeTimer.current !== null) {
      clearTimeout(closeTimer.current);
      closeTimer.current = null;
    }
  }

function InfoIcon() {
  return (
    <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-3.5 w-3.5">
      <circle cx="8" cy="8" r="6.5" />
      <path d="M8 7.5v3.5" strokeLinecap="round" />
      <circle cx="8" cy="5" r="0.5" fill="currentColor" stroke="none" />
    </svg>
  );
}

function onAnchorEnter(e: React.SyntheticEvent<HTMLElement>, id: string) {
  clearHoverTimers();
  const rect = e.currentTarget.getBoundingClientRect();
  openTimer.current = window.setTimeout(() => setHover({ id, rect }), HOVER_OPEN_MS);
}

function onAnchorLeave() {
  clearHoverTimers();
  closeTimer.current = window.setTimeout(() => setHover(null), HOVER_CLOSE_MS);
}

function onCardEnter() {
  clearHoverTimers();
}

function onCardLeave() {
  clearHoverTimers();
  closeTimer.current = window.setTimeout(() => setHover(null), HOVER_CLOSE_MS);
}

  useEffect(() => {
    if (!hover) return;
    const close = () => setHover(null);
    window.addEventListener("scroll", close, true);
    window.addEventListener("resize", close);
    return () => {
      window.removeEventListener("scroll", close, true);
      window.removeEventListener("resize", close);
    };
  }, [hover !== null]);

  useEffect(() => () => clearHoverTimers(), []);

  async function load(targetPage = page, targetPageSize = pageSize, qVal = q, deptVal = department, push = true) {
    const version = ++requestVersion.current;
    setLoading(true);
    setError(null);
    try {
      const res = await api.listPersons({ q: qVal, department: deptVal, page: targetPage, page_size: targetPageSize });
      if (version !== requestVersion.current) return;
      setPersons(res.data);
      setTotal(res.total);
      setTotalPages(res.total_pages || Math.ceil(res.total / targetPageSize) || 1);
      setPage(res.page);
      setPageSize(res.page_size);
      if (push) pushURL(qVal, deptVal, res.page, res.page_size);
    } catch (e: unknown) {
      if (version !== requestVersion.current) return;
      setError(e instanceof Error ? e.message : "request failed");
    } finally {
      if (version === requestVersion.current) setLoading(false);
    }
  }

  useEffect(() => {
    const init = readURL();
    setQ(init.q);
    setDepartment(init.department);
    setPage(init.page);
    setPageSize(init.pageSize);
    load(init.page, init.pageSize, init.q, init.department, false);

    const onPop = () => {
      const s = readURL();
      setQ(s.q);
      setDepartment(s.department);
      load(s.page, s.pageSize, s.q, s.department, false);
    };
    window.addEventListener("popstate", onPop);
    return () => window.removeEventListener("popstate", onPop);
  }, []);

  function onSearch() {
    load(1, pageSize, q, department, true);
  }

  function goTo(p: number) {
    if (p < 1 || p > totalPages) return;
    load(p, pageSize, q, department, true);
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  function onPageSizeChange(e: React.ChangeEvent<HTMLSelectElement>) {
    const s = parseInt(e.target.value, 10);
    load(1, s, q, department, true);
  }

  function pageNumbers(): (number | "...")[] {
    const pages: (number | "...")[] = [];
    if (totalPages <= 7) {
      for (let i = 1; i <= totalPages; i++) pages.push(i);
      return pages;
    }
    pages.push(1);
    if (page > 3) pages.push("...");
    for (let i = Math.max(2, page - 1); i <= Math.min(totalPages - 1, page + 1); i++) pages.push(i);
    if (page < totalPages - 2) pages.push("...");
    pages.push(totalPages);
    return pages;
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center gap-3">
        <input
          value={q}
          onChange={(e) => setQ(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && onSearch()}
          placeholder="Search name / email / job title..."
          className="w-64 rounded-md border px-3 py-2 text-sm"
        />
        <input
          value={department}
          onChange={(e) => setDepartment(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && onSearch()}
          placeholder="Department (Engineering, Finance...)"
          className="w-64 rounded-md border px-3 py-2 text-sm"
        />
        <button onClick={onSearch} className="rounded-md bg-zinc-900 px-4 py-2 text-sm font-medium text-white hover:bg-zinc-800">
          Search
        </button>
        <span className="self-center text-sm text-zinc-500">
          {total} results · page {page} of {totalPages}
        </span>
        <div className="ml-auto flex items-center gap-2 text-sm">
          <span className="text-zinc-500">Rows:</span>
          <select value={pageSize} onChange={onPageSizeChange} className="rounded-md border px-2 py-1">
            <option value={10}>10</option>
            <option value={20}>20</option>
            <option value={50}>50</option>
            <option value={100}>100</option>
          </select>
        </div>
      </div>

      {loading && <p className="text-sm text-zinc-500">Loading...</p>}
      {error && <p className="rounded-md bg-red-50 px-4 py-3 text-sm text-red-700">{error} — is the Go API running on :8080?</p>}

      <div className="overflow-x-auto rounded-lg border bg-white">
        <table className="w-full text-left text-sm">
          <thead className="bg-zinc-50 text-xs uppercase text-zinc-500">
            <tr>
              <th className="px-4 py-3">Name</th>
              <th className="px-4 py-3">Position</th>
              <th className="px-4 py-3">Department</th>
              <th className="px-4 py-3">Email</th>
              <th className="px-4 py-3">City</th>
            </tr>
          </thead>
          <tbody>
            {persons.map((p, idx) => (
              <tr
                key={p.id}
                className={`border-t ${idx % 2 === 0 ? "bg-white" : "bg-zinc-50/70"} hover:bg-zinc-100/70`}
              >
                <td className="px-4 py-3 font-medium">
                  <span className="flex items-center gap-1.5">
                    <button
                      type="button"
                      aria-label={`Show details for ${p.first_name} ${p.last_name}`}
                      onMouseEnter={(e) => onAnchorEnter(e, p.id)}
                      onMouseLeave={onAnchorLeave}
                      onFocus={(e) => onAnchorEnter(e, p.id)}
                      onBlur={onAnchorLeave}
                      className="shrink-0 rounded p-0.5 text-zinc-400 transition-colors hover:bg-zinc-200 hover:text-zinc-700 focus:outline-none focus-visible:ring-2 focus-visible:ring-zinc-500"
                    >
                      <InfoIcon />
                    </button>
                    <a href={`/person?id=${encodeURIComponent(p.id)}`} className="truncate hover:underline">
                      {p.first_name} {p.last_name}
                    </a>
                    {p.preferred_name && <span className="shrink-0 text-zinc-500">({p.preferred_name})</span>}
                  </span>
                </td>
                <td className="px-4 py-3 text-zinc-700">{p.current_job_title || "—"}</td>
                <td className="px-4 py-3">
                  {p.current_department ? (
                    <span className="rounded-full bg-zinc-900 px-2.5 py-0.5 text-xs font-medium text-white">{p.current_department}</span>
                  ) : (
                    <span className="text-zinc-400">—</span>
                  )}
                </td>
                <td className="px-4 py-3 text-zinc-600">{p.org_email || p.personal_email || "—"}</td>
                <td className="px-4 py-3">{p.city || "—"}</td>
              </tr>
            ))}
            {!loading && persons.length === 0 && (
              <tr>
                <td colSpan={5} className="px-4 py-10 text-center text-zinc-500">
                  No results. Create a person via <code className="rounded bg-zinc-100 px-1">POST /api/v1/persons</code>
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <div className="flex flex-wrap items-center justify-between gap-3">
        <p className="text-sm text-zinc-500">
          Showing {(page - 1) * pageSize + 1}–{Math.min(page * pageSize, total)} of {total}
        </p>
        <div className="flex items-center gap-1">
          <button
            onClick={() => goTo(page - 1)}
            disabled={page <= 1 || loading}
            className="rounded-md border px-3 py-1.5 text-sm disabled:opacity-40 hover:bg-zinc-50"
          >
            ← Prev
          </button>
          {pageNumbers().map((p, i) =>
            p === "..." ? (
              <span key={`dots-${i}`} className="px-2 text-zinc-400">
                …
              </span>
            ) : (
              <button
                key={p}
                onClick={() => goTo(p)}
                disabled={loading}
                className={`min-w-9 rounded-md px-3 py-1.5 text-sm ${p === page ? "bg-zinc-900 text-white" : "border hover:bg-zinc-50"}`}
              >
                {p}
              </button>
            )
          )}
          <button
            onClick={() => goTo(page + 1)}
            disabled={page >= totalPages || loading}
            className="rounded-md border px-3 py-1.5 text-sm disabled:opacity-40 hover:bg-zinc-50"
          >
            Next →
          </button>
        </div>
      </div>

      {hover && persons.some((p) => p.id === hover.id) && (
        <PersonHoverCard
          personId={hover.id}
          fallback={persons.find((p) => p.id === hover.id)!}
          anchorRect={hover.rect}
          onEnter={onCardEnter}
          onLeave={onCardLeave}
        />
      )}
    </div>
  );
}
