"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { ticketsApi } from "@/lib/api/tickets";
import type { ServiceTicket, TicketPriority, TicketStatus } from "@/types";
import { Loader2, Filter, Search, Ticket as TicketIcon, User, Package, Calendar } from "lucide-react";
import DashboardLayout from "@/components/DashboardLayout";

function StatusBadge({ status }: { status: TicketStatus }) {
  const color = {
    new: "bg-gray-100 text-gray-700",
    assigned: "bg-indigo-100 text-indigo-700",
    in_progress: "bg-blue-100 text-blue-700",
    on_hold: "bg-yellow-100 text-yellow-800",
    resolved: "bg-green-100 text-green-700",
    closed: "bg-gray-200 text-gray-800",
    cancelled: "bg-red-100 text-red-700",
  }[status];
  return <span className={`px-2 py-0.5 rounded text-xs font-medium ${color}`}>{status.replaceAll("_", " ")}</span>;
}

function PriorityBadge({ priority }: { priority: TicketPriority }) {
  const color = {
    critical: "bg-red-100 text-red-700",
    high: "bg-orange-100 text-orange-700",
    medium: "bg-amber-100 text-amber-700",
    low: "bg-green-100 text-green-700",
  }[priority];
  return <span className={`px-2 py-0.5 rounded text-xs font-medium ${color}`}>{priority}</span>;
}

export default function TicketsListPage() {
  const router = useRouter();
  const sp = useSearchParams();

  const [status, setStatus] = useState<TicketStatus | "">((sp.get("status") as TicketStatus) || "");
  const [priority, setPriority] = useState<TicketPriority | "">((sp.get("priority") as TicketPriority) || "");
  const [search, setSearch] = useState<string>(sp.get("q") || "");
  const [page, setPage] = useState<number>(Number(sp.get("page") || 1));

  // Sync URL when filters change
  useEffect(() => {
    const params = new URLSearchParams();
    if (status) params.set("status", status);
    if (priority) params.set("priority", priority);
    if (search) params.set("q", search);
    params.set("page", String(page));
    router.push(`/tickets?${params.toString()}`);
  }, [status, priority, search, page]);

  const queryParams = useMemo(() => {
    const p: Record<string, any> = { page, page_size: 20 };
    if (status) p.status = status;
    if (priority) p.priority = priority;
    // backend supports filtering by ticket number via exact match API; for now filter client-side
    return p;
  }, [status, priority, page]);

  const { data, isLoading, isFetching } = useQuery({
    queryKey: ["tickets", queryParams],
    queryFn: () => ticketsApi.list(queryParams as any),
    staleTime: 10_000,
  });

  const items: ServiceTicket[] = useMemo(() => {
    const all = (data as any)?.tickets || (data as any)?.items || [];
    if (!search) return all;
    const q = search.toLowerCase();
    return all.filter((t: ServiceTicket) =>
      t.ticket_number?.toLowerCase().includes(q) ||
      t.equipment_name?.toLowerCase().includes(q) ||
      t.customer_name?.toLowerCase().includes(q)
    );
  }, [data, search]);

  return (
    <DashboardLayout>
      <div className="bg-white border-b -mx-6 -mt-6 px-6 py-4 mb-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <TicketIcon className="h-6 w-6 text-gray-700" />
            <h1 className="text-xl font-semibold">Service Tickets</h1>
          </div>
          <Link href="/service-request" className="px-3 py-2 rounded bg-blue-600 text-white text-sm hover:bg-blue-700">Create Ticket</Link>
        </div>
      </div>

      <div>
        {/* Filters */}
        <div className="bg-white rounded-md border p-4 mb-4">
          <div className="flex flex-col md:flex-row gap-3 items-stretch md:items-center">
            <div className="flex-1 flex items-center gap-2 border rounded px-3 py-2">
              <Search className="h-4 w-4 text-gray-400" />
              <input
                value={search}
                onChange={(e) => { setSearch(e.target.value); setPage(1); }}
                placeholder="Search by ticket #, equipment, customer"
                className="w-full outline-none text-sm"
              />
            </div>
            <div className="flex items-center gap-2">
              <Filter className="h-4 w-4 text-gray-500" />
              <select value={status} onChange={(e) => { setStatus(e.target.value as any); setPage(1); }} className="border rounded px-2 py-2 text-sm">
                <option value="">All Status</option>
                <option value="new">New</option>
                <option value="assigned">Assigned</option>
                <option value="in_progress">In Progress</option>
                <option value="on_hold">On Hold</option>
                <option value="resolved">Resolved</option>
                <option value="closed">Closed</option>
                <option value="cancelled">Cancelled</option>
              </select>
              <select value={priority} onChange={(e) => { setPriority(e.target.value as any); setPage(1); }} className="border rounded px-2 py-2 text-sm">
                <option value="">All Priorities</option>
                <option value="critical">Critical</option>
                <option value="high">High</option>
                <option value="medium">Medium</option>
                <option value="low">Low</option>
              </select>
            </div>
          </div>
        </div>

        {/* Table */}
        <div className="bg-white rounded-md border overflow-x-auto">
          <table className="min-w-full text-sm">
            <thead className="bg-gray-50 text-left">
              <tr>
                <th className="px-4 py-3">Ticket #</th>
                <th className="px-4 py-3">Status</th>
                <th className="px-4 py-3">Priority</th>
                <th className="px-4 py-3">Equipment</th>
                <th className="px-4 py-3">Customer</th>
                <th className="px-4 py-3">Engineer</th>
                <th className="px-4 py-3">Created</th>
                <th className="px-4 py-3"></th>
              </tr>
            </thead>
            <tbody>
              {isLoading ? (
                <tr><td colSpan={8} className="px-4 py-10 text-center text-gray-500"><Loader2 className="inline h-5 w-5 animate-spin" /> Loading tickets...</td></tr>
              ) : items?.length ? (
                items.map((t) => (
                  <tr key={t.id} className="border-t hover:bg-gray-50">
                    <td className="px-4 py-3 font-medium"><Link href={`/tickets/${t.id}`} className="text-blue-600 hover:underline">{t.ticket_number}</Link></td>
                    <td className="px-4 py-3"><StatusBadge status={t.status} /></td>
                    <td className="px-4 py-3"><PriorityBadge priority={t.priority} /></td>
                    <td className="px-4 py-3"><div className="flex items-center gap-2"><Package className="h-4 w-4 text-gray-400" />{t.equipment_name}</div></td>
                    <td className="px-4 py-3">{t.customer_name}</td>
                    <td className="px-4 py-3"><div className="flex items-center gap-2"><User className="h-4 w-4 text-gray-400" />{t.assigned_engineer_name || "—"}</div></td>
                    <td className="px-4 py-3"><div className="flex items-center gap-2"><Calendar className="h-4 w-4 text-gray-400" />{new Date(t.created_at).toLocaleString()}</div></td>
                    <td className="px-4 py-3 text-right"><Link href={`/tickets/${t.id}`} className="text-sm text-blue-600 hover:underline">View</Link></td>
                  </tr>
                ))
              ) : (
                <tr><td colSpan={8} className="px-4 py-10 text-center text-gray-500">No tickets found</td></tr>
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        <div className="flex items-center justify-between mt-4">
          <p className="text-xs text-gray-500">{isFetching ? "Refreshing..." : `Total: ${data?.total ?? 0}`}</p>
          <div className="flex items-center gap-2">
            <button disabled={page <= 1} onClick={() => setPage((p) => Math.max(1, p - 1))} className="px-3 py-1.5 text-sm border rounded disabled:opacity-50">Prev</button>
            <span className="text-sm">Page {page}</span>
            <button disabled={(data?.items?.length ?? 0) < 20} onClick={() => setPage((p) => p + 1)} className="px-3 py-1.5 text-sm border rounded disabled:opacity-50">Next</button>
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
}
