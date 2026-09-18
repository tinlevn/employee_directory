import { useEffect, useRef, useState } from "react";
import { api, type EmergencyContact, type EmploymentRecord, type Person, type StatusChangeEvent } from "../lib/api";

interface Props {
  personId: string | null;
  onClose: () => void;
}

function Field({ label, value }: { label: string; value?: string | null }) {
  return (
    <div>
      <span className="block text-xs text-[#5A6578] dark:text-slate-400 mb-0.5">{label}</span>
      <p className="text-sm font-medium text-[#141E46] dark:text-slate-100">{value || "—"}</p>
    </div>
  );
}

export default function PersonDrawer({ personId, onClose }: Props) {
  const [person, setPerson] = useState<Person | null>(null);
  const [employment, setEmployment] = useState<EmploymentRecord | null>(null);
  const [contact, setContact] = useState<EmergencyContact | null>(null);
  const [events, setEvents] = useState<StatusChangeEvent[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const panelRef = useRef<HTMLDivElement>(null);
  const closeButtonRef = useRef<HTMLButtonElement>(null);

  // Reset and fetch whenever personId changes
  useEffect(() => {
    if (!personId) return;
    let cancelled = false;
    setPerson(null);
    setEmployment(null);
    setContact(null);
    setEvents([]);
    setError(null);
    setLoading(true);

    Promise.all([
      api.getPerson(personId),
      api.getCurrentEmployment(personId).catch(() => null),
      api.getEmergencyContact(personId).catch(() => null),
      api.listEvents(personId, { page_size: 20 }).then((r) => r.data).catch(() => []),
    ])
      .then(([p, emp, ec, evts]) => {
        if (cancelled) return;
        setPerson(p);
        setEmployment(emp);
        setContact(ec);
        setEvents(evts);
      })
      .catch((err: unknown) => {
        if (!cancelled) setError(err instanceof Error ? err.message : "Failed to load employee.");
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });

    return () => { cancelled = true; };
  }, [personId]);

  // Trap focus and handle Escape
  useEffect(() => {
    if (!personId) return;
    closeButtonRef.current?.focus();

    function onKey(e: KeyboardEvent) {
      if (e.key === "Escape") onClose();
    }
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [personId, onClose]);


  const isOpen = !!personId;

  return (
    <div
      ref={panelRef}
      role="complementary"
      aria-label="Employee details"
      className={`fixed inset-y-0 right-0 z-40 flex w-full max-w-md flex-col bg-[#FFF5E0] dark:bg-[#131519] shadow-[−4px_0_32px_rgba(0,0,0,0.12)] dark:shadow-[−4px_0_32px_rgba(0,0,0,0.5)] border-l border-[#E6DBC5] dark:border-[#2b303c] transition-transform duration-300 ease-in-out ${isOpen ? "translate-x-0" : "translate-x-full"}`}
    >
        {/* Header */}
        <div className="flex items-center justify-between border-b border-[#E6DBC5] dark:border-[#2b303c] px-5 py-4 bg-[#FFF5E0] dark:bg-[#1c1f26]">
          <h2 className="font-semibold text-[#141E46] dark:text-white text-base">
            Employee Details
          </h2>
          <button
            ref={closeButtonRef}
            type="button"
            onClick={onClose}
            aria-label="Close panel"
            className="rounded-md p-1.5 text-[#5A6578] dark:text-slate-400 hover:bg-[#E6DBC5]/60 dark:hover:bg-[#2b303c] hover:text-[#141E46] dark:hover:text-white transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-[#41B06E] dark:focus-visible:ring-[#1DCD9F]"
          >
            <svg viewBox="0 0 20 20" fill="currentColor" className="h-5 w-5" aria-hidden="true">
              <path d="M6.28 5.22a.75.75 0 0 0-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 1 0 1.06 1.06L10 11.06l3.72 3.72a.75.75 0 1 0 1.06-1.06L11.06 10l3.72-3.72a.75.75 0 0 0-1.06-1.06L10 8.94 6.28 5.22Z" />
            </svg>
          </button>
        </div>

        {/* Scrollable body */}
        <div className="flex-1 overflow-y-auto">
          {loading && (
            <div className="flex items-center justify-center py-16 text-sm text-[#5A6578] dark:text-slate-400">
              <svg className="mr-2 h-4 w-4 animate-spin text-[#41B06E] dark:text-[#1DCD9F]" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
              </svg>
              Loading…
            </div>
          )}

          {error && (
            <p className="m-5 rounded-md bg-red-500/10 border border-red-500/20 p-4 text-sm text-red-600 dark:text-red-400">
              {error}
            </p>
          )}

          {person && (
            <div className="divide-y divide-[#E6DBC5] dark:divide-[#2b303c]">

              {/* Identity */}
              <div className="px-5 py-5">
                <div className="flex items-center gap-4">
                  {person.profile_photo_url ? (
                    <img
                      src={person.profile_photo_url}
                      alt={`${person.first_name} ${person.last_name}`}
                      loading="lazy"
                      decoding="async"
                      width={56}
                      height={56}
                      className="h-14 w-14 rounded-full object-cover border border-[#E6DBC5] dark:border-[#2b303c] shrink-0"
                    />
                  ) : (
                    <div className="h-14 w-14 shrink-0 rounded-full bg-[#8DECB4]/30 dark:bg-[#1DCD9F]/20 flex items-center justify-center text-[#141E46] dark:text-[#1DCD9F] font-bold text-xl border border-[#41B06E]/30 dark:border-[#1DCD9F]/30">
                      {person.first_name.charAt(0)}{person.last_name.charAt(0)}
                    </div>
                  )}
                  <div className="min-w-0">
                    <h3 className="text-lg font-semibold text-[#141E46] dark:text-white leading-tight">
                      {person.first_name} {person.last_name}
                    </h3>
                    {person.preferred_name && (
                      <p className="text-xs text-[#5A6578] dark:text-slate-400">"{person.preferred_name}"</p>
                    )}
                    <p className="mt-0.5 text-sm text-[#5A6578] dark:text-slate-400">
                      {employment?.job_title || "—"} · {employment?.department || "—"}
                    </p>
                  </div>
                </div>

                {/* Status badges */}
                <div className="mt-3 flex flex-wrap gap-2">
                  {employment?.employment_status && (
                    <span className="rounded-full bg-[#8DECB4]/30 dark:bg-[#1DCD9F]/20 px-2.5 py-0.5 text-xs font-semibold text-[#141E46] dark:text-[#1DCD9F] border border-[#41B06E]/30 dark:border-[#1DCD9F]/40">
                      {employment.employment_status}
                    </span>
                  )}
                  {employment?.work_arrangement && (
                    <span className="rounded-full bg-[#F8EFE0] dark:bg-[#252a34] px-2.5 py-0.5 text-xs text-[#5A6578] dark:text-slate-300 border border-[#E6DBC5] dark:border-[#2b303c]">
                      {employment.work_arrangement}
                    </span>
                  )}
                  {!person.is_active && (
                    <span className="rounded-full bg-red-500/10 px-2.5 py-0.5 text-xs text-red-600 dark:text-red-400 border border-red-500/20">
                      Inactive
                    </span>
                  )}
                </div>
              </div>

              {/* Contact */}
              <div className="px-5 py-5">
                <h4 className="mb-3 text-xs font-semibold uppercase tracking-wider text-[#5A6578] dark:text-slate-500">Contact</h4>
                <div className="grid grid-cols-2 gap-3">
                  <Field label="Org email" value={person.org_email} />
                  <Field label="Personal email" value={person.personal_email} />
                  <Field label="Phone" value={person.phone_primary} />
                  <Field label="Location" value={[person.city, person.country].filter(Boolean).join(", ")} />
                </div>
              </div>

              {/* Employment */}
              <div className="px-5 py-5">
                <h4 className="mb-3 text-xs font-semibold uppercase tracking-wider text-[#5A6578] dark:text-slate-500">Employment</h4>
                <div className="grid grid-cols-2 gap-3">
                  <Field label="Job level" value={employment?.job_level} />
                  <Field label="Team" value={employment?.team} />
                  <Field label="Office" value={employment?.office_location} />
                  <Field label="Hire date" value={employment?.hire_date?.slice(0, 10)} />
                  <Field label="Employment type" value={employment?.employment_type} />
                </div>
              </div>

              {/* Tags */}
              {person.tags?.length > 0 && (
                <div className="px-5 py-5">
                  <h4 className="mb-3 text-xs font-semibold uppercase tracking-wider text-[#5A6578] dark:text-slate-500">Tags</h4>
                  <div className="flex flex-wrap gap-1.5">
                    {person.tags.map((t) => (
                      <span key={t} className="rounded-full bg-[#F8EFE0] dark:bg-[#252a34] px-2.5 py-0.5 text-xs text-[#5A6578] dark:text-slate-300 border border-[#E6DBC5] dark:border-[#2b303c]">
                        {t}
                      </span>
                    ))}
                  </div>
                </div>
              )}

              {/* Emergency contact */}
              <div className="px-5 py-5">
                <h4 className="mb-3 text-xs font-semibold uppercase tracking-wider text-[#5A6578] dark:text-slate-500">Emergency Contact</h4>
                {contact ? (
                  <div className="grid grid-cols-2 gap-3">
                    <Field label="Name" value={contact.name} />
                    <Field label="Relationship" value={contact.relationship} />
                    <Field label="Phone" value={contact.phone} />
                    <Field label="Email" value={contact.email} />
                  </div>
                ) : (
                  <p className="text-sm text-[#5A6578] dark:text-slate-400">No emergency contact recorded.</p>
                )}
              </div>

              {/* Timeline */}
              <div className="px-5 py-5">
                <h4 className="mb-3 text-xs font-semibold uppercase tracking-wider text-[#5A6578] dark:text-slate-500">Timeline</h4>
                {events.length > 0 ? (
                  <div className="space-y-1">
                    {events.map((ev) => (
                      <div key={ev.id} className="flex items-center justify-between rounded-md px-3 py-2.5 text-sm bg-white dark:bg-[#1c1f26] border border-[#E6DBC5] dark:border-[#2b303c] hover:bg-[#8DECB4]/10 dark:hover:bg-[#1DCD9F]/5 transition-colors">
                        <span className="font-medium text-[#141E46] dark:text-slate-200">{ev.event_type}</span>
                        <span className="text-xs font-mono text-[#5A6578] dark:text-slate-400 bg-[#F8EFE0] dark:bg-[#252a34] px-2 py-0.5 rounded">
                          {ev.effective_date.slice(0, 10)}
                        </span>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-sm text-[#5A6578] dark:text-slate-400">No events recorded.</p>
                )}
              </div>

            </div>
          )}
        </div>
    </div>
  );
}
