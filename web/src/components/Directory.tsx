import { useEffect, useRef, useState } from "react";
import { api, type Person } from "../lib/api";
import PersonHoverCard from "./PersonHoverCard";

const HOVER_OPEN_MS = 150;
const HOVER_CLOSE_MS = 150;

function readURL() {
  if (typeof window === "undefined") return { q: "", department: "", page: 1, pageSize: 20, sort: "" };
  const sp = new URLSearchParams(window.location.search);
  return {
    q: sp.get("q") || "",
    department: sp.get("department") || "",
    page: Math.max(1, parseInt(sp.get("page") || "1", 10)),
    pageSize: Math.min(100, Math.max(1, parseInt(sp.get("page_size") || "20", 10))),
    sort: sp.get("sort") || "",
  };
}

function pushURL(q: string, department: string, page: number, pageSize: number, sort: string) {
  const sp = new URLSearchParams();
  if (q) sp.set("q", q);
  if (department) sp.set("department", department);
  if (page !== 1) sp.set("page", String(page));
  if (pageSize !== 20) sp.set("page_size", String(pageSize));
  if (sort) sp.set("sort", sort);
  const qs = sp.toString();
  const url = `${window.location.pathname}${qs ? `?${qs}` : ""}`;
  window.history.pushState(null, "", url);
}

export default function Directory() {
  const [q, setQ] = useState("");
  const [department, setDepartment] = useState("");
  const [sort, setSort] = useState("");
  const [persons, setPersons] = useState<Person[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [totalPages, setTotalPages] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [departments, setDepartments] = useState<string[]>([]);
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

  async function load(targetPage = page, targetPageSize = pageSize, qVal = q, deptVal = department, sortVal = sort, push = true) {
    const version = ++requestVersion.current;
    setLoading(true);
    setError(null);
    try {
      const res = await api.listPersons({ q: qVal, department: deptVal, page: targetPage, page_size: targetPageSize, sort: sortVal });
      if (version !== requestVersion.current) return;

      let list = [...res.data];
      if (sortVal === "first_name") {
        list.sort((a, b) => `${a.first_name} ${a.last_name}`.localeCompare(`${b.first_name} ${b.last_name}`));
      } else if (sortVal === "-first_name") {
        list.sort((a, b) => `${b.first_name} ${b.last_name}`.localeCompare(`${a.first_name} ${a.last_name}`));
      } else if (sortVal === "city") {
        list.sort((a, b) => (a.city || "").localeCompare(b.city || ""));
      } else if (sortVal === "-city") {
        list.sort((a, b) => (b.city || "").localeCompare(a.city || ""));
      }

      setPersons(list);
      setTotal(res.total);
      setTotalPages(res.total_pages || Math.ceil(res.total / targetPageSize) || 1);
      setPage(res.page);
      setPageSize(res.page_size);
      if (push) pushURL(qVal, deptVal, res.page, res.page_size, sortVal);
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
    setSort(init.sort);
    setPage(init.page);
    setPageSize(init.pageSize);
    load(init.page, init.pageSize, init.q, init.department, init.sort, false);

    // fetch departments for dropdown
    api.headcount()
      .then((data) => setDepartments(data.map((d) => d.department).sort()))
      .catch(() => {});

    const onPop = () => {
      const s = readURL();
      setQ(s.q);
      setDepartment(s.department);
      setSort(s.sort);
      load(s.page, s.pageSize, s.q, s.department, s.sort, false);
    };
    window.addEventListener("popstate", onPop);
    return () => window.removeEventListener("popstate", onPop);
  }, []);

  function onSearch() {
    load(1, pageSize, q, department, sort, true);
  }

  function goTo(p: number) {
    if (p < 1 || p > totalPages) return;
    load(p, pageSize, q, department, sort, true);
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  function onPageSizeChange(e: React.ChangeEvent<HTMLSelectElement>) {
    const s = parseInt(e.target.value, 10);
    load(1, s, q, department, sort, true);
  }

  function toggleSort(col: "name" | "city") {
    let nextSort = "";
    if (col === "name") {
      if (sort === "first_name") {
        nextSort = "-first_name";
      } else if (sort === "-first_name") {
        nextSort = "";
      } else {
        nextSort = "first_name";
      }
    } else if (col === "city") {
      if (sort === "city") {
        nextSort = "-city";
      } else if (sort === "-city") {
        nextSort = "";
      } else {
        nextSort = "city";
      }
    }
    setSort(nextSort);
    load(1, pageSize, q, department, nextSort, true);
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
          className="w-64 rounded-md border border-[#E6DBC5] dark:border-[#2b303c] bg-white dark:bg-[#1c1f26] px-3 py-2 text-sm text-[#141E46] dark:text-slate-100 placeholder-[#7A869A] dark:placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-[#41B06E]/50 dark:focus:ring-[#1DCD9F]/50 focus:border-[#41B06E] dark:focus:border-[#1DCD9F] transition-colors"
        />
        <select
          value={department}
          onChange={(e) => setDepartment(e.target.value)}
          className="w-64 rounded-md border border-[#E6DBC5] dark:border-[#2b303c] bg-white dark:bg-[#1c1f26] px-3 py-2 text-sm text-[#141E46] dark:text-slate-100 focus:outline-none focus:ring-2 focus:ring-[#41B06E]/50 dark:focus:ring-[#1DCD9F]/50 focus:border-[#41B06E] dark:focus:border-[#1DCD9F] transition-colors"
        >
          <option value="">All Departments</option>
          {departments.map((d) => (
            <option key={d} value={d}>
              {d}
            </option>
          ))}
        </select>
        <button onClick={onSearch} className="rounded-md bg-[#41B06E] hover:bg-[#329057] text-white dark:bg-[#1DCD9F] dark:hover:bg-[#169976] dark:text-slate-950 px-4 py-2 text-sm font-semibold transition-colors shadow-sm focus:outline-none focus:ring-2 focus:ring-[#41B06E] dark:focus:ring-[#1DCD9F]">
          Search
        </button>
        <span className="self-center text-sm text-[#5A6578] dark:text-slate-400">
          {total} results · page {page} of {totalPages}
        </span>
        <div className="ml-auto flex items-center gap-2 text-sm">
          <span className="text-[#5A6578] dark:text-slate-400">Rows:</span>
          <select value={pageSize} onChange={onPageSizeChange} className="rounded-md border border-[#E6DBC5] dark:border-[#2b303c] bg-white dark:bg-[#1c1f26] text-[#141E46] dark:text-slate-100 px-2 py-1 focus:outline-none focus:ring-2 focus:ring-[#41B06E]/50 dark:focus:ring-[#1DCD9F]/50">
            <option value={10}>10</option>
            <option value={20}>20</option>
            <option value={50}>50</option>
            <option value={100}>100</option>
          </select>
        </div>
      </div>

      {loading && <p className="text-sm text-[#5A6578] dark:text-slate-400">Loading...</p>}
      {error && <p className="rounded-md bg-red-500/10 border border-red-500/20 px-4 py-3 text-sm text-red-600 dark:text-red-400">{error} — is the Go API running on :8080?</p>}

      <div className="overflow-x-auto rounded-lg border border-[#E6DBC5] dark:border-[#2b303c] bg-white dark:bg-[#1c1f26] shadow-sm">
        <table className="w-full text-left text-sm">
          <thead className="bg-[#F8EFE0] dark:bg-[#252a34]/60 text-xs uppercase font-semibold text-[#141E46] dark:text-slate-400 border-b border-[#E6DBC5] dark:border-[#2b303c]">
            <tr>
              <th
                onClick={() => toggleSort("name")}
                className="px-4 py-3 cursor-pointer select-none hover:text-[#41B06E] dark:hover:text-[#1DCD9F] transition-colors"
                title="Click to sort by Name (A-Z / Z-A)"
              >
                <div className="flex items-center gap-1.5">
                  <span>Name</span>
                  <span className="text-xs font-bold">
                    {sort === "first_name" && "▲"}
                    {sort === "-first_name" && "▼"}
                    {sort !== "first_name" && sort !== "-first_name" && <span className="opacity-30">↕</span>}
                  </span>
                </div>
              </th>
              <th className="px-4 py-3">Position</th>
              <th className="px-4 py-3">Department</th>
              <th className="px-4 py-3">Email</th>
              <th
                onClick={() => toggleSort("city")}
                className="px-4 py-3 cursor-pointer select-none hover:text-[#41B06E] dark:hover:text-[#1DCD9F] transition-colors"
                title="Click to sort by City (A-Z / Z-A)"
              >
                <div className="flex items-center gap-1.5">
                  <span>City</span>
                  <span className="text-xs font-bold">
                    {sort === "city" && "▲"}
                    {sort === "-city" && "▼"}
                    {sort !== "city" && sort !== "-city" && <span className="opacity-30">↕</span>}
                  </span>
                </div>
              </th>
            </tr>
          </thead>
          <tbody>
            {persons.map((p, idx) => (
              <tr
                key={p.id}
                className={`border-t border-[#E6DBC5] dark:border-[#2b303c] ${idx % 2 === 0 ? "bg-white dark:bg-[#1c1f26]" : "bg-[#FFF9EE] dark:bg-[#20242d]"} hover:bg-[#8DECB4]/25 dark:hover:bg-[#1DCD9F]/10 transition-colors`}
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
                      className="shrink-0 rounded p-0.5 text-[#5A6578] dark:text-slate-500 transition-colors hover:bg-[#E6DBC5]/60 dark:hover:bg-[#252a34] hover:text-[#141E46] dark:hover:text-slate-200 focus:outline-none focus-visible:ring-2 focus-visible:ring-[#41B06E]"
                    >
                      <InfoIcon />
                    </button>
                    <a href={`/person?id=${encodeURIComponent(p.id)}`} className="truncate text-[#141E46] dark:text-slate-100 hover:text-[#41B06E] dark:hover:text-[#1DCD9F] transition-colors font-medium">
                      {p.first_name} {p.last_name}
                    </a>
                    {p.preferred_name && <span className="shrink-0 text-[#7A869A] dark:text-slate-500">({p.preferred_name})</span>}
                  </span>
                </td>
                <td className="px-4 py-3 text-[#141E46]/90 dark:text-slate-300">{p.current_job_title || "—"}</td>
                <td className="px-4 py-3">
                  {p.current_department ? (
                    <span className="rounded-full bg-[#8DECB4]/30 dark:bg-[#1DCD9F]/20 px-2.5 py-0.5 text-xs font-semibold text-[#141E46] dark:text-[#1DCD9F] border border-[#41B06E]/30 dark:border-[#1DCD9F]/40">{p.current_department}</span>
                  ) : (
                    <span className="text-[#94a0b2] dark:text-slate-600">—</span>
                  )}
                </td>
                <td className="px-4 py-3 text-[#5A6578] dark:text-slate-400">{p.org_email || p.personal_email || "—"}</td>
                <td className="px-4 py-3 text-[#141E46]/90 dark:text-slate-300">{p.city || "—"}</td>
              </tr>
            ))}
            {!loading && persons.length === 0 && (
              <tr>
                <td colSpan={5} className="px-4 py-10 text-center text-[#5A6578] dark:text-slate-400">
                  No results. Create a person via <code className="rounded bg-[#F8EFE0] dark:bg-[#252a34] text-[#141E46] dark:text-[#1DCD9F] px-1 py-0.5 font-mono text-xs">POST /api/v1/persons</code>
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <div className="flex flex-wrap items-center justify-between gap-3">
        <p className="text-sm text-[#5A6578] dark:text-slate-400">
          Showing {(page - 1) * pageSize + 1}–{Math.min(page * pageSize, total)} of {total}
        </p>
        <div className="flex items-center gap-1">
          <button
            onClick={() => goTo(page - 1)}
            disabled={page <= 1 || loading}
            className="rounded-md border border-[#E6DBC5] dark:border-[#2b303c] bg-white dark:bg-[#1c1f26] text-[#141E46] dark:text-slate-300 px-3 py-1.5 text-sm disabled:opacity-40 hover:bg-[#8DECB4]/20 hover:text-[#41B06E] dark:hover:bg-[#1DCD9F]/10 dark:hover:text-[#1DCD9F] transition-colors"
          >
            ← Prev
          </button>
          {pageNumbers().map((p, i) =>
            p === "..." ? (
              <span key={`dots-${i}`} className="px-2 text-[#94a0b2] dark:text-slate-600">
                …
              </span>
            ) : (
              <button
                key={p}
                onClick={() => goTo(p)}
                disabled={loading}
                className={`min-w-9 rounded-md px-3 py-1.5 text-sm transition-colors ${p === page ? "bg-[#41B06E] text-white font-semibold border border-[#41B06E] dark:bg-[#1DCD9F] dark:text-slate-950 dark:border-[#1DCD9F]" : "border border-[#E6DBC5] dark:border-[#2b303c] bg-white dark:bg-[#1c1f26] text-[#141E46] dark:text-slate-300 hover:bg-[#8DECB4]/20 hover:text-[#41B06E] dark:hover:bg-[#1DCD9F]/10 dark:hover:text-[#1DCD9F]"}`}
              >
                {p}
              </button>
            )
          )}
          <button
            onClick={() => goTo(page + 1)}
            disabled={page >= totalPages || loading}
            className="rounded-md border border-[#E6DBC5] dark:border-[#2b303c] bg-white dark:bg-[#1c1f26] text-[#141E46] dark:text-slate-300 px-3 py-1.5 text-sm disabled:opacity-40 hover:bg-[#8DECB4]/20 hover:text-[#41B06E] dark:hover:bg-[#1DCD9F]/10 dark:hover:text-[#1DCD9F] transition-colors"
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
