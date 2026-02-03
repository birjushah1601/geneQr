"use client";

import { EngineerSuggestion } from "@/lib/api/tickets";
import { Activity, Award } from "lucide-react";

interface EngineerCardProps {
  engineer: EngineerSuggestion;
  onAssign: () => void;
  isAssigning?: boolean;
}

export default function EngineerCard({ engineer, onAssign, isAssigning }: EngineerCardProps) {
  const getLevelBadgeColor = (level: number) => {
    switch (level) {
      case 3: return "bg-purple-100 text-purple-700 border-purple-200";
      case 2: return "bg-blue-100 text-blue-700 border-blue-200";
      default: return "bg-gray-100 text-gray-700 border-gray-200";
    }
  };

  const getLevelLabel = (level: number) => {
    switch (level) {
      case 3: return "Senior";
      case 2: return "Mid-Level";
      default: return "Junior";
    }
  };

  return (
    <div className="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-lg hover:border-blue-300 transition-all">
      {/* Header with Avatar and Name */}
      <div className="flex items-center gap-3 mb-3">
        <div className="w-12 h-12 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white font-bold text-base shadow-md">
          {engineer.name.split(' ').map(n => n[0]).join('').substring(0, 2)}
        </div>
        <div className="flex-1 min-w-0">
          <h3 className="font-semibold text-gray-900 truncate">{engineer.name}</h3>
          <span className={`inline-block px-2 py-0.5 rounded text-xs font-medium border ${getLevelBadgeColor(engineer.engineer_level)}`}>
            Level {engineer.engineer_level} - {getLevelLabel(engineer.engineer_level)}
          </span>
        </div>
        {engineer.match_score && (
          <div className="flex-shrink-0">
            <div className="text-2xl font-bold text-green-600">
              {engineer.match_score}%
            </div>
            <div className="text-xs text-gray-500 text-center">Match</div>
          </div>
        )}
      </div>

      {/* Workload and Availability */}
      {engineer.workload && (
        <div className="mb-3 flex items-center justify-between text-sm bg-gray-50 rounded-md p-2">
          <div className="flex items-center gap-2 text-gray-600">
            <Activity className="h-4 w-4 text-gray-400" />
            <span className="font-medium">Workload:</span>
          </div>
          <span className={`font-semibold ${engineer.workload.active_tickets === 0 ? 'text-green-600' : 'text-orange-600'}`}>
            {engineer.workload.active_tickets} active
          </span>
        </div>
      )}

      {/* Top Match Reason */}
      {engineer.match_reasons && engineer.match_reasons.length > 0 && (
        <div className="mb-3 text-sm text-gray-700 bg-blue-50 rounded-md p-2 border border-blue-100">
          <span className="font-medium text-blue-900">Ã¢Å“â€œ</span> {engineer.match_reasons[0]}
        </div>
      )}

      {/* Certifications Badge */}
      {engineer.certifications && engineer.certifications.filter(c => c.is_certified).length > 0 && (
        <div className="mb-3 flex items-center gap-2 text-xs">
          <Award className="h-4 w-4 text-blue-500" />
          <span className="text-gray-600">
            {engineer.certifications.filter(c => c.is_certified).length} certification{engineer.certifications.filter(c => c.is_certified).length !== 1 ? 's' : ''}
          </span>
        </div>
      )}

      {/* Assign Button */}
      <button
        onClick={onAssign}
        disabled={isAssigning}
        className="w-full py-2.5 px-4 bg-blue-600 text-white rounded-md text-sm font-semibold hover:bg-blue-700 active:bg-blue-800 disabled:bg-gray-300 disabled:cursor-not-allowed transition-all shadow-sm hover:shadow-md"
      >
        {isAssigning ? 'Assigning...' : 'Assign Engineer'}
      </button>
    </div>
  );
}
