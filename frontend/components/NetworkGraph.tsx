'use client';
import { useMemo, useState } from 'react';

export interface GraphEdge {
  relationship_id: string;
  source_type: string;
  source_id: string;
  target_type: string;
  target_id: string;
  relationship_type: string;
  weight: number;
  fraud_linked: boolean;
}

interface LaidOutNode {
  id: string;
  type: string;
  fraudLinked: boolean;
  x: number;
  y: number;
}

const TYPE_COLOR: Record<string, { fill: string; stroke: string; text: string }> = {
  USER: { fill: '#e8f1fc', stroke: '#0071e3', text: '#0071e3' },
  DEVICE: { fill: '#f3e8ff', stroke: '#7c3aed', text: '#6d28d9' },
  IP: { fill: '#fef3c7', stroke: '#d97706', text: '#b45309' },
  MERCHANT: { fill: '#dcfce7', stroke: '#16a34a', text: '#15803d' },
};

function layout(edges: GraphEdge[], width: number, height: number): LaidOutNode[] {
  const map = new Map<string, LaidOutNode>();
  const add = (id: string, type: string, fraud: boolean) => {
    const existing = map.get(id);
    if (!existing) {
      map.set(id, { id, type, fraudLinked: fraud, x: 0, y: 0 });
    } else if (fraud) {
      existing.fraudLinked = true;
    }
  };
  edges.forEach(e => {
    add(e.source_id, e.source_type, e.fraud_linked);
    add(e.target_id, e.target_type, e.fraud_linked);
  });

  const byType: Record<string, LaidOutNode[]> = { USER: [], DEVICE: [], IP: [], MERCHANT: [] };
  map.forEach(n => {
    (byType[n.type] ?? byType.USER).push(n);
  });

  const columns = ['USER', 'DEVICE', 'IP', 'MERCHANT'];
  const padX = 70;
  const padY = 36;
  const colW = (width - padX * 2) / Math.max(columns.length - 1, 1);

  columns.forEach((type, col) => {
    const nodes = byType[type];
    const n = nodes.length;
    nodes.forEach((node, i) => {
      node.x = padX + col * colW;
      node.y = n <= 1 ? height / 2 : padY + (i * (height - padY * 2)) / Math.max(n - 1, 1);
    });
  });

  return Array.from(map.values());
}

interface Props {
  edges: GraphEdge[];
  height?: number;
  onSelectNode?: (id: string) => void;
  selectedId?: string | null;
}

export default function NetworkGraph({ edges, height = 320, onSelectNode, selectedId }: Props) {
	const [hoverId, setHoverId] = useState<string | null>(null);
  const width = 860;

  const nodes = useMemo(() => layout(edges, width, height), [edges, height]);
  const nodeById = useMemo(() => {
    const m = new Map<string, LaidOutNode>();
    nodes.forEach(n => m.set(n.id, n));
    return m;
  }, [nodes]);

	
	
	const activeId = selectedId ?? hoverId ?? edges[0]?.source_id ?? null;
	const visibleEdges = useMemo(() => edges.filter(e => e.source_id === activeId || e.target_id === activeId), [edges, activeId]);
	const visibleNodeIDs = useMemo(() => new Set(visibleEdges.flatMap(e => [e.source_id, e.target_id])), [visibleEdges]);
	const visibleNodes = useMemo(() => nodes.filter(n => visibleNodeIDs.has(n.id)), [nodes, visibleNodeIDs]);

  if (edges.length === 0) {
    return (
      <p className="text-sm text-center text-[#a1a1a6] py-8">
        No relationships extracted yet. Start a device_farm or ip_abuse scenario to generate graph data.
      </p>
    );
  }

  return (
    <div className="w-full overflow-x-auto">
      <svg viewBox={`0 0 ${width} ${height}`} className="w-full h-auto" role="img" aria-label="Fraud network graph">
		{visibleEdges.map(e => {
          const s = nodeById.get(e.source_id);
          const t = nodeById.get(e.target_id);
          if (!s || !t) return null;
          const lit = !activeId || e.source_id === activeId || e.target_id === activeId;
          return (
            <line
              key={e.relationship_id}
              x1={s.x} y1={s.y} x2={t.x} y2={t.y}
              stroke={e.fraud_linked ? '#ef4444' : '#c7c7cc'}
              strokeWidth={e.fraud_linked ? 2 : 1.25}
              strokeOpacity={lit ? 0.9 : 0.15}
            />
          );
        })}
		{visibleNodes.map(n => {
          const colors = TYPE_COLOR[n.type] ?? TYPE_COLOR.USER;
          const selected = n.id === activeId;
			return (
            <g
              key={n.id}
              transform={`translate(${n.x}, ${n.y})`}
              className="cursor-pointer"
              onClick={() => onSelectNode?.(n.id)}
              onMouseEnter={() => setHoverId(n.id)}
              onMouseLeave={() => setHoverId(null)}
			>
              <circle
                r={selected ? 22 : 18}
                fill={n.fraudLinked ? '#fef2f2' : colors.fill}
                stroke={n.fraudLinked ? '#dc2626' : colors.stroke}
                strokeWidth={selected ? 3 : 2}
              />
              <text textAnchor="middle" y={4} fontSize={9} fontWeight={600} fill={n.fraudLinked ? '#dc2626' : colors.text}>
                {n.id}
              </text>
            </g>
          );
        })}
		</svg>
		<p className="mt-2 text-center text-xs text-[#6e6e73]">Showing active entity <span className="font-mono font-medium text-[#1d1d1f]">{activeId}</span> and its direct relationships. Select a connected node to focus it.</p>
      <div className="flex flex-wrap gap-3 justify-center mt-2 text-[11px] text-[#6e6e73]">
        {Object.entries(TYPE_COLOR).map(([type, c]) => (
          <span key={type} className="flex items-center gap-1.5">
            <span className="h-2.5 w-2.5 rounded-full border" style={{ background: c.fill, borderColor: c.stroke }} />
            {type}
          </span>
        ))}
        <span className="flex items-center gap-1.5">
          <span className="h-2.5 w-2.5 rounded-full bg-red-50 border border-red-500" />
          Fraud-linked
        </span>
      </div>
    </div>
  );
}
