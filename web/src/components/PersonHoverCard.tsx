import { useEffect, useState, type ReactNode } from "react";
import { api, type EmploymentRecord, type Person } from "../lib/api";

interface Details {
  person: Person;
  employment: EmploymentRecord | null;
}

const MAX_CACHE_SIZE = 50;
const CACHE_TTL_MS = 5 * 60 * 1000; // 5 minutes
const cache = new Map<string, { details: Details; timestamp: number }>();

function getCached(id: string): Details | null {
  const entry = cache.get(id);
  if (!entry) return null;
  if (Date.now() - entry.timestamp > CACHE_TTL_MS) {
    cache.delete(id);
    return null;
  }
  return entry.details;
}

function setCached(id: string, details: Details) {
  if (cache.size >= MAX_CACHE_SIZE) {
    const oldestKey = cache.keys().next().value;
    if (oldestKey) cache.delete(oldestKey);
  }
  cache.set(id, { details, timestamp: Date.now() });
}

const CARD_WIDTH = 340;
const EST_HEIGHT = 380;

interface Props {
  personId: string;
  fallback: Person;
  anchorRect: DOMRect;
  onEnter: () => void;
  onLeave: () => void;
  onOpenDrawer: (id: string) => void;
}

function initials(p: Person) {
  return `${p.first_name?.[0] ?? ""}${p.last_name?.[0] ?? ""}`.toUpperCase();
}

function Row({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="flex justify-between gap-3">
      <span className="shrink-0 text-slate-500 dark:text-slate-400">{label}</span>
      <span className="min-w-0 truncate text-right text-navy dark:text-slate-200 font-medium">{children}</span>
    </div>
  );
}

export default function PersonHoverCard({ personId, fallback, anchorRect, onEnter, onLeave, onOpenDrawer }: Props) {
  const [details, setDetails] = useState<Details | null>(getCached(personId));

  useEffect(() => {
    const cached = getCached(personId);
    if (cached) {
      setDetails(cached);
      return;
    }
    const controller = new AbortController();
    setDetails(null);
    Promise.all([
      api.getPerson(personId, { signal: controller.signal }),
      api.getCurrentEmployment(personId, { signal: controller.signal }).catch(() => null),
    ])
      .then(([person, employment]) => {
        const entry: Details = { person, employment };
        setCached(personId, entry);
        if (!controller.signal.aborted) setDetails(entry);
      })
      .catch(() => {
        if (!controller.signal.aborted) setDetails(null);
      });
    return () => {
      controller.abort();
    };
  }, [personId]);

  const person = details?.person ?? fallback;
  const employment = details?.employment ?? null;
  const loading = !details;

  const vw = window.innerWidth;
  const vh = window.innerHeight;
  const left = Math.min(Math.max(anchorRect.left, 8), vw - CARD_WIDTH - 8);
  const below = vh - anchorRect.bottom >= EST_HEIGHT;

  return (
    <div
      onMouseEnter={onEnter}
      onMouseLeave={onLeave}
      style={{
        position: "fixed",
        left,
        width: CARD_WIDTH,
        ...(below ? { top: anchorRect.bottom + 8 } : { bottom: vh - anchorRect.top + 8 }),
      }}
      className="z-50 rounded-lg border border-cream-border dark:border-dark-border bg-cream-card dark:bg-dark-card p-4 shadow-2xl transition-all"
    >
      <div className="flex items-start gap-3">
        {person.profile_photo_url ? (
          <img src={person.profile_photo_url} alt="" className="h-12 w-12 shrink-0 rounded-full object-cover border border-cream-border dark:border-dark-border" />
        ) : (
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-pastel-mint/30 dark:bg-mint/20 text-navy dark:text-mint border border-forest/30 dark:border-mint/40 font-bold text-sm">
            {initials(person)}
          </div>
        )}
        <div className="min-w-0 flex-1">
          <p className="truncate font-semibold text-navy dark:text-slate-100">
            {person.first_name} {person.last_name}
            {loading && <span className="ml-2 text-xs font-normal text-slate-400">loading…</span>}
          </p>
          {person.preferred_name && <p className="text-xs text-slate-500 dark:text-slate-400">&ldquo;{person.preferred_name}&rdquo;</p>}
          <p className="mt-0.5 truncate text-sm text-slate-500 dark:text-slate-300">
            {employment?.job_title || person.current_job_title || "—"}
          </p>
          <div className="mt-1.5 flex flex-wrap gap-1">
            {(employment?.department || person.current_department) && (
              <span className="rounded-full bg-pastel-mint/30 dark:bg-mint/20 px-2.5 py-0.5 text-xs font-semibold text-navy dark:text-mint border border-forest/30 dark:border-mint/40">
                {employment?.department || person.current_department}
              </span>
            )}
            {employment?.employment_status && (
              <span className="rounded-full bg-cream-hover dark:bg-dark-hover px-2 py-0.5 text-xs text-slate-600 dark:text-slate-300 border border-cream-border dark:border-dark-border">
                {employment.employment_status}
              </span>
            )}
            {!person.is_active && (
              <span className="rounded-full bg-red-500/10 px-2 py-0.5 text-xs text-red-600 dark:text-red-400 border border-red-500/20">inactive</span>
            )}
          </div>
        </div>
      </div>

      <div className="mt-3 space-y-1.5 border-t border-cream-border dark:border-dark-border pt-3 text-sm">
        <Row label="Email">{person.org_email || person.personal_email || "—"}</Row>
        <Row label="Phone">{person.phone_primary || "—"}</Row>
        <Row label="Location">{[person.city, person.country].filter(Boolean).join(", ") || "—"}</Row>
        <Row label="Office">{person.current_office_location || "—"}</Row>
        <Row label="Team">{employment?.team || person.current_team || "—"}</Row>
        <Row label="Hired">{(employment?.hire_date || person.current_hire_date || "").slice(0, 10) || "—"}</Row>
        {employment?.job_level && <Row label="Level">{employment.job_level}</Row>}
      </div>

      {person.tags.length > 0 && (
        <div className="mt-3 flex flex-wrap gap-1 border-t border-cream-border dark:border-dark-border pt-3">
          {person.tags.map((t) => (
            <span key={t} className="rounded-full bg-cream-hover dark:bg-dark-hover px-2 py-0.5 text-xs text-slate-600 dark:text-slate-300 border border-cream-border dark:border-dark-border">
              {t}
            </span>
          ))}
        </div>
      )}

      <button
        type="button"
        onClick={() => onOpenDrawer(personId)}
        className="mt-3 block w-full rounded-md bg-forest hover:bg-forest-hover text-white dark:bg-mint dark:hover:bg-mint-hover dark:text-slate-950 px-3 py-2 text-center text-sm font-semibold transition-colors shadow-sm cursor-pointer"
      >
        View full profile →
      </button>
    </div>
  );
}
