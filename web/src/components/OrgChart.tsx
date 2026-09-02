import React, { useEffect, useState } from "react";
import { api, type OrgChartNode } from "../lib/api";

interface TreeNode extends OrgChartNode {
  children: TreeNode[];
}

function buildTree(nodes: OrgChartNode[]): TreeNode[] {
  const map = new Map<string, TreeNode>();
  const roots: TreeNode[] = [];

  nodes.forEach((n) => {
    map.set(n.id, { ...n, children: [] });
  });

  nodes.forEach((n) => {
    const node = map.get(n.id)!;
    if (n.reports_to && map.has(n.reports_to)) {
      map.get(n.reports_to)!.children.push(node);
    } else {
      roots.push(node);
    }
  });

  return roots;
}

const NodeComponent: React.FC<{ node: TreeNode }> = ({ node }) => {
  const [expanded, setExpanded] = useState(true);

  return (
    <div className="flex flex-col items-center">
      <div className="relative border border-emerald-100 rounded-lg bg-white p-4 shadow-sm w-56 text-center m-2 flex flex-col items-center hover:border-emerald-300 transition-colors cursor-default">
        {node.profile_photo_url ? (
            <img src={node.profile_photo_url} alt={node.name} className="w-12 h-12 rounded-full mb-2 object-cover border border-zinc-200" width="48" height="48" />
        ) : (
            <div className="w-12 h-12 rounded-full mb-2 bg-emerald-50 flex items-center justify-center text-emerald-700 font-medium border border-emerald-100">
                {node.name.charAt(0)}
            </div>
        )}
        <h3 className="font-semibold text-sm text-zinc-900 truncate w-full">{node.name}</h3>
        <p className="text-xs text-zinc-600 truncate w-full mt-0.5">{node.job_title || "No Title"}</p>
        {node.department && (
          <span className="mt-2 inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-zinc-100 text-zinc-800">
            {node.department}
          </span>
        )}

        {node.children.length > 0 && (
          <button
            onClick={() => setExpanded(!expanded)}
            aria-expanded={expanded}
            aria-label={`Toggle children for ${node.name}`}
            className="absolute -bottom-3 bg-white border border-emerald-200 text-emerald-700 rounded-full w-6 h-6 flex items-center justify-center hover:bg-emerald-50 hover:border-emerald-300 transition-colors z-10 text-xs shadow-sm"
          >
            {expanded ? "−" : "+"}
          </button>
        )}
      </div>

      {expanded && node.children.length > 0 && (
        <div className="flex flex-row relative mt-4 pt-4 border-t-2 border-emerald-100/50">
           {/* Connecting vertical line down to children */}
          <div className="absolute top-0 left-1/2 w-0.5 h-4 bg-emerald-100/50 -translate-x-1/2"></div>
          {node.children.map((child, _index) => (
            <div key={child.id} className="relative flex flex-col items-center px-2">
                {/* Connecting horizontal line between children */}
                <div className="absolute top-0 w-full h-0.5 bg-emerald-100/50 z-0"></div>
                {/* Connecting vertical line to each child */}
                <div className="absolute top-0 w-0.5 h-4 bg-emerald-100/50 z-0"></div>
              <div className="mt-4">
                  <NodeComponent node={child} />
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default function OrgChart() {
  const [data, setData] = useState<TreeNode[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchChart() {
      try {
        const response = await api.getOrgChart();
        const tree = buildTree(response);
        setData(tree);
      } catch (err: any) {
        setError(err.message || "Failed to load org chart.");
      } finally {
        setLoading(false);
      }
    }
    fetchChart();
  }, []);

  if (loading) {
    return <div className="text-center py-10 text-zinc-500">Loading Org Chart...</div>;
  }

  if (error) {
    return (
        <div className="rounded bg-red-50 p-4 text-sm text-red-700 max-w-2xl mx-auto mt-4">
            {error} — ensure you are authenticated.
        </div>
    );
  }

  if (data.length === 0) {
    return <div className="text-center py-10 text-zinc-500">No organizational data available.</div>;
  }

  return (
    <div className="overflow-auto p-8 w-full h-full min-h-[600px] flex justify-center bg-zinc-50/50 rounded-xl border border-zinc-200 pb-24 shadow-inner">
        {data.map(root => (
            <NodeComponent key={root.id} node={root} />
        ))}
    </div>
  );
}
