'use client';

import { useState } from 'react';
import { ChevronDown, ChevronUp } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface CollapsibleSearchPanelProps {
  children: React.ReactNode;
  title?: string;
  defaultExpanded?: boolean;
}

export function CollapsibleSearchPanel({
  children,
  title = 'Search Filters',
  defaultExpanded = true,
}: CollapsibleSearchPanelProps) {
  const [isExpanded, setIsExpanded] = useState(defaultExpanded);

  return (
    <div className="bg-zinc-800/50 backdrop-blur-sm rounded-lg border border-zinc-700 shadow-lg">
      {/* Header */}
      <button
        onClick={() => setIsExpanded(!isExpanded)}
        className="w-full px-4 py-3 flex items-center justify-between hover:bg-zinc-700/30 transition-colors rounded-t-lg"
      >
        <h3 className="text-sm font-medium text-white">{title}</h3>
        <div className="text-zinc-400">
          {isExpanded ? (
            <ChevronUp className="w-4 h-4" />
          ) : (
            <ChevronDown className="w-4 h-4" />
          )}
        </div>
      </button>

      {/* Collapsible Content */}
      {isExpanded && (
        <div className="px-4 pb-4 border-t border-zinc-700">
          {children}
        </div>
      )}
    </div>
  );
}
